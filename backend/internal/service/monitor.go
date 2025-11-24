package service

import (
	"context"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/device"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"go.uber.org/zap"
)

type PersistenceProvider interface {
	RegisterDevice(
		ctx context.Context,
		status types.DeviceHealthStatus,
		reg types.DeviceRegistration,
	) (types.Device, error)
	ListDevices(ctx context.Context) ([]types.Device, error)
	GetDiagnostics(ctx context.Context, device string) (types.Diagnostics, error)
}

type DeviceProvider interface {
	CreateClient(config device.Config) (device.Client, error)
}

type MonitorService struct {
	persistence PersistenceProvider
	device      DeviceProvider
	config      types.Config
	logger      *zap.Logger
}

func NewMonitorService(
	persistence PersistenceProvider,
	device DeviceProvider,
	config types.Config,
	logger *zap.Logger,
) *MonitorService {
	return &MonitorService{
		persistence: persistence,
		device:      device,
		config:      config,
		logger:      logger,
	}
}

func (s *MonitorService) RegisterDevice(
	ctx context.Context,
	reg types.DeviceRegistration,
) (types.Device, error) {
	port := reg.Port
	if reg.Protocol.IsHttp() {
		port = reg.GatewayPort
	}
	client, err := s.device.CreateClient(device.Config{
		Protocol: reg.Protocol,
		Host:     reg.Host,
		Port:     port,
	})
	if err != nil {
		return types.Device{}, err
	}
	health, err := client.GetHealth(ctx)
	if err != nil {
		return types.Device{}, err
	}
	device, err := s.persistence.RegisterDevice(ctx, *health, reg)
	if err != nil {
		return types.Device{}, err
	}
	return device, nil
}

func (s *MonitorService) ListDevices(ctx context.Context) ([]types.Device, error) {
	return s.persistence.ListDevices(ctx)
}

func (s *MonitorService) GetDiagnostics(ctx context.Context, device string) (types.Diagnostics, error) {
	return s.persistence.GetDiagnostics(ctx, device)
}

func (s *MonitorService) StreamDiagnostics(ctx context.Context, device string) <-chan types.Diagnostics {
	ch := make(chan types.Diagnostics)
	interval := s.config.StreamInterval
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				diagnostics, err := s.persistence.GetDiagnostics(ctx, device)
				if err != nil {
					s.logger.Error(
						"failed to get diagnostics for streaming",
						zap.String("device", device),
						zap.Error(err),
					)
					continue
				}
				ch <- diagnostics
			}
		}
	}()
	return ch
}
