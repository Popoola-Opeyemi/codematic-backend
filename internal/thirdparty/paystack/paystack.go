package paystack

import (
	"context"

	"github.com/shopspring/decimal"
)

type PaystackProvider struct{}

func (p *PaystackProvider) InitiateDeposit(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error) {
	return "paystack-ref-123", nil
}

func (p *PaystackProvider) InitiateWithdrawal(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error) {
	return "paystack-ref-456", nil
}

func (p *PaystackProvider) HandleWebhook(ctx context.Context, payload []byte) error {
	return nil
}

func (p *PaystackProvider) GetTransactionStatus(ctx context.Context, reference string) (string, error) {
	return "success", nil
}
