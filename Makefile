GO_TEST_CMD := $(if $(shell which gotest), gotest, go test)
DOCKER_COMPOSE_CMD := $(if $(shell PATH=$(PATH) command -v docker-compose), docker-compose, docker compose)

.PHONY: fmt/global fmt
fmt/global: $(GOPATH)/bin/goimports $(GOPATH)/bin/golines
	goimports -w .; \
	golines -w --chain-split-dots --ignore-generated -m 110 .

fmt: $(GOPATH)/bin/goimports $(GOPATH)/bin/golines
	CHANGED_FILES=$$(git diff --name-only --diff-filter=AM | grep '\.go$$'); \
	goimports -w $$CHANGED_FILES; \
	golines -w --chain-split-dots --ignore-generated -m 110 $$CHANGED_FILES

.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	golangci-lint run

.PHONY: generate/env generate
generate/env:
	cp ./.env.example ./.env

generate:
	go generate ./...

.PHONY: git/hooks
git/hooks:
	ln -sf $(PWD)/.github/hooks/* ./.git/hooks

$(GOPATH)/bin/goimports:
	go install golang.org/x/tools/cmd/goimports@latest

$(GOPATH)/bin/golines:
	go install github.com/segmentio/golines@latest

$(GOPATH)/bin/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(GOPATH)/bin/go-enum:
	go install github.com/abice/go-enum@latest