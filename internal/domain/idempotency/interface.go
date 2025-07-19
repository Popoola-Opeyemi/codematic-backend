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
	Get(ctx context.Context, tenantID, key, endpoint, requestHash string) (*db.IdempotencyKey, error)
	UpdateResponse(ctx context.Context, arg UpdateResponseParams) (*db.IdempotencyKey, error)
}
