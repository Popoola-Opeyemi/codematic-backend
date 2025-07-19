package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"codematic/internal/domain/provider"
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
	Producer *kafka.KafkaProducer // add this field
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

func (s *WalletService) Deposit(ctx context.Context, data DepositForm) error {
	s.logger.Sugar().Infof("Deposit started: wallet_id=%s, tenant_id=%s, amount=%s, provider=%s",
		data.WalletID, data.TenantID, data.Amount.String(), data.Provider)

	if data.Amount.LessThanOrEqual(decimal.Zero) {
		s.logger.Sugar().Errorf("Invalid deposit amount: %s", data.Amount.String())
		return errors.New("amount must be positive")
	}

	return s.withTx(ctx, func(repo Repository) error {
		s.logger.Sugar().Info("Fetching wallet:", data.WalletID)
		wallet, err := repo.GetWallet(ctx, data.WalletID)
		if err != nil {
			s.logger.Sugar().Errorw("Failed to get wallet", "wallet_id", data.WalletID, "error", err)
			return err
		}

		wallet.Balance = wallet.Balance.Add(data.Amount)
		s.logger.Sugar().Infof("Updated wallet balance: new_balance=%s", wallet.Balance.String())

		if err := repo.UpdateWalletBalance(ctx, wallet.ID, wallet.Balance); err != nil {
			s.logger.Sugar().Errorw("Failed to update wallet balance", "wallet_id", wallet.ID, "error", err)
			return err
		}

		tx := &Transaction{
			ID:        uuid.NewString(),
			WalletID:  wallet.ID,
			Type:      "deposit",
			Status:    "pending",
			Amount:    data.Amount,
			Fee:       decimal.Zero,
			Provider:  data.Provider,
			Reference: uuid.NewString(),
			Metadata:  data.Metadata,
		}

		if err := repo.CreateTransaction(ctx, tx); err != nil {
			s.logger.Sugar().Errorw("Failed to create transaction", "wallet_id", wallet.ID, "tx_id", tx.ID, "error", err)
			return err
		}

		// Emit Kafka event
		depositEvent := struct {
			TenantID  string                 `json:"tenant_id"`
			WalletID  string                 `json:"wallet_id"`
			Amount    string                 `json:"amount"`
			Provider  string                 `json:"provider"`
			Metadata  map[string]interface{} `json:"metadata"`
			Timestamp time.Time              `json:"timestamp"`
		}{
			TenantID:  data.TenantID,
			WalletID:  data.WalletID,
			Amount:    data.Amount.String(),
			Provider:  data.Provider,
			Metadata:  data.Metadata,
			Timestamp: time.Now().UTC(),
		}

		payload, err := json.Marshal(depositEvent)
		if err != nil {
			s.logger.Sugar().Errorw("Failed to marshal deposit event", "error", err)
		} else {
			if err := s.Producer.Publish(ctx, kafka.WalletDepositSuccessTopic, data.WalletID, payload); err != nil {
				s.logger.Sugar().Errorw("Failed to publish deposit event", "wallet_id", data.WalletID, "error", err)
			} else {
				s.logger.Sugar().Info("Deposit event published successfully")
			}
		}

		return nil
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

// InitiateDeposit handles the complete deposit flow with all required checks
func (s *WalletService) InitiateDeposit(ctx context.Context, data DepositRequest) (string, error) {
	// Convert amount string to decimal
	amount, err := decimal.NewFromString(data.Amount)
	if err != nil {
		return "", errors.New("invalid amount format")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return "", errors.New("amount must be positive")
	}

	// 1. Check user existence
	_, err = s.User.GetUserByID(ctx, data.UserID)
	if err != nil {
		s.logger.Error("user not found", zap.String("user_id", data.UserID), zap.Error(err))
		return "", errors.New("user not found")
	}

	// 2. Check for wallet existence based on currency
	wallet, err := s.Repo.GetWalletByUserAndCurrency(ctx, data.UserID, data.Currency)
	if err != nil {
		s.logger.Error("wallet not found for user and currency",
			zap.String("user_id", data.UserID),
			zap.String("currency", data.Currency),
			zap.Error(err))
		return "", errors.New("wallet not found for the specified currency")
	}

	// 3. Create a pending transaction record
	transactionID := uuid.NewString()
	reference := uuid.NewString()

	tx := &Transaction{
		ID:        transactionID,
		WalletID:  wallet.ID,
		Type:      "deposit",
		Status:    "pending",
		Amount:    amount,
		Fee:       decimal.Zero,
		Provider:  data.Channel, // Using channel as provider for now
		Reference: reference,
		Metadata:  data.Metadata,
	}

	err = s.Repo.CreateTransaction(ctx, tx)
	if err != nil {
		s.logger.Error("failed to create pending transaction", zap.Error(err))
		return "", errors.New("failed to create transaction record")
	}

	// 4. Call the provider service to initiate the deposit
	providerReq := provider.DepositRequest{
		UserID:   data.UserID,
		WalletID: wallet.ID,
		Amount:   amount,
		Metadata: map[string]interface{}{
			"provider":  data.Channel,
			"currency":  data.Currency,
			"reference": reference,
		},
	}

	providerRef, err := s.Provider.InitiateDeposit(ctx, providerReq)
	if err != nil {
		s.logger.Error("failed to initiate deposit with provider", zap.Error(err))
		// Update transaction status to failed
		tx.Status = "failed"
		tx.Error = err.Error()
		_ = s.Repo.CreateTransaction(ctx, tx) // This will create a new record or we need an update method
		return "", err
	}

	s.logger.Info("deposit initiated successfully",
		zap.String("user_id", data.UserID),
		zap.String("wallet_id", wallet.ID),
		zap.String("currency", data.Currency),
		zap.String("amount", amount.String()),
		zap.String("reference", reference),
		zap.String("provider_ref", providerRef))

	return reference, nil
}
