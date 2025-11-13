# Ubiquiti Backend Monitor Service

A high-performance backend monitoring service that manages network device registration, collects diagnostics data, and provides real-time streaming capabilities via multiple protocols (gRPC, HTTP) with PostgreSQL persistence and event-driven architecture.

## Supported Protocols

The service supports four different communication protocols:

- **`PROTOCOL_HTTP`** (1) - Standard HTTP REST API
- **`PROTOCOL_HTTP_STREAM`** (2) - HTTP Server-Sent Events (SSE) for streaming diagnostics
- **`PROTOCOL_GRPC`** (3) - gRPC unary calls for device management and diagnostics
- **`PROTOCOL_GRPC_STREAM`** (4) - gRPC server streaming for real-time diagnostics

Configure protocols via device registration. Protocols are defined in [`proto/monitor/v1/monitor.proto`](proto/monitor/v1/monitor.proto).

## Supported Platforms

The monitor service runs on multiple platforms:

| Platform | Architecture | Base Image | Container Port |
|----------|-------------|------------|----------------|
| Alpine Linux | `arm64` | `alpine:3.22.2` | 8080-8081 |
| Alpine Linux | `amd64` | `alpine:3.22.2` | 8082-8083 |

Platform configurations are defined in [`docker-compose.yaml`](../docker-compose.yaml) and [`Dockerfile`](Dockerfile).

## API Endpoints

### gRPC Service

The Monitor service implements the [`monitorv1.MonitorServer`](backend/proto/monitor/v1/monitor_grpc.pb.go) interface with the following methods:

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `GetHealth` | [`Empty`](proto/monitor/v1/monitor.pb.go) | [`Empty`](proto/monitor/v1/monitor.pb.go) | Get monitor service health status |
| `RegisterDevice` | [`RegisterDeviceRequest`](proto/monitor/v1/monitor.pb.go) | [`RegisterDeviceResponse`](proto/monitor/v1/monitor.pb.go) | Register a new device |
| `ListDevices` | [`Empty`](proto/monitor/v1/monitor.pb.go) | [`ListDevicesResponse`](proto/monitor/v1/monitor.pb.go) | List all registered devices |
| `UpdateDevice` | [`UpdateDeviceRequest`](proto/monitor/v1/monitor.pb.go) | [`Empty`](proto/monitor/v1/monitor.pb.go) | Update device status |
| `GetDiagnostics` | [`DiagnosticsRequest`](proto/monitor/v1/monitor.pb.go) | [`DiagnosticsResponse`](proto/monitor/v1/monitor.pb.go) | Get device diagnostics |
| `StreamDiagnostics` | [`DiagnosticsRequest`](proto/monitor/v1/monitor.pb.go) | [`DiagnosticsResponse`](proto/monitor/v1/monitor.pb.go) | Stream diagnostics in real-time |

### HTTP/REST Gateway

The service exposes an HTTP gateway via `grpc-gateway` that proxies to the gRPC service:

| Method | Endpoint | Description | Response Type |
|--------|----------|-------------|---------------|
| `GET` | `/v1/health` | Get monitor service health | JSON |
| `POST` | `/v1/devices/{device_id}` | Register a new device | JSON |
| `GET` | `/v1/devices` | List all registered devices | JSON |
| `PATCH` | `/v1/devices/{device_id}` | Update device status | JSON |
| `GET` | `/v1/diagnostics/{device_id}` | Get device diagnostics | JSON |
| `GET` | `/v1/diagnostics/{device_id}/stream` | Stream device diagnostics (SSE) | Server-Sent Events |


### Useful Commands

```bash
grpcurl -plaintext localhost:8080 monitor.v1.Monitor/GetHealth
grpcurl -plaintext -d '{"device_id": "ubiquiti-device-router-3c2d"}' localhost:8080 monitor.v1.Monitor/GetDiagnostics
grpcurl -plaintext -d '{"device_id": "ubiquiti-device-switch-b87f"}' localhost:8080 monitor.v1.Monitor/StreamDiagnostics
grpcurl -plaintext -d '{"device_id": "ubiquiti-device-switch-b87f", "device_status": "DEVICE_STATUS_ERROR"}' localhost:8080 monitor.v1.Monitor/UpdateDevice
grpcurl -plaintext -d '{"device_id": "ubiquiti-device-access-point-05da", "alias": "U7 Pro Max Ultimate", "host": "ubiquiti-device-access-point", "port": "8080", "port_gateway": "8081", "protocol": "PROTOCOL_HTTP"}' localhost:8080 monitor.v1.Monitor/RegisterDevice
```