package server

import (
	"context"
	"errors"
	"time"

	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ monitorv1.MonitorServer = (*Server)(nil)

const (
	DefaultContextTimeout time.Duration = 3 * time.Second
)

var (
	ErrorSendStream = errors.New("failed to send data over stream")
)

// type Provider interface {
// }

type Server struct {
	monitorv1.UnimplementedMonitorServer
	logger *zap.Logger
	// provider Provider
}

func NewMonitorServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) GetHealth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	return &emptypb.Empty{}, nil
}
