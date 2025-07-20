// @title Codematic API
// @version 1.0
// @description This is the Codematic API documentation.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email code.popoola@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:9082
// @BasePath /
package main

import (
	"codematic/internal/config"
	"codematic/internal/handler"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/router"
	"codematic/internal/shared/model"

	"codematic/internal/shared/utils"
	"context"
	"os"
	"os/signal"
	"syscall"

	"codematic/internal/consumers"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
)

func main() {
	cfg := config.LoadAppConfig()

	zapLogger := config.InitLogger()
	defer zapLogger.Close()

	redisCache := cache.InitRedis(cfg)

	defer redisCache.Close()

	store := db.InitDB(cfg, zapLogger.Logger)

	JWTManager := utils.NewJWTManager(
		cfg.JwtSecret,
		cfg.RefreshTokenSecret,
	)

	cacheManager := cache.NewRedisCacheManager(
		redisCache,
	)

	// setup for Kafka
	broker := os.Getenv("KAFKA_BROKER")
	kafkaProducer := kafka.NewKafkaProducer(broker)
	events.Init(kafkaProducer)

	// App environment
	app := router.InitRouterWithConfig(cfg, redisCache, zapLogger.Logger)

	providers := &model.Providers{}

	env := handler.NewEnvironment(
		cfg,
		app,
		store,
		redisCache,
		zapLogger.Logger,
		providers,
		JWTManager,
		cacheManager,
		kafkaProducer,
	)

	// Set up Kafka consumer for Paystack wallet events
	walletProvider := provider.NewService(store, cacheManager, zapLogger.Logger, kafkaProducer)
	userService := user.NewService(store, JWTManager, zapLogger.Logger)
	walletService := wallet.NewService(walletProvider, userService, store, zapLogger.Logger, kafkaProducer)

	consumers.StartWalletPaystackConsumer(context.Background(), broker, walletService, zapLogger.Logger)

	router.InitHandlers(env, []handler.IHandler{
		&handler.Auth{},
		&handler.Tenants{},
		&handler.Wallet{},
		&handler.Webhook{},
	})

	// Graceful shutdown support
	go func() {
		router.RunWithGracefulShutdown(app, cfg.PORT)
	}()

	// Wait for termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
