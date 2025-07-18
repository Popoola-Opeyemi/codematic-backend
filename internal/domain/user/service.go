package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *userService) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByEmailAndTenantID(ctx context.Context, email string, tenantID string) (db.User, error) {
	return s.repo.GetUserByEmailAndTenantID(ctx, email, tenantID)
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (db.User, error) {
	// Email validation
	email := req.Email
	email = utils.StringOrEmpty(&email)
	email = regexp.MustCompile(`\s+`).ReplaceAllString(email, "")
	if !utils.IsValidEmail(email) {
		return db.User{}, model.ErrInvalidEmailFormat
	}

	// Check for duplicate
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existing.ID != (pgtype.UUID{}) {
		return db.User{}, errors.New("user already exists")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return db.User{}, errors.New("failed to hash password")
	}

	userID := uuid.New()
	params := db.CreateUserParams{
		ID:           utils.ToUUID(userID),
		TenantID:     utils.ToUUID(uuid.MustParse(req.TenantID)),
		Email:        email,
		Phone:        utils.ToDBString(&req.Phone),
		PasswordHash: hash,
		IsActive:     pgtype.Bool{Bool: req.IsActive, Valid: true},
	}
	created, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return created, nil
}
