package invoices

import (
	"context"
	"database/sql"
	"strings"

	"desafio/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, invoices *domain.Invoice) (int64, error)
	ReadAll(ctx context.Context) ([]*domain.Invoice, error)
	CreateMany(ctx context.Context, invoices []*domain.Invoice) error
	UpdateTotals(ctx context.Context) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, invoices *domain.Invoice) (int64, error) {
	query := `INSERT INTO invoices (customer_id, datetime, total) VALUES (?, ?, ?);`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &invoices.CustomerId, &invoices.Datetime, &invoices.Total)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) ReadAll(ctx context.Context) ([]*domain.Invoice, error) {
	query := `SELECT id, customer_id, datetime, total FROM invoices;`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]*domain.Invoice, 0)
	for rows.Next() {
		invoice := domain.Invoice{}
		err := rows.Scan(&invoice.Id, &invoice.CustomerId, &invoice.Datetime, &invoice.Total)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, nil
}

func (r *repository) CreateMany(ctx context.Context, invoices []*domain.Invoice) error {
	query := `INSERT INTO invoices (id, customer_id, datetime, total) VALUES`
	values := []any{}
	for _, invoice := range invoices {
		query += " (?, ?, ?, ?),"
		values = append(values, invoice.Id, invoice.CustomerId, invoice.Datetime, invoice.Total)
	}
	query = strings.TrimSuffix(query, ",")
	query += ";"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateTotals(ctx context.Context) error {
	query := `
				UPDATE invoices i
				JOIN sales s ON i.id = s.invoice_id
				JOIN products p ON s.product_id = p.id
				SET i.total = p.price * s.quantity
				WHERE i.total = 0;
			`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
