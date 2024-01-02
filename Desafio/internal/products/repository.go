package products

import (
	"context"
	"database/sql"
	"strings"

	"desafio/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, product *domain.Product) (int64, error)
	ReadAll(ctx context.Context) ([]*domain.Product, error)
	CreateMany(ctx context.Context, products []*domain.Product) error
	GetQtySaledGroupedByDescription(ctx context.Context) ([]map[string]any, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, product *domain.Product) (int64, error) {
	query := `INSERT INTO products (description, price) VALUES (?, ?);`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &product.Description, &product.Price)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) ReadAll(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, description, price FROM products;`

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

	products := make([]*domain.Product, 0)
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(&product.Id, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *repository) CreateMany(ctx context.Context, products []*domain.Product) error {
	query := `INSERT INTO products (id, description, price) VALUES`
	values := []any{}
	for _, product := range products {
		query += " (?, ?, ?),"
		values = append(values, product.Id, product.Description, product.Price)
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

func (r *repository) GetQtySaledGroupedByDescription(ctx context.Context) ([]map[string]any, error) {
	query := `
		SELECT p.description, SUM(s.quantity) AS total FROM products p
		JOIN sales s ON p.id = s.product_id
		GROUP BY p.description
		ORDER BY total DESC
		LIMIT 5;
	`

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

	qtySaledGrouped := make([]map[string]any, 0)
	for rows.Next() {
		register := map[string]any{}
		description := ""
		total := 0.0

		err := rows.Scan(&description, &total)
		if err != nil {
			return nil, err
		}

		register["description"] = description
		register["total"] = total
		qtySaledGrouped = append(qtySaledGrouped, register)
	}

	return qtySaledGrouped, nil
}
