package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil-j-olsson/ubiquiti/backend/internal/database/exceptions"
	"github.com/emil-j-olsson/ubiquiti/backend/internal/types"
	"github.com/jackc/pgx/v5"
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

func (r *StateRetrieverRepository) ListDevices(ctx context.Context) ([]types.Device, error) {
	rows, err := r.pool.Query(ctx, `select * from devices order by created_at desc`)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to query devices (postgres): %w", exceptions.ErrorInternal, err)
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Device])
	if err != nil {
		return nil, fmt.Errorf(
			"%w: failed to collect device rows (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	return result, nil
}

func (r *StateRetrieverRepository) GetDiagnostics(
	ctx context.Context,
	device string,
) (types.Diagnostics, error) {
	rows, err := r.pool.Query(ctx, `select * from device_diagnostics_snapshot where device_id = $1`, device)
	if err != nil {
		return types.Diagnostics{}, fmt.Errorf(
			"%w: failed to query diagnostics (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.Diagnostics])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.Diagnostics{}, fmt.Errorf(
				"%w: failed to retrieve diagnostic data for device '%s': %w",
				exceptions.ErrorNotFound,
				device,
				err,
			)
		}
		return types.Diagnostics{}, fmt.Errorf(
			"%w: failed to collect diagnostic rows (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	return result, nil
}
