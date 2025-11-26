VERSION  = $(shell git describe --always --tags --dirty=-dirty)
REVISION = $(shell git rev-parse --short=8 HEAD)
BRANCH   = $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS ?= -X main.Version=${VERSION}

.PHONY: build
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/checksum .