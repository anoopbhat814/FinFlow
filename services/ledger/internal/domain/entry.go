package domain

import (
	"errors"
	"time"
)

type Direction string

const (
	Debit  Direction = "DEBIT"
	Credit Direction = "CREDIT"
)

var ErrUnbalanced = errors.New("ledger transaction is unbalanced")

type Entry struct {
	ID           string    `json:"id"`
	TxID         string    `json:"transactionId"`
	AccountID    string    `json:"accountId"`
	Direction    Direction `json:"direction"`
	AmountMinor  int64     `json:"amountMinor"`
	Currency     string    `json:"currency"`
	BalanceAfter int64     `json:"balanceAfterMinor"`
	PostedAt     time.Time `json:"postedAt"`
}

type Transaction struct {
	ID          string  `json:"id"`
	Reference   string  `json:"reference"`
	Description string  `json:"description"`
	Entries     []Entry `json:"entries"`
}

// Validate enforces the double-entry invariant: debits == credits per currency.
func (t Transaction) Validate() error {
	if len(t.Entries) < 2 {
		return ErrUnbalanced
	}
	net := map[string]int64{}
	for _, e := range t.Entries {
		if e.AmountMinor <= 0 {
			return ErrUnbalanced
		}
		if e.Direction == Debit {
			net[e.Currency] += e.AmountMinor
		} else {
			net[e.Currency] -= e.AmountMinor
		}
	}
	for _, v := range net {
		if v != 0 {
			return ErrUnbalanced
		}
	}
	return nil
}
