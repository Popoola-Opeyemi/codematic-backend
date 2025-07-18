-- +goose up
-- +goose statementbegin

create table "currencies" (
  "code" varchar primary key, 
  "name" varchar not null,
  "symbol" varchar not null, -- e.g. '$', '₦', '€'
  "is_active" boolean default true,
  "created_at" timestamptz default now() not null,
  "updated_at" timestamptz default now() not null
);

create table "providers" (
  "id" uuid primary key,
  "name" varchar not null,
  "code" varchar unique not null, -- 'paystack', 'flutterwave'
  "config" jsonb,
  "is_active" boolean default true,
  "created_at" timestamptz default now() not null,
  "updated_at" timestamptz default now() not null
);

create table "wallets" (
  "id" uuid primary key,
  "user_id" uuid not null references users(id) on delete cascade,
  "wallet_type_id" uuid not null references wallet_types(id),
  "currency_code" varchar not null references currencies(code),
  "balance" decimal(18, 2) default 0,
  "created_at" timestamptz default now() not null,
  "updated_at" timestamptz default now() not null
);

create table "transactions" (
  "id" uuid primary key,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  "wallet_id" uuid not null references wallets(id),
  "provider_id" uuid not null references providers(id),
  "currency_code" varchar not null references currencies(code),
  "reference" varchar unique not null,
  "type" varchar not null check (
    type in ('deposit', 'withdrawal', 'transfer')
  ),
  "status" varchar not null check (
    status in ('pending', 'completed', 'failed')
  ),
  "amount" decimal(18, 2) not null,
  "fee" decimal(18, 2) default 0,
  "metadata" jsonb,
  "error_reason" varchar,
  "created_at" timestamptz default now() not null,
  "updated_at" timestamptz default now() not null
);

create table "transfers" (
  "id" uuid primary key,
  "tenant_id" uuid not null references tenants(id) on delete cascade,
  "sender_wallet_id" uuid not null references wallets(id),
  "receiver_wallet_id" uuid not null references wallets(id),
  "transaction_id" uuid not null references transactions(id) on delete cascade,
  "amount" decimal(18, 2) not null,
  "status" varchar not null,
  "created_at" timestamptz default now() not null,
  "updated_at" timestamptz default now() not null
);

INSERT INTO currencies (code, name, symbol)
VALUES 
  ('USD', 'US Dollar', '$'),
  ('NGN', 'Nigerian Naira', '₦'),
  ('EUR', 'Euro', '€');


-- +goose statementend
-- +goose down
-- +goose statementbegin

drop table if exists "transfers" cascade;
drop table if exists "transactions" cascade;
drop table if exists "wallets" cascade;
drop table if exists "providers" cascade;
drop table if exists "currencies" cascade;

-- +goose statementend
