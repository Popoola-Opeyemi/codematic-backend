package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisWalletCacheStore struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisWalletCacheStore(client *redis.Client,
	logger *zap.Logger) WalletCacheStore {
	return &RedisWalletCacheStore{client: client, logger: logger}
}

func (r *RedisWalletCacheStore) SetWalletBalance(ctx context.Context,
	walletID string, balance float64, ttl time.Duration) error {
	key := fmt.Sprintf("wallet:balance:%s", walletID)
	return r.client.Set(ctx, key, balance, ttl).Err()
}

func (r *RedisWalletCacheStore) GetWalletBalance(ctx context.Context,
	walletID string) (float64, error) {
	key := fmt.Sprintf("wallet:balance:%s", walletID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Not found
		}
		return 0, err
	}
	var balance float64
	if _, err := fmt.Sscanf(val, "%f", &balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *RedisWalletCacheStore) DeleteWalletBalance(ctx context.Context,
	walletID string) error {
	key := fmt.Sprintf("wallet:balance:%s", walletID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisWalletCacheStore) SetWalletTransactions(ctx context.Context,
	walletID string, txns interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("wallet:transactions:%s", walletID)
	data, err := json.Marshal(txns)
	if err != nil {
		r.logger.Error("Failed to marshal transactions", zap.Error(err))
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisWalletCacheStore) GetWalletTransactions(ctx context.Context,
	walletID string, result interface{}) error {
	key := fmt.Sprintf("wallet:transactions:%s", walletID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Not found
		}
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

func (r *RedisWalletCacheStore) DeleteWalletTransactions(ctx context.Context,
	walletID string) error {
	key := fmt.Sprintf("wallet:transactions:%s", walletID)
	return r.client.Del(ctx, key).Err()
}
