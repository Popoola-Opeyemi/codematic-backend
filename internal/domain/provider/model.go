package provider

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	CreateProviderParams struct {
		Name   string
		Code   string
		Config map[string]interface{}
	}

	DepositRequest struct {
		Amount   decimal.Decimal
		Channel  string
		Currency string
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

	PaystackConfig struct {
		BaseURL       string `json:"base_url"`
		SecretKey     string `json:"secret_key"`
		PublicKey     string `json:"public_key"`
		WebhookSecret string `json:"webhook_secret"`
	}

	FlutterwaveConfig struct {
		BaseURL       string `json:"base_url"`
		SecretKey     string `json:"secret_key"`
		PublicKey     string `json:"public_key"`
		WebhookSecret string `json:"webhook_secret"`
		EncryptionKey string `json:"encryption_key"`
	}

	ProviderDetails struct {
		ID        string
		Name      string
		Code      string
		IsActive  bool
		CreatedAt time.Time
		UpdatedAt time.Time

		Config interface{}
	}
)
