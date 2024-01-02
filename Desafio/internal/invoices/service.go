package invoices

import (
	"context"
	"desafio/internal/domain"
	"desafio/pkg/filemanager"
)

type Service interface {
	Create(ctx context.Context, invoices *domain.Invoice) error
	ReadAll(ctx context.Context) ([]*domain.Invoice, error)
	CreateManyFromJSON(ctx context.Context) ([]*domain.Invoice, error)
	UpdateTotals(ctx context.Context) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, invoices *domain.Invoice) error {
	insertedId, err := s.r.Create(ctx, invoices)
	if err != nil {
		return err
	}

	invoices.Id = int(insertedId)

	return nil
}

func (s *service) ReadAll(ctx context.Context) ([]*domain.Invoice, error) {
	invoices, err := s.r.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *service) CreateManyFromJSON(ctx context.Context) ([]*domain.Invoice, error) {
	invoices := make([]*domain.Invoice, 0)
	if err := filemanager.LoadDataFromJSON("datos/invoices.json", &invoices); err != nil {
		return nil, err
	}

	if err := s.r.CreateMany(ctx, invoices); err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *service) UpdateTotals(ctx context.Context) error {
	err := s.r.UpdateTotals(ctx)
	if err != nil {
		return err
	}

	return nil
}
