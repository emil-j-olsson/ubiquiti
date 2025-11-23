package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Device Worker Repository
type DeviceWorkerRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewDeviceWorkerRepository(pool *pgxpool.Pool, logger *zap.Logger) *DeviceWorkerRepository {
	return &DeviceWorkerRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *DeviceWorkerRepository) GetSomething() {}

// State Retriever Repository
type StateRetrieverRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewStateRetrieverRepository(pool *pgxpool.Pool, logger *zap.Logger) *StateRetrieverRepository {
	return &StateRetrieverRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *StateRetrieverRepository) GetSomethingElse() {}
