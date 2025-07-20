
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

-- name: GetIdempotencyByKeyAndEndpoint :one
SELECT
  id,
  tenant_id,
  user_id,
  idempotency_key,
  endpoint,
  request_hash,
  response_body,
  status_code,
  created_at,
  updated_at
FROM idempotency_keys
WHERE tenant_id      = $1
  AND idempotency_key = $2
  AND endpoint        = $3
LIMIT 1;


