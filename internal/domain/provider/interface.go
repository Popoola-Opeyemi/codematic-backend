package provider

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service interface {
	InitiateDeposit(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error)
	InitiateWithdrawal(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error)
	HandleWebhook(ctx context.Context, payload []byte) error
	GetTransactionStatus(ctx context.Context, reference string) (string, error)
}
