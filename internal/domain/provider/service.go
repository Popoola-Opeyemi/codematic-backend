package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"codematic/internal/domain/provider/gateways"
	"codematic/internal/infrastructure/cache"
	dbconn "codematic/internal/infrastructure/db"
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/thirdparty/flutterwave"
	"codematic/internal/thirdparty/paystack"

	"go.uber.org/zap"
)

type providerService struct {
	DB           *dbconn.DBConn
	Repo         Repository
	cacheManager cache.CacheManager
	Logger       *zap.Logger
	Producer     *kafka.KafkaProducer
}

func NewService(
	db *dbconn.DBConn,
	cacheManager cache.CacheManager,
	logger *zap.Logger,
	producer *kafka.KafkaProducer,
) Service {
	return &providerService{
		DB:           db,
		Repo:         NewRepository(db.Queries, db.Pool),
		cacheManager: cacheManager,
		Logger:       logger,
		Producer:     producer,
	}
}

func (s *providerService) InitiateDeposit(ctx context.Context,
	req DepositRequest) (gateways.GatewayResponse, error) {
	email, _ := req.Metadata["email"].(string)

	providerRow, err := s.Repo.SelectBestProviderByCurrencyAndChannel(ctx, req.Currency, req.Channel)
	if err != nil {
		s.Logger.Error("No provider available", zap.Error(err))
		return gateways.GatewayResponse{}, fmt.Errorf("no provider available for currency %s and channel %s", req.Currency, req.Channel)
	}

	provider, err := s.GetProviderByID(ctx, providerRow.ID.String())
	if err != nil {
		s.Logger.Error("Failed to retrieve provider details", zap.Error(err))
		return gateways.GatewayResponse{}, err
	}

	code := strings.ToLower(provider.Code)

	if code == paystack.ProviderPaystack {
		var cfg PaystackConfig
		if err := json.Unmarshal(provider.Config, &cfg); err != nil {
			s.Logger.Error("Failed to decode paystack config", zap.Error(err))
			return gateways.GatewayResponse{}, err
		}

		gateway := gateways.NewPaystackProvider(s.Logger, cfg.BaseURL, cfg.SecretKey)
		return gateway.InitDeposit(ctx, gateways.DepositRequest{
			Email:      email,
			Amount:     req.Amount,
			ProviderID: provider.ID.String(),
		})
	}

	if code == flutterwave.ProviderFlutterwave {
		return gateways.GatewayResponse{}, fmt.Errorf("flutterwave not yet implemented")
	}

	return gateways.GatewayResponse{}, fmt.Errorf("unsupported provider: %s", provider.Code)
}

func (s *providerService) GetProviderByCode(ctx context.Context,
	code string) (*db.Provider, error) {
	code = strings.ToLower(code)

	if provider, err := s.cacheManager.GetProviderCacheByCode(ctx, code); err == nil && provider != nil {
		return provider, nil
	}

	provider, err := s.Repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	_ = s.cacheManager.SetProviderCache(ctx, provider)
	return provider, nil
}

func (s *providerService) GetProviderByID(ctx context.Context,
	id string) (*db.Provider, error) {

	if provider, err := s.cacheManager.GetProviderCacheByID(ctx, id); err == nil && provider != nil {
		return provider, nil
	}

	provider, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cacheManager.SetProviderCache(ctx, provider)
	return provider, nil
}

func (s *providerService) UpdateProvider(ctx context.Context,
	arg db.UpdateProviderConfigParams) (*db.Provider, error) {

	updated, err := s.Repo.Update(ctx, arg)
	if err != nil {
		s.Logger.Error("Failed to update provider", zap.Error(err))
		return nil, err
	}

	if err := s.cacheManager.SetProviderCache(ctx, updated); err != nil {
		s.Logger.Warn("Failed to update provider cache", zap.Error(err))
	}

	return updated, nil
}

func (s *providerService) InvalidateProviderCache(ctx context.Context,
	id, code string) {
	_ = s.cacheManager.InvalidateProviderCache(ctx, id, code)

	provider, err := s.Repo.GetByID(ctx, id)
	if err == nil && provider != nil {
		_ = s.cacheManager.SetProviderCache(ctx, provider)
	}
}

func (s *providerService) VerifyWebhookSignature(
	ctx context.Context,
	providerCode, signatureHeader string,
	body []byte,
) (bool, error) {

	provider, err := s.Repo.GetByCode(ctx, providerCode)
	if err != nil {
		s.Logger.Error("Failed to get provider", zap.Error(err))
		return false, err
	}

	code := strings.ToLower(provider.Code)
	if code == paystack.ProviderPaystack {
		return s.verifyPaystackSignature(provider, body, signatureHeader)
	}

	if code == flutterwave.ProviderFlutterwave {
		return false, fmt.Errorf("flutterwave webhook verification not implemented")
	}

	return false, fmt.Errorf("unsupported provider for webhook verification")
}

func (s *providerService) verifyPaystackSignature(provider *db.Provider,
	body []byte, signatureHeader string) (bool, error) {

	var cfg PaystackConfig
	if err := json.Unmarshal(provider.Config, &cfg); err != nil {
		s.Logger.Error("Failed to decode paystack config", zap.Error(err))
		return false, err
	}

	gateway := gateways.NewPaystackProvider(s.Logger, cfg.BaseURL, cfg.SecretKey)

	return gateway.VerifyWebhookSignature(body, signatureHeader)
}

func (s *providerService) VerifyPaystackTransaction(ctx context.Context, reference string) (*gateways.VerifyResponse, error) {
	// Find the best Paystack provider config (for now, just get by code)
	provider, err := s.GetProviderByCode(ctx, "paystack")
	if err != nil {
		return nil, err
	}
	var cfg PaystackConfig
	if err := json.Unmarshal(provider.Config, &cfg); err != nil {
		return nil, err
	}
	gateway := gateways.NewPaystackProvider(s.Logger, cfg.BaseURL, cfg.SecretKey)
	return gateway.VerifyTransaction(ctx, reference)
}
