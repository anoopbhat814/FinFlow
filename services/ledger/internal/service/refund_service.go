package service

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

// RefundStore is the persistence surface RefundService needs.
type RefundStore interface {
	TransferByID(ctx context.Context, id string) (*domain.Transfer, error)
	ListByTransfer(ctx context.Context, transferID string) ([]domain.Refund, error)
	Insert(ctx context.Context, f domain.Refund) (*domain.Refund, error)
}

// AccountReader checks ownership of the transfer's source account.
type AccountReader interface {
	Get(ctx context.Context, id string) (*domain.Account, error)
}

// Auditor records security-relevant actions.
type Auditor interface {
	Record(ctx context.Context, e domain.AuditEvent)
}

type RefundService struct {
	store    RefundStore
	accounts AccountReader
	audit    Auditor
	now      func() time.Time
}

func NewRefundService(s RefundStore, a AccountReader, au Auditor) *RefundService {
	return &RefundService{store: s, accounts: a, audit: au, now: time.Now}
}

// Refund creates a pending refund for part or all of a completed transfer.
// Only the owner of the source account may request one. A zero amount refunds
// whatever is still refundable.
func (s *RefundService) Refund(ctx context.Context, userID, transferID string, amount int64, reason string) (*domain.Refund, error) {
	t, err := s.store.TransferByID(ctx, transferID)
	if err != nil {
		return nil, err
	}
	src, err := s.accounts.Get(ctx, t.SourceID)
	if err != nil {
		return nil, err
	}
	if src.OwnerID != userID {
		return nil, domain.ErrTransferNotFound
	}
	existing, err := s.store.ListByTransfer(ctx, transferID)
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		amount = domain.RefundableRemaining(*t, existing)
	}
	if err := domain.ValidateRefund(*t, existing, amount); err != nil {
		return nil, err
	}
	ref, err := s.store.Insert(ctx, domain.Refund{
		TransferID:  transferID,
		AmountMinor: amount,
		Currency:    t.Currency,
		Reason:      reason,
		Status:      domain.RefundPending,
	})
	if err != nil {
		return nil, err
	}
	s.audit.Record(ctx, domain.AuditEvent{
		ActorID: userID, Action: "refund.created", EntityType: "transfer", EntityID: transferID,
		Metadata: map[string]string{"refundId": ref.ID}, CreatedAt: s.now(),
	})
	return ref, nil
}
