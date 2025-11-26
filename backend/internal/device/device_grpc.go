package device

import (
	"context"
	"fmt"
	"io"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Device Client (gRPC)
type ClientGrpc struct {
	conn   *grpc.ClientConn
	client devicev1.DeviceClient
	config Config
}

func NewClientGrpc(config Config) (*ClientGrpc, error) {
	endpoint := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w (grpc): %w", ErrorClientCreation, err)
	}
	return &ClientGrpc{
		conn:   conn,
		client: devicev1.NewDeviceClient(conn),
		config: config,
	}, nil
}

func (d *ClientGrpc) GetHealth(ctx context.Context) (*types.DeviceHealthStatus, error) {
	res, err := d.client.GetHealth(ctx, &devicev1.GetHealthRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, fmt.Errorf("%w: device is not available (grpc): %w", ErrorNotFound, err)
			}
		}
		return nil, fmt.Errorf("failed to perform health request (grpc): %w", err)
	}
	return &types.DeviceHealthStatus{
		Identifier:         res.DeviceId,
		SupportedProtocols: types.ProtocolFromDevice(res.SupportedProtocols),
		Architecture:       res.Architecture,
		OS:                 res.Os,
		Updated:            res.UpdatedAt.AsTime(),
	}, nil
}

func (d *ClientGrpc) GetDiagnostics(ctx context.Context) (*types.DeviceDiagnostics, error) {
	res, err := d.client.GetDiagnostics(ctx, &devicev1.DiagnosticsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to perform diagnostics request (grpc): %w", err)
	}
	return d.diagnostics(res), nil
}

func (d *ClientGrpc) StreamDiagnostics(ctx context.Context) (<-chan *types.DeviceDiagnostics, <-chan error) {
	ch := make(chan *types.DeviceDiagnostics)
	errCh := make(chan error, 1)
	go func() {
		defer close(ch)
		defer close(errCh)
		stream, err := d.client.StreamDiagnostics(ctx, &devicev1.DiagnosticsRequest{})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.NotFound {
					errCh <- fmt.Errorf("%w: device is not available (grpc): %w", ErrorNotFound, err)
					return
				}
			}
			errCh <- fmt.Errorf("failed to create diagnostics stream (grpc): %w", err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						return
					}
					errCh <- fmt.Errorf("failed to receive from diagnostics stream (grpc): %w", err)
					return
				}
				ch <- d.diagnostics(res)
			}
		}
	}()
	return ch, errCh
}

func (d *ClientGrpc) UpdateDevice(ctx context.Context, status types.DeviceStatus) error {
	_, err := d.client.UpdateDevice(ctx, &devicev1.UpdateDeviceRequest{DeviceStatus: status.DeviceProto()})
	if err != nil {
		return fmt.Errorf("failed to perform device update request (grpc): %w", err)
	}
	return nil
}

func (d *ClientGrpc) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

func (d *ClientGrpc) diagnostics(diag *devicev1.DiagnosticsResponse) *types.DeviceDiagnostics {
	return &types.DeviceDiagnostics{
		Identifier: diag.DeviceId,
		DeviceVersions: types.DeviceVersions{
			Hardware: diag.HardwareVersion,
			Software: diag.SoftwareVersion,
			Firmware: diag.FirmwareVersion,
		},
		CPU:          diag.CpuUsage,
		Memory:       diag.MemoryUsage,
		DeviceStatus: types.DeviceStatusFromString(diag.DeviceStatus.String()),
		Checksum:     diag.Checksum,
		Timestamp:    diag.Timestamp.AsTime(),
	}
}
