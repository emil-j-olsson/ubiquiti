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
    'DEVICE_STATUS_BOOTING',
    'DEVICE_STATUS_OFFLINE' 
);

-- Tables
create table if not exists devices (
    id uuid primary key default gen_random_uuid(),
    device_id varchar(255) not null unique,
    alias varchar(255),
    host varchar(255) not null,
    port integer,
    port_gateway integer,
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
    d.host,
    d.port,
    d.port_gateway,
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

-- Function to notify on device changes
create or replace function notify_device_change()
returns trigger AS $$
declare
    payload JSON;
begin
    payload = json_build_object(
        'device_id', new.device_id,
        'operation', TG_OP
    );
    perform pg_notify('device_changes', payload::text);
    return new;
end;
$$ language plpgsql;

-- Trigger to notify on device insertions and deletions
create trigger trigger_notify_device_change
    after insert or delete on devices
    for each row
    execute function notify_device_change();

-- Default device state (demo)
insert into devices (device_id, alias, host, port, port_gateway, architecture, os, supported_protocols) values
    ('ubiquiti-device-router-3c2d', 'Dream Machine Pro Max', 'ubiquiti-device-router', 8080, 8081, 'arm64', 'linux', array['PROTOCOL_GRPC'::device_protocol]),
    ('ubiquiti-device-switch-b87f', 'Pro Max 24 PoE', 'ubiquiti-device-switch', 8080, 8081, 'amd64', 'linux', array['PROTOCOL_GRPC_STREAM'::device_protocol]);
