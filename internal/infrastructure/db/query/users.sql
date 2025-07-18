-- name: CreateUser :one
INSERT INTO users (id, tenant_id, email, phone, password_hash, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByEmailAndTenantID :one
SELECT * FROM users
WHERE email = $1 AND tenant_id = $2;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = now()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
