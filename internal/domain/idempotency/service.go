package idempotency

import (
	"codematic/internal/config"
	"codematic/internal/infrastructure/db"
	"codematic/internal/shared/utils"

	"context"

	"go.uber.org/zap"
)

type idempotencyService struct {
	DB   *db.DBConn
	Repo Repository

	JwtManager *utils.JWTManager
	cfg        *config.Config
	logger     *zap.Logger
}

func NewService(
	db *db.DBConn,
	cfg *config.Config,
	logger *zap.Logger) Service {
	return &idempotencyService{
		DB:     db,
		Repo:   NewRepository(db.Queries, db.Pool),
		cfg:    cfg,
		logger: logger,
	}
}

func (s *idempotencyService) Create(ctx context.Context, userId string) error {
	return nil
}
