package user

import (
	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"

	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// userService implements the business logic for user-related operations.
type userService struct {
	DB         *db.DBConn
	Repo       Repository
	JwtManager *utils.JWTManager
	logger     *zap.Logger
}

// NewService initializes and returns a new instance of the user service.
func NewService(db *db.DBConn,
	jwtManager *utils.JWTManager,
	logger *zap.Logger,
) Service {
	return &userService{
		DB:         db,
		Repo:       NewRepository(db.Queries, db.Pool),
		JwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *userService) WithTx(q *dbsqlc.Queries) Service {
	return &userService{
		DB:         s.DB,
		Repo:       NewRepository(q, s.DB.Pool),
		JwtManager: s.JwtManager,
		logger:     s.logger,
	}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dbsqlc.User, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByEmailAndTenantID(ctx context.Context, email string, tenantID string) (dbsqlc.User, error) {
	return s.Repo.GetUserByEmailAndTenantID(ctx, email, tenantID)
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (dbsqlc.User, error) {
	return s.Repo.GetUserByID(ctx, userID)
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (dbsqlc.User, error) {
	// Email validation
	email := req.Email
	email = utils.StringOrEmpty(&email)
	email = regexp.MustCompile(`\s+`).ReplaceAllString(email, "")
	if !utils.IsValidEmail(email) {
		return dbsqlc.User{}, model.ErrInvalidEmailFormat
	}

	// Check for duplicate
	existing, err := s.Repo.GetUserByEmail(ctx, email)
	if err == nil && existing.ID != (pgtype.UUID{}) {
		return dbsqlc.User{}, errors.New("user already exists")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return dbsqlc.User{}, errors.New("failed to hash password")
	}

	userID := uuid.New()
	role := string(model.RoleUser)
	if req.Role != "" {
		role = string(req.Role)
	}
	params := dbsqlc.CreateUserParams{
		ID:           utils.ToUUID(userID),
		TenantID:     utils.ToUUID(uuid.MustParse(req.TenantID)),
		Email:        email,
		Phone:        utils.ToDBString(&req.Phone),
		PasswordHash: hash,
		IsActive:     pgtype.Bool{Bool: req.IsActive, Valid: true},
		Role:         utils.ToDBString(&role),
	}
	created, err := s.Repo.CreateUser(ctx, params)
	if err != nil {
		return dbsqlc.User{}, err
	}
	return created, nil
}
