ARG RUNTIME_IMAGE=alpine:latest

FROM --platform=$BUILDPLATFORM golang:1.24-alpine as build

ARG LDFLAGS="-X 'github.com/emil-j-olsson/ubiquiti/device/internal/types.Version=unknown'"
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

RUN apk add --no-cache git make

COPY device/go.mod device/go.sum ./device/
RUN cd device && go mod download

COPY device/ ./device/
RUN cd device && GOOS=$TARGETOS GOARCH=$TARGETARCH make build LDFLAGS="${LDFLAGS}"

COPY checksum/ ./checksum/
RUN cd checksum && GOOS=$TARGETOS GOARCH=$TARGETARCH make build LDFLAGS="${LDFLAGS}"

FROM ${RUNTIME_IMAGE}

WORKDIR /app

COPY --from=build /workspace/device/bin/server /app/server
COPY --from=build /workspace/checksum/bin/checksum /usr/local/bin/checksum

RUN chown -R 1000:1000 /app && \
    chmod +x /usr/local/bin/checksum

USER 1000

EXPOSE 8080 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/v1/health || exit 1

ENTRYPOINT [ "/app/server" ]
