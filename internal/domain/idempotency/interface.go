package idempotency

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	Create(ctx context.Context, userId string) error
}

type Repository interface {
	WithTx(q *db.Queries) Repository

	Create(ctx context.Context, arg CreateParams) error
	GetByKeyAndEndpoint(
		ctx context.Context,
		tenantID string,
		key string,
		endpoint string,
	) (*db.IdempotencyKey, error)
}
