-- +goose up
-- +goose statementbegin

create table "webhook_events" (
  "id" uuid primary key,
  "provider_id" uuid not null references providers(id) on delete cascade,
  "provider_event_id" varchar not null,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  is_outgoing boolean DEFAULT false,
  "event_type" varchar not null,
  "payload" jsonb not null,
  "status" varchar not null default 'received', -- 'received', 'processed', 'failed'
  "attempts" integer default 0,
  "last_error" varchar,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null,
  unique ("provider_id", "provider_event_id")
);

create table "audit_logs" (
  "id" uuid primary key,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  "user_id" uuid not null references users(id) on delete set null,
  "action" varchar not null,
  "metadata" jsonb,
  "created_at" timestamp with time zone default now() not null
);

create table "virtual_accounts" (
  "id" uuid primary key,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  "user_id" uuid not null references users(id) on delete cascade,
  "provider_id" uuid not null references providers(id),
  "account_number" varchar not null,
  "bank_name" varchar not null,
  "currency" varchar not null,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null
);

-- Deposits table
CREATE TABLE IF NOT EXISTS deposits (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    external_txid VARCHAR(255),
    amount NUMERIC(20, 2) NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Withdrawals table
CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    external_txid VARCHAR(255),
    amount NUMERIC(20, 2) NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose statementend


-- +goose down
-- +goose statementbegin

drop table if exists "virtual_accounts" cascade;
drop table if exists "audit_logs" cascade;
drop table if exists "webhook_events" cascade;

-- +goose statementend
