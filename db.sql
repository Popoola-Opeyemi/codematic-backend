CREATE TABLE "tenants" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "slug" varchar UNIQUE,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid,
  "email" varchar UNIQUE,
  "phone" varchar,
  "password_hash" varchar,
  "is_active" boolean,
  "role" varchar DEFAULT 'USER',
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "wallet_types" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "currency" varchar,
  "description" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "wallets" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "wallet_type_id" uuid,
  "balance" decimal,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "providers" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "code" varchar UNIQUE,
  "config" jsonb,
  "is_active" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "transactions" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid,
  "wallet_id" uuid,
  "provider_id" uuid,
  "reference" varchar UNIQUE,
  "type" varchar,
  "status" varchar,
  "amount" decimal,
  "fee" decimal,
  "metadata" jsonb,
  "error_reason" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "transfers" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid,
  "sender_wallet_id" uuid,
  "receiver_wallet_id" uuid,
  "transaction_id" uuid,
  "amount" decimal,
  "status" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "webhook_events" (
  "id" uuid PRIMARY KEY,
  "provider_id" uuid,
  "event_type" varchar,
  "payload" jsonb,
  "status" varchar,
  "attempts" int,
  "last_error" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "audit_logs" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid,
  "user_id" uuid,
  "action" varchar,
  "metadata" jsonb,
  "created_at" timestamp
);

CREATE TABLE "virtual_accounts" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid,
  "user_id" uuid,
  "provider_id" uuid,
  "account_number" varchar,
  "bank_name" varchar,
  "currency" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("wallet_type_id") REFERENCES "wallet_types" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("provider_id") REFERENCES "providers" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("sender_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("receiver_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "webhook_events" ADD FOREIGN KEY ("provider_id") REFERENCES "providers" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "virtual_accounts" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "virtual_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "virtual_accounts" ADD FOREIGN KEY ("provider_id") REFERENCES "providers" ("id");
