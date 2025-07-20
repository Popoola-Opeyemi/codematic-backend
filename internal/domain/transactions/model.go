package transactions

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	WalletID     string                 `json:"wallet_id"`
	ProviderID   string                 `json:"provider_id"`
	CurrencyCode string                 `json:"currency_code"`
	Reference    string                 `json:"reference"`
	Type         string                 `json:"type"`
	Status       string                 `json:"status"`
	Amount       decimal.Decimal        `json:"amount"`
	Fee          decimal.Decimal        `json:"fee"`
	Metadata     map[string]interface{} `json:"metadata"`
	ErrorReason  string                 `json:"error_reason"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

const (
	RoleUser        = "user"
	RoleTenantAdmin = "tenant_admin"
	RoleAdmin       = "admin"
)
