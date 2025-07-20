package auth

import (
	"codematic/internal/config"
	"codematic/internal/domain/tenants"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
	"codematic/internal/infrastructure/cache"
	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"
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
	DB   *db.DBConn
	Repo Repository

	userService   user.Service
	walletService wallet.Service
	tenantService tenants.Service
	cacheManager  cache.CacheManager
	JwtManager    *utils.JWTManager
	cfg           *config.Config
	logger        *zap.Logger
}

// NewService initializes and returns a new instance of the auth service.
func NewService(
	db *db.DBConn,
	userService user.Service,
	walletService wallet.Service,
	tenantService tenants.Service,
	cacheManager cache.CacheManager,
	jwtManager *utils.JWTManager,
	cfg *config.Config,
	logger *zap.Logger) Service {
	return &authService{
		DB:            db,
		Repo:          NewRepository(db.Queries, db.Pool),
		tenantService: tenantService,
		userService:   userService,
		walletService: walletService,
		cacheManager:  cacheManager,
		JwtManager:    jwtManager,
		cfg:           cfg,
		logger:        logger,
	}
}

// / externalServicesWithTx returns user and wallet services bound to a transaction.
func (s *authService) externalServicesWithTx(q *dbsqlc.Queries) (user.Service, wallet.Service) {
	return s.userService.WithTx(q), s.walletService.WithTx(q)
}

func (s *authService) Signup(ctx context.Context, req *SignupRequest) (User, error) {
	var result User

	s.logger.Debug("Starting Signup", zap.String("email", req.Email), zap.String("tenantID", req.TenantID))

	err := utils.WithTX(ctx, s.DB.Pool, func(q *dbsqlc.Queries) error {
		tenant, err := s.tenantService.GetTenantByID(ctx, req.TenantID)
		if err != nil {
			s.logger.Error("Invalid tenant", zap.Error(err))
			return errors.New("invalid tenant")
		}

		userTx, walletTx := s.externalServicesWithTx(q)

		userReq := &user.CreateUserRequest{
			TenantID: tenant.ID,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: req.Password,
			IsActive: true,
		}

		created, err := userTx.CreateUser(ctx, userReq)
		if err != nil {
			s.logger.Error("Create user failed", zap.Error(err))
			return err
		}

		if _, err := walletTx.CreateWalletForNewUser(ctx, created.ID.String()); err != nil {
			s.logger.Error("Create wallet failed", zap.Error(err))
			return err
		}

		result = User{
			ID:        created.ID.String(),
			Email:     created.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			TenantID:  created.TenantID.String(),
			Role:      created.Role.String,
		}
		return nil
	})

	if err != nil {
		s.logger.Error("Signup transaction failed", zap.Error(err))
		return User{}, err
	}

	s.logger.Info("Signup successful", zap.String("userID", result.ID))
	return result, nil
}

// Login authenticates a user, manages session caching, and returns JWT tokens
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

	existingTokenID, err := s.cacheManager.GetTokenIDForUser(ctx, user.ID.String())
	var session *model.UserSessionInfo
	if err == nil && existingTokenID != "" {
		session, _ = s.cacheManager.GetSession(ctx, existingTokenID)
		_ = s.cacheManager.DeleteSession(ctx, existingTokenID)
	}

	tokenID := uuid.New().String()
	jwtData := model.JWTData{
		UserID:   user.ID.String(),
		Email:    user.Email,
		TenantID: user.TenantID.String(),
		TokenID:  tokenID,
		Role:     user.Role.String,
	}

	jwt, err := s.JwtManager.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	var sessionToStore model.UserSessionInfo
	if session != nil {
		sessionToStore = *session
		sessionToStore.TokenID = tokenID
		sessionToStore.LoginTime = time.Now()
		sessionToStore.LastSeen = time.Now()
		sessionToStore.IsActive = true
	} else {
		sessionToStore = sessionInfo
		sessionToStore.UserID = user.ID.String()
		sessionToStore.TokenID = tokenID
		sessionToStore.LoginTime = time.Now()
		sessionToStore.LastSeen = time.Now()
		sessionToStore.IsActive = true
	}

	s.cacheManager.SetSession(ctx, tokenID, &sessionToStore, utils.SessionExpiry)

	return LoginResponse{
		Auth: JwtAuthData{
			AccessToken:  jwt,
			RefreshToken: refresh,
			ExpiresIn:    int(utils.SessionExpiry.Seconds()),
			TokenType:    "Bearer",
		},
		User: User{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: "",
			LastName:  "",
			TenantID:  user.TenantID.String(),
			Role:      user.Role.String,
		},
	}, nil
}

