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

	"codematic/internal/shared/utils"
	"context"
	"os"
	"os/signal"
	"syscall"

	"codematic/internal/app"
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

	kafkaProducer := kafka.NewKafkaProducer(cfg.KAFKA_BROKER)
	events.Init(kafkaProducer)

	appEnv := router.InitRouterWithConfig(cfg, redisCache, zapLogger.Logger)

	// Initialize services
	services := app.InitServices(
		zapLogger.Logger,
		store,
		cacheManager,
		JWTManager,
		kafkaProducer,
		cfg,
	)

	// Initialize scheduler
	app.InitScheduler(zapLogger.Logger)

	// Start consumers
	app.StartConsumers(context.Background(), cfg.KAFKA_BROKER, services, zapLogger.Logger)

	env := handler.NewEnvironment(
		cfg,
		appEnv,
		store,
		redisCache,
		zapLogger.Logger,
		JWTManager,
		cacheManager,
		kafkaProducer,
		services,
	)

	router.InitHandlers(env, []handler.IHandler{
		&handler.Auth{},
		&handler.Tenants{},
		&handler.Wallet{},
		&handler.Webhook{},
		&handler.Transactions{},
	})

	go func() {
		router.RunWithGracefulShutdown(appEnv, cfg.PORT, zapLogger.Logger)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
