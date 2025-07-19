package tenants

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	GetTenantByID(ctx context.Context, id string) (Tenant, error)
	CreateTenant(ctx context.Context, req CreateTenantRequest) (Tenant, error)
	ListTenants(ctx context.Context) ([]Tenant, error)
	GetTenantBySlug(ctx context.Context, slug string) (Tenant, error)
	UpdateTenant(ctx context.Context, id, name, slug string) (Tenant, error)
	DeleteTenant(ctx context.Context, id string) error
	WithTx(q *db.Queries) Service
}

type Repository interface {
	GetTenantByID(ctx context.Context, id string) (db.Tenant, error)
	CreateTenant(ctx context.Context, name, slug string) (db.Tenant, error)
	ListTenants(ctx context.Context) ([]db.Tenant, error)
	GetTenantBySlug(ctx context.Context, slug string) (db.Tenant, error)
	UpdateTenant(ctx context.Context, id, name, slug string) (db.Tenant, error)
	DeleteTenant(ctx context.Context, id string) error
	WithTx(q *db.Queries) Repository
}
