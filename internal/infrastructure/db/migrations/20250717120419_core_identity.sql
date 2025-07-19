-- +goose Up
-- +goose StatementBegin

CREATE TABLE "tenants" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "slug" VARCHAR UNIQUE NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY,
  "tenant_id" UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "email" VARCHAR UNIQUE NOT NULL,
  "phone" VARCHAR,
  "password_hash" VARCHAR NOT NULL,
  "role" VARCHAR DEFAULT 'USER' CHECK (
    role IN ('PLATFORM_ADMIN', 'TENANT_ADMIN', 'USER')
  ),
  "is_active" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMPTZ DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- Seed Tenants
INSERT INTO tenants (id, name, slug) VALUES
  ('5c5c14bb-47f5-479f-ba19-01f311cbdd87', 'Dangote Group', 'dangote-group'),
  ('b91f85eb-d7df-4e0c-b82b-a1de968c0264', 'Jumia', 'jumia'),
  ('d03db4dc-3406-45e0-86b5-7542c6cffd61', 'Banqroll Payments', 'banqroll-payments');

-- Seed Platform Admin User
INSERT INTO users (id, tenant_id, email, phone, password_hash, role) VALUES (
  '73db18b6-1b9e-4aa3-aed2-46688a5b3f7e',
  'd03db4dc-3406-45e0-86b5-7542c6cffd61',
  'admin@codematic.com',
  '+2348012345678',
  '$2a$12$JYBNhAsvs5kKDGHJR71/r.f5mhNC9uRdaFtvYBnVVuMvlm7V9aSci',
  'PLATFORM_ADMIN'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tenants" CASCADE;

-- +goose StatementEnd
