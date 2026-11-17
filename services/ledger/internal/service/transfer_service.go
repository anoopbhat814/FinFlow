package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/finflow/ledger/internal/domain"
	"github.com/finflow/ledger/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// FraudChecker scores a prospective transfer.
type FraudChecker interface {
	Check(ctx context.Context, t domain.Transfer) (domain.FraudVerdict, error)
}

type TransferService struct {
	accounts *repository.AccountRepo
	ledger   *repository.LedgerRepo
	fraud    FraudChecker
	now      func() time.Time
}

func NewTransferService(a *repository.AccountRepo, l *repository.LedgerRepo, f FraudChecker) *TransferService {
	return &TransferService{accounts: a, ledger: l, fraud: f, now: time.Now}
}

// Execute validates, fraud-checks and posts a transfer atomically.
func (s *TransferService) Execute(ctx context.Context, t domain.Transfer) (*domain.Transfer, error) {
	t.ID, t.Status, t.CreatedAt = uuid.NewString(), domain.TransferPending, s.now()

	verdict, err := s.fraud.Check(ctx, t)
	if err != nil {
		// Fail open into manual review rather than blocking all payments.
		verdict = domain.FraudVerdict{Score: 0.5, Action: "REVIEW"}
	}
	t.FraudScore = verdict.Score
	switch verdict.Action {
	case "BLOCK":
		t.Status = domain.TransferFailed
		return &t, fmt.Errorf("transfer blocked by fraud policy")
	case "REVIEW":
		t.Status = domain.TransferInReview
		return &t, nil
	}

	tx, err := s.ledger.Pool().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Lock in deterministic order to avoid deadlocks.
	ids := []string{t.SourceID, t.DestinationID}
	sort.Strings(ids)
	locked := map[string]*domain.Account{}
	for _, id := range ids {
		a, err := s.accounts.GetForUpdate(ctx, tx, id)
		if err != nil {
			return nil, err
		}
		locked[id] = a
	}
	src, dst := locked[t.SourceID], locked[t.DestinationID]
	if src.Currency != t.Currency || dst.Currency != t.Currency {
		return nil, fmt.Errorf("currency mismatch")
	}
	if err := src.CanDebit(t.AmountMinor); err != nil {
		return nil, err
	}

	txn := domain.Transaction{
		ID: uuid.NewString(), Reference: t.ID, Description: t.Memo,
		Entries: []domain.Entry{
			{ID: uuid.NewString(), AccountID: src.ID, Direction: domain.Debit, AmountMinor: t.AmountMinor, Currency: t.Currency, BalanceAfter: src.BalanceMinor - t.AmountMinor},
			{ID: uuid.NewString(), AccountID: dst.ID, Direction: domain.Credit, AmountMinor: t.AmountMinor, Currency: t.Currency, BalanceAfter: dst.BalanceMinor + t.AmountMinor},
		},
	}
	if err := s.ledger.InsertTransaction(ctx, tx, txn); err != nil {
		return nil, err
	}
	if err := s.accounts.AdjustBalance(ctx, tx, src.ID, -t.AmountMinor); err != nil {
		return nil, err
	}
	if err := s.accounts.AdjustBalance(ctx, tx, dst.ID, t.AmountMinor); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	done := s.now()
	t.Status, t.CompletedAt = domain.TransferCompleted, &done
	return &t, nil
}
