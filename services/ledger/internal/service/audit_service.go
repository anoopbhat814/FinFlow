package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

// AuditSink persists audit events.
type AuditSink interface {
	Append(ctx context.Context, e domain.AuditEvent) error
}

// AuditService writes audit events best-effort: a failed write is logged but
// never fails the business operation that triggered it.
type AuditService struct {
	sink AuditSink
	log  *slog.Logger
}

func NewAuditService(sink AuditSink, log *slog.Logger) *AuditService {
	return &AuditService{sink: sink, log: log}
}

func (a *AuditService) Record(ctx context.Context, e domain.AuditEvent) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if err := a.sink.Append(ctx, e); err != nil {
		a.log.Error("audit append failed", "action", e.Action, "entity", e.EntityID, "err", err)
	}
}
