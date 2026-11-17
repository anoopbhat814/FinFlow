package domain

import (
	"errors"
	"time"
)

type InvoiceStatus string

const (
	InvoiceDraft   InvoiceStatus = "DRAFT"
	InvoiceSent    InvoiceStatus = "SENT"
	InvoicePaid    InvoiceStatus = "PAID"
	InvoiceOverdue InvoiceStatus = "OVERDUE"
	InvoiceVoid    InvoiceStatus = "VOID"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrInvoiceLocked   = errors.New("invoice can no longer be changed")
)

type InvoiceLine struct {
	Description    string `json:"description"`
	Quantity       int    `json:"quantity"`
	UnitPriceMinor int64  `json:"unitPriceMinor"`
}

type Invoice struct {
	ID            string        `json:"id"`
	Number        string        `json:"number"`
	IssuerID      string        `json:"issuerId"`
	CustomerName  string        `json:"customerName"`
	CustomerEmail string        `json:"customerEmail"`
	Currency      string        `json:"currency"`
	Status        InvoiceStatus `json:"status"`
	IssuedAt      time.Time     `json:"issuedAt"`
	DueAt         time.Time     `json:"dueAt"`
	PaidAt        *time.Time    `json:"paidAt,omitempty"`
	Lines         []InvoiceLine `json:"lines"`
}

// TotalMinor sums quantity * unit price over all lines.
func (i Invoice) TotalMinor() int64 {
	var total int64
	for _, l := range i.Lines {
		total += int64(l.Quantity) * l.UnitPriceMinor
	}
	return total
}

// EffectiveStatus derives OVERDUE for unpaid sent invoices past their due date.
func (i Invoice) EffectiveStatus(now time.Time) InvoiceStatus {
	if i.Status == InvoiceSent && now.After(i.DueAt) {
		return InvoiceOverdue
	}
	return i.Status
}

// MarkPaid transitions a sent or overdue invoice to PAID.
func (i *Invoice) MarkPaid(now time.Time) error {
	switch i.EffectiveStatus(now) {
	case InvoiceSent, InvoiceOverdue:
		i.Status = InvoicePaid
		i.PaidAt = &now
		return nil
	}
	return ErrInvoiceLocked
}
