package webhook

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	ProcessWebhook(ctx context.Context, provider string, payload []byte) error
	ReplayWebhook(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Create(ctx context.Context, event *WebhookEvent) error
	GetByProviderAndEventID(ctx context.Context, providerID uuid.UUID, providerEventID string) (*WebhookEvent, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, attempts int, lastError *string) error
	GetByID(ctx context.Context, id uuid.UUID) (*WebhookEvent, error)
}

type Handler interface {
	HandleWebhookEvent(ctx context.Context, event *WebhookEvent) error
}
