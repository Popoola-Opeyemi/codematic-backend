package cache

import (
	"codematic/internal/shared/model"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SessionStore interface {
	SetSession(ctx context.Context, key string, value *model.UserSessionInfo, ttl time.Duration) error
	GetSession(ctx context.Context, key string) (*model.UserSessionInfo, error)
	DeleteSession(ctx context.Context, key string) error
}

type RedisSessionStore struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisSessionStore(client *redis.Client, logger *zap.Logger) SessionStore {
	return &RedisSessionStore{client: client, logger: logger}
}

func (r *RedisSessionStore) SetSession(ctx context.Context, key string, value *model.UserSessionInfo, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("Failed to marshal session data", zap.Error(err))
	}

	if err := r.client.Set(context.Background(), key, data, ttl).Err(); err != nil {
		r.logger.Error("Failed to set session in redis", zap.Error(err))
	}

	return nil
}

func (r *RedisSessionStore) GetSession(ctx context.Context, key string) (*model.UserSessionInfo, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("Session not found in redis", zap.String("key", key))
		} else {
			r.logger.Error("Failed to get session from redis", zap.Error(err), zap.String("key", key))
		}
		return nil, err
	}

	var session model.UserSessionInfo
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		r.logger.Error("Failed to unmarshal session data", zap.Error(err), zap.String("key", key))
		return nil, err
	}

	return &session, nil
}

func (r *RedisSessionStore) DeleteSession(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Failed to delete session from redis", zap.Error(err), zap.String("key", key))
	}
	return err
}
