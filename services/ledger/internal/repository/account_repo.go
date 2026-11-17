package repository

import (
	"context"
	"errors"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct{ pool *pgxpool.Pool }

func NewAccountRepo(pool *pgxpool.Pool) *AccountRepo { return &AccountRepo{pool: pool} }

const accountCols = `id, owner_id, name, currency, status, balance_minor, hold_minor, created_at`

func scanAccount(row pgx.Row) (*domain.Account, error) {
	var a domain.Account
	err := row.Scan(&a.ID, &a.OwnerID, &a.Name, &a.Currency, &a.Status, &a.BalanceMinor, &a.HoldMinor, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	return &a, err
}

func (r *AccountRepo) Get(ctx context.Context, id string) (*domain.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx, `SELECT `+accountCols+` FROM accounts WHERE id = $1`, id))
}

// GetForUpdate locks the row inside tx to serialise concurrent balance changes.
func (r *AccountRepo) GetForUpdate(ctx context.Context, tx pgx.Tx, id string) (*domain.Account, error) {
	return scanAccount(tx.QueryRow(ctx, `SELECT `+accountCols+` FROM accounts WHERE id = $1 FOR UPDATE`, id))
}

func (r *AccountRepo) ListByOwner(ctx context.Context, ownerID string) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+accountCols+` FROM accounts WHERE owner_id = $1 ORDER BY created_at`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.OwnerID, &a.Name, &a.Currency, &a.Status, &a.BalanceMinor, &a.HoldMinor, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AccountRepo) AdjustBalance(ctx context.Context, tx pgx.Tx, id string, delta int64) error {
	_, err := tx.Exec(ctx, `UPDATE accounts SET balance_minor = balance_minor + $2, updated_at = now() WHERE id = $1`, id, delta)
	return err
}
