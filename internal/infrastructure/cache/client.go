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

// NewRedisCacheManager initializes the Redis client, session store, and nonce store, and returns a CacheManager.
func NewRedisCacheManager(cfg *config.Config) CacheManager {
	redisClient := InitRedis(cfg)
	logger := config.InitLogger().Logger
	return NewCacheManager(
		NewRedisSessionStore(redisClient, logger),
	)
}
