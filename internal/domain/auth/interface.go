package auth

import (
	"codematic/internal/shared/model"
	"context"
)

type Service interface {
	Signup(ctx context.Context, req *SignupRequest) (User, error)
	Login(ctx context.Context, req *LoginRequest,
		sessionInfo model.UserSessionInfo) (interface{}, error)
	Logout(ctx context.Context, userId string) error
}

type Repository interface {
}
