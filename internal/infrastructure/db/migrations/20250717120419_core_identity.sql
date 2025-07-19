-- +goose Up
-- +goose StatementBegin

CREATE TABLE "tenants" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "slug" VARCHAR UNIQUE NOT NULL,
  "webhook_url" TEXT DEFAULT '' NOT NULL,
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

-- Seed 1 Tenant
INSERT INTO tenants (id, name, slug) VALUES
  ('5c5c14bb-47f5-479f-ba19-01f311cbdd87', 'Dangote Group', 'dangote-group'),
  ('b91f85eb-d7df-4e0c-b82b-a1de968c0264', 'Jumia', 'jumia'),
  ('d03db4dc-3406-45e0-86b5-7542c6cffd61', 'Banqroll Payments', 'banqroll-payments');
-- Seed 1 Platform Admin
INSERT INTO users (id, tenant_id, email, phone, password_hash, role) VALUES (
  'a1111111-1111-1111-1111-111111111111',
  'd03db4dc-3406-45e0-86b5-7542c6cffd61',
  'admin@platform.com',
  '+2348012345678',
  '$2a$12$uYX1s4J1/H.nNwocsEQ75uQbDB.9HepoclC7vjJ1OHZ9AzIHA3VIC',
  'PLATFORM_ADMIN'
);

-- Seed 1 Tenant Admin
INSERT INTO users (id, tenant_id, email, phone, password_hash, role) VALUES (
  'b2222222-2222-2222-2222-222222222222',
  'd03db4dc-3406-45e0-86b5-7542c6cffd61',
  'admin@tenant.com',
  '+2348023456789',
  '$2a$12$uYX1s4J1/H.nNwocsEQ75uQbDB.9HepoclC7vjJ1OHZ9AzIHA3VIC',
  'TENANT_ADMIN'
);

-- Seed 1 Regular User
INSERT INTO users (id, tenant_id, email, phone, password_hash, role) VALUES (
  'c3333333-3333-3333-3333-333333333333',
  'd03db4dc-3406-45e0-86b5-7542c6cffd61',
  'user@tenant.com',
  '+2348034567890',
  '$2a$12$uYX1s4J1/H.nNwocsEQ75uQbDB.9HepoclC7vjJ1OHZ9AzIHA3VIC',
  'USER'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tenants" CASCADE;

-- +goose StatementEnd
