package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/device"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"go.uber.org/zap"
)

// Worker (Polling Strategy)
type WorkerPoll struct {
	device      types.Device
	protocol    types.Protocol
	persistence PersistenceProvider
	provider    DeviceProvider
	interval    time.Duration
	logger      *zap.Logger
}

func NewWorkerPoll(
	device types.Device,
	protocol types.Protocol,
	persistence PersistenceProvider,
	provider DeviceProvider,
	interval time.Duration,
	logger *zap.Logger,
) *WorkerPoll {
	return &WorkerPoll{
		device:      device,
		protocol:    protocol,
		persistence: persistence,
		provider:    provider,
		interval:    interval,
		logger:      logger,
	}
}

func (w *WorkerPoll) Run(ctx context.Context) error {
	port := *w.device.Port
	if w.protocol.IsHttp() {
		port = *w.device.GatewayPort
	}
	client, err := w.provider.CreateClient(device.Config{
		Protocol: w.protocol,
		Host:     *w.device.Host,
		Port:     port,
	})
	if err != nil {
		return fmt.Errorf("failed to create client (%s): %w", w.protocol.String(), err)
	}
	defer client.Close() //nolint:errcheck
	deviceID := *w.device.Identifier
	config := DefaultPollingConfig(w.interval)
	job := func(ctx context.Context) error {
		diagnostics, err := client.GetDiagnostics(ctx)
		if err != nil {
			return err
		}
		return w.persistence.SaveDiagnostics(ctx, *diagnostics)
	}
	failure := func(ctx context.Context) error {
		return w.persistence.SaveDiagnostics(ctx, types.DeviceDiagnostics{
			Identifier:   deviceID,
			DeviceStatus: types.DeviceStatusOffline,
			Timestamp:    time.Now(),
		})
	}
	return NewPollingStrategy(config, job, failure, w.logger).Run(ctx)
}

// Worker (Streaming Strategy)
type WorkerStream struct {
	device      types.Device
	protocol    types.Protocol
	persistence PersistenceProvider
	provider    DeviceProvider
	logger      *zap.Logger
}

func NewWorkerStream(
	device types.Device,
	protocol types.Protocol,
	persistence PersistenceProvider,
	provider DeviceProvider,
	logger *zap.Logger,
) *WorkerStream {
	return &WorkerStream{
		device:      device,
		protocol:    protocol,
		persistence: persistence,
		provider:    provider,
		logger:      logger,
	}
}

func (w *WorkerStream) Run(ctx context.Context) error {
	// TODO: implement diagnostics streaming (stream stragegy; define externally – inject function)
	time.Sleep(time.Second * 50000)
	return nil
}
