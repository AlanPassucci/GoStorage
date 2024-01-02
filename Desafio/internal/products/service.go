package products

import (
	"context"
	"desafio/internal/domain"
	"desafio/pkg/filemanager"
)

type Service interface {
	Create(ctx context.Context, product *domain.Product) error
	ReadAll(ctx context.Context) ([]*domain.Product, error)
	CreateManyFromJSON(ctx context.Context) ([]*domain.Product, error)
	GetQtySaledGroupedByDescription(ctx context.Context) ([]map[string]any, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, product *domain.Product) error {
	id, err := s.r.Create(ctx, product)
	if err != nil {
		return err
	}

	product.Id = int(id)

	return nil
}

func (s *service) ReadAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.r.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *service) CreateManyFromJSON(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	if err := filemanager.LoadDataFromJSON("datos/products.json", &products); err != nil {
		return nil, err
	}

	if err := s.r.CreateMany(ctx, products); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *service) GetQtySaledGroupedByDescription(ctx context.Context) ([]map[string]any, error) {
	qtySaledGrouped, err := s.r.GetQtySaledGroupedByDescription(ctx)
	if err != nil {
		return nil, err
	}

	return qtySaledGrouped, nil
}
