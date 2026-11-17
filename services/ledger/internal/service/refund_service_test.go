package service

import (
	"context"
	"errors"
	"testing"

	"github.com/finflow/ledger/internal/domain"
)

type fakeRefundStore struct {
	transfer *domain.Transfer
	existing []domain.Refund
	inserted []domain.Refund
}

func (f *fakeRefundStore) TransferByID(_ context.Context, id string) (*domain.Transfer, error) {
	if f.transfer == nil || f.transfer.ID != id {
		return nil, domain.ErrTransferNotFound
	}
	return f.transfer, nil
}
func (f *fakeRefundStore) ListByTransfer(context.Context, string) ([]domain.Refund, error) {
	return f.existing, nil
}
func (f *fakeRefundStore) Insert(_ context.Context, r domain.Refund) (*domain.Refund, error) {
	r.ID = "ref-1"
	f.inserted = append(f.inserted, r)
	return &r, nil
}

type fakeAccounts struct{ owner string }

func (f fakeAccounts) Get(_ context.Context, id string) (*domain.Account, error) {
	return &domain.Account{ID: id, OwnerID: f.owner}, nil
}

type captureAudit struct{ events []domain.AuditEvent }

func (c *captureAudit) Record(_ context.Context, e domain.AuditEvent) { c.events = append(c.events, e) }

func TestRefundService(t *testing.T) {
	completed := &domain.Transfer{ID: "t1", SourceID: "acc1", AmountMinor: 10_000, Currency: "USD", Status: domain.TransferCompleted}
	tests := []struct {
		name       string
		user       string
		transferID string
		amount     int64
		existing   []domain.Refund
		wantErr    error
		wantAmount int64
	}{
		{"partial refund", "u1", "t1", 2_500, nil, nil, 2_500},
		{"zero means remaining", "u1", "t1", 0, []domain.Refund{{AmountMinor: 3_000, Status: domain.RefundSucceeded}}, nil, 7_000},
		{"exceeds", "u1", "t1", 10_001, nil, domain.ErrRefundExceedsTransfer, 0},
		{"fully refunded already", "u1", "t1", 0, []domain.Refund{{AmountMinor: 10_000, Status: domain.RefundSucceeded}}, domain.ErrInvalidAmount, 0},
		{"not the owner", "intruder", "t1", 100, nil, domain.ErrTransferNotFound, 0},
		{"unknown transfer", "u1", "nope", 100, nil, domain.ErrTransferNotFound, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeRefundStore{transfer: completed, existing: tc.existing}
			aud := &captureAudit{}
			svc := NewRefundService(store, fakeAccounts{owner: "u1"}, aud)
			ref, err := svc.Refund(context.Background(), tc.user, tc.transferID, tc.amount, "requested")
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				if len(store.inserted) != 0 || len(aud.events) != 0 {
					t.Fatal("failed refund must not persist or audit")
				}
				return
			}
			if ref.AmountMinor != tc.wantAmount || ref.Status != domain.RefundPending || ref.Currency != "USD" {
				t.Fatalf("unexpected refund %+v", ref)
			}
			if len(aud.events) != 1 || aud.events[0].Action != "refund.created" {
				t.Fatalf("audit events = %+v", aud.events)
			}
		})
	}
}
