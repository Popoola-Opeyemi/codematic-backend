-- name: CreateProvider :one
INSERT INTO providers (
  name, code, config, is_active
) VALUES (
  $1, $2, $3, COALESCE($4, true)
)
RETURNING *;

-- name: GetProviderByID :one
SELECT * FROM providers
WHERE id = $1;

-- name: GetProviderByCode :one
SELECT * FROM providers
WHERE code = $1;

-- name: ListActiveProviders :many
SELECT * FROM providers
WHERE is_active = true
ORDER BY name;

-- name: UpdateProviderConfig :one
UPDATE providers
SET config = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeactivateProvider :exec
UPDATE providers
SET is_active = false,
    updated_at = now()
WHERE id = $1;
