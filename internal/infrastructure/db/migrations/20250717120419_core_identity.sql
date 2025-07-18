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

insert into "tenants" ("id", "name", "slug")
values
  ('5c5c14bb-47f5-479f-ba19-01f311cbdd87', 'Dangote Group', 'dangote-group'),
  ('b91f85eb-d7df-4e0c-b82b-a1de968c0264', 'Jumia', 'jumia'),
  ('d03db4dc-3406-45e0-86b5-7542c6cffd61', 'Banqroll Payments', 'banqroll-payments');


insert into "wallet_types" ("id", "name", "currency", "description")
values
  ('aabdd0a6-e35a-4788-85c4-598fbbb12d9e', 'Naira Wallet', 'NGN', 'Wallet for Nigerian Naira'),
  ('ac18e3c9-dfc3-4d2a-a8a9-cf95a0359346', 'Dollar Wallet', 'USD', 'Wallet for United States Dollar'),
  ('ca10450e-40f1-41d2-af19-23f3dcd9d5a8', 'Pound Wallet', 'GBP', 'Wallet for British Pound');

-- +goose statementend

-- +goose down
-- +goose statementbegin

-- safe down migration
drop table if exists "wallet_types" cascade;
drop table if exists "users" cascade;
drop table if exists "tenants" cascade;

-- +goose statementend
