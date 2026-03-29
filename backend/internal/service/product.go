package service

import (
	"context"

	"github.com/2yuri/pack-calculator/internal/domain"
)

type ProductDeps struct {
	Repo domain.ProductRepository
}

var _ domain.ProductService = (*Product)(nil)

type Product struct {
	repo domain.ProductRepository
}

func NewProduct(deps ProductDeps) *Product {
	return &Product{
		repo: deps.Repo,
	}
}

func (s *Product) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *Product) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Product) Create(ctx context.Context, name string) (*domain.Product, error) {
	return s.repo.Create(ctx, name)
}

func (s *Product) Update(ctx context.Context, id int, name string) (*domain.Product, error) {
	return s.repo.Update(ctx, id, name)
}

func (s *Product) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
