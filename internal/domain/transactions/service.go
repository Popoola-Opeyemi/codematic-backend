package transactions

import (
	"codematic/internal/infrastructure/cache"
	"context"
	"time"
)

type service struct {
	repo  Repository
	cache cache.TransactionCacheStore
}

func NewService(
	repo Repository,
	cacheStore cache.TransactionCacheStore,
) Service {
	return &service{repo: repo, cache: cacheStore}
}

func (s *service) GetTransactionByID(ctx context.Context, id string) (*Transaction, error) {
	var tx Transaction
	if s.cache != nil {
		err := s.cache.GetTransaction(ctx, id, &tx)
		if err == nil && tx.ID != "" {
			return &tx, nil
		}
	}
	dbTx, err := s.repo.GetTransactionByID(ctx, id)
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

	dbTxns, err := s.repo.ListTransactionsByUserID(ctx, userID, limit, offset)
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

	dbTxns, err := s.repo.ListTransactionsByTenantID(ctx, tenantID, limit, offset)
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
	return s.repo.ListAllTransactions(ctx, limit, offset)
}

func (s *service) ListTransactionsByStatus(ctx context.Context, status string,
	limit, offset int) ([]*Transaction, error) {
	return s.repo.ListTransactionsByStatus(ctx, status, limit, offset)
}
