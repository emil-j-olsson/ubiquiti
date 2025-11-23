package service

import (
	"context"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"go.uber.org/zap"
)

// Business logic
// Here we spin up all layers needed for the monitor service & respond to api calls

type StateProvider interface {
	ListDevices(ctx context.Context) ([]types.Device, error)
	GetDiagnostics(ctx context.Context, device string) (types.Diagnostics, error)
}

type DeviceProvider interface {
}

type MonitorService struct {
	state  StateProvider
	device DeviceProvider
	config types.Config
	logger *zap.Logger
}

func NewMonitorService(
	state StateProvider,
	device DeviceProvider,
	config types.Config,
	logger *zap.Logger,
) *MonitorService {
	return &MonitorService{
		state:  state,
		device: device,
		config: config,
		logger: logger,
	}
}

func (s *MonitorService) ListDevices(ctx context.Context) ([]types.Device, error) {
	return s.state.ListDevices(ctx)
}

func (s *MonitorService) GetDiagnostics(ctx context.Context, device string) (types.Diagnostics, error) {
	return s.state.GetDiagnostics(ctx, device)
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
				diagnostics, err := s.state.GetDiagnostics(ctx, device)
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
