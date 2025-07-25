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

-- name: GetTransactionByReference :one
SELECT * FROM transactions WHERE reference = $1 LIMIT 1; 

-- name: UpdateTransactionStatusAndAmount :exec
UPDATE transactions
SET status = $1, amount = $2, updated_at = now()
WHERE id = $3; 

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1;

-- name: ListTransactionsByUserID :many
SELECT * FROM transactions WHERE wallet_id IN (SELECT id FROM wallets WHERE user_id = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListTransactionsByTenantID :many
SELECT * FROM transactions WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListAllTransactions :many
SELECT * FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListTransactionsByStatus :many
SELECT * FROM transactions WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3; 