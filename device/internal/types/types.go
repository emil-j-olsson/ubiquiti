package types

import "runtime"

// potentially more elsewhere
var (
	Version   = "unset"
	Revision  = "unset"
	Branch    = "unset"
	BuildDate = "unset"
	GoVersion = runtime.Version()
)

type Config struct {
	Environment        Environment    `envconfig:"ENVIRONMENT" required:"true"`
	LogLevel           string         `envconfig:"LOG_LEVEL" default:"info"`
	LogFormat          string         `envconfig:"LOG_FORMAT" default:"json"`
	Port               int            `envconfig:"PORT" default:"8080"`
	SystemPort         int            `envconfig:"SYSTEM_PORT" default:"8081"` // health
	Identifier         string         `envconfig:"IDENTIFIER" default:"device-001"`
	SupportedProtocols []Protocol     `envconfig:"PROTOCOLS" default:"http,grpc"`
	DeviceVersions     DeviceVersions `envconfig:"DEVICE_VERSION"`
}

// Authentication

// TODO: set these via env in docker
type DeviceVersions struct {
	Hardware string `envconfig:"HARDWARE" default:"hw-1.0.0"`
	Software string `envconfig:"SOFTWARE" default:"sw-1.0.0"`
	Firmware string `envconfig:"FIRMWARE" default:"fw-1.0.0"`
}

type DeviceState struct {
	Identifier     string
	DeviceVersions DeviceVersions
	Architecture   string
	OS             string
	DeviceStatus   DeviceStatus
}

type HealthStatus struct {
	Identifier         string
	SupportedProtocols []Protocol
}

type Diagnostics struct {
	Identifier     string
	DeviceVersions DeviceVersions
	CPU            float64
	Memory         float64
	DeviceStatus   DeviceStatus
	Checksum       string
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

// ENUM(healthy, degraded, error, maintenance, booting)
type DeviceStatus string

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
