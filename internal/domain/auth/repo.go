package auth

import (
	db "codematic/internal/infrastructure/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &authRepository{
		q: q,
		p: pool,
	}
}

func (r *authRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}
