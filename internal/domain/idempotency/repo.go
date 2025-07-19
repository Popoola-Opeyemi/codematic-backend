package idempotency

import (
	"context"
	"encoding/json"

	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &repository{
		q: q,
		p: pool,
	}
}

func (r *repository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

func (r *repository) Create(ctx context.Context, arg CreateParams) error {
	bodyJSON, err := json.Marshal(arg.ResponseBody)
	if err != nil {
		return err
	}

	id, err := utils.StringToPgUUID(arg.ID)
	if err != nil {
		return err
	}
	tid, err := utils.StringToPgUUID(arg.TenantID)
	if err != nil {
		return err
	}

	var uid pgtype.UUID
	if arg.UserID != "" {
		uid, err = utils.StringToPgUUID(arg.UserID)
		if err != nil {
			return err
		}
	}

	return r.q.CreateIdempotencyKey(ctx, db.CreateIdempotencyKeyParams{
		ID:             id,
		TenantID:       tid,
		UserID:         uid,
		IdempotencyKey: arg.Key,
		Endpoint:       arg.Endpoint,
		RequestHash:    arg.RequestHash,
		ResponseBody:   bodyJSON,
		StatusCode:     utils.ToPgxInt4(int32(arg.StatusCode)),
	})
}

func (r *repository) Get(ctx context.Context, tenantID, key, endpoint, requestHash string) (*db.IdempotencyKey, error) {
	tid, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return nil, err
	}

	record, err := r.q.GetIdempotencyRecord(ctx, db.GetIdempotencyRecordParams{
		TenantID:       tid,
		IdempotencyKey: key,
		Endpoint:       endpoint,
		RequestHash:    requestHash,
	})
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *repository) UpdateResponse(ctx context.Context, arg UpdateResponseParams) (*db.IdempotencyKey, error) {
	// Marshal response body to JSON
	bodyJSON, err := json.Marshal(arg.ResponseBody)
	if err != nil {
		return nil, err
	}

	// Convert TenantID to pgtype.UUID
	tid, err := utils.StringToPgUUID(arg.TenantID)
	if err != nil {
		return nil, err
	}

	// Perform the update using SQLC-generated method
	result, err := r.q.UpdateIdempotencyKeyResponse(ctx, db.UpdateIdempotencyKeyResponseParams{
		ResponseBody:   bodyJSON,
		StatusCode:     utils.ToPgxInt4(int32(arg.StatusCode)),
		TenantID:       tid,
		IdempotencyKey: arg.Key,
		Endpoint:       arg.Endpoint,
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}
