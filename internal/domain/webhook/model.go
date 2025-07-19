package webhook

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID              uuid.UUID
	ProviderID      uuid.UUID
	ProviderEventID string
	TenantID        uuid.UUID
	EventType       string
	Payload         []byte
	Status          string
	Attempts        int
	LastError       *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
