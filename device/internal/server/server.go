package server

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/emil-j-olsson/ubiquiti/device/internal/types"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ devicev1.DeviceServer = (*Server)(nil)

const (
	DefaultContextTimeout time.Duration = 3 * time.Second
)

var (
	ErrorSendStream = errors.New("failed to send data over stream")
)

type Provider interface {
	GetHealth() *types.HealthStatus
	GetDiagnostics() *types.Diagnostics
	StreamDiagnostics(context.Context) <-chan *types.Diagnostics
	UpdateDevice(types.DeviceMutation)
}

type Server struct {
	devicev1.UnimplementedDeviceServer
	provider Provider
	logger   *zap.Logger
}

func NewDeviceServer(provider Provider, logger *zap.Logger) *Server {
	return &Server{
		provider: provider,
		logger:   logger,
	}
}

func (s *Server) GetHealth(
	ctx context.Context,
	_ *devicev1.GetHealthRequest,
) (*devicev1.GetHealthResponse, error) {
	_, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	health := s.provider.GetHealth()
	protocols := make([]devicev1.Protocol, 0, len(health.SupportedProtocols))
	for _, protocol := range health.SupportedProtocols {
		if p := protocol.Proto(); !slices.Contains(protocols, p) {
			protocols = append(protocols, p)
		}
	}
	return &devicev1.GetHealthResponse{
		DeviceId:           health.Identifier,
		SupportedProtocols: protocols,
		Architecture:       health.Architecture,
		Os:                 health.OS,
		UpdatedAt:          timestamppb.New(health.Updated),
	}, nil
}

func (s *Server) GetDiagnostics(
	ctx context.Context,
	_ *devicev1.DiagnosticsRequest,
) (*devicev1.DiagnosticsResponse, error) {
	_, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	diagnostics := s.provider.GetDiagnostics()
	return s.diagnostics(diagnostics), nil
}

func (s *Server) StreamDiagnostics(
	_ *devicev1.DiagnosticsRequest,
	stream devicev1.Device_StreamDiagnosticsServer,
) error {
	ch := s.provider.StreamDiagnostics(stream.Context())
	for diagnostics := range ch {
		response := s.diagnostics(diagnostics)
		if err := stream.Send(response); err != nil {
			s.logger.Error(ErrorSendStream.Error(), zap.Error(err))
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *Server) UpdateDevice(
	ctx context.Context,
	req *devicev1.UpdateDeviceRequest,
) (*devicev1.UpdateDeviceResponse, error) {
	_, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.provider.UpdateDevice(types.DeviceMutation{
		DeviceStatus: types.DeviceStatus(req.GetDeviceStatus().String()),
	})
	return &devicev1.UpdateDeviceResponse{}, nil
}

func (s *Server) diagnostics(diag *types.Diagnostics) *devicev1.DiagnosticsResponse {
	return &devicev1.DiagnosticsResponse{
		DeviceId:        diag.Identifier,
		HardwareVersion: diag.DeviceVersions.Hardware,
		SoftwareVersion: diag.DeviceVersions.Software,
		FirmwareVersion: diag.DeviceVersions.Firmware,
		CpuUsage:        diag.CPU,
		MemoryUsage:     diag.Memory,
		DeviceStatus:    diag.DeviceStatus.Proto(),
		Checksum:        diag.Checksum,
		Timestamp:       timestamppb.Now(),
	}
}
