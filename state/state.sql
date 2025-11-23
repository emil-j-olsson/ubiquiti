-- Enums
create type device_protocol as enum (
    'DEVICE_PROTOCOL_UNSPECIFIED',
    'DEVICE_PROTOCOL_HTTP',
    'DEVICE_PROTOCOL_HTTP_STREAM',
    'DEVICE_PROTOCOL_GRPC',
    'DEVICE_PROTOCOL_GRPC_STREAM'
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
create table if not exists device (
    id uuid primary key default gen_random_uuid(),
    device_id varchar(255) not null unique,
    alias varchar(255),
    architecture varchar(50) not null,
    os varchar(50) not null,
    supported_protocols device_protocol[] not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists device_state (
    id uuid primary key default gen_random_uuid(),
    device_id uuid not null references device(id) on delete cascade,
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
create index if not exists idx_device_device_id on device(device_id);
create index if not exists idx_device_state_device_id on device_state(device_id);
create index if not exists idx_device_state_timestamp on device_state(timestamp desc);
create index if not exists idx_device_state_device_timestamp on device_state(device_id, timestamp desc);

-- Composite index for dashboard queries (latest state per device)
create index if not exists idx_device_state_latest on device_state(device_id, timestamp desc, device_status);

-- View for latest device state (optimized for dashboard)
create or replace view device_latest_state as
select distinct on (d.device_id)
    d.id,
    d.device_id,
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
from device d
left join device_state ds on d.id = ds.device_id
order by d.device_id, ds.timestamp desc nulls last;

-- Function to update device.updated_at timestamp
create or replace function update_device_updated_at()
returns trigger as $$
begin
    update device set updated_at = now() where id = new.device_id;
    return new;
end;
$$ language plpgsql;

-- Trigger to auto-update device.updated_at when device state changes
create trigger trigger_update_device_timestamp
    after insert on device_state
    for each row
    execute function update_device_updated_at();
