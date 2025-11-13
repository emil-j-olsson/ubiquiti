package fixtures

import (
	"context"
	"fmt"
	"testing"

	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	DefaultTestPrefix = "TEST"
)

type Environment struct {
	t       *testing.T
	ctx     context.Context
	config  *Config
	device  *DeviceClient
	monitor *MonitorClient
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

func (e *Environment) Monitor(service Service) *MonitorScenario {
	return &MonitorScenario{
		env:      e,
		service:  service,
		protocol: ProtocolGrpc,
	}
}

func (e *Environment) Close() {
	if e.device != nil {
		e.device.conn.Close() // nolint:errcheck
	}
	if e.monitor != nil {
		e.monitor.conn.Close() // nolint:errcheck
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

type MonitorScenario struct {
	env      *Environment
	service  Service
	protocol Protocol
}

func (s *MonitorScenario) GetHealth() (*emptypb.Empty, error) {
	monitor := s.client(s.env.t)
	return monitor.client.GetHealth(s.env.ctx, &emptypb.Empty{})
}

func (s *MonitorScenario) RegisterDevice(service ServiceConfig) (*monitorv1.RegisterDeviceResponse, error) {
	monitor := s.client(s.env.t)
	return monitor.client.RegisterDevice(s.env.ctx, &monitorv1.RegisterDeviceRequest{
		DeviceId:    service.Identifier,
		Alias:       service.Alias,
		Host:        service.Container,
		Port:        8080,
		PortGateway: 8081,
		Protocol:    service.SupportedProtocols[0].Proto(),
	})
}

func (s *MonitorScenario) ListDevices() (*monitorv1.ListDevicesResponse, error) {
	monitor := s.client(s.env.t)
	return monitor.client.ListDevices(s.env.ctx, &emptypb.Empty{})
}

func (s *MonitorScenario) UpdateDevice(
	service ServiceConfig,
	status monitorv1.DeviceStatus,
) (*emptypb.Empty, error) {
	monitor := s.client(s.env.t)
	return monitor.client.UpdateDevice(s.env.ctx, &monitorv1.UpdateDeviceRequest{
		DeviceId:     service.Identifier,
		DeviceStatus: status,
	})
}

func (s *MonitorScenario) GetDiagnostics(service ServiceConfig) (*monitorv1.DiagnosticsResponse, error) {
	monitor := s.client(s.env.t)
	return monitor.client.GetDiagnostics(s.env.ctx, &monitorv1.DiagnosticsRequest{
		DeviceId: service.Identifier,
	})
}

func (s *MonitorScenario) StreamDiagnostics(
	service ServiceConfig,
) (grpc.ServerStreamingClient[monitorv1.DiagnosticsResponse], error) {
	monitor := s.client(s.env.t)
	return monitor.client.StreamDiagnostics(s.env.ctx, &monitorv1.DiagnosticsRequest{
		DeviceId: service.Identifier,
	})
}

type MonitorClient struct {
	conn     *grpc.ClientConn
	client   monitorv1.MonitorClient
	service  ServiceConfig
	protocol Protocol
}

func (s *MonitorScenario) client(t *testing.T) *MonitorClient {
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
	client := &MonitorClient{
		conn:     conn,
		client:   monitorv1.NewMonitorClient(conn),
		service:  service,
		protocol: s.protocol,
	}
	s.env.monitor = client
	return client
}
