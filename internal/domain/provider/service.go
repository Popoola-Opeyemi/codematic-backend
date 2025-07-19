package provider

import (
	"context"

	"codematic/internal/infrastructure/cache"
	dbconn "codematic/internal/infrastructure/db"
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/shared/model"

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

// InitiateDeposit initiates a deposit transaction in a provider-agnostic way
func (s *providerService) InitiateDeposit(ctx context.Context, req DepositRequest) (string, error) {

	providerCode, ok := req.Metadata["provider"].(string)
	if !ok || providerCode == "" {
		return "", model.ErrUnsupportedProvider
	}

	provider, err := s.GetProviderByCode(ctx, providerCode)
	if err != nil {
		return "", model.ErrUnsupportedProvider
	}

	if provider == nil || !provider.IsActive.Valid || !provider.IsActive.Bool {
		return "", model.ErrUnsupportedProvider
	}

	return "", nil

}

func (s *providerService) InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (string, error) {

	providerCode, ok := req.Metadata["provider"].(string)
	if !ok || providerCode == "" {
		return "", model.ErrUnsupportedProvider
	}

	provider, err := s.GetProviderByCode(ctx, providerCode)
	if err != nil {
		return "", model.ErrUnsupportedProvider
	}
	if provider == nil || !provider.IsActive.Valid || !provider.IsActive.Bool {
		return "", model.ErrUnsupportedProvider
	}

	providerSvc, ok := GetProvider(providerCode)
	if !ok {
		return "", model.ErrUnsupportedProvider
	}

	ref, err := providerSvc.InitiateWithdrawal(ctx, req)
	if err != nil {
		return "", err
	}

	return ref, nil
}

func (s *providerService) GetTransactionStatus(ctx context.Context, reference string) (string, error) {

	return "pending", nil
}
