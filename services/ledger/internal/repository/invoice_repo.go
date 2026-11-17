package repository

import (
	"context"
	"errors"

	"github.com/finflow/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceRepo struct{ pool *pgxpool.Pool }

func NewInvoiceRepo(pool *pgxpool.Pool) *InvoiceRepo { return &InvoiceRepo{pool: pool} }

func (r *InvoiceRepo) lines(ctx context.Context, invoiceID string) ([]domain.InvoiceLine, error) {
	rows, err := r.pool.Query(ctx, `SELECT description, quantity, "unitPriceMinor" FROM invoice_lines WHERE "invoiceId" = $1`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.InvoiceLine
	for rows.Next() {
		var l domain.InvoiceLine
		if err := rows.Scan(&l.Description, &l.Quantity, &l.UnitPriceMinor); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *InvoiceRepo) Get(ctx context.Context, issuerID, id string) (*domain.Invoice, error) {
	var i domain.Invoice
	err := r.pool.QueryRow(ctx, `
		SELECT id, number, "issuerId", "customerName", "customerEmail", currency, status, "issuedAt", "dueAt", "paidAt"
		FROM invoices WHERE id = $1 AND "issuerId" = $2`, id, issuerID).
		Scan(&i.ID, &i.Number, &i.IssuerID, &i.CustomerName, &i.CustomerEmail, &i.Currency, &i.Status, &i.IssuedAt, &i.DueAt, &i.PaidAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	i.Lines, err = r.lines(ctx, i.ID)
	return &i, err
}

func (r *InvoiceRepo) ListByIssuer(ctx context.Context, issuerID string, status string) ([]domain.Invoice, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, number, "issuerId", "customerName", "customerEmail", currency, status, "issuedAt", "dueAt", "paidAt"
		FROM invoices WHERE "issuerId" = $1 AND ($2 = '' OR status::text = $2)
		ORDER BY "issuedAt" DESC LIMIT 100`, issuerID, status)
	if err != nil {
		return nil, err
	}
	var out []domain.Invoice
	for rows.Next() {
		var i domain.Invoice
		if err := rows.Scan(&i.ID, &i.Number, &i.IssuerID, &i.CustomerName, &i.CustomerEmail, &i.Currency, &i.Status, &i.IssuedAt, &i.DueAt, &i.PaidAt); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, i)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for k := range out {
		if out[k].Lines, err = r.lines(ctx, out[k].ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *InvoiceRepo) MarkPaid(ctx context.Context, inv *domain.Invoice) error {
	_, err := r.pool.Exec(ctx, `UPDATE invoices SET status = 'PAID', "paidAt" = $2 WHERE id = $1`, inv.ID, inv.PaidAt)
	return err
}
