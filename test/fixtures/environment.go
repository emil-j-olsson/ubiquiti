package fixtures

import (
	"context"
	"fmt"
	"testing"

	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultTestPrefix = "TEST"
)

type Environment struct {
	t      *testing.T
	ctx    context.Context
	config *Config
	device *DeviceClient
}

func NewEnvironment(t *testing.T) *Environment {
	err := godotenv.Load("../.env")
	assert.NoError(t, err)

	var config Config
	err = envconfig.Process(DefaultTestPrefix, &config)
	assert.NoError(t, err)

	return &Environment{
		t:      t,
		ctx:    context.Background(),
		config: &config,
	}
}

func (e *Environment) Device(service Service) *DeviceScenario {
	return &DeviceScenario{
		env:      e,
		service:  service,
		protocol: ProtocolGrpc,
	}
}

func (e *Environment) Close() {
	if e.device != nil {
		e.device.conn.Close() // nolint:errcheck
	}
}

type DeviceScenario struct {
	env      *Environment
	service  Service
	protocol Protocol
}

func (s *DeviceScenario) GetHealth() (*devicev1.GetHealthResponse, error) {
	device := s.client(s.env.t)
	return device.client.GetHealth(s.env.ctx, &devicev1.GetHealthRequest{})
}

func (s *DeviceScenario) GetDiagnostics() (*devicev1.DiagnosticsResponse, error) {
	device := s.client(s.env.t)
	return device.client.GetDiagnostics(s.env.ctx, &devicev1.DiagnosticsRequest{})
}

func (s *DeviceScenario) StreamDiagnostics() (grpc.ServerStreamingClient[devicev1.DiagnosticsResponse], error) {
	device := s.client(s.env.t)
	return device.client.StreamDiagnostics(s.env.ctx, &devicev1.DiagnosticsRequest{})
}

func (s *DeviceScenario) UpdateDevice(status devicev1.DeviceStatus) (*devicev1.UpdateDeviceResponse, error) {
	device := s.client(s.env.t)
	return device.client.UpdateDevice(s.env.ctx, &devicev1.UpdateDeviceRequest{
		DeviceStatus: status,
	})
}

type DeviceClient struct {
	conn     *grpc.ClientConn
	client   devicev1.DeviceClient
	service  ServiceConfig
	protocol Protocol
}

func (s *DeviceScenario) client(t *testing.T) *DeviceClient {
	service, exists := Services[s.service]
	if !exists {
		t.Fatalf("service %s not found in configuration", s.service)
		return nil
	}
	port := service.Port
	if s.protocol.IsHttp() {
		port = service.GatewayPort
	}
	endpoint := fmt.Sprintf("%s:%d", s.env.config.Host, port)
	conn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	client := &DeviceClient{
		conn:     conn,
		client:   devicev1.NewDeviceClient(conn),
		service:  service,
		protocol: s.protocol,
	}
	s.env.device = client
	return client
}
