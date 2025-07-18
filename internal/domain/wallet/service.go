package wallet

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

type walletService struct {
	repo Repository
	// providerRegistry ProviderRegistry // TODO: inject provider abstraction
}

func NewService(repo Repository /*, other deps */) Service {
	return &walletService{
		repo: repo,
	}
}

func (s *walletService) Deposit(ctx context.Context, data DepositForm) error {

	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.GetWallet(ctx, data.WalletID)
	if err != nil {
		return err
	}

	// TODO: Call provider abstraction to initiate deposit
	wallet.Balance = wallet.Balance.Add(data.Amount)
	if err := s.repo.UpdateWalletBalance(ctx, data.WalletID, wallet.Balance); err != nil {
		return err
	}

	// TODO: Generate reference, provider, etc.
	tx := &Transaction{
		ID:        "", // generate UUID
		WalletID:  data.WalletID,
		Type:      "deposit",
		Status:    "success",
		Amount:    data.Amount,
		Fee:       decimal.Zero,
		Provider:  data.Provider,
		Reference: "", // generate reference
		Metadata:  data.Metadata,
	}
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *walletService) Withdraw(ctx context.Context, data WithdrawalForm) error {
	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}
	wallet, err := s.repo.GetWallet(ctx, data.WalletID)
	if err != nil {
		return err
	}
	if wallet.Balance.LessThan(data.Amount) {
		return errors.New("insufficient balance")
	}
	// TODO: Call provider abstraction to initiate withdrawal
	wallet.Balance = wallet.Balance.Sub(data.Amount)
	if err := s.repo.UpdateWalletBalance(ctx, data.WalletID, wallet.Balance); err != nil {
		return err
	}
	tx := &Transaction{
		ID:        "", // generate UUID
		WalletID:  data.WalletID,
		Type:      "withdrawal",
		Status:    "success",
		Amount:    data.Amount,
		Fee:       decimal.Zero,
		Provider:  data.Provider,
		Reference: "", // generate reference
		Metadata:  data.Metadata,
	}
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *walletService) Transfer(ctx context.Context, data TransferForm) error {
	if data.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}
	fromWallet, err := s.repo.GetWallet(ctx, data.FromWalletID)
	if err != nil {
		return err
	}
	if fromWallet.Balance.LessThan(data.Amount) {
		return errors.New("insufficient balance")
	}
	toWallet, err := s.repo.GetWallet(ctx, data.ToWalletID)
	if err != nil {
		return err
	}
	fromWallet.Balance = fromWallet.Balance.Sub(data.Amount)
	toWallet.Balance = toWallet.Balance.Add(data.Amount)
	if err := s.repo.UpdateWalletBalance(ctx, data.FromWalletID, fromWallet.Balance); err != nil {
		return err
	}
	if err := s.repo.UpdateWalletBalance(ctx, data.ToWalletID, toWallet.Balance); err != nil {
		return err
	}
	// TODO: Generate reference, provider, etc.
	tx := &Transaction{
		ID:        "", // generate UUID
		WalletID:  data.FromWalletID,
		Type:      "transfer",
		Status:    "success",
		Amount:    data.Amount,
		Fee:       decimal.Zero,
		Provider:  "internal",
		Reference: "", // generate reference
		Metadata:  data.Metadata,
	}
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *walletService) GetBalance(ctx context.Context, walletID string) (decimal.Decimal, error) {
	wallet, err := s.repo.GetWallet(ctx, walletID)
	if err != nil {
		return decimal.Zero, err
	}
	return wallet.Balance, nil
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string, limit, offset int) ([]Transaction, error) {
	return s.repo.ListTransactions(ctx, walletID, limit, offset)
}
