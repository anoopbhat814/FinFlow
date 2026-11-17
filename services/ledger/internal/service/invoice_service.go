package service

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

type InvoiceStore interface {
	Get(ctx context.Context, issuerID, id string) (*domain.Invoice, error)
	ListByIssuer(ctx context.Context, issuerID string, status string) ([]domain.Invoice, error)
	MarkPaid(ctx context.Context, inv *domain.Invoice) error
}

type InvoiceView struct {
	domain.Invoice
	TotalMinor int64                `json:"totalMinor"`
	Status     domain.InvoiceStatus `json:"status"`
}

type InvoiceService struct {
	store InvoiceStore
	now   func() time.Time
}

func NewInvoiceService(s InvoiceStore) *InvoiceService { return &InvoiceService{store: s, now: time.Now} }

func (s *InvoiceService) view(i domain.Invoice) InvoiceView {
	return InvoiceView{Invoice: i, TotalMinor: i.TotalMinor(), Status: i.EffectiveStatus(s.now())}
}

func (s *InvoiceService) List(ctx context.Context, issuerID, status string) ([]InvoiceView, error) {
	invs, err := s.store.ListByIssuer(ctx, issuerID, status)
	if err != nil {
		return nil, err
	}
	out := make([]InvoiceView, 0, len(invs))
	for _, i := range invs {
		out = append(out, s.view(i))
	}
	return out, nil
}

func (s *InvoiceService) Get(ctx context.Context, issuerID, id string) (*InvoiceView, error) {
	inv, err := s.store.Get(ctx, issuerID, id)
	if err != nil {
		return nil, err
	}
	v := s.view(*inv)
	return &v, nil
}

func (s *InvoiceService) MarkPaid(ctx context.Context, issuerID, id string) (*InvoiceView, error) {
	inv, err := s.store.Get(ctx, issuerID, id)
	if err != nil {
		return nil, err
	}
	if err := inv.MarkPaid(s.now()); err != nil {
		return nil, err
	}
	if err := s.store.MarkPaid(ctx, inv); err != nil {
		return nil, err
	}
	v := s.view(*inv)
	return &v, nil
}
