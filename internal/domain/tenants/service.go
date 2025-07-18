package tenants

import (
	"codematic/internal/shared/utils"
	"context"

	"go.uber.org/zap"
)

type tenantService struct {
	repo       Repository
	JwtManager *utils.JWTManager
	logger     *zap.Logger
}

func NewService(repo Repository, jwtManager *utils.JWTManager,
	logger *zap.Logger) Service {
	return &tenantService{
		repo:       repo,
		JwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *tenantService) GetTenantByID(ctx context.Context,
	id string) (Tenant, error) {
	dbTenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) CreateTenant(ctx context.Context,
	req CreateTenantRequest) (Tenant, error) {
	dbTenant, err := s.repo.CreateTenant(ctx, req.Name, req.Slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) ListTenants(ctx context.Context) ([]Tenant, error) {
	dbTenants, err := s.repo.ListTenants(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainTenants(dbTenants), nil
}

func (s *tenantService) GetTenantBySlug(ctx context.Context,
	slug string) (Tenant, error) {
	dbTenant, err := s.repo.GetTenantBySlug(ctx, slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) UpdateTenant(ctx context.Context,
	id, name, slug string) (Tenant, error) {
	dbTenant, err := s.repo.UpdateTenant(ctx, id, name, slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) DeleteTenant(ctx context.Context, id string) error {
	return s.repo.DeleteTenant(ctx, id)
}
