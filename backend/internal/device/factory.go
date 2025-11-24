package device

import (
	"context"
	"fmt"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
)

var _ = []Client{(*ClientGrpc)(nil), (*ClientHttp)(nil)}

const (
	DefaultClientTimeout     = 3 * time.Second
	DefaultClientIdleTimeout = 30 * time.Second
)

var (
	ErrorUnsupportedProtocol = fmt.Errorf("unsupported device client protocol")
	ErrorClientCreation      = fmt.Errorf("failed to create device client")
	ErrorNotFound            = fmt.Errorf("not found")
)

type Client interface {
	GetHealth(ctx context.Context) (*types.DeviceHealthStatus, error)
	Close() error
}

type Config struct {
	Protocol types.Protocol
	Host     string
	Port     int64
}

// Device Client Factory
type factory struct{}

func NewClientFactory() *factory {
	return &factory{}
}

func (f *factory) CreateClient(config Config) (Client, error) {
	if config.Protocol.IsGrpc() {
		return NewClientGrpc(config)
	}
	if config.Protocol.IsHttp() {
		return NewClientHttp(config), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrorUnsupportedProtocol, config.Protocol.String())
}
