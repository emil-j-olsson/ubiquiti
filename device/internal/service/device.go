package service

import (
	"context"
	"runtime"
	"time"

	"github.com/emil-j-olsson/ubiquiti/device/internal/types"
	"go.uber.org/zap"
)

type StateProvider interface {
	GetState() types.DeviceState
	UpdateState(fn func(*types.DeviceState)) types.DeviceState
}

type ChecksumGenerator interface {
	GenerateChecksum(ctx context.Context, data []byte) (string, error)
}

type Service struct {
	provider StateProvider
	checksum ChecksumGenerator
	logger   *zap.Logger
}

func NewDeviceService(provider StateProvider, checksum ChecksumGenerator, logger *zap.Logger) *Service {
	return &Service{provider: provider, checksum: checksum, logger: logger}
}

func (s *Service) GetHealth() *types.HealthStatus {
	state := s.provider.GetState()
	return &types.HealthStatus{
		Identifier:         state.Identifier,
		SupportedProtocols: state.SupportedProtocols,
		Architecture:       state.Architecture,
		OS:                 state.OS,
		Updated:            state.Updated,
	}
}

func (s *Service) GetDiagnostics() *types.Diagnostics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	state := s.provider.GetState()
	return s.diagnostics(state)
}

func (s *Service) StreamDiagnostics(ctx context.Context) <-chan *types.Diagnostics {
	ch := make(chan *types.Diagnostics)
	interval := s.provider.GetState().StreamInterval
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				state := s.provider.GetState()
				ch <- s.diagnostics(state)
			}
		}
	}()
	return ch
}

func (s *Service) UpdateDevice(mutation types.DeviceMutation) {
	s.provider.UpdateState(func(state *types.DeviceState) {
		state.DeviceStatus = mutation.DeviceStatus
	})
}

func (s *Service) GenerateChecksum(ctx context.Context, data []byte) (string, error) {
	return s.checksum.GenerateChecksum(ctx, data)
}

func (s *Service) diagnostics(state types.DeviceState) *types.Diagnostics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &types.Diagnostics{
		Identifier:     state.Identifier,
		DeviceVersions: state.DeviceVersions,
		CPU:            m.GCCPUFraction * 100.0,
		Memory:         float64(m.Alloc) / float64(m.Sys) * 100.0,
		DeviceStatus:   state.DeviceStatus,
	}
}
