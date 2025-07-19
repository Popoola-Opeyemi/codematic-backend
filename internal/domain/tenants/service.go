package tenants

import (
	"codematic/internal/infrastructure/db"
	dbsqlc "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"

	"context"

	"go.uber.org/zap"
)

type tenantService struct {
	DB         *db.DBConn
	Repo       Repository
	JwtManager *utils.JWTManager
	logger     *zap.Logger
}

func NewService(db *db.DBConn, jwtManager *utils.JWTManager,
	logger *zap.Logger) Service {
	return &tenantService{
		DB:         db,
		Repo:       NewRepository(db.Queries, db.Pool),
		JwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *tenantService) WithTx(q *dbsqlc.Queries) Service {
	return &tenantService{
		DB:         s.DB,
		Repo:       NewRepository(q, s.DB.Pool),
		JwtManager: s.JwtManager,
		logger:     s.logger,
	}
}

func (s *tenantService) GetTenantByID(ctx context.Context,
	id string) (Tenant, error) {
	dbTenant, err := s.Repo.GetTenantByID(ctx, id)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) CreateTenant(ctx context.Context,
	req CreateTenantRequest) (Tenant, error) {
	dbTenant, err := s.Repo.CreateTenant(ctx, req.Name, req.Slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) ListTenants(ctx context.Context) ([]Tenant, error) {
	dbTenants, err := s.Repo.ListTenants(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainTenants(dbTenants), nil
}

func (s *tenantService) GetTenantBySlug(ctx context.Context,
	slug string) (Tenant, error) {
	dbTenant, err := s.Repo.GetTenantBySlug(ctx, slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) UpdateTenant(ctx context.Context,
	id, name, slug string) (Tenant, error) {
	dbTenant, err := s.Repo.UpdateTenant(ctx, id, name, slug)
	if err != nil {
		return Tenant{}, err
	}
	return toDomainTenant(dbTenant), nil
}

func (s *tenantService) DeleteTenant(ctx context.Context, id string) error {
	return s.Repo.DeleteTenant(ctx, id)
}
