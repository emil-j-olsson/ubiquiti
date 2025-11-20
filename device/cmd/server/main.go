package main

import (
	"github.com/emil-j-olsson/ubiquiti/device/internal/cache"
	"github.com/emil-j-olsson/ubiquiti/device/internal/types"
	"github.com/emil-j-olsson/ubiquiti/device/logging"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

const (
	defaultDevicePrefix string = "DEVICE"
)

func main() {
	// logger... encapsulate zap and clean up startup process!
	defaultLogger, defaultClose := logging.NewDefaultLogger()
	if err := godotenv.Load(); err == nil {
		defaultLogger.Info("Successfully loaded environment variables from .env file")
	}
	config := types.Config{}
	if err := envconfig.Process(defaultDevicePrefix, &config); err != nil {
		defer defaultClose()
		defaultLogger.Fatal("Failed to process config from environment variables", zap.Error(err))
	}
	defaultClose()

	logger, close := logging.NewProductionLogger(config.LogLevel, config.LogFormat)
	defer close()
	zap.ReplaceGlobals(logger)

	logger.Info("Starting device service",
		zap.String("version", types.Version),
		zap.String("revision", types.Revision),
		zap.String("branch", types.Branch),
		zap.String("build_date", types.BuildDate),
		zap.String("go_version", types.GoVersion),
	)

	// mapping + validation of config values...

	// keep all clients outside of business logic!
	// depending on this state... enable/disable grpc or http requests
	cache.NewDeviceState(config)

	// start api servers based on config flags...
	// !! check grpc examples (repos)
	// !! create the small executable for checksum

	// hexagonal...
	// input validation
	// try out shared package later...
	// should support various protocols (activate a certain protocol depending on environment variable)
	// optimize: custom protocol?
	// optimize: kafka topic/event driven?
	// use various image bases, ubuntu vX and MacOS? Set the operating system and versions as the result of the `/health` endpoint
	// consider using `shared` code with the connection adapters
	// simulate disruptions via feature toggling / API call to device...
	// API versioning
	// document everything!!!
}
