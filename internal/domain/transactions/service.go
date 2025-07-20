package transactions

import (
	"context"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTransactionByID(ctx context.Context, id string) (*Transaction, error) {
	return s.repo.GetTransactionByID(ctx, id)
}

func (s *service) ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*Transaction, error) {
	return s.repo.ListTransactionsByUserID(ctx, userID, limit, offset)
}

func (s *service) ListTransactionsByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*Transaction, error) {
	return s.repo.ListTransactionsByTenantID(ctx, tenantID, limit, offset)
}

func (s *service) ListAllTransactions(ctx context.Context, limit, offset int) ([]*Transaction, error) {
	return s.repo.ListAllTransactions(ctx, limit, offset)
}

func (s *service) ListTransactionsByStatus(ctx context.Context, status string, limit, offset int) ([]*Transaction, error) {
	return s.repo.ListTransactionsByStatus(ctx, status, limit, offset)
}
