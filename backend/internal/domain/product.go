package domain

import (
	"context"
	"time"
)

type Product struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	GetByIDs(ctx context.Context, ids []int) ([]*Product, error)
	Create(ctx context.Context, name string) (*Product, error)
	Update(ctx context.Context, id int, name string) (*Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductService interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Create(ctx context.Context, name string) (*Product, error)
	Update(ctx context.Context, id int, name string) (*Product, error)
	Delete(ctx context.Context, id int) error
}
