package tenants

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"

	"github.com/google/uuid"
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

func (r *repository) CreateTenant(ctx context.Context,
	name string, slug string) (db.Tenant, error) {
	uid := uuid.New().String()
	uuid, err := utils.StringToPgUUID(uid)
	if err != nil {
		return db.Tenant{}, nil
	}

	arg := db.CreateTenantParams{
		ID:   uuid,
		Name: name,
		Slug: slug,
	}
	return r.q.CreateTenant(ctx, arg)
}

func (r *repository) GetTenantByID(ctx context.Context, id string) (db.Tenant, error) {
	uuid, err := utils.StringToPgUUID(id)
	if err != nil {
		return db.Tenant{}, nil
	}
	return r.q.GetTenantByID(ctx, uuid)
}

func (r *repository) GetTenantBySlug(ctx context.Context, slug string) (
	db.Tenant, error) {
	return r.q.GetTenantBySlug(ctx, slug)
}

func (r *repository) ListTenants(ctx context.Context) ([]db.Tenant, error) {
	return r.q.ListTenants(ctx)
}

func (r *repository) UpdateTenant(ctx context.Context, id, name, slug string) (
	db.Tenant, error) {
	uuid, err := utils.StringToPgUUID(id)
	if err != nil {
		return db.Tenant{}, nil
	}
	arg := db.UpdateTenantParams{
		ID:   uuid,
		Name: name,
		Slug: slug,
	}
	return r.q.UpdateTenant(ctx, arg)
}

func (r *repository) DeleteTenant(ctx context.Context, id string) error {
	uuid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil
	}
	return r.q.DeleteTenant(ctx, uuid)
}
