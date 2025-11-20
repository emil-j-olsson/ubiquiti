FROM golang:1.24-alpine

RUN apk add --no-cache git make protobuf protobuf-dev curl

# TODO: pin versions...
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

WORKDIR /tmp
RUN git clone --depth 1 https://github.com/googleapis/googleapis.git && \
    mkdir -p /usr/include/google && \
    cp -r googleapis/google/api /usr/include/google/ && \
    rm -rf googleapis

WORKDIR /workspace