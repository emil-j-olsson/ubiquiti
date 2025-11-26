package fixtures

type Config struct {
	Environment string      `envconfig:"ENVIRONMENT"  default:"test"`
	Port        int         `envconfig:"PORT"         default:"8080"`
	Host        string      `envconfig:"HOST"         default:"localhost"`
	GatewayPort int         `envconfig:"GATEWAY_PORT" default:"8081"`
	Persistence Persistence `envconfig:"PERSISTENCE"`
}

type Persistence struct {
	Postgres Postgres `envconfig:"POSTGRES"`
}

type Postgres struct {
	ConnectionString    string `envconfig:"CONNECTION_STRING"`
	MaxPoolSize         int32  `envconfig:"MAX_POOL_SIZE"        default:"25"`
	NotificationChannel string `envconfig:"NOTIFICATION_CHANNEL" default:"device_changes"`
}

type ServiceConfig struct {
	Container          string
	Identifier         string
	Port               int
	GatewayPort        int
	SupportedProtocols []Protocol
	Architecture       string
	OS                 string
}

var Services = map[Service]ServiceConfig{
	ServiceDeviceRouter: {
		Container:          "ubiquiti-device-router",
		Identifier:         "ubiquiti-device-router-3c2d",
		Port:               8084,
		GatewayPort:        8085,
		SupportedProtocols: []Protocol{ProtocolGrpc},
		Architecture:       "arm64",
		OS:                 "linux",
	},
	ServiceDeviceSwitch: {
		Container:          "ubiquiti-device-switch",
		Identifier:         "ubiquiti-device-switch-b87f",
		Port:               8086,
		GatewayPort:        8087,
		SupportedProtocols: []Protocol{ProtocolGrpcStream},
		Architecture:       "amd64",
		OS:                 "linux",
	},
	ServiceDeviceAccessPoint: {
		Container:          "ubiquiti-device-access-point",
		Identifier:         "ubiquiti-device-access-point-05da",
		Port:               8088,
		GatewayPort:        8089,
		SupportedProtocols: []Protocol{ProtocolHttp},
		Architecture:       "arm",
		OS:                 "linux",
	},
	ServiceBackendMonitorArm: {
		Container:          "ubiquiti-monitor-arm",
		Identifier:         "ubiquiti-monitor-arm",
		Port:               8080,
		GatewayPort:        8081,
		SupportedProtocols: []Protocol{ProtocolGrpc, ProtocolHttp, ProtocolGrpcStream, ProtocolHttpStream},
		Architecture:       "arm64",
		OS:                 "linux",
	},
	ServiceBackendMonitorAmd: {
		Container:          "ubiquiti-monitor-amd",
		Identifier:         "ubiquiti-monitor-amd",
		Port:               8082,
		GatewayPort:        8083,
		SupportedProtocols: []Protocol{ProtocolGrpc, ProtocolHttp, ProtocolGrpcStream, ProtocolHttpStream},
		Architecture:       "amd64",
		OS:                 "linux",
	},
	ServiceInvalid: {
		Container:   "ubiquiti-invalid",
		Identifier:  "ubiquiti-invalid",
		Port:        8090,
		GatewayPort: 8091,
	},
}

//go:generate go-enum

// ENUM(device-router, device-switch, device-access-point, backend-monitor-arm, backend-monitor-amd, invalid)
type Service string

/*
ENUM(

	http = PROTOCOL_HTTP
	http-stream = PROTOCOL_HTTP_STREAM
	grpc = PROTOCOL_GRPC
	grpc-stream = PROTOCOL_GRPC_STREAM

)
*/
type Protocol string

func (p *Protocol) IsGrpc() bool {
	return *p == ProtocolGrpc || *p == ProtocolGrpcStream
}

func (p *Protocol) IsHttp() bool {
	return *p == ProtocolHttp || *p == ProtocolHttpStream
}
