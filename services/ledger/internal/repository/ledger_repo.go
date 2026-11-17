package repository

import (
	"context"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerRepo struct{ pool *pgxpool.Pool }

func NewLedgerRepo(pool *pgxpool.Pool) *LedgerRepo { return &LedgerRepo{pool: pool} }

func (r *LedgerRepo) Pool() *pgxpool.Pool { return r.pool }

// InsertTransaction writes the header and all entries inside tx.
func (r *LedgerRepo) InsertTransaction(ctx context.Context, tx pgx.Tx, t domain.Transaction) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO ledger_transactions (id, reference, description) VALUES ($1, $2, $3)`,
		t.ID, t.Reference, t.Description); err != nil {
		return err
	}
	for _, e := range t.Entries {
		if _, err := tx.Exec(ctx,
			`INSERT INTO ledger_entries (id, transaction_id, account_id, direction, amount_minor, currency, balance_after_minor)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			e.ID, t.ID, e.AccountID, e.Direction, e.AmountMinor, e.Currency, e.BalanceAfter); err != nil {
			return err
		}
	}
	return nil
}

func (r *LedgerRepo) EntriesForAccount(ctx context.Context, accountID string, limit int) ([]domain.Entry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, transaction_id, account_id, direction, amount_minor, currency, balance_after_minor, posted_at
		 FROM ledger_entries WHERE account_id = $1 ORDER BY posted_at DESC LIMIT $2`, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Entry
	for rows.Next() {
		var e domain.Entry
		if err := rows.Scan(&e.ID, &e.TxID, &e.AccountID, &e.Direction, &e.AmountMinor, &e.Currency, &e.BalanceAfter, &e.PostedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
