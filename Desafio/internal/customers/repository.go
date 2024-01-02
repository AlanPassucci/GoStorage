package customers

import (
	"context"
	"database/sql"
	"strings"

	"desafio/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, customers *domain.Customer) (int64, error)
	ReadAll(ctx context.Context) ([]*domain.Customer, error)
	CreateMany(ctx context.Context, customers []*domain.Customer) error
	GetTotalsGroupedByCondition(ctx context.Context) ([]map[string]any, error)
	GetActivesWhoSpentTheMost(ctx context.Context) ([]map[string]any, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, customers *domain.Customer) (int64, error) {
	query := `INSERT INTO customers (first_name, last_name, customers.condition) VALUES (?, ?, ?);`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &customers.FirstName, &customers.LastName, &customers.Condition)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) ReadAll(ctx context.Context) ([]*domain.Customer, error) {
	query := `SELECT id, first_name, last_name, customers.condition FROM customers;`

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

	customers := make([]*domain.Customer, 0)
	for rows.Next() {
		customer := domain.Customer{}
		err := rows.Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Condition)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	return customers, nil
}

func (r *repository) CreateMany(ctx context.Context, customers []*domain.Customer) error {
	query := `INSERT INTO customers (id, first_name, last_name, customers.condition) VALUES`
	values := []any{}
	for _, customer := range customers {
		query += " (?, ?, ?, ?),"
		values = append(values, customer.Id, customer.FirstName, customer.LastName, customer.Condition)
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

func (r *repository) GetTotalsGroupedByCondition(ctx context.Context) ([]map[string]any, error) {
	query := `
			SELECT c.condition, ROUND(SUM(i.total), 2) AS total FROM customers c
			JOIN invoices i ON c.id = i.customer_id
			GROUP BY c.condition;
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

	totalsGrouped := make([]map[string]any, 0)
	for rows.Next() {
		register := map[string]any{}
		condition := false
		total := 0.0

		err := rows.Scan(&condition, &total)
		if err != nil {
			return nil, err
		}

		register["condition"] = condition
		register["total"] = total
		totalsGrouped = append(totalsGrouped, register)
	}

	return totalsGrouped, nil
}

func (r *repository) GetActivesWhoSpentTheMost(ctx context.Context) ([]map[string]any, error) {
	query := `
			SELECT c.last_name, c.first_name, ROUND(SUM(i.total), 2) AS amount FROM customers c
			JOIN invoices i ON c.id = i.customer_id
			WHERE c.condition = 1
			GROUP BY c.last_name, c.first_name
			ORDER BY amount DESC
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

	spentTheMost := make([]map[string]any, 0)
	for rows.Next() {
		register := map[string]any{}
		lastName := ""
		firstName := ""
		amount := 0.0

		err := rows.Scan(&lastName, &firstName, &amount)
		if err != nil {
			return nil, err
		}

		register["last_name"] = lastName
		register["first_name"] = firstName
		register["amount"] = amount
		spentTheMost = append(spentTheMost, register)
	}

	return spentTheMost, nil
}
