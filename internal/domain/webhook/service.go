package webhook

import (
	"codematic/internal/config"
	"codematic/internal/domain/provider"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/shared/model"
	"codematic/internal/thirdparty/flutterwave"
	"codematic/internal/thirdparty/paystack"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const (
	HeaderPaystackSignature    = "x-paystack-signature"
	HeaderFlutterwaveSignature = "flutterwave-signature"
)

type service struct {
	DB       *db.DBConn
	Repo     Repository
	Provider provider.Service
	logger   *zap.Logger
	cfg      *config.Config
	Producer *kafka.KafkaProducer
}

func NewService(
	Provider provider.Service,
	logger *zap.Logger,
	db *db.DBConn,
	cfg *config.Config,
	producer *kafka.KafkaProducer,
) Service {
	return &service{
		DB:       db,
		Repo:     NewRepository(db.Queries, db.Pool),
		Provider: Provider,
		logger:   logger,
		cfg:      cfg,
		Producer: producer,
	}
}

func (s *service) HandleWebhook(
	ctx context.Context,
	provider string,
	headers map[string]string,
	payload []byte,
) error {
	s.logger.Sugar().Infof("Handling webhook for provider: %s", provider)

	if err := s.VerifyWebhookSignature(ctx, provider, headers, payload); err != nil {
		s.logger.Sugar().Errorf("Webhook signature verification failed: %v", err)
		return err
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		s.logger.Sugar().Errorf("Invalid webhook payload format: %v", err)
		return fmt.Errorf("invalid webhook payload: %w", err)
	}

	switch strings.ToLower(provider) {
	case paystack.ProviderPaystack:
		return s.handlePaystackEvent(ctx, event, payload)
	case flutterwave.ProviderFlutterwave:
		return s.handleFlutterwaveEvent(ctx, event, payload)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *service) handlePaystackEvent(
	ctx context.Context,
	event map[string]interface{},
	rawPayload []byte,
) error {
	s.logger.Sugar().Infof("Paystack event received: %+v", event)

	// Extract event type as key if available
	key := ""
	if evt, ok := event["event"]; ok {
		if str, ok := evt.(string); ok {
			key = str
		}
	}

	// Emit event to Kafka for wallet service to process
	return s.Producer.Publish(ctx, kafka.PaystackWalletEventTopic, key, rawPayload)
}

func (s *service) handleFlutterwaveEvent(
	ctx context.Context,
	event map[string]interface{},
	rawPayload []byte,
) error {
	s.logger.Sugar().Infof("Flutterwave event received: %+v", event)

	return nil
}

func (s *service) VerifyWebhookSignature(
	ctx context.Context,
	provider string,
	headers map[string]string,
	payload []byte,
) error {
	provider = strings.ToLower(provider)

	switch provider {
	case paystack.ProviderPaystack:
		signature := getHeader(headers, HeaderPaystackSignature)
		if signature == "" {
			return fmt.Errorf("missing %s header", HeaderPaystackSignature)
		}

		isValid, err := s.Provider.VerifyWebhookSignature(ctx,
			paystack.ProviderPaystack, signature, payload)
		if err != nil {
			return fmt.Errorf("paystack signature verification failed: %w", err)
		}
		if !isValid {
			return model.ErrInvalidSignature
		}
		return nil

	case flutterwave.ProviderFlutterwave:
		secret := s.cfg.FlwSecretHash
		if secret == "" {
			return errors.New("missing Flutterwave secret key")
		}

		if getHeader(headers, HeaderFlutterwaveSignature) != secret {
			return model.ErrInvalidSignature
		}
		return nil

	default:
		return fmt.Errorf("unsupported provider for signature verification: %s", provider)
	}
}

func getHeader(headers map[string]string, key string) string {
	for k, v := range headers {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}
