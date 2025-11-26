package test

import (
	"testing"

	devicev1 "github.com/emil-j-olsson/ubiquiti/device/proto/device/v1"
	"github.com/emil-j-olsson/ubiquiti/test/fixtures"
	"github.com/stretchr/testify/assert"
)

const DefaultInvalidChecksum = "invalid-checksum"

func TestDevice_GetHealth(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve health status of available device (router)",
			service: fixtures.ServiceDeviceRouter,
		},
		{
			name:    "should retrieve health status of available device (switch)",
			service: fixtures.ServiceDeviceSwitch,
		},
		{
			name:    "should retrieve health status of available device (access point)",
			service: fixtures.ServiceDeviceAccessPoint,
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
			res, err := env.Device(tt.service).GetHealth()
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			expected := fixtures.Services[tt.service]
			assert.Equal(t, expected.Identifier, res.DeviceId)
			assert.Equal(t, expected.Architecture, res.Architecture)
			assert.Equal(t, expected.OS, res.Os)
			for i, protocol := range expected.SupportedProtocols {
				assert.Equal(t, protocol, fixtures.Protocol(res.SupportedProtocols[i].String()))
			}
			assert.NotNil(t, res.UpdatedAt)
		})
	}
}

func TestDevice_GetDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve diagnostics of available device (router)",
			service: fixtures.ServiceDeviceRouter,
		},
		{
			name:    "should retrieve diagnostics of available device (switch)",
			service: fixtures.ServiceDeviceSwitch,
		},
		{
			name:    "should retrieve diagnostics of available device (access point)",
			service: fixtures.ServiceDeviceAccessPoint,
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
			res, err := env.Device(tt.service).GetDiagnostics()
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			expected := fixtures.Services[tt.service]
			assert.Equal(t, expected.Identifier, res.DeviceId)
			assertValidDiagnostics(t, res)
		})
	}
}

func TestDevice_StreamDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		service fixtures.Service
		err     bool
	}{
		{
			name:    "should retrieve streamed diagnostics of available device (router)",
			service: fixtures.ServiceDeviceRouter,
		},
		{
			name:    "should retrieve streamed diagnostics of available device (switch)",
			service: fixtures.ServiceDeviceSwitch,
		},
		{
			name:    "should retrieve streamed diagnostics of available device (access point)",
			service: fixtures.ServiceDeviceAccessPoint,
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
			stream, err := env.Device(tt.service).StreamDiagnostics()
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			res, err := stream.Recv()
			assert.NoError(t, err)
			expected := fixtures.Services[tt.service]
			assert.Equal(t, expected.Identifier, res.DeviceId)
			assertValidDiagnostics(t, res)
		})
	}
}

func TestDevice_UpdateDiagnostics(t *testing.T) {
	t.Run("should update status of available device (router)", func(t *testing.T) {
		env := fixtures.NewEnvironment(t)
		defer env.Close()

		// Healthy -> Maintenance
		_, err := env.Device(fixtures.ServiceDeviceRouter).UpdateDevice(
			devicev1.DeviceStatus_DEVICE_STATUS_MAINTENANCE,
		)
		assert.NoError(t, err)
		res, err := env.Device(fixtures.ServiceDeviceRouter).GetDiagnostics()
		assert.NoError(t, err)
		assert.Equal(t, devicev1.DeviceStatus_DEVICE_STATUS_MAINTENANCE, res.DeviceStatus)

		// Maintenance -> Healthy
		_, err = env.Device(fixtures.ServiceDeviceRouter).UpdateDevice(
			devicev1.DeviceStatus_DEVICE_STATUS_HEALTHY,
		)
		assert.NoError(t, err)
		res, err = env.Device(fixtures.ServiceDeviceRouter).GetDiagnostics()
		assert.NoError(t, err)
		assert.Equal(t, devicev1.DeviceStatus_DEVICE_STATUS_HEALTHY, res.DeviceStatus)
	})
	t.Run("should return error due to invalid device", func(t *testing.T) {
		env := fixtures.NewEnvironment(t)
		defer env.Close()
		_, err := env.Device(fixtures.ServiceInvalid).UpdateDevice(
			devicev1.DeviceStatus_DEVICE_STATUS_MAINTENANCE,
		)
		assert.Error(t, err)
	})
}

func assertValidDiagnostics(t *testing.T, diag *devicev1.DiagnosticsResponse) {
	assert.NotNil(t, diag.DeviceStatus)
	assert.NotNil(t, diag.HardwareVersion)
	assert.NotNil(t, diag.SoftwareVersion)
	assert.NotNil(t, diag.FirmwareVersion)
	assert.NotZero(t, diag.MemoryUsage)
	assert.NotNil(t, diag.Timestamp)
	assert.NotNil(t, diag.Checksum)
	assert.NotEqual(t, DefaultInvalidChecksum, diag.Checksum)
}
