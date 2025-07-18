package auth

import (
	db "codematic/internal/infrastructure/db/sqlc"
)

type authRepository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &authRepository{q: q}
}
