package auth

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type authRepository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &authRepository{q: q}
}

func (r *authRepository) GetTenantBySlug(ctx context.Context, slug string) (db.Tenant, error) {
	return r.q.GetTenantBySlug(ctx, slug)
}
