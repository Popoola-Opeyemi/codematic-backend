package webhook

import (
	"context"
	"time"

	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type webhookRepository struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) Repository {
	return &webhookRepository{
		q: q,
		p: pool,
	}
}

func (r *webhookRepository) WithTx(q *db.Queries) Repository {
	return NewRepository(q, r.p)
}

func (r *webhookRepository) Create(ctx context.Context, event *WebhookEvent) error {

	uid, err := utils.StringToPgUUID(event.ID)
	if err != nil {
		return err
	}
	pid, err := utils.StringToPgUUID(event.ProviderID)
	if err != nil {
		return err
	}
	tid, err := utils.StringToPgUUID(event.TenantID)
	if err != nil {
		return err
	}

	_, err = r.q.CreateWebhookEvent(ctx, db.CreateWebhookEventParams{
		ID:              uid,
		ProviderID:      pid,
		ProviderEventID: event.ProviderEventID,
		TenantID:        tid,
		EventType:       event.EventType,
		Payload:         event.Payload,
		Status:          event.Status,
		Attempts:        utils.ToPgxInt4(int32(event.Attempts)),
		LastError:       utils.ToDBString(&event.LastError),
		CreatedAt:       utils.ToPgTimestamptz(event.CreatedAt),
		UpdatedAt:       utils.ToPgTimestamptz(event.UpdatedAt),
		IsOutgoing:      utils.ToPgxBool(event.IsOutgoing),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *webhookRepository) GetByProviderAndEventID(ctx context.Context, providerID, providerEventID string) (*WebhookEvent, error) {

	pid, err := utils.StringToPgUUID(providerID)
	if err != nil {
		return nil, err
	}
	record, err := r.q.GetWebhookEventByProviderAndEventID(ctx, db.GetWebhookEventByProviderAndEventIDParams{
		ProviderID:      pid,
		ProviderEventID: providerEventID,
	})
	if err != nil {
		return nil, err
	}

	return &WebhookEvent{
		ID:              utils.FromPgUUID(record.ID),
		ProviderID:      utils.FromPgUUID(record.ProviderID),
		ProviderEventID: record.ProviderEventID,
		TenantID:        utils.FromPgUUID(record.TenantID),
		EventType:       record.EventType,
		Payload:         record.Payload,
		Status:          record.Status,
		Attempts:        int(record.Attempts.Int32),
		LastError:       utils.StringOrEmpty(utils.FromPgText(record.LastError)),
		CreatedAt:       utils.FromPgTimestamptz(record.CreatedAt),
		UpdatedAt:       utils.FromPgTimestamptz(record.UpdatedAt),
	}, nil
}

func (r *webhookRepository) UpdateStatus(ctx context.Context, id string, status string, attempts int, lastError *string) error {
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return err
	}
	return r.q.UpdateWebhookEventStatus(ctx, db.UpdateWebhookEventStatusParams{
		Status:    status,
		Attempts:  utils.ToPgxInt4(int32(attempts)),
		LastError: utils.ToDBString(lastError),
		UpdatedAt: utils.ToPgTimestamptz(time.Now()),
		ID:        uid,
	})
}

func (r *webhookRepository) GetByID(ctx context.Context, id string) (*WebhookEvent, error) {
	// TODO: Replace with a direct GetWebhookEventByID query in db.Queries
	uid, err := utils.StringToPgUUID(id)
	if err != nil {
		return nil, err
	}
	list, err := r.q.ListFailedWebhookEvents(ctx)
	if err != nil {
		return nil, err
	}
	for _, record := range list {
		if record.ID == uid {
			return &WebhookEvent{
				ID:              utils.FromPgUUID(record.ID),
				ProviderID:      utils.FromPgUUID(record.ProviderID),
				ProviderEventID: record.ProviderEventID,
				TenantID:        utils.FromPgUUID(record.TenantID),
				EventType:       record.EventType,
				Payload:         record.Payload,
				Status:          record.Status,
				Attempts:        int(record.Attempts.Int32),
				LastError:       utils.StringOrEmpty(utils.FromPgText(record.LastError)),
				CreatedAt:       utils.FromPgTimestamptz(record.CreatedAt),
				UpdatedAt:       utils.FromPgTimestamptz(record.UpdatedAt),
			}, nil
		}
	}
	return nil, err
}
