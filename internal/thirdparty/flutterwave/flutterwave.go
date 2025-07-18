package flutterwave

import (
	"context"

	"github.com/shopspring/decimal"
)

type FlutterwaveProvider struct{}

func (f *FlutterwaveProvider) InitiateDeposit(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error) {
	return "flutterwave-ref-789", nil
}

func (f *FlutterwaveProvider) InitiateWithdrawal(ctx context.Context, userID, walletID string, amount decimal.Decimal, metadata map[string]interface{}) (string, error) {
	return "flutterwave-ref-101", nil
}

func (f *FlutterwaveProvider) HandleWebhook(ctx context.Context, payload []byte) error {
	return nil
}

func (f *FlutterwaveProvider) GetTransactionStatus(ctx context.Context, reference string) (string, error) {
	return "success", nil
}
