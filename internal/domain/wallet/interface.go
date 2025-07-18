package wallet

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service interface {
	Deposit(ctx context.Context, userID, tenantID, walletID string, amount decimal.Decimal, providerCode string, metadata map[string]interface{}) error
	Withdraw(ctx context.Context, userID, tenantID, walletID string, amount decimal.Decimal, providerCode string, metadata map[string]interface{}) error
	Transfer(ctx context.Context, userID, tenantID, fromWalletID, toWalletID string, amount decimal.Decimal, metadata map[string]interface{}) error
	GetBalance(ctx context.Context, walletID string) (decimal.Decimal, error)
	GetTransactions(ctx context.Context, walletID string, limit, offset int) ([]Transaction, error)
}

type Repository interface {
	// Define repository methods for wallet persistence
}
