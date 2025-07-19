package wallet

import (
	db "codematic/internal/infrastructure/db/sqlc"
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
	CreateWallet(ctx context.Context, userID string, walletTypeID string,
		balance decimal.Decimal) (*Wallet, error)
	GetWalletTypeIDByCurrency(ctx context.Context, currency string) (string, error)
	CreateWalletForNewUser(ctx context.Context,
		userID string) ([]*Wallet, error)
	WithTx(q *db.Queries) Service
}

type Repository interface {
	CreateWalletsForUserFromAvailableWallets(ctx context.Context,
		userID string) ([]*Wallet, error)
	GetWallet(ctx context.Context, walletID string) (*Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID string,
		amount decimal.Decimal) error
	CreateTransaction(ctx context.Context, tx *Transaction) error
	ListTransactions(ctx context.Context, walletID string, limit,
		offset int) ([]Transaction, error)
	CreateWallet(ctx context.Context, userID string, walletTypeID string,
		balance decimal.Decimal) (*Wallet, error)
	GetWalletTypeIDByCurrency(ctx context.Context, currency string) (string, error)
	ListActiveCurrencyCodes(ctx context.Context) ([]string, error)
	WithTx(q *db.Queries) Repository
}
