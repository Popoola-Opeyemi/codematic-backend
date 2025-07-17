package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
}
