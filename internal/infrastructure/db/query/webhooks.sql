-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (
  id, provider_id, provider_event_id, tenant_id, event_type, payload, status, attempts, last_error, created_at, updated_at, is_outgoing
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetWebhookEventByProviderAndEventID :one
SELECT * FROM webhook_events WHERE provider_id = $1 AND provider_event_id = $2;

-- name: UpdateWebhookEventStatus :exec
UPDATE webhook_events SET status = $1, attempts = $2, last_error = $3, updated_at = $4 WHERE id = $5;

-- name: ListFailedWebhookEvents :many
SELECT * FROM webhook_events WHERE status = 'failed'; 