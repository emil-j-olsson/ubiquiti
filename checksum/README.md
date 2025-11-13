# Checksum Generator

A lightweight `SHA-256` checksum utility for generating deterministic cryptographic hashes from streaming data. Used for data integrity verification between the device monitoring service and the devices.

## Usage

```sh
# Build the binary:
make build
# Generate checksum from piped input:
echo '{"device_id":"test"}' | ./bin/checksum
# or:
echo '{"device_id":"test"}' | go run main.go
```
