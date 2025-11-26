ARG RUNTIME_IMAGE=alpine:latest

FROM --platform=$BUILDPLATFORM golang:1.24-alpine as build

ARG LDFLAGS="-X 'github.com/emil-j-olsson/ubiquiti/backend/internal/types.Version=unknown'"
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

RUN apk add --no-cache git make

COPY go.work go.work.sum* ./

COPY device/go.mod device/go.sum ./device/
COPY device/ ./device/

COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

COPY backend/ ./backend/

RUN cd backend && GOOS=$TARGETOS GOARCH=$TARGETARCH make build LDFLAGS="${LDFLAGS}"

FROM ${RUNTIME_IMAGE}

WORKDIR /app

COPY --from=build /workspace/backend/bin/server /app/server

RUN chown -R 1000:1000 /app
USER 1000

EXPOSE 8080 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/v1/health || exit 1

ENTRYPOINT [ "/app/server" ]
