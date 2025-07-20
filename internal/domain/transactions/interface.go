package transactions

import (
	"context"
)

type Repository interface {
	GetTransactionByID(ctx context.Context, id string) (*Transaction, error)
	ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*Transaction, error)
	ListTransactionsByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*Transaction, error)
	ListAllTransactions(ctx context.Context, limit, offset int) ([]*Transaction, error)
	ListTransactionsByStatus(ctx context.Context, status string, limit, offset int) ([]*Transaction, error)
}

type Service interface {
	GetTransactionByID(ctx context.Context, id string) (*Transaction, error)
	ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*Transaction, error)
	ListTransactionsByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*Transaction, error)
	ListAllTransactions(ctx context.Context, limit, offset int) ([]*Transaction, error)
	ListTransactionsByStatus(ctx context.Context, status string, limit, offset int) ([]*Transaction, error)
}
