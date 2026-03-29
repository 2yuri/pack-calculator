package domain

import (
	"context"
	"time"
)

type Pack struct {
	ID        int        `json:"id" db:"id"`
	ProductID int        `json:"product_id" db:"product_id"`
	Size      int        `json:"size" db:"size"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ProductPacks struct {
	ProductID int
	Packs     []*Pack
}

type PackRepository interface {
	GetAll(ctx context.Context) ([]*Pack, error)
	GetByID(ctx context.Context, id int) (*Pack, error)
	GetByProductID(ctx context.Context, productID int) ([]*Pack, error)
	Create(ctx context.Context, productID int, size int) (*Pack, error)
	CreateBatch(ctx context.Context, productID int, sizes []int) ([]*Pack, error)
	Update(ctx context.Context, id int, size int) (*Pack, error)
	Delete(ctx context.Context, id int) error
}

type CacheRepository interface {
	GetPacksByProductIDs(ctx context.Context, productIDs []int) ([]*ProductPacks, error)
	SetPack(ctx context.Context, pack *Pack) error
	SetPacks(ctx context.Context, packs []*Pack) error
	RemovePack(ctx context.Context, id int) error
	Invalidate(ctx context.Context) error
}

type PackService interface {
	GetByProductID(ctx context.Context, productID int) ([]*Pack, error)
	Create(ctx context.Context, productID int, size int) (*Pack, error)
	CreateBatch(ctx context.Context, productID int, sizes []int) ([]*Pack, error)
	Update(ctx context.Context, id int, size int) (*Pack, error)
	Delete(ctx context.Context, id int) error
}
