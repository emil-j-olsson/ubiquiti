package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/device"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"go.uber.org/zap"
)

var _ = []Worker{(*WorkerPoll)(nil), (*WorkerStream)(nil)}

var (
	// Protocol hierarchy preference for workers (desc)
	WorkerProtocolHierarchy = []types.Protocol{
		types.ProtocolGrpcStream,
		types.ProtocolHttpStream,
		types.ProtocolGrpc,
		types.ProtocolHttp,
	}
)

type PersistenceProvider interface {
	GetDevice(ctx context.Context, deviceID string) (types.Device, error)
	ListDevices(ctx context.Context) ([]types.Device, error)
	SaveDiagnostics(ctx context.Context, diag types.DeviceDiagnostics) error
}

type PersistenceNotifier interface {
	Subscribe(ctx context.Context) <-chan types.Event
}

type DeviceProvider interface {
	CreateClient(config device.Config) (device.Client, error)
}

type EventPayload struct {
	DeviceID  string `json:"device_id"`
	Operation string `json:"operation"`
}

// Worker Orchestrator
type orchestrator struct {
	persistence PersistenceProvider
	notifier    PersistenceNotifier
	device      DeviceProvider
	pool        *pool
	interval    time.Duration
	logger      *zap.Logger
}

func NewOrchestrator(
	persistence PersistenceProvider,
	notifier PersistenceNotifier,
	device DeviceProvider,
	interval time.Duration,
	logger *zap.Logger,
) *orchestrator {
	return &orchestrator{
		persistence: persistence,
		notifier:    notifier,
		device:      device,
		pool:        NewPool(persistence, device, interval, logger),
		interval:    interval,
		logger:      logger,
	}
}

func (o *orchestrator) Run(ctx context.Context) error {
	// Run workers for existing devices
	devices, err := o.persistence.ListDevices(ctx)
	if err != nil {
		return err
	}
	for _, dev := range devices {
		if err := o.pool.NewWorker(ctx, *dev.Identifier); err != nil {
			o.logger.Error(
				"failed to create worker",
				zap.String("device_id", *dev.Identifier),
				zap.Error(err),
			)
		}
	}
	// Orchestrate workers on device events
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-o.notifier.Subscribe(ctx):
			var payload EventPayload
			if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
				o.logger.Error("failed to unmarshal event payload", zap.Error(err))
				continue
			}
			o.logger.Info(
				"received device event",
				zap.String("device_id", payload.DeviceID),
				zap.String("operation", payload.Operation),
			)
			switch payload.Operation {
			case "INSERT":
				if err := o.pool.NewWorker(ctx, payload.DeviceID); err != nil {
					o.logger.Error(
						"failed to create worker",
						zap.String("device_id", payload.DeviceID),
						zap.Error(err),
					)
				}
			case "DELETE":
				o.pool.StopWorker(payload.DeviceID)
			}
		}
	}
}

// Worker Pool
type Worker interface {
	Run(ctx context.Context) error
}

type pool struct {
	persistence PersistenceProvider
	device      DeviceProvider
	interval    time.Duration
	workers     map[string]context.CancelFunc
	mu          sync.RWMutex
	logger      *zap.Logger
}

func NewPool(
	persistence PersistenceProvider,
	device DeviceProvider,
	interval time.Duration,
	logger *zap.Logger,
) *pool {
	return &pool{
		workers:     make(map[string]context.CancelFunc),
		persistence: persistence,
		device:      device,
		interval:    interval,
		logger:      logger,
	}
}

func (p *pool) NewWorker(ctx context.Context, deviceID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.workers[deviceID]; exists {
		p.logger.Info("worker already exists", zap.String("device_id", deviceID))
		return nil
	}
	device, err := p.persistence.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}
	protocol := p.protocol(*device.SupportedProtocols)
	if protocol == "" {
		return fmt.Errorf("no supported protocol found for device %s", deviceID)
	}
	wctx, cancel := context.WithCancel(ctx)
	p.workers[deviceID] = cancel

	var worker Worker
	switch protocol {
	case types.ProtocolGrpc:
		worker = NewWorkerPoll(device, types.ProtocolGrpc, p.persistence, p.device, p.interval, p.logger)
	case types.ProtocolGrpcStream:
		worker = NewWorkerStream(device, types.ProtocolGrpcStream, p.persistence, p.device, p.logger)
	case types.ProtocolHttp:
		worker = NewWorkerPoll(device, types.ProtocolHttp, p.persistence, p.device, p.interval, p.logger)
	case types.ProtocolHttpStream:
		worker = NewWorkerStream(device, types.ProtocolHttpStream, p.persistence, p.device, p.logger)
	}
	go func() {
		defer p.delete(deviceID)
		if err := worker.Run(wctx); err != nil {
			p.logger.Error("worker encountered an error", zap.String("device_id", deviceID), zap.Error(err))
		}
	}()
	return nil
}

func (p *pool) StopWorker(deviceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if cancel, exists := p.workers[deviceID]; exists {
		cancel()
		delete(p.workers, deviceID)
		p.logger.Info("stopped worker", zap.String("device_id", deviceID))
	}
}

func (p *pool) protocol(protocols []string) types.Protocol {
	for _, preferred := range WorkerProtocolHierarchy {
		if slices.ContainsFunc(protocols, func(proto string) bool {
			return types.Protocol(proto) == preferred
		}) {
			return preferred
		}
	}
	return ""
}

func (p *pool) delete(deviceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.workers, deviceID)
}
