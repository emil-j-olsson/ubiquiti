package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool = pgxpool.Pool

type connector struct {
	ctx        context.Context
	config     *pgxpool.Config
	connection types.PostgresConnection
	pool       *pgxpool.Pool
}

type Option func(*options)

type options struct {
	connection  types.PostgresConnection
	url         string
	connections int32
}

func (o *options) validate() error {
	if !o.connection.IsValid() {
		return errors.New("invalid connection type configuration (postgres)")
	}
	return nil
}

func WithProxyConnection(url string, connections int32) Option {
	return func(o *options) {
		o.connection = types.PostgresConnectionProxy
		o.url = url
		o.connections = connections
	}
}

func NewConnector(ctx context.Context, opts ...Option) (*connector, error) {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	if err := options.validate(); err != nil {
		return nil, err
	}
	switch options.connection {
	case types.PostgresConnectionProxy:
		config, err := pgxpool.ParseConfig(options.url)
		if err != nil {
			return nil, fmt.Errorf("failed to parse url connection string (postgres): %w", err)
		}
		config.MaxConns = options.connections
		return &connector{ctx: ctx, config: config, connection: options.connection}, nil
	default:
		return nil, fmt.Errorf("unsupported connection type (postgres): %s", options.connection)
	}
}

func (c *connector) Connect() (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(c.ctx, c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection (postgres): %w", err)
	}
	c.pool = pool
	return pool, nil
}

func (c *connector) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}
