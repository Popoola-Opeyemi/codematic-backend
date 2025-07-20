package webhook

import (
	"time"
)

type (
	WebhookEvent struct {
		ID              string
		ProviderID      string
		ProviderEventID string
		TenantID        string
		EventType       string
		Payload         []byte
		Status          string
		Attempts        int
		LastError       string
		CreatedAt       time.Time
		UpdatedAt       time.Time
		IsOutgoing      bool
	}
	event struct {
		TenantID  string                 `json:"tenant_id"`
		WalletID  string                 `json:"wallet_id"`
		Amount    string                 `json:"amount"`
		Provider  string                 `json:"provider"`
		Metadata  map[string]interface{} `json:"metadata"`
		Timestamp time.Time              `json:"timestamp"`
	}
)
