package test

import (
	"testing"
	"time"

	monitorv1 "github.com/emil-j-olsson/ubiquiti/backend/proto/monitor/v1"
	"github.com/emil-j-olsson/ubiquiti/test/fixtures"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultTickerTimeout  = 10 * time.Second
	DefaultTickerInterval = 500 * time.Millisecond
)

func TestMonitor_GetHealth(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve health status of available monitor service (arm64)",
			service: fixtures.ServiceBackendMonitorArm,
		},
		{
			name:    "should retrieve health status of available monitor service (amd64)",
			service: fixtures.ServiceBackendMonitorAmd,
		},
		{
			name:    "should return error due to invalid device",
			service: fixtures.ServiceInvalid,
			err:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := fixtures.NewEnvironment(t)
			defer env.Close()
			res, err := env.Monitor(tt.service).GetHealth()
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}

func TestMonitor_ListDevices(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve list of devices from available monitor service (arm64)",
			service: fixtures.ServiceBackendMonitorArm,
		},
		{
			name:    "should retrieve list of devices from available monitor service (amd64)",
			service: fixtures.ServiceBackendMonitorAmd,
		},
		{
			name:    "should return error due to invalid device",
			service: fixtures.ServiceInvalid,
			err:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := fixtures.NewEnvironment(t)
			defer env.Close()
			res, err := env.Monitor(tt.service).ListDevices()
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for _, device := range res.GetDevices() {
				var expected fixtures.ServiceConfig
				var found bool
				for _, svc := range fixtures.Services {
					if svc.Identifier == device.DeviceId {
						expected = svc
						found = true
						break
					}
				}
				assert.True(t, found, "device %s not found in fixtures", device.DeviceId)
				assertValidDevice(t, expected, device)
			}
		})
	}
}

func TestMonitor_GetDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		device  fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve diagnostics from available monitor service (arm64)",
			service: fixtures.ServiceBackendMonitorArm,
			device:  fixtures.ServiceDeviceRouter,
		},
		{
			name:    "should retrieve diagnostics from available monitor service (amd64)",
			service: fixtures.ServiceBackendMonitorAmd,
			device:  fixtures.ServiceDeviceSwitch,
		},
		{
			name:    "should return error due to invalid device",
			service: fixtures.ServiceInvalid,
			err:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := fixtures.NewEnvironment(t)
			defer env.Close()
			device := fixtures.Services[tt.device]
			res, err := env.Monitor(tt.service).GetDiagnostics(device)
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertValidDevice(t, device, res.GetDevice())
			assertValidMonitorDiagnostics(t, res.GetDiagnostics())
			assert.NotNil(t, res.UpdatedAt)
		})
	}
}

func TestMonitor_StreamDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		device  fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve diagnostics from available monitor service (arm64)",
			service: fixtures.ServiceBackendMonitorArm,
			device:  fixtures.ServiceDeviceRouter,
		},
		{
			name:    "should retrieve diagnostics from available monitor service (amd64)",
			service: fixtures.ServiceBackendMonitorAmd,
			device:  fixtures.ServiceDeviceSwitch,
		},
		{
			name:    "should return error due to invalid device",
			service: fixtures.ServiceInvalid,
			err:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := fixtures.NewEnvironment(t)
			defer env.Close()
			device := fixtures.Services[tt.device]
			stream, err := env.Monitor(tt.service).StreamDiagnostics(device)
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			res, err := stream.Recv()
			assert.NoError(t, err)
			assertValidDevice(t, device, res.GetDevice())
			assertValidMonitorDiagnostics(t, res.GetDiagnostics())
			assert.NotNil(t, res.UpdatedAt)
		})
	}
}

func TestMonitor_UpdateDevice(t *testing.T) {
	t.Run("should update device via available monitor service (arm64)", func(t *testing.T) {
		env := fixtures.NewEnvironment(t)
		defer env.Close()
		service := fixtures.ServiceBackendMonitorArm
		device := fixtures.Services[fixtures.ServiceDeviceRouter]

		// Healthy -> Degraded -> Healthy
		statuses := []monitorv1.DeviceStatus{
			monitorv1.DeviceStatus_DEVICE_STATUS_DEGRADED,
			monitorv1.DeviceStatus_DEVICE_STATUS_HEALTHY,
		}
		for _, status := range statuses {
			_, err := env.Monitor(service).UpdateDevice(device, status)
			assert.NoError(t, err)

			// Await status change to propagate through the system
			timeout := time.After(DefaultTickerTimeout)
			ticker := time.NewTicker(DefaultTickerInterval)
			defer ticker.Stop()
			var diag *monitorv1.DiagnosticsResponse
			for {
				select {
				case <-timeout:
					t.Fatalf("timeout waiting for device status to change to %v", status)
				case <-ticker.C:
					diag, err = env.Monitor(service).GetDiagnostics(device)
					assert.NoError(t, err)
					if diag.Diagnostics.DeviceStatus == status {
						goto statusChanged
					}
				}
			}
		statusChanged:
			assert.Equal(t, status, diag.Diagnostics.DeviceStatus)
		}
	})
}

func TestMonitor_RegisterDevice(t *testing.T) {
	t.Run("should register device via available monitor service (amd64)", func(t *testing.T) {
		env := fixtures.NewEnvironment(t)
		defer env.Close()
		service := fixtures.ServiceBackendMonitorAmd
		device := fixtures.Services[fixtures.ServiceDeviceAccessPoint]

		object, err := env.Monitor(service).RegisterDevice(device)
		assert.NoError(t, err)
		assertValidDevice(t, device, object.GetDevice())

		devices, err := env.Monitor(service).ListDevices()
		assert.NoError(t, err)
		assert.Len(t, devices.Devices, 3)
		for _, dev := range devices.GetDevices() {
			if dev.DeviceId == device.Identifier {
				assertValidDevice(t, device, dev)
				return
			}
		}
		t.Error("registered device not found in device list")
	})
}

func assertValidDevice(t *testing.T, expected fixtures.ServiceConfig, actual *monitorv1.Device) {
	assert.Equal(t, expected.Identifier, actual.DeviceId)
	assert.Equal(t, expected.Alias, actual.Alias)
	assert.Equal(t, expected.Container, actual.Host)
	assert.Equal(t, 8080, int(actual.Port))
	assert.Equal(t, 8081, int(actual.PortGateway))
	assert.Equal(t, expected.Architecture, actual.Architecture)
	assert.Equal(t, expected.OS, actual.Os)
	for i, protocol := range expected.SupportedProtocols {
		if i < len(actual.SupportedProtocols) {
			assert.Equal(t, protocol.Proto(), actual.SupportedProtocols[i])
		}
	}
	assert.NotNil(t, actual.CreatedAt)
	assert.NotNil(t, actual.UpdatedAt)
}

func assertValidMonitorDiagnostics(t *testing.T, actual *monitorv1.Diagnostics) {
	assert.NotNil(t, actual.DeviceStatus)
	assert.NotNil(t, actual.HardwareVersion)
	assert.NotNil(t, actual.SoftwareVersion)
	assert.NotNil(t, actual.FirmwareVersion)
	assert.NotZero(t, actual.MemoryUsage)
	assert.NotNil(t, actual.Checksum)
	assert.NotEqual(t, DefaultInvalidChecksum, actual.Checksum)
}
