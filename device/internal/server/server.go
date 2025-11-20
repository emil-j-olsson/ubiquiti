package server

import (
	"context"
	"time"

	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"go.uber.org/zap"
)

var _ devicev1.DeviceServer = (*Server)(nil)

const (
	DefaultContextTimeout time.Duration = 3 * time.Second
)

// could generate mocks...
type Provider interface {
	GetHealth()
}

type Server struct {
	devicev1.UnimplementedDeviceServer
	logger   *zap.Logger
	provider Provider
}

func NewDeviceServer(logger *zap.Logger, provider Provider) *Server {
	return &Server{
		logger:   logger,
		provider: provider,
	}
}

func (s *Server) GetHealth(ctx context.Context, req *devicev1.GetHealthRequest) (*devicev1.GetHealthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	s.provider.GetHealth()
	return nil, nil
}
