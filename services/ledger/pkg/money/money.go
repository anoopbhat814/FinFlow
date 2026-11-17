// Package money provides integer minor-unit arithmetic.
package money

import (
	"errors"
	"fmt"
)

var (
	ErrCurrencyMismatch = errors.New("money: currency mismatch")
	ErrOverflow         = errors.New("money: overflow")
)

// Amount is a quantity of minor units (e.g. cents) in a given currency.
type Amount struct {
	Minor    int64
	Currency string
}

func New(minor int64, currency string) Amount { return Amount{Minor: minor, Currency: currency} }

func (a Amount) Add(b Amount) (Amount, error) {
	if a.Currency != b.Currency {
		return Amount{}, ErrCurrencyMismatch
	}
	sum := a.Minor + b.Minor
	if (b.Minor > 0 && sum < a.Minor) || (b.Minor < 0 && sum > a.Minor) {
		return Amount{}, ErrOverflow
	}
	return Amount{sum, a.Currency}, nil
}

func (a Amount) Sub(b Amount) (Amount, error) {
	return a.Add(Amount{-b.Minor, b.Currency})
}

func (a Amount) IsNegative() bool { return a.Minor < 0 }
func (a Amount) IsZero() bool     { return a.Minor == 0 }

func (a Amount) String() string {
	sign := ""
	m := a.Minor
	if m < 0 {
		sign, m = "-", -m
	}
	return fmt.Sprintf("%s%d.%02d %s", sign, m/100, m%100, a.Currency)
}
