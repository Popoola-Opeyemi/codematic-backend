package wallet

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        string
	UserID    string
	Balance   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID        string
	WalletID  string
	Type      string
	Status    string
	Amount    decimal.Decimal
	Fee       decimal.Decimal
	Provider  string
	Reference string
	Metadata  map[string]interface{}
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
