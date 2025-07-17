-- +goose up
-- +goose statementbegin

create table "tenants" (
  "id" uuid primary key,
  "name" varchar not null,
  "slug" varchar unique not null,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null
);

create table "users" (
  "id" uuid primary key,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  "email" varchar unique not null,
  "phone" varchar,
  "password_hash" varchar not null,
  "is_active" boolean default true,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null
);

create table "wallet_types" (
  "id" uuid primary key,
  "name" varchar not null,
  "currency" varchar not null, -- usd, naira, gbp
  "description" varchar,
  "created_at" timestamp with time zone default now() not null,
  "updated_at" timestamp with time zone default now() not null
);

-- +goose statementend

-- +goose down
-- +goose statementbegin

-- safe down migration
drop table if exists "wallet_types" cascade;
drop table if exists "users" cascade;
drop table if exists "tenants" cascade;

-- +goose statementend
