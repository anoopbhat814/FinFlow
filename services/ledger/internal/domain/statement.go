package domain

import (
	"errors"
	"time"
)

var ErrInvalidPeriod = errors.New("statement period is invalid")

type Statement struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"accountId"`
	PeriodStart    time.Time `json:"periodStart"`
	PeriodEnd      time.Time `json:"periodEnd"`
	OpeningMinor   int64     `json:"openingBalanceMinor"`
	ClosingMinor   int64     `json:"closingBalanceMinor"`
	TotalDebits    int64     `json:"totalDebitsMinor"`
	TotalCredits   int64     `json:"totalCreditsMinor"`
	EntryCount     int       `json:"entryCount"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

// BuildStatement summarises entries posted in [start, end). Entries must be
// ordered by PostedAt ascending; opening is derived from the first entry's
// balance-after so the statement always chains with the ledger.
func BuildStatement(accountID string, start, end time.Time, entries []Entry, now time.Time) (Statement, error) {
	if !start.Before(end) {
		return Statement{}, ErrInvalidPeriod
	}
	st := Statement{AccountID: accountID, PeriodStart: start, PeriodEnd: end, GeneratedAt: now}
	first := true
	for _, e := range entries {
		if e.PostedAt.Before(start) || !e.PostedAt.Before(end) {
			continue
		}
		delta := e.AmountMinor
		if e.Direction == Debit {
			st.TotalDebits += e.AmountMinor
			delta = -delta
		} else {
			st.TotalCredits += e.AmountMinor
		}
		if first {
			st.OpeningMinor = e.BalanceAfter - delta
			first = false
		}
		st.ClosingMinor = e.BalanceAfter
		st.EntryCount++
	}
	if first {
		return st, nil
	}
	return st, nil
}

// NetChange is closing minus opening.
func (s Statement) NetChange() int64 { return s.ClosingMinor - s.OpeningMinor }
