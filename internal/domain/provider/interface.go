package provider

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	InitiateDeposit(ctx context.Context, req DepositRequest) (string, error)
	InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (string, error)
	GetTransactionStatus(ctx context.Context, reference string) (string, error)
}

type Repository interface {
	Create(ctx context.Context, arg CreateProviderParams) (*db.Provider, error)
	GetByID(ctx context.Context, id string) (*db.Provider, error)
	GetByCode(ctx context.Context, code string) (*db.Provider, error)
	ListActive(ctx context.Context) ([]db.Provider, error)
	UpdateConfig(ctx context.Context, id string, config map[string]interface{}) (*db.Provider, error)
	Deactivate(ctx context.Context, id string) error
	WithTx(q *db.Queries) Repository
	Update(ctx context.Context, arg db.UpdateProviderConfigParams) (*db.Provider, error)
}
