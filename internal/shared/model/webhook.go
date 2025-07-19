package model

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID              uuid.UUID `json:"id"`
	ProviderID      uuid.UUID `json:"provider_id"`
	ProviderEventID string    `json:"provider_event_id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	EventType       string    `json:"event_type"`
	Payload         []byte    `json:"payload"`
	Status          string    `json:"status"`
	Attempts        int       `json:"attempts"`
	LastError       *string   `json:"last_error,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
