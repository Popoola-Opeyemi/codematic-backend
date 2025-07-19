package webhook

import (
	"time"
)

type WebhookEvent struct {
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
}
