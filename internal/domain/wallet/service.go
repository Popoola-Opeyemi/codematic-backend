package wallet

import (
	"context"

	"github.com/shopspring/decimal"
)

type walletService struct {
	repo Repository
	// Add provider registry, logger, etc. as needed
}

func NewService(repo Repository /*, other deps */) Service {
	return &walletService{
		repo: repo,
	}
}

func (s *walletService) Deposit(ctx context.Context, userID, tenantID, walletID string, amount decimal.Decimal, providerCode string, metadata map[string]interface{}) error {
	// TODO: Implement deposit logic
	return nil
}

func (s *walletService) Withdraw(ctx context.Context, userID, tenantID, walletID string, amount decimal.Decimal, providerCode string, metadata map[string]interface{}) error {
	// TODO: Implement withdraw logic
	return nil
}

func (s *walletService) Transfer(ctx context.Context, userID, tenantID, fromWalletID, toWalletID string, amount decimal.Decimal, metadata map[string]interface{}) error {
	// TODO: Implement transfer logic
	return nil
}

func (s *walletService) GetBalance(ctx context.Context, walletID string) (decimal.Decimal, error) {
	// TODO: Implement get balance logic
	return decimal.Zero, nil
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string, limit, offset int) ([]Transaction, error) {
	// TODO: Implement get transactions logic
	return nil, nil
}
