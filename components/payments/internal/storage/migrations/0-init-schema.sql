create extension if not exists pgcrypto;

-- connectors
create table if not exists connectors (
    -- Mandatory fields
    id         varchar not null,
    name       text not null,
    created_at timestamp without time zone not null,
    provider   text not null,

    -- Optional fields
    config bytea,

    -- Primary key
    primary key (id)
);
create unique index connectors_unique_name on connectors (name);

-- accounts
create table if not exists accounts (
    -- Mandatory fields
    id           varchar not null,
    connector_id varchar not null,
    created_at   timestamp without time zone not null,
    reference    text not null,
    type         text not null,
    raw          json not null,

    -- Optional fields
    default_asset text,
    name          text,

    -- Optional fields with default
    metadata jsonb not null default '{}'::jsonb,

    -- Primary key
    primary key (id)
);
alter table accounts 
    add constraint accounts_connector_id_fk foreign key (connector_id) 
    references connectors (id)
    on delete cascade;

-- balances
create table if not exists balances (
    -- Mandatory fields
    account_id      varchar not null,
    connector_id    varchar not null,
    created_at      timestamp without time zone not null,
    last_updated_at timestamp without time zone not null,
    asset           text not null,
    balance         numeric not null,

    -- Primary key
    primary key (account_id, created_at, asset)
);
create index balances_account_id_created_at_asset on balances (account_id, last_updated_at desc, asset);
alter table balances
    add constraint balances_connector_id foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- bank accounts
create table if not exists bank_accounts (
    -- Mandatory fields
    id uuid    not null,
    created_at timestamp without time zone not null,
    name       text not null,

    -- Optional fields
    account_number bytea,
    iban           bytea,
    swift_bic_code bytea,
    country        text,

    -- Optional fields with default
    metadata jsonb not null default '{}'::jsonb,

    -- Primary key
    primary key (id)
);
create table if not exists bank_accounts_related_accounts (
    -- Mandatory fields
    bank_account_id uuid not null,
    account_id      varchar not null,
    connector_id    varchar not null,
    created_at      timestamp without time zone not null,

    -- Primary key
    primary key (bank_account_id, account_id)
);
alter table bank_accounts_related_accounts
    add constraint bank_accounts_related_accounts_bank_account_id_fk foreign key (bank_account_id)
    references bank_accounts (id)
    on delete cascade;
alter table bank_accounts_related_accounts
    add constraint bank_accounts_related_accounts_account_id_fk foreign key (account_id)
    references accounts (id)
    on delete cascade;
alter table bank_accounts_related_accounts
    add constraint bank_accounts_related_accounts_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- payments
create table if not exists payments (
    -- Mandatory fields
    id             varchar not null,
    connector_id   varchar not null,
    reference      text not null,
    created_at     timestamp without time zone not null,
    type           text not null,
    initial_amount numeric not null,
    amount         numeric not null,
    asset          text not null,
    scheme         text not null,

    -- Optional fields
    source_account_id      varchar,
    destination_account_id varchar,

    -- Optional fields with default
    metadata jsonb not null default '{}'::jsonb,

    -- Primary key
    primary key (id)
);
alter table payments
    add constraint payments_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- payment adjustments
create table if not exists payment_adjustments (
    -- Mandatory fields
    id          varchar not null,
    payment_id  varchar not null,
    reference   text not null,
    created_at  timestamp without time zone not null,
    status      text not null,
    raw         json not null,

    -- Optional fields
    amount numeric,
    asset  text,

    -- Optional fields with default
    metadata jsonb not null default '{}'::jsonb,

    -- Primary key
    primary key (id)
);
alter table payment_adjustments
    add constraint payment_adjustments_payment_id_fk foreign key (payment_id)
    references payments (id)
    on delete cascade;

-- pools
create table if not exists pools (
    -- Mandatory fields
    id         uuid not null,
    name       text not null,
    created_at timestamp without time zone not null,

    -- Primary key
    primary key (id)
);
create unique index pools_unique_name on pools (name);

create table if not exists pool_accounts (
    -- Mandatory fields
    pool_id     uuid not null,
    account_id  varchar not null,

    -- Primary key
    primary key (pool_id, account_id)
);
alter table pool_accounts
    add constraint pool_accounts_pool_id_fk foreign key (pool_id)
    references pools (id)
    on delete cascade;
alter table pool_accounts
    add constraint pool_accounts_account_id_fk foreign key (account_id)
    references accounts (id)
    on delete cascade;

-- schedules
create table if not exists schedules (
    -- Mandatory fields
    id text not null,
    connector_id varchar not null,
    created_at timestamp without time zone not null,
    
    -- Primary key
    primary key (id, connector_id)
);
alter table schedules
    add constraint schedules_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- states
create table if not exists states (
    -- Mandatory fields
    id           varchar not null,
    connector_id varchar not null,

    -- Optional fields with default
    state json not null default '{}'::json,

    -- Primary key
    primary key (id)
);
alter table states
    add constraint states_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- tasks
create table if not exists tasks (
    -- Mandatory fields
    connector_id varchar not null,
    tasks        json not null,

    -- Primary key
    primary key (connector_id)
);
alter table tasks
    add constraint tasks_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- Workflow instance
create table if not exists workflows_instances (
    -- Mandatory fields
    id           text not null,
    schedule_id  text not null,
    connector_id varchar not null,
    created_at   timestamp without time zone not null,
    updated_at   timestamp without time zone not null,

    -- Optional fields with default
    terminated boolean not null default false,

    -- Optional fields
    terminated_at timestamp without time zone,
    error         text,

    -- Primary key
    primary key (id, schedule_id, connector_id)
);
alter table workflows_instances
    add constraint workflows_instances_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;
alter table workflows_instances
    add constraint workflows_instances_schedule_id_fk foreign key (schedule_id, connector_id)
    references schedules (id, connector_id)
    on delete cascade;

-- Webhook configs
create table if not exists webhooks_configs (
    -- Mandatory fields
    name         text not null,
    connector_id varchar not null,
    url_path     text not null,

    -- Primary key
    primary key (name, connector_id)
);
alter table webhooks_configs
    add constraint webhooks_configs_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;

-- Webhooks
create table if not exists webhooks (
    -- Mandatory fields
    id text not null,
    connector_id varchar not null,

    -- Optional fields
    headers json,
    query_values json,
    body bytea,

    -- Primary key
    primary key (id)
);
alter table webhooks
    add constraint webhooks_connector_id_fk foreign key (connector_id)
    references connectors (id)
    on delete cascade;