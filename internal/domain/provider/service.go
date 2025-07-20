package provider

import (
	"context"
	"fmt"

	"codematic/internal/infrastructure/cache"
	dbconn "codematic/internal/infrastructure/db"
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/infrastructure/events/kafka"
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

// InitiateDeposit initiates a deposit transaction in a provider-agnostic way
func (s *providerService) InitiateDeposit(ctx context.Context, req DepositRequest) (string, error) {

	provider, err := s.Repo.SelectBestProviderByCurrencyAndChannel(ctx, req.Currency, req.Channel)
	if err != nil {
		s.Logger.Error("No provider available", zap.Error(err))
		return "", fmt.Errorf("no provider available for currency %s and channel %s", req.Currency, req.Channel)
	}

	if provider.Code == paystack.ProviderPaystack {

		paystackProvider, err := s.GetProviderByID(ctx, provider.ID.String())
		if err == nil && provider != nil {
			return "", nil
		}

		x := paystackProvider.Config
		// gateway := gateways.NewPaystackProvider(paystackProvider., apiKey string, client *paystack.Client)

	}
	s.Logger.Sugar().Info("Provider", provider.Code)

	return "", nil

}

func (s *providerService) InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (string, error) {

	return "ref", nil
}

func (s *providerService) GetProviderByCode(ctx context.Context, code string) (ProviderDetails, error) {

	cache, err := s.cacheManager.GetProviderCacheByCode(ctx, code)
	if err == nil && cache != nil {
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
		return provider, nil
	}

	provider, err = s.Repo.GetByID(ctx, id)
	if err != nil || provider == nil {
		return provider, err
	}

	_ = s.cacheManager.SetProviderCache(ctx, provider)
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
