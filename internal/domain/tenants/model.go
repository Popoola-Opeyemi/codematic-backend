package tenants

import db "codematic/internal/infrastructure/db/sqlc"

type (
	CreateTenantRequest struct {
		Name       string `json:"id" validate:"required"`
		Slug       string `json:"slug" validate:"required"`
		WebhookURL string `json:"webhook_url"`
	}

	CreateTenantResponse struct {
		ID         string `json:"id"`
		WebhookURL string `json:"webhook_url"`
	}

	Tenant struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Slug       string `json:"slug" `
		WebhookURL string `json:"webhook_url"`
	}
)

func toDomainTenant(dbTenant db.Tenant) Tenant {
	return Tenant{
		ID:         dbTenant.ID.String(),
		Name:       dbTenant.Name,
		Slug:       dbTenant.Slug,
		WebhookURL: dbTenant.WebhookUrl,
	}
}

func toDomainTenants(dbTenants []db.Tenant) []Tenant {
	ts := make([]Tenant, len(dbTenants))
	for i, t := range dbTenants {
		ts[i] = toDomainTenant(t)
	}
	return ts
}
