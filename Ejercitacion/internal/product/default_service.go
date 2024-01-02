package product

import (
	"context"
	"errors"
	"gostorage/internal/domain"
	"gostorage/pkg"
	"strings"
)

// service is the default implementation of the Service interface
type service struct {
	rp Repository
}

// NewService creates a new product service
func NewService(r Repository) Service {
	return &service{r}
}

func (sv *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	// Get the products from the repository
	products, err := sv.rp.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (sv *service) GetByID(ctx context.Context, id int) (domain.Product, error) {
	// Validate that the ID of the product is not zero or negative
	if id < 1 {
		return domain.Product{}, ErrServiceInvalidProductID
	}

	// Get the product from the repository
	product, err := sv.rp.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRepositoryProductNotFound):
			return domain.Product{}, ErrServiceProductNotFound
		default:
			return domain.Product{}, err
		}
	}
	return product, nil
}

func (sv *service) Create(ctx context.Context, p *domain.Product) (domain.Product, error) {
	// Validate that the name of the product is not empty
	if strings.TrimSpace(p.Name) == "" {
		return domain.Product{}, ErrServiceInvalidProductName
	}
	// Validate that the quantity of the product is not zero or negative
	if p.Quantity < 1 {
		return domain.Product{}, ErrServiceInvalidProductQuantity
	}
	// Validate that the code value of the product is not empty
	if strings.TrimSpace(p.CodeValue) == "" {
		return domain.Product{}, ErrServiceInvalidProductCodeValue
	}
	// Validate that the expiration of the product is not empty or invalid
	if strings.TrimSpace(p.Expiration) == "" || !pkg.IsValidDate(p.Expiration) {
		return domain.Product{}, ErrServiceInvalidProductExpiration
	}
	// Validate that the price of the product is not zero or negative
	if p.Price < 1 {
		return domain.Product{}, ErrServiceInvalidProductPrice
	}
	// Validate that the product code value does not already exist
	if exists, err := sv.rp.Exists(ctx, p.CodeValue); err != nil {
		return domain.Product{}, err
	} else if exists {
		return domain.Product{}, ErrorServiceAlreadyExistsCodeValue
	}

	// Create the product in the repository
	if err := sv.rp.Create(ctx, p); err != nil {
		return domain.Product{}, err
	}

	return *p, nil
}
func (sv *service) Update(ctx context.Context, p *domain.Product) error {
	// Validate that the ID of the product is not zero or negative
	if p.Id < 1 {
		return ErrServiceInvalidProductID
	}
	// Validate that the name of the product is not empty
	if strings.TrimSpace(p.Name) == "" {
		return ErrServiceInvalidProductID
	}
	// Validate that the quantity of the product is not zero or negative
	if p.Quantity < 1 {
		return ErrServiceInvalidProductQuantity
	}
	// Validate that the code value of the product is not empty
	if strings.TrimSpace(p.CodeValue) == "" {
		return ErrServiceInvalidProductCodeValue
	}
	// Validate that the expiration of the product is not empty or invalid
	if strings.TrimSpace(p.Expiration) == "" || !pkg.IsValidDate(p.Expiration) {
		return ErrServiceInvalidProductExpiration
	}
	// Validate that the price of the product is not zero or negative
	if p.Price < 1 {
		return ErrServiceInvalidProductPrice
	}
	// Validate that the product code value does not already exist with different ID
	if exists, err := sv.rp.ExistsWithDifferentID(ctx, p.Id, p.CodeValue); err != nil {
		return err
	} else if exists {
		return ErrorServiceAlreadyExistsCodeValue
	}

	// Update the product in the repository
	if err := sv.rp.Update(ctx, p); err != nil {
		switch {
		case errors.Is(err, ErrRepositoryProductNotFound):
			return ErrServiceProductNotFound
		default:
			return err
		}
	}

	return nil
}

func (sv *service) Delete(ctx context.Context, id int) error {
	// Validate that the ID of the product is not zero or negative
	if id < 1 {
		return ErrServiceInvalidProductID
	}

	// Delete the product from the repository
	if err := sv.rp.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, ErrRepositoryProductNotFound):
			return ErrServiceProductNotFound
		default:
			return err
		}
	}

	return nil
}
