package worker

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

type PollingConfig struct {
	Interval          time.Duration
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
	Timeout           time.Duration
}

func DefaultPollingConfig(interval time.Duration) PollingConfig {
	return PollingConfig{
		Interval:          interval,
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		Timeout:           5 * time.Second,
	}
}

type PollFunc func(ctx context.Context) error

type ErrorFunc func(ctx context.Context) error

type PollingStrategy struct {
	config    PollingConfig
	pollFunc  PollFunc
	errorFunc ErrorFunc
	logger    *zap.Logger
}

func NewPollingStrategy(
	config PollingConfig,
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
