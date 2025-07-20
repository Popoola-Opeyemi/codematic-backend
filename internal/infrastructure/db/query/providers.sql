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

-- name: AddSupportedCurrency :exec
INSERT INTO provider_supported_currencies (
  provider_id, currency_code
) VALUES (
  $1, $2
)
ON CONFLICT DO NOTHING;

-- name: AddSupportedChannel :exec
INSERT INTO provider_supported_channels (
  provider_id, channel
) VALUES (
  $1, $2
)
ON CONFLICT DO NOTHING;

-- name: RemoveSupportedCurrency :exec
DELETE FROM provider_supported_currencies
WHERE provider_id = $1 AND currency_code = $2;

-- name: ListSupportedCurrencies :many
SELECT c.*
FROM provider_supported_currencies psc
JOIN currencies c ON psc.currency_code = c.code
WHERE psc.provider_id = $1
ORDER BY c.name;

-- name: DeleteAllSupportedCurrencies :exec
DELETE FROM provider_supported_currencies
WHERE provider_id = $1;

-- name: ListProviderDetails :many
SELECT
  p.id AS provider_id,
  p.name AS provider_name,
  p.code AS provider_code,
  p.config,
  p.is_active,
  p.created_at,
  p.updated_at,
  ARRAY_AGG(DISTINCT c.code ORDER BY c.code) AS currency_codes,
  ARRAY_AGG(DISTINCT ch.channel ORDER BY ch.channel) AS supported_channels
FROM providers p
LEFT JOIN provider_supported_currencies psc ON p.id = psc.provider_id
LEFT JOIN currencies c ON psc.currency_code = c.code
LEFT JOIN provider_supported_channels ch ON p.id = ch.provider_id
GROUP BY p.id, p.name, p.code, p.config, p.is_active, p.created_at, p.updated_at
ORDER BY p.name;

-- name: GetProviderMetrics :one
SELECT * FROM provider_metrics
WHERE provider_id = $1;

-- name: CreateProviderMetrics :exec
INSERT INTO provider_metrics (provider_id)
VALUES ($1);

-- name: IncrementSuccess :exec
UPDATE provider_metrics
SET
  success_count = success_count + 1,
  priority = GREATEST(priority - 10, 0),
  last_success_at = now(),
  updated_at = now()
WHERE provider_id = $1;

-- name: IncrementFailure :exec
UPDATE provider_metrics
SET
  failure_count = failure_count + 1,
  priority = priority + 20,
  last_failure_at = now(),
  updated_at = now()
WHERE provider_id = $1;

-- name: ResetDailyMetrics :exec
UPDATE provider_metrics
SET
  success_count = 0,
  failure_count = 0,
  priority = 100,
  updated_at = now();

-- name: DecayPriority :exec
UPDATE provider_metrics
SET
  priority = GREATEST(priority - 5, 0),
  updated_at = now();

-- name: RemoveSupportedChannel :exec
DELETE FROM provider_supported_channels
WHERE provider_id = $1 AND channel = $2;

-- name: SelectBestProvider :one
SELECT 
  p.id, p.name, p.code, p.config,
  COALESCE(m.priority, 100) as priority,
  COALESCE(m.success_count, 0) as success_count,
  COALESCE(m.failure_count, 0) as failure_count
FROM providers p
LEFT JOIN provider_metrics m ON p.id = m.provider_id
WHERE p.is_active = true
ORDER BY priority ASC, success_count DESC
LIMIT 1;

-- name: SelectBestProviderByCurrencyAndChannel :one
SELECT 
  p.id, p.name, p.code, p.config,
  COALESCE(m.priority, 100) as priority,
  COALESCE(m.success_count, 0) as success_count,
  COALESCE(m.failure_count, 0) as failure_count
FROM providers p
LEFT JOIN provider_metrics m ON p.id = m.provider_id
JOIN provider_supported_currencies pc ON p.id = pc.provider_id
JOIN provider_supported_channels ch ON p.id = ch.provider_id
WHERE p.is_active = true
  AND pc.currency_code = $1
  AND ch.channel = $2
ORDER BY priority ASC, success_count DESC
LIMIT 1;
