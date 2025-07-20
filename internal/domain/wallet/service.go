package wallet

import (
	"context"
	"errors"
	"fmt"

	"codematic/internal/domain/provider"
	"codematic/internal/domain/provider/gateways"
	"codematic/internal/domain/user"
	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"
	kafka "codematic/internal/infrastructure/events/kafka"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type WalletService struct {
	DB   *db.DBConn
	Repo Repository

	Provider provider.Service
	User     user.Service

	logger   *zap.Logger
	Producer *kafka.KafkaProducer
}

func NewService(Provider provider.Service,
	User user.Service,
	db *db.DBConn, logger *zap.Logger,
	producer *kafka.KafkaProducer) Service {
	return &WalletService{
		DB:       db,
		Repo:     NewRepository(db.Queries, db.Pool),
		Provider: Provider,
		User:     User,
		logger:   logger,
		Producer: producer,
	}
}

func (s *WalletService) WithTx(q *dbsqlc.Queries) Service {
	return &WalletService{
		DB:     s.DB,
		Repo:   NewRepository(q, s.DB.Pool),
		User:   s.User,
		logger: s.logger,
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

func (s *WalletService) InitiateDeposit(ctx context.Context, data DepositForm) (gateways.GatewayResponse, error) {
	s.logger.Sugar().Infof("Deposit started: tenant_id=%s, amount=%s, channel=%s",
		data.TenantID, data.Amount.String(), data.Channel)

	var response gateways.GatewayResponse

	if data.Amount.LessThanOrEqual(decimal.Zero) {
		s.logger.Sugar().Errorf("Invalid deposit amount: %s", data.Amount.String())
		return response, errors.New("amount must be positive")
	}

	err := s.withTx(ctx, func(repo Repository) error {
		// Check wallet existence
		wallet, err := repo.GetWalletByUserAndCurrency(ctx, data.UserID, data.Currency)
		if err != nil {
			s.logger.Sugar().Errorf("Failed to get wallet for user %s and currency %s: %v", data.UserID, data.Currency, err)
			return fmt.Errorf("failed to get %s wallet for user", data.Currency)
		}

		// Call the provider service to initiate the payment first
		providerReq := provider.DepositRequest{
			Amount:   data.Amount,
			Channel:  data.Channel,
			Currency: data.Currency,
			Metadata: data.Metadata,
		}

		gateway, err := s.Provider.InitiateDeposit(ctx, providerReq)
		if err != nil {
			s.logger.Sugar().Errorf("Failed to initiate deposit with provider: %v", err)
			return err
		}

		s.logger.Sugar().Infow("Gateway response", "response", fmt.Sprintf("%+v", gateway))

		transaction := &Transaction{
			ID:           uuid.NewString(),
			WalletID:     wallet.ID,
			Type:         TransactionDeposit,
			TenantID:     data.TenantID,
			Status:       StatusPending,
			Amount:       data.Amount,
			Fee:          decimal.Zero,
			Provider:     gateway.ProviderID,
			CurrencyCode: data.Currency,
			Reference:    gateway.Reference,
			Metadata:     data.Metadata,
		}

		if err := repo.CreateTransaction(ctx, transaction); err != nil {
			s.logger.Sugar().Errorf("Failed to create transaction: %v", err)
			return err
		}

		response = gateway

		return nil
	})

	return response, err
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

func (s *WalletService) CreateWalletForNewUser(ctx context.Context,
	userID string) ([]*Wallet, error) {

	d, err := s.Repo.CreateWalletsForNewUserFromAvailableWallets(ctx, userID)
	if err != nil {
		s.logger.Sugar().Info("CreateWalletsForUserByCurrencies error occured", err)
		return nil, err
	}

	return d, nil
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
