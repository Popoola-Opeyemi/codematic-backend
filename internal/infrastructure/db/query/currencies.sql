-- name: CreateCurrency :one
INSERT INTO currencies (code, name, symbol)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCurrencyByCode :one
SELECT * FROM currencies WHERE code = $1;

-- name: ListActiveCurrencies :many
SELECT * FROM currencies WHERE is_active = true ORDER BY code ASC;

-- name: ListAllCurrencies :many
SELECT * FROM currencies ORDER BY code ASC;

-- name: ListActiveCurrencyCodes :many
SELECT code FROM currencies WHERE is_active = true ORDER BY code ASC;

-- name: UpdateCurrency :exec
UPDATE currencies 
SET name = $1, symbol = $2, is_active = $3, updated_at = now()
WHERE code = $4;

-- name: DeactivateCurrency :exec
UPDATE currencies 
SET is_active = false, updated_at = now()
WHERE code = $1;

-- name: ActivateCurrency :exec
UPDATE currencies 
SET is_active = true, updated_at = now()
WHERE code = $1;

-- name: CurrencyExists :one
SELECT EXISTS (
  SELECT 1 FROM currencies WHERE code = $1 AND is_active = true
);

-- name: CurrencyExistsAnyStatus :one
SELECT EXISTS (
  SELECT 1 FROM currencies WHERE code = $1
);
