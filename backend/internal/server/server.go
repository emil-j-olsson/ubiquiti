package server

import (
	"context"
	"errors"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/database/exceptions"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/device"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ monitorv1.MonitorServer = (*Server)(nil)

const (
	DefaultContextTimeout time.Duration = 3 * time.Second
)

var (
	ErrorSendStream = errors.New("failed to send data over stream")
)

type Provider interface {
	RegisterDevice(ctx context.Context, reg types.DeviceRegistration) (types.Device, error)
	ListDevices(ctx context.Context) ([]types.Device, error)
	UpdateDevice(ctx context.Context, device string, status types.DeviceStatus) error
	GetDiagnostics(ctx context.Context, device string) (types.Diagnostics, error)
	StreamDiagnostics(ctx context.Context, device string) <-chan types.Diagnostics
}

type Server struct {
	monitorv1.UnimplementedMonitorServer
	logger   *zap.Logger
	provider Provider
}

func NewMonitorServer(provider Provider, logger *zap.Logger) *Server {
	return &Server{
		provider: provider,
		logger:   logger,
	}
}

func (s *Server) GetHealth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	return &emptypb.Empty{}, nil
}

func (s *Server) RegisterDevice(
	ctx context.Context,
	req *monitorv1.RegisterDeviceRequest,
) (*monitorv1.RegisterDeviceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	protocol := types.Protocol(req.GetProtocol().String())
	dev, err := s.provider.RegisterDevice(ctx, types.DeviceRegistration{
		Protocol:    protocol,
		Alias:       req.GetAlias(),
		Host:        req.GetHost(),
		Port:        req.GetPort(),
		GatewayPort: req.GetPortGateway(),
	})
	if err != nil {
		if errors.Is(err, device.ErrorNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &monitorv1.RegisterDeviceResponse{Device: s.device(dev)}, nil
}

func (s *Server) ListDevices(ctx context.Context, _ *emptypb.Empty) (*monitorv1.ListDevicesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	result, err := s.provider.ListDevices(ctx)
	if err != nil {
		return nil, s.databaseError(err)
	}
	devices := make([]*monitorv1.Device, len(result))
	for i, device := range result {
		devices[i] = s.device(device)
	}
	return &monitorv1.ListDevicesResponse{Devices: devices}, nil
}

func (s *Server) UpdateDevice(
	ctx context.Context,
	req *monitorv1.UpdateDeviceRequest,
) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	status := types.DeviceStatusFromString(req.GetDeviceStatus().String())
	err := s.provider.UpdateDevice(ctx, req.GetDeviceId(), status)
	if err != nil {
		return nil, s.databaseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetDiagnostics(
	ctx context.Context,
	req *monitorv1.DiagnosticsRequest,
) (*monitorv1.DiagnosticsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultContextTimeout)
	defer cancel()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := s.provider.GetDiagnostics(ctx, req.GetDeviceId())
	if err != nil {
		return nil, s.databaseError(err)
	}
	return s.diagnostics(result), nil
}

func (s *Server) StreamDiagnostics(
	req *monitorv1.DiagnosticsRequest,
	stream monitorv1.Monitor_StreamDiagnosticsServer,
) error {
	if err := req.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	ch := s.provider.StreamDiagnostics(stream.Context(), req.GetDeviceId())
	for diagnostics := range ch {
		response := s.diagnostics(diagnostics)
		if err := stream.Send(response); err != nil {
			s.logger.Error(ErrorSendStream.Error(), zap.Error(err))
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *Server) device(device types.Device) *monitorv1.Device {
	return &monitorv1.Device{
		Id:                 deref(device.ID),
		DeviceId:           deref(device.Identifier),
		Alias:              deref(device.Alias),
		Host:               deref(device.Host),
		Port:               deref(device.Port),
		PortGateway:        deref(device.GatewayPort),
		Architecture:       deref(device.Architecture),
		Os:                 deref(device.OS),
		SupportedProtocols: types.ProtocolFromStrings(deref(device.SupportedProtocols)),
		CreatedAt:          timestamp(device.Created),
		UpdatedAt:          timestamp(device.Updated),
	}
}

func (s *Server) diagnostics(diag types.Diagnostics) *monitorv1.DiagnosticsResponse {
	status := types.DeviceStatusFromString(deref(diag.DeviceStatus))
	return &monitorv1.DiagnosticsResponse{
		Device: &monitorv1.Device{
			Id:                 deref(diag.ID),
			DeviceId:           deref(diag.Identifier),
			Alias:              deref(diag.Alias),
			Host:               deref(diag.Host),
			Port:               deref(diag.Port),
			PortGateway:        deref(diag.GatewayPort),
			Architecture:       deref(diag.Architecture),
			Os:                 deref(diag.OS),
			SupportedProtocols: types.ProtocolFromStrings(deref(diag.SupportedProtocols)),
			CreatedAt:          timestamp(diag.Created),
			UpdatedAt:          timestamp(diag.Updated),
		},
		Diagnostics: &monitorv1.Diagnostics{
			HardwareVersion: deref(diag.Hardware),
			SoftwareVersion: deref(diag.Software),
			FirmwareVersion: deref(diag.Firmware),
			CpuUsage:        deref(diag.CPU),
			MemoryUsage:     deref(diag.Memory),
			DeviceStatus:    status.Proto(),
			Checksum:        deref(diag.Checksum),
		},
		UpdatedAt: timestamp(diag.LastUpdated),
	}

}

func (s *Server) databaseError(err error) error {
	if errors.Is(err, exceptions.ErrorNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}

func deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

func timestamp(t *time.Time) *timestamppb.Timestamp {
	if t != nil {
		return timestamppb.New(*t)
	}
	return nil
}
