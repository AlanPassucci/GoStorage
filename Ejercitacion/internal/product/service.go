package product

import (
	"context"
	"errors"
	"gostorage/internal/domain"
)

var (
	// ErrServiceProductNotFound is returned when a product is not found
	ErrServiceProductNotFound = errors.New("product not found")
	// ErrServiceInvalidProductID is returned when the product ID is invalid
	ErrServiceInvalidProductID = errors.New("invalid product identifier")
	// ErrServiceInvalidProductName is returned when the product name is invalid
	ErrServiceInvalidProductName = errors.New("invalid product name")
	// ErrServiceInvalidProductQuantity is returned when the product quantity is invalid
	ErrServiceInvalidProductQuantity = errors.New("invalid product quantity")
	// ErrServiceInvalidProductCodeValue is returned when the product code value is invalid
	ErrServiceInvalidProductCodeValue = errors.New("invalid product code value")
	// ErrServiceInvalidProductExpiration is returned when the product expiration is invalid
	ErrServiceInvalidProductExpiration = errors.New("invalid product expiration")
	// ErrServiceInvalidProductPrice is returned when the product price is invalid
	ErrServiceInvalidProductPrice = errors.New("invalid product price")
	// ErrServiceAlreadyExistsCodeValue is returned when the product code value already exists
	ErrorServiceAlreadyExistsCodeValue = errors.New("product code value already exists")
)

type Service interface {
	// GetAll returns all the products
	GetAll(ctx context.Context) ([]domain.Product, error)
	// GetByID returns a product by its ID
	GetByID(ctx context.Context, id int) (domain.Product, error)
	// Create creates a new product
	Create(ctx context.Context, p *domain.Product) (domain.Product, error)
	// Update updates a product
	Update(ctx context.Context, p *domain.Product) error
	// Delete deletes a product
	Delete(ctx context.Context, id int) error
}
