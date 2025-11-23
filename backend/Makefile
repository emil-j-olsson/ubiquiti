VERSION  = $(shell git describe --always --tags --dirty=-dirty)
REVISION = $(shell git rev-parse --short=8 HEAD)
BRANCH   = $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS ?= -X $(PKG).Version=${VERSION} \
	-X $(PKG).Revision=${REVISION} \
	-X $(PKG).Branch=${BRANCH} \
	-X $(PKG).Branch=${shell date +%Y-%m-%dT%T%z}

.PHONY: ldflags
ldflags:
	@echo $(LDFLAGS)

.PHONY: build
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/server ./cmd/server

PROTOC_IMAGE := monitor-protoc
.PHONY: generate generate/proto
generate/proto:
	docker build -t $(PROTOC_IMAGE) -f tools/protoc/Dockerfile .
	docker run --rm -v $(PWD):/workspace $(PROTOC_IMAGE) \
    	protoc -I . -I /usr/include \
    	--go_out . --go_opt paths=source_relative \
    	--go-grpc_out . --go-grpc_opt paths=source_relative \
    	--grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
    	proto/monitor/v1/monitor.proto
