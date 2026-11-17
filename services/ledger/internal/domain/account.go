package domain

import (
	"errors"
	"time"
)

type AccountStatus string

const (
	AccountPending AccountStatus = "PENDING"
	AccountActive  AccountStatus = "ACTIVE"
	AccountFrozen  AccountStatus = "FROZEN"
	AccountClosed  AccountStatus = "CLOSED"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountNotActive  = errors.New("account is not active")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Account struct {
	ID           string        `json:"id"`
	OwnerID      string        `json:"ownerId"`
	Name         string        `json:"name"`
	Currency     string        `json:"currency"`
	Status       AccountStatus `json:"status"`
	BalanceMinor int64         `json:"balanceMinor"`
	HoldMinor    int64         `json:"holdMinor"`
	CreatedAt    time.Time     `json:"createdAt"`
}

func (a Account) AvailableMinor() int64 { return a.BalanceMinor - a.HoldMinor }

// CanDebit reports whether amount may be debited from the account.
func (a Account) CanDebit(amount int64) error {
	if a.Status != AccountActive {
		return ErrAccountNotActive
	}
	if a.AvailableMinor() < amount {
		return ErrInsufficientFunds
	}
	return nil
}
