# Ubiquiti Device

A lightweight network device service that simulates Ubiquiti network hardware, exposing health status and diagnostics data via multiple protocols (gRPC, HTTP) with support for various platforms and architectures.

## Overview

The Device service provides a standardized API for monitoring network device health and diagnostics. It supports multiple communication protocols and can run on different operating systems and architectures, making it suitable for simulating various Ubiquiti network devices like routers, switches, and access points.

## Supported Protocols

The service supports four different communication protocols:

- **`PROTOCOL_HTTP`** - Standard HTTP REST API.
- **`PROTOCOL_HTTP_STREAM`** - HTTP Server-Sent Events (SSE) for streaming diagnostics.
- **`PROTOCOL_GRPC`** - gRPC unary calls for health and diagnostics.
- **`PROTOCOL_GRPC_STREAM`** - gRPC server streaming for real-time diagnostics.

## Operating Systems & Architectures

The service is designed to run on multiple platforms to simulate different Ubiquiti hardware:

| Platform | Architecture | Base Image | Example Device |
|----------|-------------|------------|----------------|
| Alpine Linux | `arm64` | `alpine:3.22.2` | Dream Machine Pro Max |
| Ubuntu | `amd64` | `ubuntu:22.04` | Pro Max 24 PoE |
| Debian | `armv7` | `debian:bullseye-slim` | U7 Pro Max Ultimate |

## API Endpoints

### gRPC Service

The Device service implements the [`devicev1.DeviceServer`](proto/device/v1/device_grpc.pb.go) interface with the following RPC methods:

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `GetHealth` | [`GetHealthRequest`](proto/device/v1/device.pb.go) | [`GetHealthResponse`](proto/device/v1/device.pb.go) | Get device health status |
| `GetDiagnostics` | [`DiagnosticsRequest`](proto/device/v1/device.pb.go) | [`DiagnosticsResponse`](proto/device/v1/device.pb.go) | Get current diagnostics |
| `StreamDiagnostics` | [`DiagnosticsRequest`](proto/device/v1/device.pb.go) | [`DiagnosticsResponse`](proto/device/v1/device.pb.go) | Stream diagnostics in real-time |
| `UpdateDevice` | [`UpdateDeviceRequest`](proto/device/v1/device.pb.go) | [`UpdateDeviceResponse`](proto/device/v1/device.pb.go) | Update device status |

### HTTP/REST Gateway

The service exposes an HTTP gateway via `grpc-gateway` that proxies to the gRPC service:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/health` | Get device health status |
| `GET` | `/v1/diagnostics` | Get current diagnostics |
| `GET` | `/v1/diagnostics/stream` | Stream diagnostics (SSE) |
| `PATCH` | `/v1/device` | Update device status |

### Useful Commands

```bash
# grpcurl:
grpcurl -plaintext localhost:8086 device.v1.Device/GetHealth
grpcurl -plaintext localhost:8086 device.v1.Device/GetDiagnostics
grpcurl -plaintext localhost:8086 device.v1.Device/StreamDiagnostics
grpcurl -plaintext -d '{"device_status": "DEVICE_STATUS_MAINTENANCE"}' localhost:8086 device.v1.Device/UpdateDevice

# curl:
curl localhost:8087/v1/health
curl localhost:8087/v1/diagnostics
curl localhost:8087/v1/diagnostics/stream
```