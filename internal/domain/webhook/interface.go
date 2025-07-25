package webhook

import (
	"context"
)

type Service interface {
	VerifyWebhookSignature(
		ctx context.Context,
		provider string,
		headers map[string]string,
		payload []byte,
	) error
	HandleWebhook(
		ctx context.Context,
		provider string,
		headers map[string]string,
		payload []byte,
	) error
}

type Repository interface {
	Create(ctx context.Context, event *WebhookEvent) error
	GetByProviderAndEventID(ctx context.Context, providerID string, providerEventID string) (*WebhookEvent, error)
	UpdateStatus(ctx context.Context, id string, status string, attempts int, lastError *string) error
	GetByID(ctx context.Context, id string) (*WebhookEvent, error)
}

type Handler interface {
	HandleWebhookEvent(ctx context.Context, event *WebhookEvent) error
}
