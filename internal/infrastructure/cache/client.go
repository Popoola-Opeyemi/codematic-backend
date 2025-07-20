package cache

import (
	"codematic/internal/config"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitRedis(config *config.Config) *redis.Client {
	var redisClient *redis.Client
	for {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         config.RedisAddr,
			Password:     config.RedisPassword,
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 2,
			MaxIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		if _, err := redisClient.Ping(Ctx).Result(); err != nil {
			log.Printf("failed to connect to Redis: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return redisClient
}

func NewRedisCacheManager(rc *redis.Client) CacheManager {
	logger := config.InitLogger().Logger
	return NewCacheManager(
		NewRedisSessionStore(rc, logger),
		NewRedisProviderCacheStore(rc, logger),
		NewRedisWalletCacheStore(rc, logger),
		NewRedisTransactionCacheStore(rc, logger),
	)
}

type CacheManager interface {
	SessionStore
	ProviderCacheStore
	WalletCacheStore
	TransactionCacheStore
}

type unifiedCacheManager struct {
	SessionStore
	ProviderCacheStore
	WalletCacheStore
	TransactionCacheStore
}

func NewCacheManager(sessionStore SessionStore,
	providerCacheStore ProviderCacheStore,
	walletCacheStore WalletCacheStore,
	transactionCacheStore TransactionCacheStore) CacheManager {
	return &unifiedCacheManager{
		SessionStore:          sessionStore,
		ProviderCacheStore:    providerCacheStore,
		WalletCacheStore:      walletCacheStore,
		TransactionCacheStore: transactionCacheStore,
	}
}
