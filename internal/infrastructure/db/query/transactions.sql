-- name: CreateTransaction :one
INSERT INTO transactions (
  id,
  tenant_id,
  wallet_id,
  provider_id,
  currency_code,
  reference,
  type,
  status,
  amount,
  fee,
  metadata,
  error_reason,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, now(), now()
) RETURNING *;

-- name: ListTransactionsByWalletID :many
SELECT * FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3; 