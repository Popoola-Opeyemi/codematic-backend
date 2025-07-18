-- name: CreateTenant :one
INSERT INTO tenants (id, name, slug, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants
WHERE id = $1;

-- name: GetTenantBySlug :one
SELECT * FROM tenants
WHERE slug = $1;

-- name: ListTenants :many
SELECT * FROM tenants
ORDER BY created_at DESC;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2,
    slug = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = $1;
