package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisTransactionCacheStore struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisTransactionCacheStore(client *redis.Client, logger *zap.Logger) TransactionCacheStore {
	return &RedisTransactionCacheStore{client: client, logger: logger}
}

func (r *RedisTransactionCacheStore) SetTransaction(ctx context.Context, txID string, tx interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("transaction:%s", txID)
	data, err := json.Marshal(tx)
	if err != nil {
		r.logger.Error("Failed to marshal transaction", zap.Error(err))
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisTransactionCacheStore) GetTransaction(ctx context.Context, txID string, result interface{}) error {
	key := fmt.Sprintf("transaction:%s", txID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Not found
		}
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

func (r *RedisTransactionCacheStore) DeleteTransaction(ctx context.Context, txID string) error {
	key := fmt.Sprintf("transaction:%s", txID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisTransactionCacheStore) SetTransactionsByUser(ctx context.Context, userID string, txns interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("transactions:user:%s", userID)
	data, err := json.Marshal(txns)
	if err != nil {
		r.logger.Error("Failed to marshal user transactions", zap.Error(err))
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisTransactionCacheStore) GetTransactionsByUser(ctx context.Context, userID string, result interface{}) error {
	key := fmt.Sprintf("transactions:user:%s", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Not found
		}
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

func (r *RedisTransactionCacheStore) DeleteTransactionsByUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("transactions:user:%s", userID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisTransactionCacheStore) SetTransactionsByTenant(ctx context.Context, tenantID string, txns interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("transactions:tenant:%s", tenantID)
	data, err := json.Marshal(txns)
	if err != nil {
		r.logger.Error("Failed to marshal tenant transactions", zap.Error(err))
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisTransactionCacheStore) GetTransactionsByTenant(ctx context.Context, tenantID string, result interface{}) error {
	key := fmt.Sprintf("transactions:tenant:%s", tenantID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Not found
		}
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

func (r *RedisTransactionCacheStore) DeleteTransactionsByTenant(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("transactions:tenant:%s", tenantID)
	return r.client.Del(ctx, key).Err()
}
