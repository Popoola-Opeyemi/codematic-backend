package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Handler interface is defined in this package
type service struct {
	repo             Repository
	providerRegistry map[string]Handler // provider name -> Handler
}

func NewService(repo Repository, registry map[string]Handler) Service {
	return &service{
		repo:             repo,
		providerRegistry: registry,
	}
}

func (s *service) ProcessWebhook(ctx context.Context, provider string, payload []byte) error {
	handler, ok := s.providerRegistry[provider]
	if !ok {
		return errors.New("unknown provider")
	}
	// Parse provider-specific event (for demo, assume generic WebhookEvent)
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	// Idempotency check
	existing, err := s.repo.GetByProviderAndEventID(ctx, event.ProviderID, event.ProviderEventID)
	if err == nil && existing != nil {
		return nil // already processed
	}
	event.ID = uuid.New()
	event.Status = "received"
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, &event); err != nil {
		return err
	}
	err = handler.HandleWebhookEvent(ctx, &event)
	status := "processed"
	var lastError *string
	if err != nil {
		msg := err.Error()
		lastError = &msg
		status = "failed"
	}
	s.repo.UpdateStatus(ctx, event.ID, status, event.Attempts+1, lastError)
	return err
}

func (s *service) ReplayWebhook(ctx context.Context, id uuid.UUID) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil || event == nil {
		return errors.New("event not found")
	}
	handler, ok := s.providerRegistry[event.ProviderID.String()]
	if !ok {
		return errors.New("unknown provider")
	}
	err = handler.HandleWebhookEvent(ctx, event)
	status := "processed"
	var lastError *string
	if err != nil {
		msg := err.Error()
		lastError = &msg
		status = "failed"
	}
	s.repo.UpdateStatus(ctx, event.ID, status, event.Attempts+1, lastError)
	return err
}
