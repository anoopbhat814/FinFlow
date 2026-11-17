package repository

import (
	"context"
	"time"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatementRepo struct{ pool *pgxpool.Pool }

func NewStatementRepo(pool *pgxpool.Pool) *StatementRepo { return &StatementRepo{pool: pool} }

// EntriesBetween returns an account's entries posted in [start, end), oldest first.
func (r *StatementRepo) EntriesBetween(ctx context.Context, accountID string, start, end time.Time) ([]domain.Entry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, "transactionId", "accountId", direction, "amountMinor", currency, "balanceAfterMinor", "postedAt"
		FROM ledger_entries
		WHERE "accountId" = $1 AND "postedAt" >= $2 AND "postedAt" < $3
		ORDER BY "postedAt", id`, accountID, start, end)
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

// Save upserts a generated statement; regenerating a period replaces its totals.
func (r *StatementRepo) Save(ctx context.Context, s domain.Statement) (*domain.Statement, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO statements (id, "accountId", "periodStart", "periodEnd", "openingBalanceMinor", "closingBalanceMinor", "totalDebitsMinor", "totalCreditsMinor", status)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6, $7, 'CLOSED')
		ON CONFLICT ("accountId", "periodStart", "periodEnd") DO UPDATE SET
			"openingBalanceMinor" = EXCLUDED."openingBalanceMinor",
			"closingBalanceMinor" = EXCLUDED."closingBalanceMinor",
			"totalDebitsMinor" = EXCLUDED."totalDebitsMinor",
			"totalCreditsMinor" = EXCLUDED."totalCreditsMinor",
			"generatedAt" = now()
		RETURNING id, "generatedAt"`,
		s.AccountID, s.PeriodStart, s.PeriodEnd, s.OpeningMinor, s.ClosingMinor, s.TotalDebits, s.TotalCredits).
		Scan(&s.ID, &s.GeneratedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
