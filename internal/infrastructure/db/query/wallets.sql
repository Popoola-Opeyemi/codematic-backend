-- name: CreateWallet :one
INSERT INTO wallets (id, user_id, wallet_type_id, balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWalletByID :one
SELECT * FROM wallets WHERE id = $1;

-- name: ListActiveWalletTypes :many
SELECT * FROM wallet_types WHERE is_active ORDER BY currency ASC;

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
INSERT INTO wallets (id, user_id, wallet_type_id, balance, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListWalletsByCurrency :many
SELECT w.*
FROM wallets w
JOIN wallet_types wt ON w.wallet_type_id = wt.id
WHERE wt.currency = $1
ORDER BY w.created_at DESC;

-- name: ListWalletsByUserAndCurrency :many
SELECT w.*
FROM wallets w
JOIN wallet_types wt ON w.wallet_type_id = wt.id
WHERE w.user_id = $1 AND wt.currency = $2
ORDER BY w.created_at DESC;

-- name: GetWalletTypeIDByCurrency :one
SELECT id
FROM wallet_types
WHERE currency = $1
LIMIT 1;

-- name: GetWalletByUserAndCurrency :one
SELECT w.*
FROM wallets w
JOIN wallet_types wt ON w.wallet_type_id = wt.id
WHERE w.user_id = $1 AND wt.currency = $2
LIMIT 1;

-- name: CreateDeposit :one
INSERT INTO deposits (user_id, transaction_id, external_txid, amount, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, transaction_id, external_txid, amount, status, created_at, updated_at;

-- name: GetDepositByID :one
SELECT id, user_id, transaction_id, external_txid, amount, status, created_at, updated_at
FROM deposits
WHERE id = $1;

-- name: CreateWithdrawal :one
INSERT INTO withdrawals (user_id, transaction_id, external_txid, amount, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, transaction_id, external_txid, amount, status, created_at, updated_at;

-- name: GetWithdrawalByID :one
SELECT id, user_id, transaction_id, external_txid, amount, status, created_at, updated_at
FROM withdrawals
WHERE id = $1;
