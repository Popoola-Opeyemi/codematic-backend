-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE idempotency_keys (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid NOT NULL REFERENCES tenants(id),
  "user_id" uuid,
  "idempotency_key" VARCHAR(64) NOT NULL,
  "endpoint" varchar(128) NOT NULL,
  "request_hash" VARCHAR(128) NOT NULL,
  "response_body" jsonb,
  "status_code" int,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  UNIQUE ("tenant_id", "idempotency_key", "endpoint")
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists "idempotency_keys" cascade;

-- +goose StatementEnd
