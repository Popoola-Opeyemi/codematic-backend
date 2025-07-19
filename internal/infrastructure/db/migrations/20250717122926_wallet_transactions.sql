-- +goose up
-- +goose statementbegin

CREATE TABLE "providers" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "code" VARCHAR UNIQUE NOT NULL,
  "config" JSONB,
  "is_active" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "currencies" (
  "code" VARCHAR PRIMARY KEY, 
  "name" VARCHAR NOT NULL,
  "symbol" VARCHAR NOT NULL,
  "is_active" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "wallet_types" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "currency" VARCHAR NOT NULL REFERENCES currencies(code),
  "is_active" BOOLEAN DEFAULT true,
  "description" VARCHAR,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "wallets" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  "wallet_type_id" UUID NOT NULL REFERENCES wallet_types(id),
  "balance" DECIMAL(18, 2) DEFAULT 0 NOT NULL,
  "status" VARCHAR NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'frozen', 'closed')),
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "transactions" (
  "id" UUID PRIMARY KEY,
  "tenant_id" UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "wallet_id" UUID NOT NULL REFERENCES wallets(id),
  "provider_id" UUID NOT NULL REFERENCES providers(id),
  "currency_code" VARCHAR NOT NULL REFERENCES currencies(code),
  "reference" VARCHAR UNIQUE NOT NULL,
  "type" VARCHAR NOT NULL CHECK (
    type IN ('deposit', 'withdrawal', 'transfer')
  ),
  "status" VARCHAR NOT NULL CHECK (
    status IN ('pending', 'completed', 'failed')
  ),
  "amount" DECIMAL(18, 2) NOT NULL,
  "fee" DECIMAL(18, 2) DEFAULT 0 NOT NULL,
  "metadata" JSONB,
  "error_reason" VARCHAR,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "transfers" (
  "id" UUID PRIMARY KEY,
  "tenant_id" UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "sender_wallet_id" UUID NOT NULL REFERENCES wallets(id),
  "receiver_wallet_id" UUID NOT NULL REFERENCES wallets(id),
  "transaction_id" UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  "amount" DECIMAL(18, 2) NOT NULL,
  "status" VARCHAR NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- Seed currencies
INSERT INTO currencies (code, name, symbol) VALUES 
  ('USD', 'US Dollar', '$'),
  ('NGN', 'Nigerian Naira', '₦'),
  ('EUR', 'Euro', '€'),
  ('GBP', 'British Pound', '£');

-- Seed wallet types
INSERT INTO wallet_types (id, name, currency, description) VALUES
  ('aabdd0a6-e35a-4788-85c4-598fbbb12d9e', 'Naira Wallet', 'NGN', 'Wallet for Nigerian Naira'),
  ('ac18e3c9-dfc3-4d2a-a8a9-cf95a0359346', 'Dollar Wallet', 'USD', 'Wallet for United States Dollar'),
  ('ca10450e-40f1-41d2-af19-23f3dcd9d5a8', 'Pound Wallet', 'GBP', 'Wallet for British Pound');

INSERT INTO providers (id, name, code, config, is_active)
VALUES
  (
    gen_random_uuid(),
    'Paystack',
    'paystack',
    '{
      "base_url": "https://api.paystack.co",
      "secret_key": "sk_test_6d247fc20f3e89b2702be48d926bbd86d4e7239b",
      "public_key": "pk_test_008b51881a11a8b7c5ae04a548e0c0b65328153b",
      "webhook_secret": "PAYSTACK_WEBHOOK_SECRET"
    }'::jsonb,
    true
  ),
  (
    gen_random_uuid(),
    'Flutterwave',
    'flutterwave',
    '{
      "base_url": "https://api.flutterwave.com/v3",
      "secret_key": "FLWSECK_TEST-b77bfdd76bd39ab9bd9edc4fd33f6154-X",
      "public_key": "FLWPUBK_TEST-a52a0b9da395cbab3ad412b40cc608c5-X",
      "webhook_secret": "FLUTTERWAVE_WEBHOOK_SECRET",
      "encryption_key":"FLWSECK_TESTa320cb3aefe7"
    }'::jsonb,
    true
  ),
  (
    gen_random_uuid(),
    'Stripe',
    'stripe',
    '{
      "base_url": "https://api.stripe.com",
      "secret_key": "STRIPE_SECRET_KEY",
      "webhook_secret": "STRIPE_WEBHOOK_SECRET"
    }'::jsonb,
    true
  );

-- +goose statementend


-- +goose down
-- +goose statementbegin

DROP TABLE IF EXISTS "transfers" CASCADE;
DROP TABLE IF EXISTS "transactions" CASCADE;
DROP TABLE IF EXISTS "wallets" CASCADE;
DROP TABLE IF EXISTS "wallet_types" CASCADE;
DROP TABLE IF EXISTS "providers" CASCADE;
DROP TABLE IF EXISTS "currencies" CASCADE;

-- +goose statementend
