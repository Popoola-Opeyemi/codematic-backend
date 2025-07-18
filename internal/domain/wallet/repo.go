package wallet

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type walletRepository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &walletRepository{
		q: q,
	}
}

func (r *walletRepository) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	uid, err := uuid.Parse(walletID)
	if err != nil {
		return nil, err
	}
	w, err := r.q.GetWalletByID(ctx, pgtype.UUID{Bytes: [16]byte(uid), Valid: true})
	if err != nil {
		return nil, err
	}
	return &Wallet{
		ID:        w.ID.String(),
		UserID:    w.UserID.String(),
		Balance:   *w.Balance,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}, nil
}

func (r *walletRepository) UpdateWalletBalance(ctx context.Context,
	walletID string, amount decimal.Decimal) error {
	uid, err := uuid.Parse(walletID)
	if err != nil {
		return err
	}
	return r.q.UpdateWalletBalance(ctx, db.UpdateWalletBalanceParams{
		Balance: &amount,
		ID:      pgtype.UUID{Bytes: [16]byte(uid), Valid: true},
	})
}

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *Transaction) error {
	uid, _ := utils.StringToPgUUID(tx.ID)
	wid, _ := utils.StringToPgUUID(tx.WalletID)
	pid, _ := utils.StringToPgUUID(tx.Provider)
	meta, _ := json.Marshal(tx.Metadata)

	fee := &tx.Fee
	_, err := r.q.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:          uid,
		TenantID:    pgtype.UUID{}, // TODO: set tenant id
		WalletID:    wid,
		ProviderID:  pid,
		Reference:   tx.Reference,
		Type:        tx.Type,
		Status:      tx.Status,
		Amount:      tx.Amount,
		Fee:         fee,
		Metadata:    meta,
		ErrorReason: pgtype.Text{String: tx.Error, Valid: tx.Error != ""},
	})
	return err
}

func (r *walletRepository) ListTransactions(ctx context.Context,
	walletID string, limit, offset int) ([]Transaction, error) {

	wid, _ := utils.StringToPgUUID(walletID)

	rows, err := r.q.ListTransactionsByWalletID(ctx,
		db.ListTransactionsByWalletIDParams{
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
		if row.Fee != nil {
			fee = *row.Fee
		}

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
