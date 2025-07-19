package provider

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type idempotencyRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &idempotencyRepository{
		q: q,
		p: pool,
	}
}

func (r *idempotencyRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

func (r *idempotencyRepository) Create(ctx context.Context, arg CreateProviderParams) (*db.Provider, error) {
	configJSON, err := json.Marshal(arg.Config)
	if err != nil {
		return nil, err
	}

	p, err := r.q.CreateProvider(ctx, db.CreateProviderParams{
		Name:   arg.Name,
		Code:   arg.Code,
		Config: configJSON,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *idempotencyRepository) GetByID(ctx context.Context, id string) (*db.Provider, error) {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, err
	}
	p, err := r.q.GetProviderByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *idempotencyRepository) GetByCode(ctx context.Context, code string) (*db.Provider, error) {
	p, err := r.q.GetProviderByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *idempotencyRepository) ListActive(ctx context.Context) ([]db.Provider, error) {
	return r.q.ListActiveProviders(ctx)
}

func (r *idempotencyRepository) Update(ctx context.Context, arg db.UpdateProviderConfigParams) (*db.Provider, error) {
	p, err := r.q.UpdateProviderConfig(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *idempotencyRepository) UpdateConfig(ctx context.Context, id string, config map[string]interface{}) (*db.Provider, error) {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p, err := r.q.UpdateProviderConfig(ctx, db.UpdateProviderConfigParams{
		ID:     uid,
		Config: configJSON,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *idempotencyRepository) Deactivate(ctx context.Context, id string) error {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return err
	}
	return r.q.DeactivateProvider(ctx, uid)
}
