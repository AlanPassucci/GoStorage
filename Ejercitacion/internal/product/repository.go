package product

import (
	"context"
	"errors"

	"gostorage/internal/domain"
)

var (
	// ErrNotFound is returned when a product is not found
	ErrRepositoryProductNotFound = errors.New("product not found")
)

// Repository is an interface that defines the methods that a product repository must implement
type Repository interface {
	// GetAll returns all the products
	GetAll(ctx context.Context) ([]domain.Product, error)
	// GetByID returns a product by its ID
	GetByID(ctx context.Context, id int) (domain.Product, error)
	// Create creates a new product
	Create(ctx context.Context, product *domain.Product) error
	// Update updates a product
	Update(ctx context.Context, product *domain.Product) error
	// Delete deletes a product
	Delete(ctx context.Context, id int) error
	// Exists verify the existence of a product with the given product code value
	Exists(ctx context.Context, codeValue string) (bool, error)
	// ExistsWithDifferentID verify the existence of a product with the given product code value and different ID
	ExistsWithDifferentID(ctx context.Context, id int, codeValue string) (bool, error)
}
