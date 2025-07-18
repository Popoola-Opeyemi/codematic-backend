package tenants

import db "codematic/internal/infrastructure/db/sqlc"

type (
	CreateTenantRequest struct {
		Name string `json:"id" validate:"required"`
		Slug string `json:"slug" validate:"required"`
	}

	CreateTenantResponse struct {
		ID string `json:"id"`
	}

	Tenant struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug" `
	}
)

func toDomainTenant(dbTenant db.Tenant) Tenant {
	return Tenant{
		ID:   dbTenant.ID.String(),
		Name: dbTenant.Name,
		Slug: dbTenant.Slug,
	}
}

func toDomainTenants(dbTenants []db.Tenant) []Tenant {
	ts := make([]Tenant, len(dbTenants))
	for i, t := range dbTenants {
		ts[i] = toDomainTenant(t)
	}
	return ts
}
