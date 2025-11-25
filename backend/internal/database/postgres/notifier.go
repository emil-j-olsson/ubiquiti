package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	DefaultBackoffDuration = 1 * time.Second
)

type Notifier struct {
	pool    *pgxpool.Pool
	channel string
	events  chan types.Event
	logger  *zap.Logger
}

func NewNotifier(pool *pgxpool.Pool, channel string, logger *zap.Logger) *Notifier {
	return &Notifier{
		pool:    pool,
		channel: channel,
		events:  make(chan types.Event, 100),
		logger:  logger,
	}
}

func (n *Notifier) Subscribe(ctx context.Context) <-chan types.Event {
	return n.events
}

func (n *Notifier) Listen(ctx context.Context) error {
	defer close(n.events)
	backoff := DefaultBackoffDuration
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := n.listener(ctx); err != nil {
			n.logger.Error("notifier listener error", zap.Error(err))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff = min(2*backoff, 30*time.Second)
			}
			continue
		}
		return nil
	}
}

func (n *Notifier) listener(ctx context.Context) error {
	conn, err := n.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, fmt.Sprintf("listen %s", n.channel)); err != nil {
		return err
	}
	n.logger.Info("listening for notifications", zap.String("channel", n.channel))
	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}
		select {
		case n.events <- types.Event{
			Channel: notification.Channel,
			Payload: notification.Payload,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
