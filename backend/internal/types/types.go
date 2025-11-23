package types

import (
	"runtime"
	"time"
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
	Persistence    Persistence
}

type Persistence struct {
	Postgres Postgres `envconfig:"POSTGRES"`
}

type Postgres struct {
	ConnectionString string `envconfig:"CONNECTION_STRING"`
	MaxPoolSize      int32  `envconfig:"MAX_POOL_SIZE"     default:"10"`
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
