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

func (r *walletRepository) GetWallet(ctx context.Context,
	walletID string) (*Wallet, error) {
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
		ID:           uid,
		TenantID:     tid,
		WalletID:     wid,
		ProviderID:   pid,
		Reference:    tx.Reference,
		Type:         tx.Type,
		CurrencyCode: tx.CurrencyCode,
		Status:       tx.Status,
		Amount:       tx.Amount,
		Fee:          tx.Fee,
		Metadata:     meta,
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

func (r *walletRepository) GetWalletByUserAndCurrency(ctx context.Context,
	userID, currency string) (*Wallet, error) {
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

func (r *walletRepository) GetTransactionByReference(ctx context.Context,
	reference string) (*Transaction, error) {
	tx, err := r.q.GetTransactionByReference(ctx, reference)
	if err != nil {
		return nil, err
	}
	var meta map[string]interface{}
	_ = json.Unmarshal(tx.Metadata, &meta)
	return &Transaction{
		ID:           tx.ID.String(),
		WalletID:     tx.WalletID.String(),
		Type:         tx.Type,
		TenantID:     tx.TenantID.String(),
		Status:       tx.Status,
		CurrencyCode: tx.CurrencyCode,
		Amount:       tx.Amount,
		Fee:          tx.Fee,
		Provider:     tx.ProviderID.String(),
		Reference:    tx.Reference,
		Metadata:     meta,
		Error:        tx.ErrorReason.String,
		CreatedAt:    tx.CreatedAt.Time,
		UpdatedAt:    tx.UpdatedAt.Time,
	}, nil
}

func (r *walletRepository) UpdateTransactionStatusAndAmount(
	ctx context.Context, id,
	status string,
	amount decimal.Decimal,
) error {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return err
	}
	return r.q.UpdateTransactionStatusAndAmount(ctx, db.UpdateTransactionStatusAndAmountParams{
		Status: status,
		Amount: amount,
		ID:     uid,
	})
}

func (r *walletRepository) CreateDeposit(ctx context.Context, deposit *Deposit) error {
	userUUID, _ := utils.StringToPgUUID(deposit.UserID)
	transactionUUID, _ := utils.StringToPgUUID(deposit.TransactionID)
	extTxid := utils.ToPgxText(deposit.ExternalTxID)
	amount := decimal.NewFromFloat(deposit.Amount)
	params := db.CreateDepositParams{
		UserID:        userUUID,
		TransactionID: transactionUUID,
		ExternalTxid:  extTxid,
		Amount:        amount,
		Status:        deposit.Status,
		CreatedAt:     utils.ToPgxTstamp(deposit.CreatedAt),
		UpdatedAt:     utils.ToPgxTstamp(deposit.UpdatedAt),
	}
	row, err := r.q.CreateDeposit(ctx, params)
	if err != nil {
		return err
	}
	deposit.ID = int(row.ID)
	return nil
}

func (r *walletRepository) GetDepositByID(ctx context.Context, id int) (*Deposit, error) {
	row, err := r.q.GetDepositByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	extTxid := utils.FromPgText(row.ExternalTxid)
	var extTxidStr string
	if extTxid != nil {
		extTxidStr = *extTxid
	}
	return &Deposit{
		ID:            int(row.ID),
		UserID:        utils.FromPgUUID(row.UserID),
		TransactionID: utils.FromPgUUID(row.TransactionID),
		ExternalTxID:  extTxidStr,
		Amount:        row.Amount.InexactFloat64(),
		Status:        row.Status,
		CreatedAt:     utils.FromPgTimestamp(row.CreatedAt),
		UpdatedAt:     utils.FromPgTimestamp(row.UpdatedAt),
	}, nil
}

func (r *walletRepository) CreateWithdrawal(ctx context.Context, withdrawal *Withdrawal) error {
	userUUID, _ := utils.StringToPgUUID(withdrawal.UserID)
	transactionUUID, _ := utils.StringToPgUUID(withdrawal.TransactionID)
	extTxid := utils.ToPgxText(withdrawal.ExternalTxID)
	amount := decimal.NewFromFloat(withdrawal.Amount)
	params := db.CreateWithdrawalParams{
		UserID:        userUUID,
		TransactionID: transactionUUID,
		ExternalTxid:  extTxid,
		Amount:        amount,
		Status:        withdrawal.Status,
		CreatedAt:     utils.ToPgxTstamp(withdrawal.CreatedAt),
		UpdatedAt:     utils.ToPgxTstamp(withdrawal.UpdatedAt),
	}
	row, err := r.q.CreateWithdrawal(ctx, params)
	if err != nil {
		return err
	}
	withdrawal.ID = int(row.ID)
	return nil
}

func (r *walletRepository) GetWithdrawalByID(ctx context.Context, id int) (*Withdrawal, error) {
	row, err := r.q.GetWithdrawalByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	extTxid := utils.FromPgText(row.ExternalTxid)
	var extTxidStr string
	if extTxid != nil {
		extTxidStr = *extTxid
	}
	return &Withdrawal{
		ID:            int(row.ID),
		UserID:        utils.FromPgUUID(row.UserID),
		TransactionID: utils.FromPgUUID(row.TransactionID),
		ExternalTxID:  extTxidStr,
		Amount:        row.Amount.InexactFloat64(),
		Status:        row.Status,
		CreatedAt:     utils.FromPgTimestamp(row.CreatedAt),
		UpdatedAt:     utils.FromPgTimestamp(row.UpdatedAt),
	}, nil
}

func (r *walletRepository) UpdateDepositStatus(ctx context.Context, transactionID,
	status string) error {
	tid, err := utils.StringToPgUUID(transactionID)
	if err != nil {
		return err
	}

	return r.q.UpdateDepositStatusByTransactionID(ctx,
		db.UpdateDepositStatusByTransactionIDParams{
			Status:        status,
			TransactionID: tid,
		})
}
