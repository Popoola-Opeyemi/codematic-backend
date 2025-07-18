-- name: GetWalletByID :one
SELECT * FROM wallets WHERE id = $1;

-- name: UpdateWalletBalance :exec
UPDATE wallets SET balance = $1, updated_at = now() WHERE id = $2;

-- name: ListWalletsByUserID :many
SELECT * FROM wallets WHERE user_id = $1 ORDER BY created_at DESC; 