// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: transactions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

const createTransaction = `-- name: CreateTransaction :one
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
) RETURNING id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at
`

type CreateTransactionParams struct {
	ID           pgtype.UUID
	TenantID     pgtype.UUID
	WalletID     pgtype.UUID
	ProviderID   pgtype.UUID
	CurrencyCode string
	Reference    string
	Type         string
	Status       string
	Amount       decimal.Decimal
	Fee          decimal.Decimal
	Metadata     []byte
	ErrorReason  pgtype.Text
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, createTransaction,
		arg.ID,
		arg.TenantID,
		arg.WalletID,
		arg.ProviderID,
		arg.CurrencyCode,
		arg.Reference,
		arg.Type,
		arg.Status,
		arg.Amount,
		arg.Fee,
		arg.Metadata,
		arg.ErrorReason,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.WalletID,
		&i.ProviderID,
		&i.CurrencyCode,
		&i.Reference,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Fee,
		&i.Metadata,
		&i.ErrorReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE id = $1
`

func (q *Queries) GetTransactionByID(ctx context.Context, id pgtype.UUID) (Transaction, error) {
	row := q.db.QueryRow(ctx, getTransactionByID, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.WalletID,
		&i.ProviderID,
		&i.CurrencyCode,
		&i.Reference,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Fee,
		&i.Metadata,
		&i.ErrorReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTransactionByReference = `-- name: GetTransactionByReference :one
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE reference = $1 LIMIT 1
`

func (q *Queries) GetTransactionByReference(ctx context.Context, reference string) (Transaction, error) {
	row := q.db.QueryRow(ctx, getTransactionByReference, reference)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.WalletID,
		&i.ProviderID,
		&i.CurrencyCode,
		&i.Reference,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Fee,
		&i.Metadata,
		&i.ErrorReason,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllTransactions = `-- name: ListAllTransactions :many
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

type ListAllTransactionsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListAllTransactions(ctx context.Context, arg ListAllTransactionsParams) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, listAllTransactions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.WalletID,
			&i.ProviderID,
			&i.CurrencyCode,
			&i.Reference,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Fee,
			&i.Metadata,
			&i.ErrorReason,
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

const listTransactionsByStatus = `-- name: ListTransactionsByStatus :many
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type ListTransactionsByStatusParams struct {
	Status string
	Limit  int32
	Offset int32
}

func (q *Queries) ListTransactionsByStatus(ctx context.Context, arg ListTransactionsByStatusParams) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, listTransactionsByStatus, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.WalletID,
			&i.ProviderID,
			&i.CurrencyCode,
			&i.Reference,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Fee,
			&i.Metadata,
			&i.ErrorReason,
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

const listTransactionsByTenantID = `-- name: ListTransactionsByTenantID :many
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type ListTransactionsByTenantIDParams struct {
	TenantID pgtype.UUID
	Limit    int32
	Offset   int32
}

func (q *Queries) ListTransactionsByTenantID(ctx context.Context, arg ListTransactionsByTenantIDParams) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, listTransactionsByTenantID, arg.TenantID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.WalletID,
			&i.ProviderID,
			&i.CurrencyCode,
			&i.Reference,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Fee,
			&i.Metadata,
			&i.ErrorReason,
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

const listTransactionsByUserID = `-- name: ListTransactionsByUserID :many
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE wallet_id IN (SELECT id FROM wallets WHERE user_id = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type ListTransactionsByUserIDParams struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) ListTransactionsByUserID(ctx context.Context, arg ListTransactionsByUserIDParams) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, listTransactionsByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.WalletID,
			&i.ProviderID,
			&i.CurrencyCode,
			&i.Reference,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Fee,
			&i.Metadata,
			&i.ErrorReason,
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

const listTransactionsByWalletID = `-- name: ListTransactionsByWalletID :many
SELECT id, tenant_id, wallet_id, provider_id, currency_code, reference, type, status, amount, fee, metadata, error_reason, created_at, updated_at FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type ListTransactionsByWalletIDParams struct {
	WalletID pgtype.UUID
	Limit    int32
	Offset   int32
}

func (q *Queries) ListTransactionsByWalletID(ctx context.Context, arg ListTransactionsByWalletIDParams) ([]Transaction, error) {
	rows, err := q.db.Query(ctx, listTransactionsByWalletID, arg.WalletID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.WalletID,
			&i.ProviderID,
			&i.CurrencyCode,
			&i.Reference,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Fee,
			&i.Metadata,
			&i.ErrorReason,
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

const updateTransactionStatusAndAmount = `-- name: UpdateTransactionStatusAndAmount :exec
UPDATE transactions
SET status = $1, amount = $2, updated_at = now()
WHERE id = $3
`

type UpdateTransactionStatusAndAmountParams struct {
	Status string
	Amount decimal.Decimal
	ID     pgtype.UUID
}

func (q *Queries) UpdateTransactionStatusAndAmount(ctx context.Context, arg UpdateTransactionStatusAndAmountParams) error {
	_, err := q.db.Exec(ctx, updateTransactionStatusAndAmount, arg.Status, arg.Amount, arg.ID)
	return err
}
