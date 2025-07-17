package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
)

type userRepository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &userRepository{q: q}
}
