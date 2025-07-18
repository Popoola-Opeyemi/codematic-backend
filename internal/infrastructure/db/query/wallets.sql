-- name: CreateWallet :one
INSERT INTO wallets (id, user_id, wallet_type_id, balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWalletByID :one
SELECT * FROM wallets WHERE id = $1;

-- name: ListWalletsByUserID :many
SELECT * FROM wallets WHERE user_id = $1 ORDER BY created_at DESC; 

-- name: UpdateWalletBalance :exec
UPDATE wallets SET balance = $1, updated_at = now() WHERE id = $2;

-- name: UpdateWalletType :exec
UPDATE wallets SET wallet_type_id = $1, updated_at = now() WHERE id = $2;

-- name: ListWalletsByType :many
SELECT * FROM wallets WHERE wallet_type_id = $1 ORDER BY created_at DESC;

-- name: IncrementWalletBalance :exec
UPDATE wallets SET balance = balance + $1, updated_at = now() WHERE id = $2;

-- name: DecrementWalletBalance :exec
UPDATE wallets SET balance = balance - $1, updated_at = now() WHERE id = $2;

-- name: DeleteWallet :exec
DELETE FROM wallets WHERE id = $1;


-- name: CreateWalletWithCurrency :one
INSERT INTO wallets (id, user_id, wallet_type_id, currency_code, balance)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListWalletsByCurrency :many
SELECT * FROM wallets WHERE currency_code = $1 ORDER BY created_at DESC;

-- name: ListWalletsByUserAndCurrency :many
SELECT * FROM wallets WHERE user_id = $1 AND currency_code = $2 ORDER BY created_at DESC;
