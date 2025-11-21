package types

import (
	"runtime"
	"time"

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
	Environment        Environment    `envconfig:"ENVIRONMENT"     default:"development"`
	LogLevel           string         `envconfig:"LOG_LEVEL"       default:"info"`
	LogFormat          string         `envconfig:"LOG_FORMAT"      default:"json"`
	Port               int            `envconfig:"PORT"            default:"8080"`
	GatewayPort        int            `envconfig:"GATEWAY_PORT"    default:"8081"`
	GatewayHost        string         `envconfig:"GATEWAY_HOST"    default:"localhost"`
	Identifier         string         `envconfig:"IDENTIFIER"      default:"device-001"`
	SupportedProtocols []Protocol     `envconfig:"PROTOCOLS"       default:"http,grpc"`
	DeviceVersions     DeviceVersions `envconfig:"DEVICE_VERSION"`
	StreamInterval     time.Duration  `envconfig:"STREAM_INTERVAL" default:"500ms"`
}

type DeviceVersions struct {
	Hardware string `envconfig:"HARDWARE" default:"hw-1.0.0"`
	Software string `envconfig:"SOFTWARE" default:"sw-1.0.0"`
	Firmware string `envconfig:"FIRMWARE" default:"fw-1.0.0"`
}

type DeviceState struct {
	Identifier         string
	SupportedProtocols []Protocol
	DeviceVersions     DeviceVersions
	Architecture       string
	OS                 string
	DeviceStatus       DeviceStatus
	StreamInterval     time.Duration
	Updated            time.Time
}

type HealthStatus struct {
	Identifier         string
	SupportedProtocols []Protocol
	Architecture       string
	OS                 string
	Updated            time.Time
}

type Diagnostics struct {
	Identifier     string
	DeviceVersions DeviceVersions
	CPU            float64
	Memory         float64
	DeviceStatus   DeviceStatus
	Checksum       string
}

type DeviceMutation struct {
	DeviceStatus DeviceStatus
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

func (d *DeviceStatus) Proto() devicev1.DeviceStatus {
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

// ENUM(http, grpc)
type Protocol string

func (p *Protocol) Decode(value string) error {
	parsed, err := ParseProtocol(value)
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}

func (p *Protocol) Proto() devicev1.Protocol {
	switch *p {
	case ProtocolHttp:
		return devicev1.Protocol_PROTOCOL_HTTP
	case ProtocolGrpc:
		return devicev1.Protocol_PROTOCOL_GRPC
	default:
		return devicev1.Protocol_PROTOCOL_UNSPECIFIED
	}
}
