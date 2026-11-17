package domain

import (
	"testing"
	"time"
)

func TestValidateRefund(t *testing.T) {
	done := Transfer{ID: "t1", AmountMinor: 10_000, Status: TransferCompleted}
	tests := []struct {
		name     string
		transfer Transfer
		existing []Refund
		amount   int64
		want     error
	}{
		{"full refund", done, nil, 10_000, nil},
		{"partial then rest", done, []Refund{{AmountMinor: 4_000, Status: RefundSucceeded}}, 6_000, nil},
		{"failed refund frees budget", done, []Refund{{AmountMinor: 9_000, Status: RefundFailed}}, 10_000, nil},
		{"over remaining", done, []Refund{{AmountMinor: 4_000, Status: RefundPending}}, 6_001, ErrRefundExceedsTransfer},
		{"zero amount", done, nil, 0, ErrInvalidAmount},
		{"negative amount", done, nil, -5, ErrInvalidAmount},
		{"pending transfer", Transfer{AmountMinor: 100, Status: TransferPending}, nil, 50, ErrTransferNotRefundable},
		{"reversed transfer", Transfer{AmountMinor: 100, Status: TransferReversed}, nil, 50, ErrTransferNotRefundable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidateRefund(tc.transfer, tc.existing, tc.amount); got != tc.want {
				t.Fatalf("ValidateRefund() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDisputeAdvance(t *testing.T) {
	now := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name         string
		from, to     DisputeStatus
		wantErr      bool
		wantResolved bool
	}{
		{"open to review", DisputeOpen, DisputeUnderReview, false, false},
		{"open to evidence", DisputeOpen, DisputeEvidenceRequired, false, false},
		{"review to won", DisputeUnderReview, DisputeWon, false, true},
		{"review to lost", DisputeUnderReview, DisputeLost, false, true},
		{"open to withdrawn", DisputeOpen, DisputeWithdrawn, false, true},
		{"open straight to won", DisputeOpen, DisputeWon, true, false},
		{"won is final", DisputeWon, DisputeLost, true, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := Dispute{Status: tc.from}
			err := d.Advance(tc.to, now)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && (d.ResolvedAt != nil) != tc.wantResolved {
				t.Fatalf("ResolvedAt set = %v, want %v", d.ResolvedAt != nil, tc.wantResolved)
			}
			if err != nil && d.Status != tc.from {
				t.Fatalf("status changed on failed transition: %s", d.Status)
			}
		})
	}
}
