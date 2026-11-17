package service

import (
	"errors"
	"testing"

	"github.com/finflow/ledger/internal/domain"
)

func TestCanDebit(t *testing.T) {
	acct := domain.Account{Status: domain.AccountActive, BalanceMinor: 10_000, HoldMinor: 2_500}
	cases := []struct {
		name   string
		amount int64
		want   error
	}{
		{"within available", 7_500, nil},
		{"exceeds available", 7_501, domain.ErrInsufficientFunds},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := acct.CanDebit(c.amount); !errors.Is(got, c.want) {
				t.Fatalf("got %v want %v", got, c.want)
			}
		})
	}
}

func TestTransactionValidate(t *testing.T) {
	ok := domain.Transaction{Entries: []domain.Entry{
		{Direction: domain.Debit, AmountMinor: 500, Currency: "USD"},
		{Direction: domain.Credit, AmountMinor: 500, Currency: "USD"},
	}}
	if err := ok.Validate(); err != nil {
		t.Fatalf("balanced tx rejected: %v", err)
	}
	bad := ok
	bad.Entries = []domain.Entry{ok.Entries[0], {Direction: domain.Credit, AmountMinor: 400, Currency: "USD"}}
	if err := bad.Validate(); !errors.Is(err, domain.ErrUnbalanced) {
		t.Fatalf("expected ErrUnbalanced, got %v", err)
	}
}
