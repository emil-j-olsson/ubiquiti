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

// Persistence Repository
type PersistenceRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPersistenceRepository(pool *pgxpool.Pool, logger *zap.Logger) *PersistenceRepository {
	return &PersistenceRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *PersistenceRepository) RegisterDevice(
	ctx context.Context,
	status types.DeviceHealthStatus,
	reg types.DeviceRegistration,
) (types.Device, error) {
	protocols := make([]string, len(status.SupportedProtocols))
	for i, p := range status.SupportedProtocols {
		protocols[i] = p.String()
	}
	rows, err := r.pool.Query(ctx, `
		insert into devices (
			device_id, alias, host, port, port_gateway, architecture, os, supported_protocols
		) values (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		on conflict (device_id) do update set
			alias = excluded.alias,
			host = excluded.host,
			port = excluded.port,
			port_gateway = excluded.port_gateway,
			architecture = excluded.architecture,
			os = excluded.os,
			supported_protocols = excluded.supported_protocols,
			updated_at = excluded.updated_at
		returning *
	`,
		status.Identifier,
		reg.Alias,
		reg.Host,
		reg.Port,
		reg.GatewayPort,
		status.Architecture,
		status.OS,
		protocols,
	)
	if err != nil {
		return types.Device{}, fmt.Errorf(
			"%w: failed to insert device rows (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.Device])
	if err != nil {
		return types.Device{}, fmt.Errorf(
			"%w: failed to collect device rows (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	return result, nil
}

func (r *PersistenceRepository) GetDevice(ctx context.Context, deviceID string) (types.Device, error) {
	rows, err := r.pool.Query(ctx, `select * from devices where device_id = $1`, deviceID)
	if err != nil {
		return types.Device{}, fmt.Errorf(
			"%w: failed to query device (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.Device])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.Device{}, fmt.Errorf(
				"%w: failed to retrieve device '%s': %w",
				exceptions.ErrorNotFound,
				deviceID,
				err,
			)
		}
		return types.Device{}, fmt.Errorf(
			"%w: failed to collect device rows (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	return result, nil
}

func (r *PersistenceRepository) ListDevices(ctx context.Context) ([]types.Device, error) {
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

func (r *PersistenceRepository) GetDiagnostics(
	ctx context.Context,
	deviceID string,
) (types.Diagnostics, error) {
	rows, err := r.pool.Query(ctx, `select * from device_diagnostics_snapshot where device_id = $1`, deviceID)
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
				deviceID,
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

func (r *PersistenceRepository) SaveDiagnostics(ctx context.Context, diag types.DeviceDiagnostics) error {
	_, err := r.pool.Exec(ctx, `
		insert into device_diagnostics (
			device_id, cpu_usage, memory_usage, device_status,
            hardware_version, software_version, firmware_version,
            checksum, timestamp
		) values (
			(select id from devices where device_id = $1),
            $2, $3, $4, $5, $6, $7, $8, $9
		)
	`,
		diag.Identifier,
		diag.CPU,
		diag.Memory,
		diag.DeviceStatus,
		diag.DeviceVersions.Hardware,
		diag.DeviceVersions.Software,
		diag.DeviceVersions.Firmware,
		diag.Checksum,
		diag.Timestamp,
	)
	if err != nil {
		return fmt.Errorf(
			"%w: failed to insert diagnostics (postgres): %w",
			exceptions.ErrorInternal,
			err,
		)
	}
	return nil
}
