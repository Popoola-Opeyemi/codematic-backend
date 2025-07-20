package transactions

import (
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"

	"context"
	"time"
)

// service implements the business logic for transaction-related operations.
type service struct {
	DB   *db.DBConn
	Repo Repository

	cache cache.TransactionCacheStore
}

// NewService initializes and returns a new instance of the transaction service.
func NewService(
	db *db.DBConn,
	cacheStore cache.TransactionCacheStore,
) Service {
	return &service{
		Repo:  NewRepository(db.Queries, db.Pool),
		cache: cacheStore,
	}
}

func (s *service) GetTransactionByID(ctx context.Context, id string) (*Transaction, error) {
	var tx Transaction
	if s.cache != nil {
		err := s.cache.GetTransaction(ctx, id, &tx)
		if err == nil && tx.ID != "" {
			return &tx, nil
		}
	}
	dbTx, err := s.Repo.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.cache != nil && dbTx != nil {
		_ = s.cache.SetTransaction(ctx, id, dbTx, 5*time.Minute)
	}
	return dbTx, nil
}

func (s *service) ListTransactionsByUserID(ctx context.Context,
	userID string, limit, offset int) ([]*Transaction, error) {

	var txns []*Transaction
	if s.cache != nil {
		err := s.cache.GetTransactionsByUser(ctx, userID, &txns)
		if err == nil && len(txns) > 0 {
			return txns, nil
		}
	}

	dbTxns, err := s.Repo.ListTransactionsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	if s.cache != nil && len(dbTxns) > 0 {
		_ = s.cache.SetTransactionsByUser(ctx, userID, dbTxns, 5*time.Minute)
	}

	return dbTxns, nil
}

func (s *service) ListTransactionsByTenantID(ctx context.Context,
	tenantID string, limit, offset int) ([]*Transaction, error) {

	var txns []*Transaction
	if s.cache != nil {
		err := s.cache.GetTransactionsByTenant(ctx, tenantID, &txns)
		if err == nil && len(txns) > 0 {
			return txns, nil
		}
	}

	dbTxns, err := s.Repo.ListTransactionsByTenantID(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}

	if s.cache != nil && len(dbTxns) > 0 {
		_ = s.cache.SetTransactionsByTenant(ctx, tenantID, dbTxns, 5*time.Minute)
	}
	return dbTxns, nil
}

func (s *service) ListAllTransactions(ctx context.Context,
	limit, offset int) ([]*Transaction, error) {
	return s.Repo.ListAllTransactions(ctx, limit, offset)
}

func (s *service) ListTransactionsByStatus(ctx context.Context, status string,
	limit, offset int) ([]*Transaction, error) {
	return s.Repo.ListTransactionsByStatus(ctx, status, limit, offset)
}
