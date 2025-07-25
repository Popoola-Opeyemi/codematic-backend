// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: providers.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addSupportedChannel = `-- name: AddSupportedChannel :exec
INSERT INTO provider_supported_channels (
  provider_id, channel
) VALUES (
  $1, $2
)
ON CONFLICT DO NOTHING
`

type AddSupportedChannelParams struct {
	ProviderID pgtype.UUID
	Channel    string
}

func (q *Queries) AddSupportedChannel(ctx context.Context, arg AddSupportedChannelParams) error {
	_, err := q.db.Exec(ctx, addSupportedChannel, arg.ProviderID, arg.Channel)
	return err
}

const addSupportedCurrency = `-- name: AddSupportedCurrency :exec
INSERT INTO provider_supported_currencies (
  provider_id, currency_code
) VALUES (
  $1, $2
)
ON CONFLICT DO NOTHING
`

type AddSupportedCurrencyParams struct {
	ProviderID   pgtype.UUID
	CurrencyCode string
}

func (q *Queries) AddSupportedCurrency(ctx context.Context, arg AddSupportedCurrencyParams) error {
	_, err := q.db.Exec(ctx, addSupportedCurrency, arg.ProviderID, arg.CurrencyCode)
	return err
}

const createProvider = `-- name: CreateProvider :one
INSERT INTO providers (
  name, code, config, is_active
) VALUES (
  $1, $2, $3, COALESCE($4, true)
)
RETURNING id, name, code, config, is_active, created_at, updated_at
`

type CreateProviderParams struct {
	Name    string
	Code    string
	Config  []byte
	Column4 interface{}
}

func (q *Queries) CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error) {
	row := q.db.QueryRow(ctx, createProvider,
		arg.Name,
		arg.Code,
		arg.Config,
		arg.Column4,
	)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createProviderMetrics = `-- name: CreateProviderMetrics :exec
INSERT INTO provider_metrics (provider_id)
VALUES ($1)
`

func (q *Queries) CreateProviderMetrics(ctx context.Context, providerID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, createProviderMetrics, providerID)
	return err
}

const deactivateProvider = `-- name: DeactivateProvider :exec
UPDATE providers
SET is_active = false,
    updated_at = now()
WHERE id = $1
`

func (q *Queries) DeactivateProvider(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deactivateProvider, id)
	return err
}

const decayPriority = `-- name: DecayPriority :exec
UPDATE provider_metrics
SET
  priority = GREATEST(priority - 5, 0),
  updated_at = now()
`

func (q *Queries) DecayPriority(ctx context.Context) error {
	_, err := q.db.Exec(ctx, decayPriority)
	return err
}

const deleteAllSupportedCurrencies = `-- name: DeleteAllSupportedCurrencies :exec
DELETE FROM provider_supported_currencies
WHERE provider_id = $1
`

func (q *Queries) DeleteAllSupportedCurrencies(ctx context.Context, providerID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAllSupportedCurrencies, providerID)
	return err
}

const getProviderByCode = `-- name: GetProviderByCode :one
SELECT id, name, code, config, is_active, created_at, updated_at FROM providers
WHERE code = $1
`

func (q *Queries) GetProviderByCode(ctx context.Context, code string) (Provider, error) {
	row := q.db.QueryRow(ctx, getProviderByCode, code)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProviderByID = `-- name: GetProviderByID :one
SELECT id, name, code, config, is_active, created_at, updated_at FROM providers
WHERE id = $1
`

func (q *Queries) GetProviderByID(ctx context.Context, id pgtype.UUID) (Provider, error) {
	row := q.db.QueryRow(ctx, getProviderByID, id)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProviderMetrics = `-- name: GetProviderMetrics :one
SELECT provider_id, priority, success_count, failure_count, last_success_at, last_failure_at, updated_at FROM provider_metrics
WHERE provider_id = $1
`

func (q *Queries) GetProviderMetrics(ctx context.Context, providerID pgtype.UUID) (ProviderMetric, error) {
	row := q.db.QueryRow(ctx, getProviderMetrics, providerID)
	var i ProviderMetric
	err := row.Scan(
		&i.ProviderID,
		&i.Priority,
		&i.SuccessCount,
		&i.FailureCount,
		&i.LastSuccessAt,
		&i.LastFailureAt,
		&i.UpdatedAt,
	)
	return i, err
}

const incrementFailure = `-- name: IncrementFailure :exec
UPDATE provider_metrics
SET
  failure_count = failure_count + 1,
  priority = priority + 20,
  last_failure_at = now(),
  updated_at = now()
WHERE provider_id = $1
`

func (q *Queries) IncrementFailure(ctx context.Context, providerID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, incrementFailure, providerID)
	return err
}

const incrementSuccess = `-- name: IncrementSuccess :exec
UPDATE provider_metrics
SET
  success_count = success_count + 1,
  priority = GREATEST(priority - 10, 0),
  last_success_at = now(),
  updated_at = now()
