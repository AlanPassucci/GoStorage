package sales

import (
	"context"
	"desafio/internal/domain"
	"desafio/pkg/filemanager"
)

type Service interface {
	Create(ctx context.Context, sales *domain.Sale) error
	ReadAll(ctx context.Context) ([]*domain.Sale, error)
	CreateManyFromJSON(ctx context.Context) ([]*domain.Sale, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, sales *domain.Sale) error {
	id, err := s.r.Create(ctx, sales)
	if err != nil {
		return err
	}

	sales.Id = int(id)

	return nil
}

func (s *service) ReadAll(ctx context.Context) ([]*domain.Sale, error) {
	sales, err := s.r.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *service) CreateManyFromJSON(ctx context.Context) ([]*domain.Sale, error) {
	sales := make([]*domain.Sale, 0)
	if err := filemanager.LoadDataFromJSON("datos/sales.json", &sales); err != nil {
		return nil, err
	}

	if err := s.r.CreateMany(ctx, sales); err != nil {
		return nil, err
	}

	return sales, nil
}
