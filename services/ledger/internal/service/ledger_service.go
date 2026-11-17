package service

import (
	"context"

	"github.com/finflow/ledger/internal/domain"
	"github.com/finflow/ledger/internal/repository"
)

type LedgerService struct {
	accounts *repository.AccountRepo
	ledger   *repository.LedgerRepo
}

func NewLedgerService(a *repository.AccountRepo, l *repository.LedgerRepo) *LedgerService {
	return &LedgerService{accounts: a, ledger: l}
}

// Statement returns recent entries, verifying the caller owns the account.
func (s *LedgerService) Statement(ctx context.Context, ownerID, accountID string, limit int) ([]domain.Entry, error) {
	acct, err := s.accounts.Get(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if acct.OwnerID != ownerID {
		return nil, domain.ErrAccountNotFound
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.ledger.EntriesForAccount(ctx, accountID, limit)
}
