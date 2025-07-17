package auth

import (
	"codematic/internal/config"
	"codematic/internal/domain/user"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/shared/utils"

	"go.uber.org/zap"
)

type authService struct {
	userRepo     user.Repository
	authRepo     Repository
	cacheManager cache.CacheManager
	JwtManager   *utils.JWTManager
	cfg          *config.Config
	logger       *zap.Logger
}

func NewService(userRepo user.Repository, authRepo Repository,
	cacheManager cache.CacheManager, jwtManager *utils.JWTManager,
	cfg *config.Config, logger *zap.Logger) Service {
	return &authService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		cacheManager: cacheManager,
		JwtManager:   jwtManager,
		cfg:          cfg,
		logger:       logger,
	}
}
