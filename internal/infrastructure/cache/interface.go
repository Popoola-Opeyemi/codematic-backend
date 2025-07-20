package cache

import (
	"context"
	"time"
)

type WalletCacheStore interface {
	SetWalletBalance(ctx context.Context, walletID string, balance float64, ttl time.Duration) error
	GetWalletBalance(ctx context.Context, walletID string) (float64, error)
	DeleteWalletBalance(ctx context.Context, walletID string) error

	SetWalletTransactions(ctx context.Context, walletID string, txns interface{}, ttl time.Duration) error
	GetWalletTransactions(ctx context.Context, walletID string, result interface{}) error
	DeleteWalletTransactions(ctx context.Context, walletID string) error
}

// TransactionCacheStore provides caching for transactions
// (single transaction and lists by user/tenant)
type TransactionCacheStore interface {
	SetTransaction(ctx context.Context, txID string, tx interface{}, ttl time.Duration) error
	GetTransaction(ctx context.Context, txID string, result interface{}) error
	DeleteTransaction(ctx context.Context, txID string) error

	SetTransactionsByUser(ctx context.Context, userID string, txns interface{}, ttl time.Duration) error
	GetTransactionsByUser(ctx context.Context, userID string, result interface{}) error
	DeleteTransactionsByUser(ctx context.Context, userID string) error

	SetTransactionsByTenant(ctx context.Context, tenantID string, txns interface{}, ttl time.Duration) error
	GetTransactionsByTenant(ctx context.Context, tenantID string, result interface{}) error
	DeleteTransactionsByTenant(ctx context.Context, tenantID string) error
}
