package service

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

type CardStore interface {
	GetOwned(ctx context.Context, ownerID, id string) (*domain.Card, error)
	SetStatus(ctx context.Context, id string, status domain.CardStatus) error
	RecordAuthorization(ctx context.Context, cardID, merchant, mcc, currency string, amount int64, declineCode string) error
}

type CardService struct {
	store CardStore
	audit Auditor
	now   func() time.Time
}

func NewCardService(s CardStore, a Auditor) *CardService {
	return &CardService{store: s, audit: a, now: time.Now}
}

// declineCode maps a domain error to the short code stored with the authorization.
func declineCode(err error) string {
	switch err {
	case domain.ErrLimitExceeded:
		return "limit_exceeded"
	case domain.ErrCardExpired:
		return "expired_card"
	case domain.ErrCardNotActive:
		return "card_inactive"
	case domain.ErrInvalidAmount:
		return "invalid_amount"
	}
	return "declined"
}

// Authorize evaluates and records a purchase attempt. A decline is returned
// as the domain error after it has been stored.
func (s *CardService) Authorize(ctx context.Context, userID, cardID, merchant, mcc, currency string, amount int64) error {
	card, err := s.store.GetOwned(ctx, userID, cardID)
	if err != nil {
		return err
	}
	decision := card.Authorize(amount, s.now())
	code := ""
	if decision != nil {
		code = declineCode(decision)
	}
	if err := s.store.RecordAuthorization(ctx, cardID, merchant, mcc, currency, amount, code); err != nil {
		return err
	}
	return decision
}

// Freeze and Unfreeze toggle an active card; cancelled cards stay cancelled.
func (s *CardService) Freeze(ctx context.Context, userID, cardID string) error {
	return s.toggle(ctx, userID, cardID, domain.CardActive, domain.CardFrozen, "card.frozen")
}

func (s *CardService) Unfreeze(ctx context.Context, userID, cardID string) error {
	return s.toggle(ctx, userID, cardID, domain.CardFrozen, domain.CardActive, "card.unfrozen")
}

func (s *CardService) toggle(ctx context.Context, userID, cardID string, from, to domain.CardStatus, action string) error {
	card, err := s.store.GetOwned(ctx, userID, cardID)
	if err != nil {
		return err
	}
	if card.Status != from {
		return domain.ErrCardNotActive
	}
	if err := s.store.SetStatus(ctx, cardID, to); err != nil {
		return err
	}
	s.audit.Record(ctx, domain.AuditEvent{ActorID: userID, Action: action, EntityType: "card", EntityID: cardID, CreatedAt: s.now()})
	return nil
}
