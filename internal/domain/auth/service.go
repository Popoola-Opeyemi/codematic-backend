package auth

import (
	"codematic/internal/config"
	"codematic/internal/domain/user"
	"codematic/internal/infrastructure/cache"
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *authService) Signup(ctx context.Context, req *SignupRequest) (User, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	existing, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existing.ID != (pgtype.UUID{}) {
		return User{}, errors.New("user already exists")
	}

	// Get tenant by slug
	tenant, err := s.authRepo.GetTenantBySlug(ctx, req.TenantSlug)
	if err != nil {
		return User{}, errors.New("invalid tenant")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, errors.New("failed to hash password")
	}

	userID := uuid.New()
	params := db.CreateUserParams{
		ID:           utils.ToUUID(userID),
		TenantID:     tenant.ID,
		Email:        email,
		Phone:        utils.ToDBString(&req.Phone),
		PasswordHash: hash,
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	}
	created, err := s.userRepo.CreateUser(ctx, params)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        created.ID.String(),
		Email:     created.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		TenantID:  created.TenantID.String(),
		Role:      "user",
	}, nil
}

func (s *authService) Login(ctx context.Context, req *LoginRequest,
	sessionInfo interface{}) (interface{}, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil || !user.IsActive.Bool {
		return nil, errors.New(model.InvalidCredentials)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New(model.InvalidCredentials)
	}

	tokenID := uuid.New().String()
	jwt, err := s.JwtManager.GenerateJWT(user.ID.String(), tokenID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(user.ID.String(), tokenID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Store session in cache
	if s.cacheManager != nil {
		sess := &model.UserSessionInfo{
			UserAgent: "",
			IPAddress: "",
			TokenID:   tokenID,
		}
		s.cacheManager.SetSession(ctx, tokenID, sess, 24*time.Hour)
	}

	return LoginResponse{
		AccessToken:  jwt,
		RefreshToken: refresh,
		ExpiresIn:    86400,
		TokenType:    "Bearer",
		User: User{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: "",
			LastName:  "",
			TenantID:  user.TenantID.String(),
			Role:      "user",
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (interface{}, error) {
	claims, err := s.JwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	userID := claims.UserID
	tokenID := uuid.New().String()
	jwt, err := s.JwtManager.GenerateJWT(userID, tokenID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refresh, err := s.JwtManager.GenerateRefreshToken(userID, tokenID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}
	return map[string]interface{}{
		"access_token":  jwt,
		"refresh_token": refresh,
		"expires_in":    86400,
		"token_type":    "Bearer",
	}, nil
}

func (s *authService) Logout(ctx context.Context, userId string) error {
	return nil
}
