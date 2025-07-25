package auth

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/model"
	"context"
)

type Service interface {
	Signup(ctx context.Context, req *SignupRequest) (User, error)
	Login(ctx context.Context, req *LoginRequest,
		sessionInfo model.UserSessionInfo) (interface{}, error)
	AdminLogin(ctx context.Context, req *LoginRequest,
		sessionInfo model.UserSessionInfo) (interface{}, error)
	Logout(ctx context.Context, userId string) error
	RefreshToken(ctx context.Context, refreshToken string) (JwtAuthData, error)
}

type Repository interface {
	WithTx(q *db.Queries) Repository
}
