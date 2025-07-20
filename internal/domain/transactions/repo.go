package transactions

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &repository{
		q: q,
		p: pool,
	}
}

func (r *repository) GetTransactionByID(ctx context.Context, id string) (*Transaction, error) {
	uuid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to uuid: %w", err)
	}

	tx, err := r.q.GetTransactionByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}

	return dbToDomainTransaction(&tx), nil
}

func (r *repository) ListTransactionsByUserID(ctx context.Context,
	userID string, limit, offset int) ([]*Transaction, error) {

	uuid, err := utils.StringToPgUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user id to uuid: %w", err)
	}

	params := db.ListTransactionsByUserIDParams{
		UserID: uuid,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	txs, err := r.q.ListTransactionsByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by user id: %w", err)
	}

	return dbToDomainTransactions(txs), nil
}

func (r *repository) ListTransactionsByTenantID(ctx context.Context,
	tenantID string, limit, offset int) ([]*Transaction, error) {
	uuid, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tenant id to uuid: %w", err)
	}
	params := db.ListTransactionsByTenantIDParams{
		TenantID: uuid,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	txs, err := r.q.ListTransactionsByTenantID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by tenant id: %w", err)
	}
	return dbToDomainTransactions(txs), nil
}

func (r *repository) ListAllTransactions(ctx context.Context, limit,
	offset int) ([]*Transaction, error) {

	params := db.ListAllTransactionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	txs, err := r.q.ListAllTransactions(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}

	return dbToDomainTransactions(txs), nil
}

func (r *repository) ListTransactionsByStatus(ctx context.Context, status string,
	limit, offset int) ([]*Transaction, error) {

	params := db.ListTransactionsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	txs, err := r.q.ListTransactionsByStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by status: %w", err)
	}

	return dbToDomainTransactions(txs), nil
}

func dbToDomainTransaction(tx *db.Transaction) *Transaction {
	return &Transaction{
		ID:           utils.FromPgUUID(tx.ID),
		TenantID:     utils.FromPgUUID(tx.TenantID),
		WalletID:     utils.FromPgUUID(tx.WalletID),
		ProviderID:   utils.FromPgUUID(tx.ProviderID),
		CurrencyCode: tx.CurrencyCode,
		Reference:    tx.Reference,
		Type:         tx.Type,
		Status:       tx.Status,
		Amount:       tx.Amount,
		Fee:          tx.Fee,
		Metadata:     utils.JSONBToMap(tx.Metadata),
		ErrorReason:  tx.ErrorReason.String,
		CreatedAt:    utils.FromPgTimestamptz(tx.CreatedAt),
		UpdatedAt:    utils.FromPgTimestamptz(tx.UpdatedAt),
	}
}

func dbToDomainTransactions(txs []db.Transaction) []*Transaction {
	result := make([]*Transaction, len(txs))
	for i, tx := range txs {
		result[i] = dbToDomainTransaction(&tx)
	}
	return result
}
