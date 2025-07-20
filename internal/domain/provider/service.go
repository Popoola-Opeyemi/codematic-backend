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

	Logger   *zap.Logger
	Producer *kafka.KafkaProducer
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
	s.Logger.Sugar().Info("User Email:", email)

	var response gateways.GatewayResponse

	providerRow, err := s.Repo.SelectBestProviderByCurrencyAndChannel(ctx, req.Currency, req.Channel)
	if err != nil {
		s.Logger.Error("No provider available", zap.Error(err))
		return response, fmt.Errorf("no provider available for currency %s and channel %s", req.Currency, req.Channel)
	}

	s.Logger.Sugar().Info("Selected provider:", providerRow.Name)

	// Get full provider with config
	provider, err := s.GetProviderByID(ctx, providerRow.ID.String())
	if err != nil {
		s.Logger.Error("Failed to retrieve provider details", zap.Error(err))
		return response, err
	}

	switch strings.ToLower(provider.Code) {
	case paystack.ProviderPaystack:

		s.Logger.Sugar().Info("Handling Paystack provider")

		var cfg PaystackConfig
		if err := json.Unmarshal(provider.Config, &cfg); err != nil {
			s.Logger.Error("Failed to decode paystack config", zap.Error(err))
			return response, err
		}

		gateway := gateways.NewPaystackProvider(s.Logger, cfg.BaseURL, cfg.SecretKey)

		gatewayReq := gateways.DepositRequest{
			Email:      email,
			Amount:     req.Amount,
			ProviderID: provider.ID.String(),
		}

		res, err := gateway.InitDeposit(ctx, gatewayReq)
		if err != nil {
			s.Logger.Error("Paystack InitDeposit failed", zap.Error(err))
			return response, err
		}
		response = res

	case flutterwave.ProviderFlutterwave:
		return response, fmt.Errorf("flutterwave not yet implemented")

	default:
		return response, fmt.Errorf("unsupported provider: %s", provider.Code)
	}

	s.Logger.Sugar().Info("Paystack deposit initiated successfully", response)

	return response, nil
}

func (s *providerService) InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (string, error) {

	return "ref", nil
}

func (s *providerService) GetProviderByCode(ctx context.Context, code string) (*db.Provider, error) {
	provider, err := s.cacheManager.GetProviderCacheByCode(ctx, code)
	if err == nil && provider != nil {
		return provider, nil
	}
	provider, err = s.Repo.GetByCode(ctx, code)
	if err != nil || provider == nil {
		return provider, err
	}
	_ = s.cacheManager.SetProviderCache(ctx, provider)
	return provider, nil
}

func (s *providerService) GetProviderByID(ctx context.Context, id string) (*db.Provider, error) {
	provider, err := s.cacheManager.GetProviderCacheByID(ctx, id)
	if err == nil && provider != nil {
		s.Logger.Sugar().Info("Provider found in cache ", "id ", id)
		return provider, nil
	}

	if err != nil {
		s.Logger.Warn("Cache lookup failed", zap.Error(err))
	}

	// Fallback to DB
	provider, err = s.Repo.GetByID(ctx, id)
	if err != nil || provider == nil {
		return provider, err
	}

	// Populate cache
	if cacheErr := s.cacheManager.SetProviderCache(ctx, provider); cacheErr != nil {
		s.Logger.Warn("Failed to set provider in cache", zap.Error(cacheErr))
	}

	return provider, nil
}

func (s *providerService) UpdateProvider(ctx context.Context,
	arg db.UpdateProviderConfigParams) (*db.Provider, error) {
	updated, err := s.Repo.Update(ctx, arg)
	if err != nil {
		s.Logger.Error("failed to update provider", zap.Error(err))
		return nil, err
	}

	err = s.cacheManager.SetProviderCache(ctx, updated)
	if err != nil {
		s.Logger.Warn("failed to update provider cache", zap.Error(err))
	}
	return updated, nil
}

func (s *providerService) InvalidateProviderCache(ctx context.Context, id, code string) {
	_ = s.cacheManager.InvalidateProviderCache(ctx, id, code)

	provider, err := s.Repo.GetByID(ctx, id)
	if err == nil && provider != nil {
		_ = s.cacheManager.SetProviderCache(ctx, provider)
	}
}