WHERE provider_id = $1
`

func (q *Queries) IncrementSuccess(ctx context.Context, providerID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, incrementSuccess, providerID)
	return err
}

const listActiveProviders = `-- name: ListActiveProviders :many
SELECT id, name, code, config, is_active, created_at, updated_at FROM providers
WHERE is_active = true
ORDER BY name
`

func (q *Queries) ListActiveProviders(ctx context.Context) ([]Provider, error) {
	rows, err := q.db.Query(ctx, listActiveProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Provider
	for rows.Next() {
		var i Provider
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
			&i.Config,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProviderDetails = `-- name: ListProviderDetails :many
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
ORDER BY p.name
`

type ListProviderDetailsRow struct {
	ProviderID        pgtype.UUID
	ProviderName      string
	ProviderCode      string
	Config            []byte
	IsActive          pgtype.Bool
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
	CurrencyCodes     interface{}
	SupportedChannels interface{}
}

func (q *Queries) ListProviderDetails(ctx context.Context) ([]ListProviderDetailsRow, error) {
	rows, err := q.db.Query(ctx, listProviderDetails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProviderDetailsRow
	for rows.Next() {
		var i ListProviderDetailsRow
		if err := rows.Scan(
			&i.ProviderID,
			&i.ProviderName,
			&i.ProviderCode,
			&i.Config,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CurrencyCodes,
			&i.SupportedChannels,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSupportedCurrencies = `-- name: ListSupportedCurrencies :many
SELECT c.code, c.name, c.symbol, c.is_active, c.created_at, c.updated_at
FROM provider_supported_currencies psc
JOIN currencies c ON psc.currency_code = c.code
WHERE psc.provider_id = $1
ORDER BY c.name
`

func (q *Queries) ListSupportedCurrencies(ctx context.Context, providerID pgtype.UUID) ([]Currency, error) {
	rows, err := q.db.Query(ctx, listSupportedCurrencies, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Currency
	for rows.Next() {
		var i Currency
		if err := rows.Scan(
			&i.Code,
			&i.Name,
			&i.Symbol,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeSupportedChannel = `-- name: RemoveSupportedChannel :exec
DELETE FROM provider_supported_channels
WHERE provider_id = $1 AND channel = $2
`

type RemoveSupportedChannelParams struct {
	ProviderID pgtype.UUID
	Channel    string
}

func (q *Queries) RemoveSupportedChannel(ctx context.Context, arg RemoveSupportedChannelParams) error {
	_, err := q.db.Exec(ctx, removeSupportedChannel, arg.ProviderID, arg.Channel)
	return err
}

const removeSupportedCurrency = `-- name: RemoveSupportedCurrency :exec
DELETE FROM provider_supported_currencies
WHERE provider_id = $1 AND currency_code = $2
`

type RemoveSupportedCurrencyParams struct {
	ProviderID   pgtype.UUID
	CurrencyCode string
}

func (q *Queries) RemoveSupportedCurrency(ctx context.Context, arg RemoveSupportedCurrencyParams) error {
	_, err := q.db.Exec(ctx, removeSupportedCurrency, arg.ProviderID, arg.CurrencyCode)
	return err
}

const resetDailyMetrics = `-- name: ResetDailyMetrics :exec
UPDATE provider_metrics
SET
  success_count = 0,
  failure_count = 0,
  priority = 100,
  updated_at = now()
`

func (q *Queries) ResetDailyMetrics(ctx context.Context) error {
	_, err := q.db.Exec(ctx, resetDailyMetrics)
	return err
}

const selectBestProvider = `-- name: SelectBestProvider :one
SELECT 
  p.id, p.name, p.code, p.config,
  COALESCE(m.priority, 100) as priority,
  COALESCE(m.success_count, 0) as success_count,
  COALESCE(m.failure_count, 0) as failure_count
FROM providers p
LEFT JOIN provider_metrics m ON p.id = m.provider_id
WHERE p.is_active = true
ORDER BY priority ASC, success_count DESC
LIMIT 1
`

type SelectBestProviderRow struct {
	ID           pgtype.UUID
	Name         string
	Code         string
	Config       []byte
	Priority     int32
	SuccessCount int32
	FailureCount int32
}

func (q *Queries) SelectBestProvider(ctx context.Context) (SelectBestProviderRow, error) {
	row := q.db.QueryRow(ctx, selectBestProvider)
	var i SelectBestProviderRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.Priority,
		&i.SuccessCount,
		&i.FailureCount,
	)
	return i, err
}

const selectBestProviderByCurrencyAndChannel = `-- name: SelectBestProviderByCurrencyAndChannel :one
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
LIMIT 1
`

type SelectBestProviderByCurrencyAndChannelParams struct {
	CurrencyCode string
	Channel      string
}

type SelectBestProviderByCurrencyAndChannelRow struct {
	ID           pgtype.UUID
	Name         string
	Code         string
	Config       []byte
	Priority     int32
	SuccessCount int32
	FailureCount int32
}

func (q *Queries) SelectBestProviderByCurrencyAndChannel(ctx context.Context, arg SelectBestProviderByCurrencyAndChannelParams) (SelectBestProviderByCurrencyAndChannelRow, error) {
	row := q.db.QueryRow(ctx, selectBestProviderByCurrencyAndChannel, arg.CurrencyCode, arg.Channel)
	var i SelectBestProviderByCurrencyAndChannelRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.Priority,
		&i.SuccessCount,
		&i.FailureCount,
	)
	return i, err
}

const updateProviderConfig = `-- name: UpdateProviderConfig :one
UPDATE providers
SET config = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, name, code, config, is_active, created_at, updated_at
`

type UpdateProviderConfigParams struct {
	ID     pgtype.UUID
	Config []byte
}

func (q *Queries) UpdateProviderConfig(ctx context.Context, arg UpdateProviderConfigParams) (Provider, error) {
	row := q.db.QueryRow(ctx, updateProviderConfig, arg.ID, arg.Config)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
