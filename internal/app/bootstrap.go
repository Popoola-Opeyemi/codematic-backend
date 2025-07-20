package app

import (
	"codematic/internal/config"
	"codematic/internal/consumers"
	"codematic/internal/domain/auth"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/tenants"
	"codematic/internal/domain/transactions"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
	"codematic/internal/domain/webhook"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"
	"codematic/internal/infrastructure/events/kafka"
	"codematic/internal/scheduler"
	"codematic/internal/scheduler/jobs"
	"codematic/internal/shared/utils"
	"context"

	"go.uber.org/zap"
)

type Services struct {
	Wallet       wallet.Service
	User         user.Service
	Provider     provider.Service
	Transactions transactions.Service
	Tenants      tenants.Service
	Auth         auth.Service
	Webhook      webhook.Service
}

func InitServices(
	logger *zap.Logger,
	store *db.DBConn,
	cacheManager cache.CacheManager,
	jwtManager *utils.JWTManager,
	kafkaProducer *kafka.KafkaProducer,
	cfg *config.Config,
) *Services {

	logger.Info("initializing services...")

	providerService := provider.NewService(
		store,
		cacheManager,
		logger,
		kafkaProducer,
	)

	userService := user.NewService(store, jwtManager, logger)

	tenantsService := tenants.NewService(store, jwtManager, logger)

	walletService := wallet.NewService(
		logger,
		providerService,
		userService,
		store,
		kafkaProducer,
	)

	authService := auth.NewService(
		store,
		userService,
		walletService,
		tenantsService,
		cacheManager,
		jwtManager,
		cfg, logger,
	)

	webhookService := webhook.NewService(
		providerService,
		tenantsService,
		logger, store, cfg,
		kafkaProducer,
	)

	logger.Info("services initialized.")

	return &Services{
		Wallet:   walletService,
		User:     userService,
		Provider: providerService,
		Tenants:  tenantsService,
		Auth:     authService,
		Webhook:  webhookService,
	}
}

func InitScheduler(logger *zap.Logger) *scheduler.Scheduler {

	logger.Info("initializing scheduler...")

	sched, err := scheduler.New(logger)
	if err != nil {
		logger.Fatal("failed to create scheduler", zap.Error(err))
	}

	jobList := []scheduler.Job{
		jobs.HelloJob{},
	}

	if err := sched.RegisterJobs(context.Background(), jobList); err != nil {
		logger.Fatal("failed to register jobs", zap.Error(err))
	}

	sched.Start()

	logger.Info("scheduler started.")

	return sched
}

func StartConsumers(ctx context.Context, broker string, services *Services, logger *zap.Logger) {

	logger.Info("starting Kafka consumers...")

	consumers.StartWalletPaystackConsumer(ctx, broker, services.Wallet, logger)

	logger.Info("wallet Paystack consumer started.", zap.String("consumer", "wallet_paystack"))
}
