package wallet

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	DepositRequest struct {
		UserID   string                 `json:"user_id"`
		TenantID string                 `json:"tenant_id"`
		WalletID string                 `json:"wallet_id"`
		Amount   string                 `json:"amount"`
		Provider string                 `json:"provider"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	DepositForm struct {
		UserID   string                 `json:"user_id"`
		TenantID string                 `json:"tenant_id"`
		WalletID string                 `json:"wallet_id"`
		Amount   decimal.Decimal        `json:"amount"`
		Provider string                 `json:"provider"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	WithdrawalForm struct {
		UserID   string                 `json:"user_id"`
		TenantID string                 `json:"tenant_id"`
		WalletID string                 `json:"wallet_id"`
		Amount   decimal.Decimal        `json:"amount"`
		Provider string                 `json:"provider"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	WithdrawalRequest struct {
		UserID   string                 `json:"user_id"`
		TenantID string                 `json:"tenant_id"`
		WalletID string                 `json:"wallet_id"`
		Amount   string                 `json:"amount"`
		Provider string                 `json:"provider"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	TransferRequest struct {
		UserID       string                 `json:"user_id"`
		TenantID     string                 `json:"tenant_id"`
		FromWalletID string                 `json:"from_wallet_id"`
		ToWalletID   string                 `json:"to_wallet_id"`
		Amount       string                 `json:"amount"`
		Metadata     map[string]interface{} `json:"metadata"`
	}

	TransferForm struct {
		UserID       string                 `json:"user_id"`
		TenantID     string                 `json:"tenant_id"`
		FromWalletID string                 `json:"from_wallet_id"`
		ToWalletID   string                 `json:"to_wallet_id"`
		Amount       decimal.Decimal        `json:"amount"`
		Metadata     map[string]interface{} `json:"metadata"`
	}

	Wallet struct {
		ID        string          `json:"id"`
		UserID    string          `json:"user_id"`
		Balance   decimal.Decimal `json:"balance"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
	Transaction struct {
		ID        string                 `json:"id"`
		WalletID  string                 `json:"wallet_id"`
		Type      string                 `json:"type"`
		Status    string                 `json:"status"`
		Amount    decimal.Decimal        `json:"amount"`
		Fee       decimal.Decimal        `json:"fee"`
		Provider  string                 `json:"provider"`
		Reference string                 `json:"reference"`
		Metadata  map[string]interface{} `json:"metadata"`
		Error     string                 `json:"error"`
		CreatedAt time.Time              `json:"created_at"`
		UpdatedAt time.Time              `json:"updated_at"`
	}
)
