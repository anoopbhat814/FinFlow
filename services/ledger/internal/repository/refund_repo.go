package repository

import (
	"context"
	"errors"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefundRepo struct{ pool *pgxpool.Pool }

func NewRefundRepo(pool *pgxpool.Pool) *RefundRepo { return &RefundRepo{pool: pool} }

func (r *RefundRepo) Insert(ctx context.Context, f domain.Refund) (*domain.Refund, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO refunds (id, "transferId", "amountMinor", currency, reason, status)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5)
		RETURNING id, "createdAt"`,
		f.TransferID, f.AmountMinor, f.Currency, f.Reason, f.Status)
	if err := row.Scan(&f.ID, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *RefundRepo) ListByTransfer(ctx context.Context, transferID string) ([]domain.Refund, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, "transferId", "amountMinor", currency, reason, status, "createdAt", "settledAt"
		FROM refunds WHERE "transferId" = $1 ORDER BY "createdAt"`, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Refund
	for rows.Next() {
		var f domain.Refund
		if err := rows.Scan(&f.ID, &f.TransferID, &f.AmountMinor, &f.Currency, &f.Reason, &f.Status, &f.CreatedAt, &f.SettledAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// MarkSettled records the final outcome of a refund inside tx.
func (r *RefundRepo) MarkSettled(ctx context.Context, tx pgx.Tx, id string, status domain.RefundStatus) error {
	tag, err := tx.Exec(ctx, `UPDATE refunds SET status = $2, "settledAt" = now() WHERE id = $1 AND status = 'PENDING'`, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("refund not pending")
	}
	return nil
}

func (r *RefundRepo) TransferByID(ctx context.Context, id string) (*domain.Transfer, error) {
	var t domain.Transfer
	err := r.pool.QueryRow(ctx, `
		SELECT id, "idempotencyKey", "sourceAccountId", "destinationAccountId", "amountMinor", currency, status, "createdAt", "completedAt"
		FROM transfers WHERE id = $1`, id).
		Scan(&t.ID, &t.IdempotencyKey, &t.SourceID, &t.DestinationID, &t.AmountMinor, &t.Currency, &t.Status, &t.CreatedAt, &t.CompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTransferNotFound
	}
	return &t, err
}
