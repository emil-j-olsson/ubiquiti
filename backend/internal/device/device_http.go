package device

import (
	"bytes"
	"context"
	"encoding/json"
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
	defer response.Body.Close() // nolint:errcheck
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

func (d *ClientHttp) GetDiagnostics(ctx context.Context) (*types.DeviceDiagnostics, error) {
	endpoint, err := url.JoinPath(d.url, "/v1/diagnostics")
	if err != nil {
		return nil, fmt.Errorf("failed to join url path (http): %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create diagnostics request (http): %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform diagnostics request (http): %w", err)
	}
	defer response.Body.Close() // nolint:errcheck
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
	var res devicev1.DiagnosticsResponse
	if err := protojson.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode diagnostics response (http): %w", err)
	}
	return d.diagnostics(&res), nil
}

func (d *ClientHttp) StreamDiagnostics(ctx context.Context) (<-chan *types.DeviceDiagnostics, <-chan error) {
	ch := make(chan *types.DeviceDiagnostics)
	errCh := make(chan error, 1)
	go func() {
		defer close(ch)
		defer close(errCh)
		endpoint, err := url.JoinPath(d.url, "/v1/diagnostics/stream")
		if err != nil {
			errCh <- fmt.Errorf("failed to join url path (http): %w", err)
			return
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			errCh <- fmt.Errorf("failed to create stream diagnostics request (http): %w", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := d.client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("failed to perform stream diagnostics request (http): %w", err)
			return
		}
		defer response.Body.Close() // nolint:errcheck
		if response.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(response.Body)
			if response.StatusCode == http.StatusNotFound {
				errCh <- fmt.Errorf("%w: device is not available (http): %s", ErrorNotFound, endpoint)
				return
			}
			errCh <- fmt.Errorf("failed request with status %d (http): %s", response.StatusCode, string(body))
			return
		}
		decoder := json.NewDecoder(response.Body)
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				var res devicev1.DiagnosticsResponse
				if err := decoder.Decode(&res); err != nil {
					if err == io.EOF {
						return
					}
					errCh <- fmt.Errorf("failed to decode stream diagnostics response (http): %w", err)
					return
				}
				ch <- d.diagnostics(&res)
			}
		}
	}()
	return ch, errCh
}

func (d *ClientHttp) UpdateDevice(ctx context.Context, status types.DeviceStatus) error {
	endpoint, err := url.JoinPath(d.url, "/v1/device")
	if err != nil {
		return fmt.Errorf("failed to join url path (http): %w", err)
	}
	body := &devicev1.UpdateDeviceRequest{DeviceStatus: status.DeviceProto()}
	marshaled, err := protojson.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal update device request (http): %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewReader(marshaled))
	if err != nil {
		return fmt.Errorf("failed to create update device request (http): %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform update device request (http): %w", err)
	}
	defer response.Body.Close() // nolint:errcheck
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed request with status %d (http): %s", response.StatusCode, string(body))
	}
	return nil
}

func (d *ClientHttp) Close() error {
	d.client.CloseIdleConnections()
	return nil
}

func (d *ClientHttp) diagnostics(diag *devicev1.DiagnosticsResponse) *types.DeviceDiagnostics {
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
