package service

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

type StatementStore interface {
	EntriesBetween(ctx context.Context, accountID string, start, end time.Time) ([]domain.Entry, error)
	Save(ctx context.Context, s domain.Statement) (*domain.Statement, error)
}

type StatementService struct {
	accounts AccountReader
	store    StatementStore
	now      func() time.Time
}

func NewStatementService(a AccountReader, s StatementStore) *StatementService {
	return &StatementService{accounts: a, store: s, now: time.Now}
}

// Monthly builds and stores the statement for a calendar month (UTC). Months
// that have not yet ended are rejected so closed statements stay immutable.
func (s *StatementService) Monthly(ctx context.Context, userID, accountID string, year int, month time.Month) (*domain.Statement, error) {
	acct, err := s.accounts.Get(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if acct.OwnerID != userID {
		return nil, domain.ErrAccountNotFound
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	now := s.now()
	if end.After(now) {
		return nil, domain.ErrInvalidPeriod
	}
	entries, err := s.store.EntriesBetween(ctx, accountID, start, end)
	if err != nil {
		return nil, err
	}
	st, err := domain.BuildStatement(accountID, start, end, entries, now)
	if err != nil {
		return nil, err
	}
	if st.EntryCount == 0 {
		st.OpeningMinor, st.ClosingMinor = acct.BalanceMinor, acct.BalanceMinor
	}
	return s.store.Save(ctx, st)
}
