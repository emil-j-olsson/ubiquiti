FROM golang:1.24-alpine

ARG PROTOC_GEN_GO_VERSION=v1.36.10
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.5.1
ARG GRPC_GATEWAY_VERSION=v2.27.3

RUN apk add --no-cache git make protobuf protobuf-dev curl

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION} && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${GRPC_GATEWAY_VERSION} && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${GRPC_GATEWAY_VERSION}

WORKDIR /tmp
RUN git clone --depth 1 https://github.com/googleapis/googleapis.git && \
    mkdir -p /usr/include/google && \
    cp -r googleapis/google/api /usr/include/google/ && \
    rm -rf googleapis

WORKDIR /workspace
