package webhook

import (
	"codematic/internal/config"
	"codematic/internal/infrastructure/db"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"errors"
	"strings"
)

type service struct {
	DB   *db.DBConn
	Repo Repository

	cfg *config.Config
}

func NewService(db *db.DBConn, cfg *config.Config) Service {
	return &service{
		DB:   db,
		Repo: NewRepository(db.Queries, db.Pool),

		cfg: cfg,
	}
}

func (s *service) ProcessWebhook(
	ctx context.Context,
	providerCode string,
	headers map[string]string,
	payload []byte,
) error {

	// err := s.VerifyWebhookSignature(providerCode, headers, payload)
	// if err != nil {
	// 	return err
	// }

	// // 4. Extract webhook event ID (implement per provider)
	// eventID, err := extractWebhookEventID(providerCode, payload)
	// if err != nil {
	// 	return err
	// }

	// // 5. Check if already processed (idempotency)
	// exists, err := s.DB.WebhookEventRepo.Exists(ctx, provider.ID, eventID)
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return nil // already processed, skip
	// }

	// // 6. Save raw event to DB
	// err = s.DB.WebhookEventRepo.Save(ctx, &model.WebhookEvent{
	// 	ProviderID: provider.ID,
	// 	EventID:    eventID,
	// 	Payload:    payload,
	// 	Status:     "pending",
	// 	Attempts:   0,
	// })
	// if err != nil {
	// 	return err
	// }

	// 7. Optionally enqueue to async processor
	// s.queue.Enqueue(eventID)

	return nil
}

func (s *service) ReplayWebhook(ctx context.Context, id string) error {
	return nil
}

// VerifyWebhookSignature verifies the webhook signature based on the provider
func (s *service) VerifyWebhookSignature(
	provider string,
	headers map[string]string,
	payload []byte,
) error {
	switch strings.ToLower(provider) {
	case "paystack":
		secret := s.cfg.PstkSecretHash
		if secret == "" {
			return errors.New("missing paystack secret key")
		}
		hash := utils.ComputeHMACSHA512(payload, secret)
		if hash != headers["x-paystack-signature"] {
			return model.ErrInvalidSignature
		}

	case "flutterwave":
		expected := s.cfg.FlwSecretHash
		if expected == "" {
			return errors.New("missing FLW_SECRET_HASH in environment")
		}
		if headers["flutterwave-signature"] != expected {
			return model.ErrInvalidSignature
		}

	default:
		return errors.New("unsupported provider for signature verification")
	}

	return nil
}
