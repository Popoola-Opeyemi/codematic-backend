package wallet

import (
	"context"
	"errors"

	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type WalletService struct {
	DB   *db.DBConn
	Repo Repository
}

func NewService(db *db.DBConn) Service {
	return &WalletService{
		DB:   db,
		Repo: NewRepository(db.Queries, db.Pool),
	}
}

// Transactional wrapper
func (s *WalletService) withTx(ctx context.Context, fn func(repo Repository) error) error {
	tx, err := s.DB.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	txRepo := s.Repo.WithTx(dbsqlc.New(tx)) // use tx-bound version of the repo

	if err := fn(txRepo); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (s *WalletService) Deposit(ctx context.Context, data DepositForm) error {

	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	return s.withTx(ctx, func(repo Repository) error {
		wallet, err := repo.GetWallet(ctx, data.WalletID)
		if err != nil {
			return err
		}
		wallet.Balance = wallet.Balance.Add(data.Amount)

		if err := repo.UpdateWalletBalance(ctx, wallet.ID, wallet.Balance); err != nil {
			return err
		}

		tx := &Transaction{
			ID:        uuid.NewString(),
			WalletID:  wallet.ID,
			Type:      "deposit",
			Status:    "success",
			Amount:    data.Amount,
			Fee:       decimal.Zero,
			Provider:  data.Provider,
			Reference: uuid.NewString(),
			Metadata:  data.Metadata,
		}
		return repo.CreateTransaction(ctx, tx)
	})
}

func (s *WalletService) Withdraw(ctx context.Context, data WithdrawalForm) error {
	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	return s.withTx(ctx, func(repo Repository) error {
		wallet, err := repo.GetWallet(ctx, data.WalletID)
		if err != nil {
			return err
		}
		if wallet.Balance.LessThan(data.Amount) {
			return errors.New("insufficient balance")
		}

		wallet.Balance = wallet.Balance.Sub(data.Amount)

		if err := repo.UpdateWalletBalance(ctx, wallet.ID, wallet.Balance); err != nil {
			return err
		}

		tx := &Transaction{
			ID:        uuid.NewString(),
			WalletID:  wallet.ID,
			Type:      "withdrawal",
			Status:    "success",
			Amount:    data.Amount,
			Fee:       decimal.Zero,
			Provider:  data.Provider,
			Reference: uuid.NewString(),
			Metadata:  data.Metadata,
		}
		return repo.CreateTransaction(ctx, tx)
	})
}

func (s *WalletService) Transfer(ctx context.Context, data TransferForm) error {
	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	return s.withTx(ctx, func(repo Repository) error {
		from, err := repo.GetWallet(ctx, data.FromWalletID)
		if err != nil {
			return err
		}
		if from.Balance.LessThan(data.Amount) {
			return errors.New("insufficient balance")
		}

		to, err := repo.GetWallet(ctx, data.ToWalletID)
		if err != nil {
			return err
		}

		from.Balance = from.Balance.Sub(data.Amount)
		to.Balance = to.Balance.Add(data.Amount)

		if err := repo.UpdateWalletBalance(ctx, from.ID, from.Balance); err != nil {
			return err
		}
		if err := repo.UpdateWalletBalance(ctx, to.ID, to.Balance); err != nil {
			return err
		}

		tx := &Transaction{
			ID:        uuid.NewString(),
			WalletID:  from.ID,
			Type:      "transfer",
			Status:    "success",
			Amount:    data.Amount,
			Fee:       decimal.Zero,
			Provider:  "internal",
			Reference: uuid.NewString(),
			Metadata:  data.Metadata,
		}
		return repo.CreateTransaction(ctx, tx)
	})
}

func (s *WalletService) CreateWalletsForUserByCurrencies(ctx context.Context,
	userID string, currencies []string) ([]*Wallet, error) {

	return s.Repo.CreateWalletsForUserByCurrencies(ctx, userID, currencies)
}

func (s *WalletService) CreateWalletForNewUser(ctx context.Context,
	userID string) ([]*Wallet, error) {
	currencies, err := s.Repo.ListActiveCurrencyCodes(ctx)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateWalletsForUserByCurrencies(ctx, userID, currencies)
}

func (s *WalletService) CreateWallet(ctx context.Context, userID,
	walletTypeID string, balance decimal.Decimal) (*Wallet, error) {
	return s.Repo.CreateWallet(ctx, userID, walletTypeID, balance)
}

func (s *WalletService) GetBalance(ctx context.Context,
	walletID string) (decimal.Decimal, error) {
	w, err := s.Repo.GetWallet(ctx, walletID)
	if err != nil {
		return decimal.Zero, err
	}
	return w.Balance, nil
}

func (s *WalletService) GetTransactions(ctx context.Context,
	walletID string, limit, offset int) ([]Transaction, error) {
	return s.Repo.ListTransactions(ctx, walletID, limit, offset)
}

func (s *WalletService) GetWalletTypeIDByCurrency(ctx context.Context,
	currency string) (string, error) {
	return s.Repo.GetWalletTypeIDByCurrency(ctx, currency)
}
