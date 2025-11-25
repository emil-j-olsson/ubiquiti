package types

import (
	"runtime"
	"time"

	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
)

var (
	Version   = "unset"
	Revision  = "unset"
	Branch    = "unset"
	BuildDate = "unset"
	GoVersion = runtime.Version()
)

type Config struct {
	Environment    Environment   `envconfig:"ENVIRONMENT"     default:"development"`
	LogLevel       string        `envconfig:"LOG_LEVEL"       default:"info"`
	LogFormat      string        `envconfig:"LOG_FORMAT"      default:"json"`
	Port           int           `envconfig:"PORT"            default:"8080"`
	GatewayPort    int           `envconfig:"GATEWAY_PORT"    default:"8081"`
	GatewayHost    string        `envconfig:"GATEWAY_HOST"    default:"localhost"`
	Identifier     string        `envconfig:"IDENTIFIER"      default:"monitor-001"`
	StreamInterval time.Duration `envconfig:"STREAM_INTERVAL" default:"500ms"`
	Persistence    Persistence   `envconfig:"PERSISTENCE"`
}

type Persistence struct {
	Postgres Postgres `envconfig:"POSTGRES"`
}

type Postgres struct {
	ConnectionString string `envconfig:"CONNECTION_STRING"`
	MaxPoolSize      int32  `envconfig:"MAX_POOL_SIZE"     default:"10"`
}

type Device struct {
	ID                 *string    `db:"id"`
	Identifier         *string    `db:"device_id"`
	Alias              *string    `db:"alias"`
	Host               *string    `db:"host"`
	Port               *int64     `db:"port"`
	GatewayPort        *int64     `db:"port_gateway"`
	Architecture       *string    `db:"architecture"`
	OS                 *string    `db:"os"`
	SupportedProtocols *[]string  `db:"supported_protocols"`
	Created            *time.Time `db:"created_at"`
	Updated            *time.Time `db:"updated_at"`
}

type Diagnostics struct {
	ID                 *string    `db:"id"`
	Identifier         *string    `db:"device_id"`
	Alias              *string    `db:"alias"`
	Host               *string    `db:"host"`
	Port               *int64     `db:"port"`
	GatewayPort        *int64     `db:"port_gateway"`
	Architecture       *string    `db:"architecture"`
	OS                 *string    `db:"os"`
	SupportedProtocols *[]string  `db:"supported_protocols"`
	Hardware           *string    `db:"hardware_version"`
	Software           *string    `db:"software_version"`
	Firmware           *string    `db:"firmware_version"`
	CPU                *float64   `db:"cpu_usage"`
	Memory             *float64   `db:"memory_usage"`
	DeviceStatus       *string    `db:"device_status"`
	Checksum           *string    `db:"checksum"`
	LastUpdated        *time.Time `db:"last_updated"`
	Created            *time.Time `db:"created_at"`
	Updated            *time.Time `db:"updated_at"`
}

type DeviceHealthStatus struct {
	Identifier         string
	SupportedProtocols []Protocol
	Architecture       string
	OS                 string
	Updated            time.Time
}

type DeviceDiagnostics struct {
	Identifier     string
	DeviceVersions DeviceVersions
	CPU            float64
	Memory         float64
	DeviceStatus   DeviceStatus
	Checksum       string
	Timestamp      time.Time
}

type DeviceVersions struct {
	Hardware string
	Software string
	Firmware string
}

type DeviceRegistration struct {
	Protocol    Protocol
	Alias       string
	Host        string
	Port        int64
	GatewayPort int64
}

//go:generate go-enum

// ENUM(test, development, staging, production)
type Environment string

func (e *Environment) Decode(value string) error {
	parsed, err := ParseEnvironment(value)
	if err != nil {
		return err
	}
	*e = parsed
	return nil
}

// ENUM(proxy)
type PostgresConnection string

// ENUM(postgres)
type Database string

