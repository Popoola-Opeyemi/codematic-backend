package wallet

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type walletRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &walletRepository{
		q: q,
		p: pool,
	}
}

func (r *walletRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

func (r *walletRepository) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	uid, err := utils.StringToPgUUID(walletID)
	if err != nil {
		return nil, err
	}
	w, err := r.q.GetWalletByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		ID:        w.ID.String(),
		UserID:    w.UserID.String(),
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}, nil
}

func (r *walletRepository) UpdateWalletBalance(ctx context.Context,
	walletID string, amount decimal.Decimal) error {
	uid, err := utils.StringToPgUUID(walletID)
	if err != nil {
		return err
	}
	return r.q.UpdateWalletBalance(ctx, db.UpdateWalletBalanceParams{
		Balance: amount,
		ID:      uid,
	})
}

func (r *walletRepository) CreateWallet(ctx context.Context, userID,
	walletTypeID string, balance decimal.Decimal) (*Wallet, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	userUUID, err := utils.StringToPgUUID(userID)
	if err != nil {
		return nil, err
	}
	walletTypeUUID, err := utils.StringToPgUUID(walletTypeID)
	if err != nil {
		return nil, err
	}
	w, err := r.q.CreateWallet(ctx, db.CreateWalletParams{
		ID:           pgtype.UUID{Bytes: [16]byte(uid), Valid: true},
		UserID:       userUUID,
		WalletTypeID: walletTypeUUID,
		Balance:      balance,
	})
	if err != nil {
		return nil, err
	}
	return &Wallet{
		ID:        w.ID.String(),
		UserID:    w.UserID.String(),
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}, nil
}

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *Transaction) error {
	uid, _ := utils.StringToPgUUID(tx.ID)
	tid, _ := utils.StringToPgUUID(tx.TenantID)
	wid, _ := utils.StringToPgUUID(tx.WalletID)
	pid, _ := utils.StringToPgUUID(tx.Provider)
	meta, _ := json.Marshal(tx.Metadata)

	_, err := r.q.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:         uid,
		TenantID:   tid,
		WalletID:   wid,
		ProviderID: pid,
		Reference:  tx.Reference,
		Type:       tx.Type,
		Status:     tx.Status,
		Amount:     tx.Amount,
		Fee:        tx.Fee,
		Metadata:   meta,
	})
	return err
}

func (r *walletRepository) ListTransactions(ctx context.Context,
	walletID string, limit, offset int) ([]Transaction, error) {
	wid, _ := utils.StringToPgUUID(walletID)

	rows, err := r.q.ListTransactionsByWalletID(ctx, db.ListTransactionsByWalletIDParams{
		WalletID: wid,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var txs []Transaction
	for _, row := range rows {
		var meta map[string]interface{}
		_ = json.Unmarshal(row.Metadata, &meta)

		fee := decimal.Zero

		txs = append(txs, Transaction{
			ID:        row.ID.String(),
			WalletID:  row.WalletID.String(),
			Type:      row.Type,
			Status:    row.Status,
			Amount:    row.Amount,
			Fee:       fee,
			Provider:  row.ProviderID.String(),
			Reference: row.Reference,
			Metadata:  meta,
			Error:     row.ErrorReason.String,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return txs, nil
}

func (r *walletRepository) ListActiveCurrencyCodes(ctx context.Context) ([]string, error) {
	rows, err := r.q.ListActiveCurrencyCodes(ctx)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *walletRepository) GetWalletTypeIDByCurrency(ctx context.Context,
	currency string) (string, error) {
	wtID, err := r.q.GetWalletTypeIDByCurrency(ctx, currency)
	if err != nil {
		return "", err
	}
	return wtID.String(), nil
}

func (r *walletRepository) CreateWalletsForNewUserFromAvailableWallets(ctx context.Context,
	userID string) ([]*Wallet, error) {
	userUUID, err := utils.StringToPgUUID(userID)
	if err != nil {
		return nil, err
	}

	types, err := r.q.ListActiveWalletTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wallet types: %w", err)
	}

	var wallets []*Wallet
	for _, wt := range types {
		walletID := uuid.New()

		w, err := r.q.CreateWallet(ctx, db.CreateWalletParams{
			ID:           pgtype.UUID{Bytes: walletID, Valid: true},
			UserID:       userUUID,
			WalletTypeID: wt.ID,
			Balance:      decimal.Zero,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create wallet for wallet type %s: %w", wt.Name, err)
		}

		wallets = append(wallets, &Wallet{
			ID:        w.ID.String(),
			UserID:    w.UserID.String(),
			Balance:   w.Balance,
			CreatedAt: w.CreatedAt.Time,
			UpdatedAt: w.UpdatedAt.Time,
		})
	}

	return wallets, nil
}

func (r *walletRepository) GetWalletByUserAndCurrency(ctx context.Context, userID string, currency string) (*Wallet, error) {
	userUUID, err := utils.StringToPgUUID(userID)
	if err != nil {
		return nil, err
	}

	w, err := r.q.GetWalletByUserAndCurrency(ctx, db.GetWalletByUserAndCurrencyParams{
		UserID:   userUUID,
		Currency: currency,
	})
	if err != nil {
		return nil, err
	}

	return &Wallet{
		ID:        w.ID.String(),
		UserID:    w.UserID.String(),
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}, nil
}
