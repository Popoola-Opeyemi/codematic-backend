package handler

import (
	"codematic/internal/app"
	"codematic/internal/config"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/shared/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var validate = validator.New()

type IHandler interface {
	Init(string, *Environment) error
}

type Environment struct {
	Config        *config.Config
	Fiber         *fiber.App
	DB            *db.DBConn
	Logger        *zap.Logger
	Cache         *redis.Client
	JWTManager    *utils.JWTManager
	CacheManager  cache.CacheManager
	KafkaProducer *kafka.KafkaProducer

	Services *app.Services
}

func NewEnvironment(
	config *config.Config,
	fiber *fiber.App,
	db *db.DBConn,
	Cache *redis.Client,
	logger *zap.Logger,
	jwtManager *utils.JWTManager,
	cacheManager cache.CacheManager,
	kafkaProducer *kafka.KafkaProducer,
	services *app.Services,
) *Environment {
	return &Environment{
		Config:        config,
		Fiber:         fiber,
		DB:            db,
		Cache:         Cache,
		Logger:        logger,
		JWTManager:    jwtManager,
		CacheManager:  cacheManager,
		KafkaProducer: kafkaProducer,
		Services:      services,
	}
}
