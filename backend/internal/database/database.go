package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/database/postgres"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"go.uber.org/zap"
)

var (
	ErrorNotFound         = errors.New("not found")
	ErrorInternal         = errors.New("internal error")
	ErrorUnspecifiedField = errors.New("unspecified field")
	ErrorInvalidType      = errors.New("invalid type")
)

type Pool interface {
	Close()
}

type model struct {
	ctx         context.Context
	logger      *zap.Logger
	connections map[types.DatabaseInstance]*connection
}

type connection struct {
	database types.Database
	pool     Pool
}

type Config struct {
	Database        types.Database
	Instance        types.DatabaseInstance
	PostgresOptions []postgres.Option
}

func NewDatabaseModel(ctx context.Context, logger *zap.Logger) *model {
	return &model{
		ctx:         ctx,
		logger:      logger,
		connections: make(map[types.DatabaseInstance]*connection),
	}
}

func (m *model) AddDatabaseConnections(config ...Config) (*model, error) {
	for _, cfg := range config {
		if !cfg.Database.IsValid() {
			return nil, fmt.Errorf("%w: database", ErrorUnspecifiedField)
		}
		if !cfg.Instance.IsValid() {
			return nil, fmt.Errorf("%w: instance", ErrorUnspecifiedField)
		}
		switch cfg.Database {
		case types.DatabasePostgres:
			connector, err := postgres.NewConnector(m.ctx, cfg.PostgresOptions...)
			if err != nil {
				return nil, err
			}
			pool, err := connector.Connect()
			if err != nil {
				return nil, err
			}
			m.connections[cfg.Instance] = &connection{database: cfg.Database, pool: pool}
		}
	}
	return m, nil
}

func (m *model) GetPostgresPool(instance types.DatabaseInstance) (*postgres.Pool, error) {
	return GetDatabaseConnection[*postgres.Pool](m, instance)
}

func (m *model) Close() {
	for _, connection := range m.connections {
		connection.pool.Close()
	}
}

func GetDatabaseConnection[T Pool](m *model, instance types.DatabaseInstance) (T, error) {
	var zero T
	if conn, exists := m.connections[instance]; exists {
		if typedPool, ok := conn.pool.(T); ok {
			return typedPool, nil
		}
		return zero, fmt.Errorf("%w: connection exists but is not of type %T", ErrorInvalidType, zero)
	}
	return zero, ErrorNotFound
}
