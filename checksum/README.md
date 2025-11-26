## Checksum Generator

Simple program that returns a SHA-256 encoded checksum via:

```sh
make build
echo '{"device_id":"test"}' | ./bin/checksum
# or:
echo '{"device_id":"test"}' | go run main.go
```
