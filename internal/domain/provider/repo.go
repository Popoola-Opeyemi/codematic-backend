package provider

import (
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type providerRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &providerRepository{
		q: q,
		p: pool,
	}
}

func (r *providerRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

// Create a new provider
func (r *providerRepository) CreateProvider(ctx context.Context,
	arg CreateProviderParams) (*db.Provider, error) {
	configJSON, err := json.Marshal(arg.Config)
	if err != nil {
		return nil, err
	}

	p, err := r.q.CreateProvider(ctx, db.CreateProviderParams{
		Name:   arg.Name,
		Code:   arg.Code,
		Config: configJSON,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Get provider by ID
func (r *providerRepository) GetByID(ctx context.Context, id string) (*db.Provider, error) {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, err
	}
	p, err := r.q.GetProviderByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Get provider by code
func (r *providerRepository) GetByCode(ctx context.Context,
	code string) (*db.Provider, error) {
	p, err := r.q.GetProviderByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List all active providers
func (r *providerRepository) ListActiveProviders(ctx context.Context) ([]db.Provider, error) {
	return r.q.ListActiveProviders(ctx)
}

// Update provider config
func (r *providerRepository) Update(ctx context.Context, arg db.UpdateProviderConfigParams) (*db.Provider, error) {
	p, err := r.q.UpdateProviderConfig(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Update config using map[string]interface{}
func (r *providerRepository) UpdateConfig(ctx context.Context,
	id string, config map[string]interface{}) (*db.Provider, error) {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p, err := r.q.UpdateProviderConfig(ctx, db.UpdateProviderConfigParams{
		ID:     uid,
		Config: configJSON,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Deactivate a provider
func (r *providerRepository) Deactivate(ctx context.Context, id string) error {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return err
	}
	return r.q.DeactivateProvider(ctx, uid)
}

// Add supported currency
func (r *providerRepository) AddSupportedCurrency(ctx context.Context,
	providerID string, currency string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.AddSupportedCurrency(ctx, db.AddSupportedCurrencyParams{
		ProviderID:   uid,
		CurrencyCode: currency,
	})
}

// Remove supported currency
func (r *providerRepository) RemoveSupportedCurrency(ctx context.Context,
	providerID string, currency string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.RemoveSupportedCurrency(ctx, db.RemoveSupportedCurrencyParams{
		ProviderID:   uid,
		CurrencyCode: currency,
	})
}

// Delete all supported currencies
func (r *providerRepository) DeleteAllSupportedCurrencies(ctx context.Context,
	providerID string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.DeleteAllSupportedCurrencies(ctx, uid)
}

// List supported currencies for a provider
func (r *providerRepository) ListSupportedCurrencies(ctx context.Context,
	providerID string) ([]db.Currency, error) {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return nil, err
	}
	return r.q.ListSupportedCurrencies(ctx, uid)
}

// Add supported channel
func (r *providerRepository) AddSupportedChannel(ctx context.Context,
	providerID string, channel string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.AddSupportedChannel(ctx, db.AddSupportedChannelParams{
		ProviderID: uid,
		Channel:    channel,
	})
}

// List provider with currencies + channels (full details view)
func (r *providerRepository) ListProviderDetails(ctx context.Context) (
	[]db.ListProviderDetailsRow, error) {
	return r.q.ListProviderDetails(ctx)
}

// Remove supported channel
func (r *providerRepository) RemoveSupportedChannel(ctx context.Context,
	providerID string, channel string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	// Assuming you define this SQLC query
	return r.q.RemoveSupportedChannel(ctx, db.RemoveSupportedChannelParams{
		ProviderID: uid,
		Channel:    channel,
	})
}

// Get provider metrics
func (r *providerRepository) GetProviderMetrics(ctx context.Context,
	providerID string) (*db.ProviderMetric, error) {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return nil, err
	}
	m, err := r.q.GetProviderMetrics(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Create provider metrics
func (r *providerRepository) CreateProviderMetrics(ctx context.Context,
	providerID string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.CreateProviderMetrics(ctx, uid)
}

// Increment success
func (r *providerRepository) IncrementSuccess(ctx context.Context,
	providerID string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.IncrementSuccess(ctx, uid)
}

// Increment failure
func (r *providerRepository) IncrementFailure(ctx context.Context,
	providerID string) error {
	uid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return err
	}
	return r.q.IncrementFailure(ctx, uid)
}

// Reset daily metrics
func (r *providerRepository) ResetDailyMetrics(ctx context.Context) error {
	return r.q.ResetDailyMetrics(ctx)
}

// Decay priority for all providers
func (r *providerRepository) DecayPriority(ctx context.Context) error {
	return r.q.DecayPriority(ctx)
}

// Select best provider
func (r *providerRepository) SelectBestProvider(ctx context.Context) (*db.SelectBestProviderRow, error) {
	p, err := r.q.SelectBestProvider(ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Select best provider by currency and channel
func (r *providerRepository) SelectBestProviderByCurrencyAndChannel(ctx context.Context,
	currency, channel string) (*db.SelectBestProviderByCurrencyAndChannelRow, error) {
	p, err := r.q.SelectBestProviderByCurrencyAndChannel(ctx, db.SelectBestProviderByCurrencyAndChannelParams{
		CurrencyCode: currency,
		Channel:      channel,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}
