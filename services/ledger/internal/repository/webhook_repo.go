package repository

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookRepo struct{ pool *pgxpool.Pool }

func NewWebhookRepo(pool *pgxpool.Pool) *WebhookRepo { return &WebhookRepo{pool: pool} }

func (r *WebhookRepo) CreateEndpoint(ctx context.Context, e domain.WebhookEndpoint) (*domain.WebhookEndpoint, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO webhook_endpoints (id, "userId", url, secret, events)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4)
		RETURNING id, active, "createdAt"`, e.UserID, e.URL, e.Secret, e.Events).
		Scan(&e.ID, &e.Active, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *WebhookRepo) ListEndpoints(ctx context.Context, userID string) ([]domain.WebhookEndpoint, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, "userId", url, secret, events, active, "createdAt"
		FROM webhook_endpoints WHERE "userId" = $1 ORDER BY "createdAt"`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.WebhookEndpoint
	for rows.Next() {
		var e domain.WebhookEndpoint
		if err := rows.Scan(&e.ID, &e.UserID, &e.URL, &e.Secret, &e.Events, &e.Active, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *WebhookRepo) EnqueueDelivery(ctx context.Context, endpointID, event string, payload []byte) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO webhook_deliveries (id, "endpointId", event, payload, "nextAttemptAt")
		VALUES (gen_random_uuid()::text, $1, $2, $3::jsonb, now())`, endpointID, event, payload)
	return err
}

// DueDeliveries claims up to limit pending deliveries whose retry time has passed.
func (r *WebhookRepo) DueDeliveries(ctx context.Context, now time.Time, limit int) ([]domain.WebhookDelivery, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, "endpointId", event, payload, status, attempts, "nextAttemptAt"
		FROM webhook_deliveries
		WHERE status = 'PENDING' AND "nextAttemptAt" <= $1
		ORDER BY "nextAttemptAt" LIMIT $2 FOR UPDATE SKIP LOCKED`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.WebhookDelivery
	for rows.Next() {
		var d domain.WebhookDelivery
		if err := rows.Scan(&d.ID, &d.EndpointID, &d.Event, &d.Payload, &d.Status, &d.Attempts, &d.NextAttemptAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *WebhookRepo) RecordAttempt(ctx context.Context, id, status string, code int, next *time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE webhook_deliveries
		SET status = $2, attempts = attempts + 1, "responseCode" = $3, "nextAttemptAt" = $4
		WHERE id = $1`, id, status, code, next)
	return err
}

// GetEndpoint loads one endpoint including its signing secret.
func (r *WebhookRepo) GetEndpoint(ctx context.Context, id string) (*domain.WebhookEndpoint, error) {
	var e domain.WebhookEndpoint
	err := r.pool.QueryRow(ctx, `
		SELECT id, "userId", url, secret, events, active, "createdAt" FROM webhook_endpoints WHERE id = $1`, id).
		Scan(&e.ID, &e.UserID, &e.URL, &e.Secret, &e.Events, &e.Active, &e.CreatedAt)
	return &e, err
}
