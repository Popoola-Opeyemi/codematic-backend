package user

import (
	"codematic/internal/shared/utils"

	"go.uber.org/zap"
)

type userService struct {
	repo       Repository
	JwtManager *utils.JWTManager
	logger     *zap.Logger
}

func NewService(repo Repository, jwtManager *utils.JWTManager,
	logger *zap.Logger) Service {
	return &userService{
		repo:       repo,
		JwtManager: jwtManager,
		logger:     logger,
	}
}
