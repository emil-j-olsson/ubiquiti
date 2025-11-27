FROM node:20-alpine

ARG GRPC_WEB_VERSION=1.5.0

RUN apk add --no-cache git curl protobuf protobuf-dev

# Install ts-proto for modern protobuf generation
RUN npm install -g \
    ts-proto@^1.181.0 \
    google-protobuf@^3.21.0 \
    grpc-web@^1.5.0

# Download protoc-gen-grpc-web binary
RUN curl -Lo /usr/local/bin/protoc-gen-grpc-web \
    https://github.com/grpc/grpc-web/releases/download/${GRPC_WEB_VERSION}/protoc-gen-grpc-web-${GRPC_WEB_VERSION}-linux-x86_64 && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web

# Clone googleapis for google/api annotations
WORKDIR /tmp
RUN git clone --depth 1 https://github.com/googleapis/googleapis.git && \
    mkdir -p /usr/include/google && \
    cp -r googleapis/google/api /usr/include/google/ && \
    rm -rf googleapis

WORKDIR /workspace