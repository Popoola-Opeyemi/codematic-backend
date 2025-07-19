package user

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &userRepository{
		q: q,
		p: pool,
	}
}

func (r *userRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *userRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, params)
}

func (r *userRepository) GetUserByEmailAndTenantID(ctx context.Context,
	email string, tenantID string) (db.User, error) {
	uuidTenant, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return db.User{}, err
	}
	return r.q.GetUserByEmailAndTenantID(ctx, db.GetUserByEmailAndTenantIDParams{
		Email:    email,
		TenantID: uuidTenant,
	})
}
