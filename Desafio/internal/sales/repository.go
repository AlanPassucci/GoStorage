package sales

import (
	"context"
	"database/sql"
	"strings"

	"desafio/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, sales *domain.Sale) (int64, error)
	ReadAll(ctx context.Context) ([]*domain.Sale, error)
	CreateMany(ctx context.Context, sales []*domain.Sale) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, sales *domain.Sale) (int64, error) {
	query := `INSERT INTO sales (product_id, invoice_id, quantity) VALUES (?, ?, ?);`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &sales.ProductId, &sales.InvoicesId, &sales.Quantity)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) ReadAll(ctx context.Context) ([]*domain.Sale, error) {
	query := `SELECT id, product_id, invoice_id, quantity FROM sales`

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

	sales := make([]*domain.Sale, 0)
	for rows.Next() {
		sale := domain.Sale{}
		err := rows.Scan(&sale.Id, &sale.ProductId, &sale.InvoicesId, &sale.Quantity)
		if err != nil {
			return nil, err
		}
		sales = append(sales, &sale)
	}
	return sales, nil
}

func (r *repository) CreateMany(ctx context.Context, sales []*domain.Sale) error {
	query := `INSERT INTO sales (id, product_id, invoice_id, quantity) VALUES`
	values := []any{}
	for _, sale := range sales {
		query += " (?, ?, ?, ?),"
		values = append(values, sale.Id, sale.ProductId, sale.InvoicesId, sale.Quantity)
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
