package repository

import (
	"context"
	"encoding/json"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct{ pool *pgxpool.Pool }

func NewAuditRepo(pool *pgxpool.Pool) *AuditRepo { return &AuditRepo{pool: pool} }

func (r *AuditRepo) Append(ctx context.Context, e domain.AuditEvent) error {
	meta, err := json.Marshal(e.Metadata)
	if err != nil {
		return err
	}
	var actor *string
	if e.ActorID != "" {
		actor = &e.ActorID
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO audit_logs (id, "actorId", action, "entityType", "entityId", metadata, ip)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5::jsonb, NULLIF($6, ''))`,
		actor, e.Action, e.EntityType, e.EntityID, meta, e.IP)
	return err
}

func (r *AuditRepo) ForEntity(ctx context.Context, entityType, entityID string, limit int) ([]domain.AuditEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE("actorId", ''), action, "entityType", "entityId", COALESCE(ip, ''), "createdAt"
		FROM audit_logs WHERE "entityType" = $1 AND "entityId" = $2
		ORDER BY "createdAt" DESC LIMIT $3`, entityType, entityID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AuditEvent
	for rows.Next() {
		var e domain.AuditEvent
		if err := rows.Scan(&e.ID, &e.ActorID, &e.Action, &e.EntityType, &e.EntityID, &e.IP, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
