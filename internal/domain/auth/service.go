package auth

import (
	"codematic/internal/config"
	"codematic/internal/domain/tenants"
	"codematic/internal/domain/user"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type authService struct {
	userService  user.Service
	authRepo     Repository
	tenantRepo   tenants.Repository
	cacheManager cache.CacheManager
	JwtManager   *utils.JWTManager
	cfg          *config.Config
	logger       *zap.Logger
}

func NewService(
	userService user.Service,
	authRepo Repository,
	tenantRepo tenants.Repository,
	cacheManager cache.CacheManager,
	jwtManager *utils.JWTManager,
	cfg *config.Config,
	logger *zap.Logger) Service {
	return &authService{
		userService:  userService,
		authRepo:     authRepo,
		tenantRepo:   tenantRepo,
		cacheManager: cacheManager,
		JwtManager:   jwtManager,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *authService) Signup(ctx context.Context, req *SignupRequest) (User, error) {
	// Get tenant by ID
	tenant, err := s.tenantRepo.GetTenantByID(ctx, req.TenantID)
	if err != nil {
		return User{}, errors.New("invalid tenant")
	}

	userReq := &user.CreateUserRequest{
		TenantID: tenant.ID.String(),
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		IsActive: true,
	}
	created, err := s.userService.CreateUser(ctx, userReq)
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
	sessionInfo model.UserSessionInfo) (interface{}, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userService.GetUserByEmailAndTenantID(ctx, email, req.TenantID)
	if err != nil || !user.IsActive.Bool {
		return nil, errors.New(model.InvalidCredentials)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New(model.InvalidCredentials)
	}

	tokenID := uuid.New().String()
	jwtData := model.JWTData{
		UserID:   user.ID.String(),
		Email:    user.Email,
		TenantID: user.TenantID.String(),
		TokenID:  tokenID,
	}
	jwt, err := s.JwtManager.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	s.cacheManager.SetSession(ctx, tokenID, &sessionInfo, 24*time.Hour)

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

	jwtData := model.JWTData{
		UserID:   claims.UserID,
		Email:    claims.Email,
		TenantID: claims.TenantID,
		TokenID:  claims.TenantID,
	}

	jwt, err := s.JwtManager.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(jwtData)
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
