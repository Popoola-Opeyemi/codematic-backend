package webhook

import (
	"bytes"
	"codematic/internal/config"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/tenants"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/shared/model"
	"codematic/internal/thirdparty/flutterwave"
	"codematic/internal/thirdparty/paystack"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

const (
	HeaderPaystackSignature    = "x-paystack-signature"
	HeaderFlutterwaveSignature = "flutterwave-signature"
)

type service struct {
	DB            *db.DBConn
	Repo          Repository
	Provider      provider.Service
	tenantService tenants.Service
	logger        *zap.Logger
	cfg           *config.Config
	Producer      *kafka.KafkaProducer
}

func NewService(
	Provider provider.Service,
	tenantService tenants.Service,
	logger *zap.Logger,
	db *db.DBConn,
	cfg *config.Config,
	producer *kafka.KafkaProducer,
) Service {
	return &service{
		DB:            db,
		Repo:          NewRepository(db.Queries, db.Pool),
		Provider:      Provider,
		tenantService: tenantService,
		logger:        logger,
		cfg:           cfg,
		Producer:      producer,
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

func (s *service) StartWalletDepositSuccessConsumer(ctx context.Context,
	broker string) {
	go func() {
		groupID := "webhook-wallet-deposit-success-group"
		err := kafka.Subscribe(ctx, broker, kafka.WalletDepositSuccessTopic, groupID, func(key, value []byte) {
			s.logger.Sugar().Infof("Received wallet deposit success event: %s", string(value))

			// Parse event
			var event event
			if err := json.Unmarshal(value, &event); err != nil {
				s.logger.Sugar().Errorf("Failed to parse deposit event: %v", err)
				return
			}

			// Get tenant webhook URL
			tenant, err := s.tenantService.GetTenantByID(ctx, event.TenantID)
			if err != nil {
				s.logger.Sugar().Errorf("Failed to get tenant for webhook: %v", err)
				return
			}
			if tenant.WebhookURL == "" {
				s.logger.Sugar().Warnf("No webhook URL for tenant %s", event.TenantID)
				return
			}

			// Save outgoing webhook event to DB
			webhookEvent := &WebhookEvent{
				ID:         uuid.NewString(),
				TenantID:   event.TenantID,
				EventType:  "wallet.deposit.success",
				Payload:    value,
				Status:     "pending",
				Attempts:   0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				IsOutgoing: true,
			}
			err = s.Repo.Create(ctx, webhookEvent)
			if err != nil {
				s.logger.Sugar().Errorf("Failed to save outgoing webhook event: %v", err)
				return
			}

			// Send HTTP POST to tenant webhook URL
			resp, err := http.Post(tenant.WebhookURL, "application/json", bytes.NewReader(value))
			if err != nil {
				s.logger.Sugar().Errorf("Failed to send webhook to tenant: %v", err)
				webhookEvent.Status = "failed"
				webhookEvent.Attempts++
				webhookEvent.UpdatedAt = time.Now()
				s.Repo.Create(ctx, webhookEvent) // Optionally update status
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				webhookEvent.Status = "success"
			} else {
				webhookEvent.Status = "failed"
			}
			webhookEvent.Attempts++
			webhookEvent.UpdatedAt = time.Now()
			s.Repo.Create(ctx, webhookEvent) // Optionally update status
		})
		if err != nil {
			s.logger.Sugar().Errorf("Failed to subscribe to wallet deposit success events: %v", err)
		}
	}()
}
