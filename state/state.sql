-- Enums
create type device_protocol as enum (
    'PROTOCOL_UNSPECIFIED',
    'PROTOCOL_HTTP',
    'PROTOCOL_HTTP_STREAM',
    'PROTOCOL_GRPC',
    'PROTOCOL_GRPC_STREAM'
);

create type device_status as enum (
    'DEVICE_STATUS_UNSPECIFIED',
    'DEVICE_STATUS_HEALTHY',
    'DEVICE_STATUS_DEGRADED',
    'DEVICE_STATUS_ERROR',
    'DEVICE_STATUS_MAINTENANCE',
    'DEVICE_STATUS_BOOTING' 
);

-- Tables
create table if not exists devices (
    id uuid primary key default gen_random_uuid(),
    device_id varchar(255) not null unique,
    alias varchar(255),
    architecture varchar(50) not null,
    os varchar(50) not null,
    supported_protocols device_protocol[] not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists device_diagnostics (
    id uuid primary key default gen_random_uuid(),
    device_id uuid not null references devices(id) on delete cascade,
    cpu_usage double precision not null,
    memory_usage double precision not null,
    device_status device_status not null,
    hardware_version varchar(50) not null,
    software_version varchar(50) not null,
    firmware_version varchar(50) not null,
    checksum varchar(255),
    timestamp timestamptz not null,
    created_at timestamptz not null default now()
);

-- Indexes for efficient queries
create index if not exists idx_device_device_id on devices(device_id);
create index if not exists idx_device_diagnostics_device_id on device_diagnostics(device_id);
create index if not exists idx_device_diagnostics_timestamp on device_diagnostics(timestamp desc);
create index if not exists idx_device_diagnostics_device_timestamp on device_diagnostics(device_id, timestamp desc);

-- Composite index for dashboard queries (latest state per device)
create index if not exists idx_device_diagnostics_latest on device_diagnostics(device_id, timestamp desc, device_status);

-- View for latest device diagnostics snapshot (optimized for dashboard)
create or replace view device_diagnostics_snapshot as
select distinct on (d.device_id)
    d.id,
    d.device_id,
    d.alias,
    d.architecture,
    d.os,
    d.supported_protocols,
    ds.hardware_version,
    ds.software_version,
    ds.firmware_version,
    ds.cpu_usage,
    ds.memory_usage,
    ds.device_status,
    ds.checksum,
    ds.timestamp as last_updated,
    d.created_at,
    d.updated_at
from devices d
left join device_diagnostics ds on d.id = ds.device_id
order by d.device_id, ds.timestamp desc nulls last;

-- Function to update device.updated_at timestamp
create or replace function update_device_updated_at()
returns trigger as $$
begin
    update devices set updated_at = now() where id = new.device_id;
    return new;
end;
$$ language plpgsql;

-- Trigger to auto-update device.updated_at when device state changes
create trigger trigger_update_device_timestamp
    after insert on device_diagnostics
    for each row
    execute function update_device_updated_at();

-- Dummy data (to be removed later)
insert into devices (device_id, alias, architecture, os, supported_protocols) values
    ('device-001', 'Gateway Router', 'arm64', 'Linux', array['PROTOCOL_HTTP'::device_protocol, 'PROTOCOL_GRPC'::device_protocol]),
    ('device-002', 'Access Point', 'amd64', 'Linux', array['PROTOCOL_HTTP'::device_protocol, 'PROTOCOL_GRPC_STREAM'::device_protocol]),
    ('device-003', 'Switch', 'arm', 'Linux', array['PROTOCOL_HTTP_STREAM'::device_protocol]);

insert into device_diagnostics (device_id, cpu_usage, memory_usage, device_status, hardware_version, software_version, firmware_version, checksum, timestamp) values
    ((select id from devices where device_id = 'device-001'), 45.2, 62.8, 'DEVICE_STATUS_HEALTHY', 'HW:1.0.0', 'SW:2.1.3', 'FW:3.0.1', 'abc123def456', now() - interval '1 hour'),
    ((select id from devices where device_id = 'device-001'), 52.3, 68.1, 'DEVICE_STATUS_HEALTHY', 'HW:1.0.0', 'SW:2.1.3', 'FW:3.0.1', 'abc123def457', now() - interval '30 minutes'),
    ((select id from devices where device_id = 'device-001'), 38.7, 59.2, 'DEVICE_STATUS_HEALTHY', 'HW:1.0.0', 'SW:2.1.3', 'FW:3.0.1', 'abc123def458', now()),
    ((select id from devices where device_id = 'device-002'), 78.9, 85.3, 'DEVICE_STATUS_DEGRADED', 'HW:1.2.0', 'SW:1.8.5', 'FW:2.5.0', 'def789ghi012', now() - interval '2 hours'),
    ((select id from devices where device_id = 'device-002'), 82.1, 87.6, 'DEVICE_STATUS_ERROR', 'HW:1.2.0', 'SW:1.8.5', 'FW:2.5.0', 'def789ghi013', now()),
    ((select id from devices where device_id = 'device-003'), 15.3, 42.7, 'DEVICE_STATUS_BOOTING', 'HW:2.0.0', 'SW:3.0.0', 'FW:4.1.2', 'ghi345jkl678', now() - interval '5 minutes'),
    ((select id from devices where device_id = 'device-003'), 22.1, 48.3, 'DEVICE_STATUS_HEALTHY', 'HW:2.0.0', 'SW:3.0.0', 'FW:4.1.2', 'ghi345jkl679', now());