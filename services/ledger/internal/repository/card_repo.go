package repository

import (
	"context"
	"errors"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepo struct{ pool *pgxpool.Pool }

func NewCardRepo(pool *pgxpool.Pool) *CardRepo { return &CardRepo{pool: pool} }

const cardCols = `c.id, c."accountId", c.type, c.status, c.brand, c.last4, c."expMonth", c."expYear", c."monthlyLimitMinor", c."spentThisMonthMinor", c."createdAt"`

func scanCard(row pgx.Row) (*domain.Card, error) {
	var c domain.Card
	err := row.Scan(&c.ID, &c.AccountID, &c.Type, &c.Status, &c.Brand, &c.Last4, &c.ExpMonth, &c.ExpYear, &c.MonthlyLimitMinor, &c.SpentMinor, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCardNotFound
	}
	return &c, err
}

// GetOwned loads a card only if its account belongs to ownerID.
func (r *CardRepo) GetOwned(ctx context.Context, ownerID, id string) (*domain.Card, error) {
	return scanCard(r.pool.QueryRow(ctx, `
		SELECT `+cardCols+` FROM cards c JOIN accounts a ON a.id = c."accountId"
		WHERE c.id = $1 AND a."ownerId" = $2`, id, ownerID))
}

func (r *CardRepo) ListByOwner(ctx context.Context, ownerID string) ([]domain.Card, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+cardCols+` FROM cards c JOIN accounts a ON a.id = c."accountId"
		WHERE a."ownerId" = $1 ORDER BY c."createdAt"`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Card
	for rows.Next() {
		var c domain.Card
		if err := rows.Scan(&c.ID, &c.AccountID, &c.Type, &c.Status, &c.Brand, &c.Last4, &c.ExpMonth, &c.ExpYear, &c.MonthlyLimitMinor, &c.SpentMinor, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CardRepo) SetStatus(ctx context.Context, id string, status domain.CardStatus) error {
	tag, err := r.pool.Exec(ctx, `UPDATE cards SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCardNotFound
	}
	return nil
}

// RecordAuthorization stores the decision and, when approved, bumps month-to-date spend.
func (r *CardRepo) RecordAuthorization(ctx context.Context, cardID, merchant, mcc, currency string, amount int64, declineCode string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status := "APPROVED"
	if declineCode != "" {
		status = "DECLINED"
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO card_authorizations (id, "cardId", merchant, mcc, "amountMinor", currency, status, "declineCode")
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6, NULLIF($7, ''))`,
		cardID, merchant, mcc, amount, currency, status, declineCode); err != nil {
		return err
	}
	if declineCode == "" {
		if _, err := tx.Exec(ctx, `UPDATE cards SET "spentThisMonthMinor" = "spentThisMonthMinor" + $2 WHERE id = $1`, cardID, amount); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
