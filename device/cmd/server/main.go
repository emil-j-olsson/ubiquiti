package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emil-j-olsson/ubiquiti/device/internal/cache"
	"github.com/emil-j-olsson/ubiquiti/device/internal/server"
	"github.com/emil-j-olsson/ubiquiti/device/internal/service"
	"github.com/emil-j-olsson/ubiquiti/device/internal/types"
	"github.com/emil-j-olsson/ubiquiti/device/logging"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultDevicePrefix    = "DEVICE"
	DefaultShutdownTimeout = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		zap.L().Fatal("application startup error", zap.Error(err))
		os.Exit(1)
	}
}

func run() error {
	// Configuration & Logging
	_ = godotenv.Load()

	var config types.Config
	if err := envconfig.Process(DefaultDevicePrefix, &config); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logger, close := logging.NewProductionLogger(config.LogLevel, config.LogFormat)
	defer close()
	zap.ReplaceGlobals(logger)

	logger.Info("starting device service",
		zap.String("version", types.Version),
		zap.String("revision", types.Revision),
		zap.String("branch", types.Branch),
		zap.String("build_date", types.BuildDate),
		zap.String("go_version", types.GoVersion),
	)

	// Context & Signal Handling
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Application Layer
	deviceState := cache.NewDeviceState(config)
	deviceService := service.NewDeviceService(deviceState)
	deviceServer := server.NewDeviceServer(logger, deviceService)

	// Server Lifecycle
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return startServer(gctx, config, deviceServer, logger)
	})
	g.Go(func() error {
		return startGateway(gctx, config, logger)
	})
	deviceState.UpdateState(func(ds *types.DeviceState) {
		ds.DeviceStatus = types.DeviceStatusHealthy
	})
	return g.Wait()
}

func startServer(ctx context.Context, config types.Config, srv *server.Server, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return fmt.Errorf("failed to register listener on port %d: %w", config.Port, err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			server.UnaryLoggingInterceptor(logger),
			server.UnaryRecoveryInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			server.StreamLoggingInterceptor(logger),
			server.StreamRecoveryInterceptor(logger),
		),
	)
	devicev1.RegisterDeviceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	logger.Info("server started (grpc)", zap.Int("port", config.Port))

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server (grpc)")
		grpcServer.GracefulStop()
	}()
	return grpcServer.Serve(lis)
}

func startGateway(ctx context.Context, config types.Config, logger *zap.Logger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("%s:%d", config.GatewayHost, config.Port)

	if err := devicev1.RegisterDeviceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register handler (gateway): %w", err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GatewayPort),
		Handler: mux,
	}
	logger.Info("server started (gateway)", zap.Int("port", config.GatewayPort))

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server (gateway)")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown error (gateway)", zap.Error(err))
		}
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("startup error (gateway): %w", err)
	}
	return nil
}
