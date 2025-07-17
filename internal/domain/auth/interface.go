package auth

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	Signup(ctx context.Context, req *SignupRequest) (User, error)
	Login(ctx context.Context, req *LoginRequest, sessionInfo interface{}) (interface{}, error)
	Logout(ctx context.Context, userId string) error
}

type Repository interface {
	GetTenantBySlug(ctx context.Context, slug string) (db.Tenant, error)
}
