package worker

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

type RetryConfig struct {
	Interval          time.Duration
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
	Timeout           time.Duration
	HeartbeatTimeout  time.Duration
}

func DefaultPollingConfig(interval time.Duration) RetryConfig {
	return RetryConfig{
		Interval:          interval,
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		Timeout:           5 * time.Second,
	}
}

func DefaultStreamingConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    2 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		HeartbeatTimeout:  5 * time.Second,
	}
}

type PollFunc func(ctx context.Context) error

type StreamFunc func(ctx context.Context, messageCh chan<- struct{}) error

type ErrorFunc func(ctx context.Context) error

// Polling Strategy
type PollingStrategy struct {
	config    RetryConfig
	pollFunc  PollFunc
	errorFunc ErrorFunc
	logger    *zap.Logger
}

func NewPollingStrategy(
	config RetryConfig,
	pollFunc PollFunc,
	errorFunc ErrorFunc,
	logger *zap.Logger,
) *PollingStrategy {
	return &PollingStrategy{
		config:    config,
		pollFunc:  pollFunc,
		errorFunc: errorFunc,
		logger:    logger,
	}
}

func (s *PollingStrategy) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.poll(ctx); err != nil {
				s.logger.Error("poll failed", zap.Error(err))
			}
		}
	}
}

func (s *PollingStrategy) poll(ctx context.Context) error {
	var lastErr error
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := s.calculateBackoff(attempt)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		pctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
		err := s.pollFunc(pctx)
		cancel()
		if err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	if err := s.errorFunc(ctx); err != nil {
		return err
	}
	return fmt.Errorf("poll failed after %d attempts: %w", s.config.MaxRetries+1, lastErr)
}

func (s *PollingStrategy) calculateBackoff(attempt int) time.Duration {
	backoff := float64(s.config.InitialBackoff) * math.Pow(s.config.BackoffMultiplier, float64(attempt-1))
	duration := time.Duration(backoff)
	if duration > s.config.MaxBackoff {
		duration = s.config.MaxBackoff
	}
	return duration
}

// Streaming Strategy
type StreamingStrategy struct {
	config     RetryConfig
	streamFunc StreamFunc
	errorFunc  ErrorFunc
	logger     *zap.Logger
}

func NewStreamingStrategy(
	config RetryConfig,
	streamFunc StreamFunc,
	errorFunc ErrorFunc,
	logger *zap.Logger,
) *StreamingStrategy {
	return &StreamingStrategy{
		config:     config,
		streamFunc: streamFunc,
		errorFunc:  errorFunc,
		logger:     logger,
	}
}

func (s *StreamingStrategy) Run(ctx context.Context) error {
	var lastErr error
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := s.calculateBackoff(attempt)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		err := s.stream(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		lastErr = err
	}
	if err := s.errorFunc(ctx); err != nil {
		s.logger.Error("failed to handle stream error", zap.Error(err))
	}
	return fmt.Errorf("stream failed after %d attempts: %w", s.config.MaxRetries+1, lastErr)
}

func (s *StreamingStrategy) stream(ctx context.Context) error {
	ch := make(chan struct{}, 1)
	sctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.streamFunc(sctx, ch)
	}()
	heartbeat := time.NewTimer(s.config.HeartbeatTimeout)
	defer heartbeat.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		case <-ch:
			if !heartbeat.Stop() {
				<-heartbeat.C
			}
			heartbeat.Reset(s.config.HeartbeatTimeout)
		case <-heartbeat.C:
			return fmt.Errorf("heartbeat timeout after %v", s.config.HeartbeatTimeout)
		}
	}
}

func (s *StreamingStrategy) calculateBackoff(attempt int) time.Duration {
	backoff := float64(s.config.InitialBackoff) * math.Pow(s.config.BackoffMultiplier, float64(attempt-1))
	duration := time.Duration(backoff)
	if duration > s.config.MaxBackoff {
		duration = s.config.MaxBackoff
	}
	return duration
}
