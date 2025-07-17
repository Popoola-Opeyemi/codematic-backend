package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type userRepository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &userRepository{q: q}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *userRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, params)
}
