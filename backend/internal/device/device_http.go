package device

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Device Client (HTTP)
type ClientHttp struct {
	url    string
	client *http.Client
	config Config
}

func NewClientHttp(config Config) *ClientHttp {
	url := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	return &ClientHttp{
		url: url,
		client: &http.Client{
			Timeout: DefaultClientTimeout,
			Transport: &http.Transport{
				IdleConnTimeout: DefaultClientIdleTimeout,
			},
		},
		config: config,
	}
}

func (d *ClientHttp) GetHealth(ctx context.Context) (*types.DeviceHealthStatus, error) {
	endpoint, err := url.JoinPath(d.url, "/v1/health")
	if err != nil {
		return nil, fmt.Errorf("failed to join url path (http): %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create health request (http): %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform health request (http): %w", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body (http): %w", err)
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: device is not available (http): %s", ErrorNotFound, endpoint)
		}
		return nil, fmt.Errorf("failed request with status %d (http): %s", response.StatusCode, string(body))
	}
	var res devicev1.GetHealthResponse
	if err := protojson.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode health response (http): %w", err)
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

// UpdateDevice

func (d *ClientHttp) Close() error {
	d.client.CloseIdleConnections()
	return nil
}
