package provider

import "github.com/shopspring/decimal"

type (
	CreateProviderParams struct {
		Name   string
		Code   string
		Config map[string]interface{}
	}

	DepositRequest struct {
		UserID   string
		WalletID string
		Amount   decimal.Decimal
		Metadata map[string]interface{}
	}

	InitDepositResponse struct {
		AuthorizationURL string
		Reference        string
		Provider         string
	}

	VerifyResponse struct {
		Provider  string
		Status    string
		Amount    int64
		Currency  string
		Reference string
		Raw       interface{}
	}

	WithdrawalRequest struct {
		UserID   string
		WalletID string
		Amount   decimal.Decimal
		Metadata map[string]interface{}
	}
)
