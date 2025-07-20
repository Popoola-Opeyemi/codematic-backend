package app

import (
	"codematic/internal/consumers"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
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
	Wallet wallet.Service
	// Add other services as needed
}

func InitServices(
	logger *zap.Logger,
	store *db.DBConn,
	cacheManager cache.CacheManager,
	jwtManager *utils.JWTManager,
	kafkaProducer *kafka.KafkaProducer,
) *Services {

	walletProvider := provider.NewService(store, cacheManager, logger, kafkaProducer)

	userService := user.NewService(store, jwtManager, logger)

	walletService := wallet.NewService(logger, walletProvider, userService, store, kafkaProducer)

	return &Services{
		Wallet: walletService,
	}
}

func InitScheduler(logger *zap.Logger) *scheduler.Scheduler {

	logger.Info("Initializing scheduler...")

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

	logger.Info("Scheduler started.")
	return sched
}

func StartConsumers(ctx context.Context, broker string, services *Services, logger *zap.Logger) {
	logger.Info("Starting Kafka consumers...")
	consumers.StartWalletPaystackConsumer(ctx, broker, services.Wallet, logger)
	logger.Info("Wallet Paystack consumer started.", zap.String("consumer", "wallet_paystack"))
}
