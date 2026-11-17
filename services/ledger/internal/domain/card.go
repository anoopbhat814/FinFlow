package domain

import (
	"errors"
	"time"
)

type CardStatus string

const (
	CardInactive  CardStatus = "INACTIVE"
	CardActive    CardStatus = "ACTIVE"
	CardFrozen    CardStatus = "FROZEN"
	CardCancelled CardStatus = "CANCELLED"
)

var (
	ErrCardNotFound  = errors.New("card not found")
	ErrCardNotActive = errors.New("card is not active")
	ErrCardExpired   = errors.New("card has expired")
	ErrLimitExceeded = errors.New("monthly card limit exceeded")
)

type Card struct {
	ID                string     `json:"id"`
	AccountID         string     `json:"accountId"`
	Type              string     `json:"type"`
	Status            CardStatus `json:"status"`
	Brand             string     `json:"brand"`
	Last4             string     `json:"last4"`
	ExpMonth          int        `json:"expMonth"`
	ExpYear           int        `json:"expYear"`
	MonthlyLimitMinor int64      `json:"monthlyLimitMinor"`
	SpentMinor        int64      `json:"spentThisMonthMinor"`
	CreatedAt         time.Time  `json:"createdAt"`
}

// Expired is true once the last day of the expiry month has passed.
func (c Card) Expired(now time.Time) bool {
	end := time.Date(c.ExpYear, time.Month(c.ExpMonth)+1, 1, 0, 0, 0, 0, time.UTC)
	return !now.Before(end)
}

// Authorize decides whether a purchase of amount may be approved.
func (c Card) Authorize(amount int64, now time.Time) error {
	switch {
	case amount <= 0:
		return ErrInvalidAmount
	case c.Status != CardActive:
		return ErrCardNotActive
	case c.Expired(now):
		return ErrCardExpired
	case c.SpentMinor+amount > c.MonthlyLimitMinor:
		return ErrLimitExceeded
	}
	return nil
}

// Remaining is the unspent monthly allowance, never negative.
func (c Card) Remaining() int64 {
	if r := c.MonthlyLimitMinor - c.SpentMinor; r > 0 {
		return r
	}
	return 0
}
