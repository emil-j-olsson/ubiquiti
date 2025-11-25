package device

import (
	"context"
	"fmt"

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

// GetDiagnostics

// StreamDiagnostics

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
