-- name: CreateTenant :one
INSERT INTO tenants (id, name, slug, webhook_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants WHERE id = $1;

-- name: GetTenantBySlug :one
SELECT * FROM tenants WHERE slug = $1;

-- name: ListTenants :many
SELECT * FROM tenants;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2, slug = $3, webhook_url = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants WHERE id = $1;
