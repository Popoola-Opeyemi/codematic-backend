package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "codematic/internal/infrastructure/db/sqlc"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	providerCacheKeyByID   = "provider:id:%s"
	providerCacheKeyByCode = "provider:code:%s"
	providerCacheTTL       = 10 * time.Minute
)

type ProviderCacheStore interface {
	SetProviderCache(ctx context.Context, provider *db.Provider) error
	GetProviderCacheByID(ctx context.Context, id string) (*db.Provider, error)
	GetProviderCacheByCode(ctx context.Context, code string) (*db.Provider, error)
	InvalidateProviderCache(ctx context.Context, id, code string) error
}

type RedisProviderCacheStore struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisProviderCacheStore(client *redis.Client, logger *zap.Logger) ProviderCacheStore {
	return &RedisProviderCacheStore{
		client: client,
		logger: logger,
	}
}

func (r *RedisProviderCacheStore) SetProviderCache(ctx context.Context, provider *db.Provider) error {
	if provider == nil {
		return nil
	}

	data, err := json.Marshal(provider)
	if err != nil {
		r.logger.Error("Failed to marshal provider", zap.Error(err))
		return err
	}

	idKey := fmt.Sprintf(providerCacheKeyByID, provider.ID.String())
	codeKey := fmt.Sprintf(providerCacheKeyByCode, provider.Code)

	pipe := r.client.Pipeline()
	pipe.Set(ctx, idKey, data, providerCacheTTL)
	pipe.Set(ctx, codeKey, data, providerCacheTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Error("Failed to set provider cache", zap.Error(err))
		return err
	}

	return nil
}

func (r *RedisProviderCacheStore) GetProviderCacheByID(ctx context.Context, id string) (*db.Provider, error) {
	key := fmt.Sprintf(providerCacheKeyByID, id)
	data, err := r.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("Provider not found in cache by ID", zap.String("id", id))
		} else {
			r.logger.Error("Failed to get provider by ID from cache", zap.Error(err), zap.String("id", id))
		}
		return nil, err
	}

	var provider db.Provider
	if err := json.Unmarshal([]byte(data), &provider); err != nil {
		r.logger.Error("Failed to unmarshal provider data by ID", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return &provider, nil
}

func (r *RedisProviderCacheStore) GetProviderCacheByCode(ctx context.Context, code string) (*db.Provider, error) {
	key := fmt.Sprintf(providerCacheKeyByCode, code)
	data, err := r.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("Provider not found in cache by code", zap.String("code", code))
		} else {
			r.logger.Error("Failed to get provider by code from cache", zap.Error(err), zap.String("code", code))
		}
		return nil, err
	}

	var provider db.Provider
	if err := json.Unmarshal([]byte(data), &provider); err != nil {
		r.logger.Error("Failed to unmarshal provider data by code", zap.Error(err), zap.String("code", code))
		return nil, err
	}

	return &provider, nil
}

func (r *RedisProviderCacheStore) InvalidateProviderCache(ctx context.Context, id, code string) error {
	idKey := fmt.Sprintf(providerCacheKeyByID, id)
	codeKey := fmt.Sprintf(providerCacheKeyByCode, code)

	if err := r.client.Del(ctx, idKey, codeKey).Err(); err != nil {
		r.logger.Error("Failed to invalidate provider cache", zap.Error(err), zap.String("id", id), zap.String("code", code))
		return err
	}
	return nil
}
