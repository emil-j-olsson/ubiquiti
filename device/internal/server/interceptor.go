package server

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
		var resp any
		return resp, withLogging(l, info.FullMethod, "grpc request", func() (err error) {
			resp, err = h(ctx, req)
			return
		})
	}
}

func UnaryRecoveryInterceptor(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
		var resp any
		return resp, withRecovery(l, info.FullMethod, func() (err error) {
			resp, err = h(ctx, req)
			return
		})
	}
}

func StreamLoggingInterceptor(l *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		return withLogging(l, info.FullMethod, "grpc stream", func() error { return h(srv, ss) })
	}
}

func StreamRecoveryInterceptor(l *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		return withRecovery(l, info.FullMethod, func() error { return h(srv, ss) })
	}
}

func withLogging(l *zap.Logger, method, msg string, fn func() error) error {
	start := time.Now()
	err := fn()
	l.Info(msg,
		zap.String("method", method),
		zap.Duration("duration", time.Since(start)),
		zap.String("code", status.Code(err).String()),
	)
	return err
}

func withRecovery(l *zap.Logger, method string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			l.Error("panic recovered", zap.String("method", method), zap.Any("panic", r))
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return fn()
}
