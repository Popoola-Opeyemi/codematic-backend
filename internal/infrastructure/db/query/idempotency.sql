
-- name: CreateIdempotencyKey :exec
INSERT INTO idempotency_keys (
  id,
  tenant_id,
  user_id,
  idempotency_key,
  endpoint,
  request_hash,
  response_body,
  status_code
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (tenant_id, idempotency_key, endpoint) DO NOTHING;


-- name: GetIdempotencyRecord :one
SELECT * FROM idempotency_keys
WHERE tenant_id = $1
  AND idempotency_key = $2
  AND endpoint = $3
  AND request_hash = $4
LIMIT 1;

-- name: UpdateIdempotencyKeyResponse :one
UPDATE idempotency_keys
SET
  response_body = $1,
  status_code = $2,
  updated_at = now()
WHERE tenant_id = $3
  AND idempotency_key = $4
  AND endpoint = $5
RETURNING *;

