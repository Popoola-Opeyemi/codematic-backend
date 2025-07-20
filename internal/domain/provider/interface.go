package provider

import (
	"codematic/internal/domain/provider/gateways"
	db "codematic/internal/infrastructure/db/sqlc"
	"context"
)

type Service interface {
	InitiateDeposit(ctx context.Context,
		req DepositRequest) (gateways.GatewayResponse, error)
	// InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (string, error)
}

type Repository interface {
	CreateProvider(ctx context.Context, arg CreateProviderParams) (*db.Provider, error)
	GetByID(ctx context.Context, id string) (*db.Provider, error)
	GetByCode(ctx context.Context, code string) (*db.Provider, error)
	ListActiveProviders(ctx context.Context) ([]db.Provider, error)
	UpdateConfig(ctx context.Context, id string, config map[string]interface{}) (*db.Provider, error)
	Update(ctx context.Context, arg db.UpdateProviderConfigParams) (*db.Provider, error)
	Deactivate(ctx context.Context, id string) error
	WithTx(q *db.Queries) Repository

	AddSupportedCurrency(ctx context.Context, providerID, currencyCode string) error
	RemoveSupportedCurrency(ctx context.Context, providerID, currencyCode string) error
	DeleteAllSupportedCurrencies(ctx context.Context, providerID string) error
	ListSupportedCurrencies(ctx context.Context, providerID string) ([]db.Currency, error)
	AddSupportedChannel(ctx context.Context, providerID, channel string) error
	ListProviderDetails(ctx context.Context) ([]db.ListProviderDetailsRow, error)
	SelectBestProviderByCurrencyAndChannel(ctx context.Context,
		currency, channel string) (*db.SelectBestProviderByCurrencyAndChannelRow, error)
	SelectBestProvider(ctx context.Context) (*db.SelectBestProviderRow, error)
	DecayPriority(ctx context.Context) error
	ResetDailyMetrics(ctx context.Context) error
	IncrementFailure(ctx context.Context,
		providerID string) error
	IncrementSuccess(ctx context.Context,
		providerID string) error
	CreateProviderMetrics(ctx context.Context,
		providerID string) error
	GetProviderMetrics(ctx context.Context,
		providerID string) (*db.ProviderMetric, error)
	RemoveSupportedChannel(ctx context.Context,
		providerID string, channel string) error
}
