package wallet

import (
	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type walletRepository struct {
	q *dbsqlc.Queries
	p *pgxpool.Pool
}

func NewRepository(db *db.DBConn) Repository {
	return &walletRepository{
		q: db.Queries,
		p: db.Pool,
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
	return r.q.UpdateWalletBalance(ctx, dbsqlc.UpdateWalletBalanceParams{
		Balance: &amount,
		ID:      pgtype.UUID{Bytes: [16]byte(uid), Valid: true},
	})
}

func (r *walletRepository) CreateWallet(ctx context.Context, userID string,
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
	w, err := r.q.CreateWallet(ctx, dbsqlc.CreateWalletParams{
		ID:           pgtype.UUID{Bytes: [16]byte(uid), Valid: true},
		UserID:       userUUID,
		WalletTypeID: walletTypeUUID,
		Balance:      &balance,
	})
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

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *Transaction) error {
	uid, _ := utils.StringToPgUUID(tx.ID)
	wid, _ := utils.StringToPgUUID(tx.WalletID)
	pid, _ := utils.StringToPgUUID(tx.Provider)
	meta, _ := json.Marshal(tx.Metadata)

	fee := &tx.Fee
	_, err := r.q.CreateTransaction(ctx, dbsqlc.CreateTransactionParams{
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
		dbsqlc.ListTransactionsByWalletIDParams{
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

func (r *walletRepository) GetWalletTypeIDByCurrency(ctx context.Context, currency string) (string, error) {
	var id uuid.UUID
	row := r.p.QueryRow(ctx, "SELECT id FROM wallet_types WHERE currency = $1 LIMIT 1", currency)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id.String(), nil
}

func (r *walletRepository) CreateWalletsForUserByCurrencies(ctx context.Context, userID string, currencies []string) ([]*Wallet, error) {
	var wallets []*Wallet
	for _, currency := range currencies {
		walletTypeID, err := r.GetWalletTypeIDByCurrency(ctx, currency)
		if err != nil {
			return nil, err
		}
		w, err := r.CreateWallet(ctx, userID, walletTypeID, decimal.NewFromInt(0))
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *walletRepository) ListActiveCurrencyCodes(ctx context.Context) ([]string, error) {
	rows, err := r.p.Query(ctx, "SELECT code FROM currencies WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}
