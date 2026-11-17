package domain

import (
	"testing"
	"time"
)

func TestInvoiceTotalAndStatus(t *testing.T) {
	now := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	inv := Invoice{
		Status: InvoiceSent,
		DueAt:  now.AddDate(0, 0, 14),
		Lines: []InvoiceLine{
			{Description: "Design retainer", Quantity: 1, UnitPriceMinor: 120_000},
			{Description: "Extra revisions", Quantity: 3, UnitPriceMinor: 7_500},
		},
	}
	if got := inv.TotalMinor(); got != 142_500 {
		t.Fatalf("TotalMinor = %d", got)
	}
	tests := []struct {
		name   string
		status InvoiceStatus
		at     time.Time
		want   InvoiceStatus
	}{
		{"sent before due", InvoiceSent, now, InvoiceSent},
		{"sent after due", InvoiceSent, now.AddDate(0, 0, 15), InvoiceOverdue},
		{"paid stays paid", InvoicePaid, now.AddDate(1, 0, 0), InvoicePaid},
		{"draft never overdue", InvoiceDraft, now.AddDate(1, 0, 0), InvoiceDraft},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			i := inv
			i.Status = tc.status
			if got := i.EffectiveStatus(tc.at); got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}

func TestInvoiceMarkPaid(t *testing.T) {
	now := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		status  InvoiceStatus
		wantErr error
	}{
		{InvoiceSent, nil}, {InvoiceOverdue, nil},
		{InvoiceDraft, ErrInvoiceLocked}, {InvoicePaid, ErrInvoiceLocked}, {InvoiceVoid, ErrInvoiceLocked},
	}
	for _, tc := range tests {
		t.Run(string(tc.status), func(t *testing.T) {
			i := Invoice{Status: tc.status, DueAt: now.AddDate(0, 0, 5)}
			if err := i.MarkPaid(now); err != tc.wantErr {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && (i.Status != InvoicePaid || i.PaidAt == nil) {
				t.Fatalf("invoice not marked paid: %+v", i)
			}
		})
	}
}
