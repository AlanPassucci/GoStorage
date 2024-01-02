package product

import (
	"context"
	"database/sql"
	"gostorage/internal/domain"
)

// MySQLRepository is a repository that implements the Repository interface
type MySQLRepository struct {
	// db is the underlying MySQL database instance
	db *sql.DB
}

// NewRepository creates a new MySQLRepository
func NewRepository(db *sql.DB) Repository {
	return &MySQLRepository{db}
}

// GetAll returns all the products
func (rp *MySQLRepository) GetAll(ctx context.Context) (products []domain.Product, err error) {
	// Create the query
	query := `
		SELECT id, name, quantity, code_value, is_published, expiration, price
		FROM products
	`
	// Execute the query
	rows, err := rp.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	// Iterate over the rows
	for rows.Next() {
		// Create a product
		product := domain.Product{}

		// Scan the row into the product
		if err = rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price); err != nil {
			return
		}

		// Append the product to the products slice
		products = append(products, product)
	}

	return
}

// GetByID returns a product by its ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int) (product domain.Product, err error) {
	// Create the query
	query := `
		SELECT id, name, quantity, code_value, is_published, expiration, price
		FROM products WHERE id = ?
	`

	// Execute the query
	row := r.db.QueryRowContext(ctx, query, id)

	// Scan the row into the product
	if err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price); err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.Product{}, ErrRepositoryProductNotFound
		default:
			return domain.Product{}, err
		}
	}

	return
}

// Create creates a new product
func (r *MySQLRepository) Create(ctx context.Context, product *domain.Product) (err error) {
	// Create the query
	query := `
		INSERT INTO products (name, quantity, code_value, is_published, expiration, price)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	// Create the statement
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.ExecContext(ctx, product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted product
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the product
	product.Id = int(id)

	return
}

// Update updates a product
func (r *MySQLRepository) Update(ctx context.Context, product *domain.Product) (err error) {
	// Create the query
	query := `
		UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ?
		WHERE id = ?
	`
	// Create the statement
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.ExecContext(ctx, product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price, product.Id)
	if err != nil {
		return err
	}

	// Get the rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Check if the product was not found
	if rowsAffected == 0 {
		return ErrRepositoryProductNotFound
	}

	return
}

// Delete deletes a product
func (r *MySQLRepository) Delete(ctx context.Context, id int) (err error) {
	// Create the query
	query := `
		DELETE FROM products WHERE id = ?
	`
	// Create the statement
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted product
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Check if the product was not found
	if rowsAffected == 0 {
		return ErrRepositoryProductNotFound
	}

	return
}

func (r *MySQLRepository) Exists(ctx context.Context, codeValue string) (exists bool, err error) {
	// Create the query
	query := `
		SELECT EXISTS (
			SELECT code_value
			FROM products
			WHERE code_value = ?
		)
	`

	// Execute the query
	row := r.db.QueryRowContext(ctx, query, codeValue)

	// Scan the row into the exists variable
	if err = row.Scan(&exists); err != nil {
		return false, err
	}

	return
}

func (r *MySQLRepository) ExistsWithDifferentID(ctx context.Context, id int, codeValue string) (exists bool, err error) {
	// Create the query
	query := `
		SELECT EXISTS (
			SELECT code_value
			FROM products
			WHERE code_value = ? AND id != ?
		)
	`

	// Execute the query
	row := r.db.QueryRowContext(ctx, query, codeValue, id)

	// Scan the row into the exists variable
	if err = row.Scan(&exists); err != nil {
		return false, err
	}

	return
}
