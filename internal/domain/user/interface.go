package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByEmailAndTenantID(ctx context.Context, email string, tenantID string) (db.User, error)
	WithTx(q *db.Queries) Service
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetUserByEmailAndTenantID(ctx context.Context, email string, tenantID string) (db.User, error)
	WithTx(q *db.Queries) Repository
}
