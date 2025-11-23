package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewDefaultLogger() (*zap.Logger, func()) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger, func() { _ = logger.Sync() }
}

func NewProductionLogger(level, format string) (*zap.Logger, func()) {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = zapcore.InfoLevel
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch format {
	case "json":
		encoder = zapcore.NewJSONEncoder(config)
	case "console":
		encoder = zapcore.NewConsoleEncoder(config)
	default:
		encoder = zapcore.NewConsoleEncoder(config)
	}
	logger := zap.New(
		zapcore.NewCore(encoder, os.Stdout, lvl),
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCaller(),
	)
	return logger, func() { _ = logger.Sync() }
}
