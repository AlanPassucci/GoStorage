package customers

import (
	"context"
	"desafio/internal/domain"
	"desafio/pkg/filemanager"
)

type Service interface {
	Create(ctx context.Context, customers *domain.Customer) error
	ReadAll(ctx context.Context) ([]*domain.Customer, error)
	CreateManyFromJSON(ctx context.Context) ([]*domain.Customer, error)
	GetTotalsGroupedByCondition(ctx context.Context) ([]map[string]any, error)
	GetActivesWhoSpentTheMost(ctx context.Context) ([]map[string]any, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, customer *domain.Customer) error {
	insertedId, err := s.r.Create(ctx, customer)
	if err != nil {
		return err
	}

	customer.Id = int(insertedId)

	return nil
}

func (s *service) ReadAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := s.r.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *service) CreateManyFromJSON(ctx context.Context) ([]*domain.Customer, error) {
	customers := make([]*domain.Customer, 0)
	if err := filemanager.LoadDataFromJSON("datos/customers.json", &customers); err != nil {
		return nil, err
	}

	if err := s.r.CreateMany(ctx, customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *service) GetTotalsGroupedByCondition(ctx context.Context) ([]map[string]any, error) {
	totalGrouped, err := s.r.GetTotalsGroupedByCondition(ctx)
	if err != nil {
		return nil, err
	}

	return totalGrouped, nil
}

func (s *service) GetActivesWhoSpentTheMost(ctx context.Context) ([]map[string]any, error) {
	activesWhoSpentTheMost, err := s.r.GetActivesWhoSpentTheMost(ctx)
	if err != nil {
		return nil, err
	}

	return activesWhoSpentTheMost, nil
}