// ENUM(ubiquiti)
type DatabaseInstance string

/*
ENUM(

	healthy = DEVICE_STATUS_HEALTHY
	degraded = DEVICE_STATUS_DEGRADED
	error = DEVICE_STATUS_ERROR
	maintenance = DEVICE_STATUS_MAINTENANCE
	booting = DEVICE_STATUS_BOOTING

)
*/
type DeviceStatus string

func (d *DeviceStatus) Proto() monitorv1.DeviceStatus {
	switch *d {
	case DeviceStatusHealthy:
		return monitorv1.DeviceStatus_DEVICE_STATUS_HEALTHY
	case DeviceStatusDegraded:
		return monitorv1.DeviceStatus_DEVICE_STATUS_DEGRADED
	case DeviceStatusError:
		return monitorv1.DeviceStatus_DEVICE_STATUS_ERROR
	case DeviceStatusMaintenance:
		return monitorv1.DeviceStatus_DEVICE_STATUS_MAINTENANCE
	case DeviceStatusBooting:
		return monitorv1.DeviceStatus_DEVICE_STATUS_BOOTING
	default:
		return monitorv1.DeviceStatus_DEVICE_STATUS_UNSPECIFIED
	}
}

func (d *DeviceStatus) DeviceProto() devicev1.DeviceStatus {
	switch *d {
	case DeviceStatusHealthy:
		return devicev1.DeviceStatus_DEVICE_STATUS_HEALTHY
	case DeviceStatusDegraded:
		return devicev1.DeviceStatus_DEVICE_STATUS_DEGRADED
	case DeviceStatusError:
		return devicev1.DeviceStatus_DEVICE_STATUS_ERROR
	case DeviceStatusMaintenance:
		return devicev1.DeviceStatus_DEVICE_STATUS_MAINTENANCE
	case DeviceStatusBooting:
		return devicev1.DeviceStatus_DEVICE_STATUS_BOOTING
	default:
		return devicev1.DeviceStatus_DEVICE_STATUS_UNSPECIFIED
	}
}

func DeviceStatusFromString(value string) DeviceStatus {
	parsed, err := ParseDeviceStatus(value)
	if err != nil {
		return DeviceStatus("")
	}
	return parsed
}

/*
ENUM(

	http = PROTOCOL_HTTP
	http-stream = PROTOCOL_HTTP_STREAM
	grpc = PROTOCOL_GRPC
	grpc-stream = PROTOCOL_GRPC_STREAM

)
*/
type Protocol string

func (p *Protocol) Proto() monitorv1.Protocol {
	switch *p {
	case ProtocolHttp:
		return monitorv1.Protocol_PROTOCOL_HTTP
	case ProtocolHttpStream:
		return monitorv1.Protocol_PROTOCOL_HTTP_STREAM
	case ProtocolGrpc:
		return monitorv1.Protocol_PROTOCOL_GRPC
	case ProtocolGrpcStream:
		return monitorv1.Protocol_PROTOCOL_GRPC_STREAM
	default:
		return monitorv1.Protocol_PROTOCOL_UNSPECIFIED
	}
}

func (p *Protocol) IsGrpc() bool {
	return *p == ProtocolGrpc || *p == ProtocolGrpcStream
}

func (p *Protocol) IsHttp() bool {
	return *p == ProtocolHttp || *p == ProtocolHttpStream
}

func ProtocolFromStrings(values []string) []monitorv1.Protocol {
	result := make([]monitorv1.Protocol, 0, len(values))
	for _, value := range values {
		if p, err := ParseProtocol(value); err == nil {
			result = append(result, p.Proto())
		}
	}
	return result
}

func ProtocolFromDevice(values []devicev1.Protocol) []Protocol {
	result := make([]Protocol, 0, len(values))
	for _, value := range values {
		if p, err := ParseProtocol(value.String()); err == nil {
			result = append(result, p)
		}
	}
	return result
}
