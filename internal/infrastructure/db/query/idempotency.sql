-- name: GetIdempotencyRecord :one
SELECT * FROM idempotency_keys
WHERE tenant_id = $1
  AND idempotency_key = $2
  AND endpoint = $3
  AND request_hash = $4;

-- name: SaveIdempotencyRecord :one
INSERT INTO idempotency_keys (
  id, tenant_id, user_id, idempotency_key, endpoint, request_hash, response_body, status_code, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, now(), now()
)
ON CONFLICT (tenant_id, idempotency_key, endpoint)
DO UPDATE SET
  response_body = EXCLUDED.response_body,
  status_code = EXCLUDED.status_code,
  updated_at = now()
RETURNING *; 