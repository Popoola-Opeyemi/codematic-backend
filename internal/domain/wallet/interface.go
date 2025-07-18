package wallet

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service interface {
	Deposit(ctx context.Context, data DepositForm) error
	Withdraw(ctx context.Context, data WithdrawalForm) error
	Transfer(ctx context.Context, data TransferForm) error
	GetBalance(ctx context.Context, walletID string) (decimal.Decimal, error)
	GetTransactions(ctx context.Context, walletID string,
		limit, offset int) ([]Transaction, error)
}

type Repository interface {
	GetWallet(ctx context.Context, walletID string) (*Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID string,
		amount decimal.Decimal) error
	CreateTransaction(ctx context.Context, tx *Transaction) error
	ListTransactions(ctx context.Context, walletID string, limit,
		offset int) ([]Transaction, error)
}
