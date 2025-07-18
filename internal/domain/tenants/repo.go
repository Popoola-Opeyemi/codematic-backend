package tenants

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"

	"github.com/google/uuid"
)

type repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &repository{q: q}
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
