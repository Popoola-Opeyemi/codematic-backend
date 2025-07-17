-- +goose up
-- +goose statementbegin

create table "webhook_events" (
  "id" uuid primary key,
  "provider_id" uuid not null references providers(id) on delete cascade,
  "event_type" varchar not null,
  "payload" jsonb not null,
  "status" varchar not null default 'received', -- 'received', 'processed', 'failed'
  "attempts" integer default 0,
  "last_error" varchar,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null
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

-- +goose statementend


-- +goose down
-- +goose statementbegin

drop table if exists "virtual_accounts" cascade;
drop table if exists "audit_logs" cascade;
drop table if exists "webhook_events" cascade;

-- +goose statementend