// AdminLogin authenticates a platform admin and returns tokens
func (s *authService) AdminLogin(ctx context.Context, req *LoginRequest,
	sessionInfo model.UserSessionInfo) (interface{}, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil || !user.IsActive.Bool {
		return nil, errors.New(model.InvalidCredentials)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New(model.InvalidCredentials)
	}

	if user.Role.String != model.RolePlatformAdmin.String() {
		return nil, errors.New("not an admin user")
	}

	existingTokenID, err := s.cacheManager.GetTokenIDForUser(ctx, user.ID.String())
	var session *model.UserSessionInfo
	if err == nil && existingTokenID != "" {
		session, _ = s.cacheManager.GetSession(ctx, existingTokenID)
		_ = s.cacheManager.DeleteSession(ctx, existingTokenID)
	}

	tokenID := uuid.New().String()
	jwtData := model.JWTData{
		UserID:   user.ID.String(),
		Email:    user.Email,
		TenantID: "",
		TokenID:  tokenID,
		Role:     user.Role.String,
	}

	jwt, err := s.JwtManager.GenerateJWT(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(jwtData)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	var sessionToStore model.UserSessionInfo
	if session != nil {
		sessionToStore = *session
		sessionToStore.TokenID = tokenID
		sessionToStore.LoginTime = time.Now()
		sessionToStore.LastSeen = time.Now()
		sessionToStore.IsActive = true
	} else {
		sessionToStore = sessionInfo
		sessionToStore.UserID = user.ID.String()
		sessionToStore.TokenID = tokenID
		sessionToStore.LoginTime = time.Now()
		sessionToStore.LastSeen = time.Now()
		sessionToStore.IsActive = true
	}

	s.cacheManager.SetSession(ctx, tokenID, &sessionToStore, utils.SessionExpiry)

	return LoginResponse{
		Auth: JwtAuthData{
			AccessToken:  jwt,
			RefreshToken: refresh,
			ExpiresIn:    int(utils.SessionExpiry.Seconds()),
			TokenType:    "Bearer",
		},
		User: User{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role.String,
		},
	}, nil
}

// RefreshToken generates a new access and refresh token pair from a valid refresh token.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (JwtAuthData, error) {
	claims, err := s.JwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return JwtAuthData{}, errors.New("invalid or expired refresh token")
	}

	jwtData := model.JWTData{
		UserID:   claims.UserID,
		Email:    claims.Email,
		TenantID: claims.TenantID,
		TokenID:  claims.ID,
		Role:     claims.Role,
	}

	jwt, err := s.JwtManager.GenerateJWT(jwtData)
	if err != nil {
		return JwtAuthData{}, errors.New("failed to generate token")
	}

	refresh, err := s.JwtManager.GenerateRefreshToken(jwtData)
	if err != nil {
		return JwtAuthData{}, errors.New("failed to generate refresh token")
	}

	return JwtAuthData{
		AccessToken:  jwt,
		RefreshToken: refresh,
		ExpiresIn:    int(utils.SessionExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// Logout handles user logout
func (s *authService) Logout(ctx context.Context, userId string) error {
	tokenID, err := s.cacheManager.GetTokenIDForUser(ctx, userId)
	if err != nil {
		return errors.New("failed to get token id for user")
	}

	if tokenID == "" {
		return errors.New("no active session found for user")
	}

	if err := s.cacheManager.DeleteSession(ctx, tokenID); err != nil {
		return errors.New("failed to delete session")
	}

	return nil
}
