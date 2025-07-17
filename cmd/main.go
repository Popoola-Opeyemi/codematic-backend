package main

import (
	"codematic/internal/config"
	"codematic/internal/handler"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"
	"codematic/internal/router"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"os"
	"os/signal"
	"syscall"
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
	)

	router.InitHandlers(env, []handler.IHandler{
		&handler.Auth{},
	})

	// Graceful shutdown support
	go func() {
		router.RunWithGracefulShutdown(app, cfg.PORT)
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
