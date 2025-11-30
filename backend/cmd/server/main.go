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

	"github.com/emil-j-olsson/ubiquiti/backend/internal/checksum"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/database"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/database/postgres"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/device"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/logging"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/server"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/service"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/worker"
	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
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
	DefaultMonitorPrefix   = "MONITOR"
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
	if err := envconfig.Process(DefaultMonitorPrefix, &config); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logger, close := logging.NewProductionLogger(config.LogLevel, config.LogFormat)
	defer close()
	zap.ReplaceGlobals(logger)

	logger.Info("starting monitor service",
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

	// Checksum
	generator := checksum.NewGenerator(config.ChecksumBinaryPath)

	// Persistence Layer
	model, err := database.NewDatabaseModel(ctx, logger).AddDatabaseConnections(
		database.Config{
			Database: types.DatabasePostgres,
			Instance: types.DatabaseInstanceUbiquiti,
			PostgresOptions: []postgres.Option{
				postgres.WithProxyConnection(
					config.Persistence.Postgres.ConnectionString,
					config.Persistence.Postgres.MaxPoolSize,
				),
			},
		},
	)
	if err != nil {
		return err
	}
	defer model.Close()
	pool, err := model.GetPostgresPool(types.DatabaseInstanceUbiquiti)
	if err != nil {
		return err
	}
	persistence := postgres.NewPersistenceRepository(pool, logger)
	notifier := postgres.NewNotifier(pool, config.Persistence.Postgres.NotificationChannel, logger)

	// External Clients
	factory := device.NewClientFactory(generator)

	// Application Layer
	monitorService := service.NewMonitorService(persistence, factory, config, logger)
	monitorServer := server.NewMonitorServer(monitorService, logger)

	// Worker Lifecycle
	interval := config.StreamInterval
	orchestrator := worker.NewOrchestrator(persistence, notifier, factory, interval, logger)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return notifier.Listen(gctx)
	})
	g.Go(func() error {
		return orchestrator.Run(gctx)
	})

	// Server Lifecycle
	g.Go(func() error {
		return startServer(gctx, config, monitorServer, logger)
	})
	g.Go(func() error {
		return startGateway(gctx, config, logger)
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
	monitorv1.RegisterMonitorServer(grpcServer, srv)
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

	if err := monitorv1.RegisterMonitorHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register handler (gateway): %w", err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GatewayPort),
		Handler: server.CORSMiddleware(mux),
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
