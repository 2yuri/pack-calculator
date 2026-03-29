package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/jmoiron/sqlx"
)

var _ domain.ProductRepository = (*Product)(nil)

type Product struct {
	db *sqlx.DB
}

func NewProduct(db *sqlx.DB) *Product {
	return &Product{db: db}
}

func (r *Product) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product

	query := `
		SELECT id, name, created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Product) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	if err := r.db.GetContext(ctx, &product, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &product, nil
}

func (r *Product) GetByIDs(ctx context.Context, ids []int) ([]*domain.Product, error) {
	var products []*domain.Product

	queryStr := `
		SELECT id, name, created_at, updated_at
		FROM products
		WHERE id IN (?) AND deleted_at IS NULL
		ORDER BY id
	`

	query, args, err := sqlx.In(queryStr, ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Product) Create(ctx context.Context, name string) (*domain.Product, error) {
	var product domain.Product

	query := `
		INSERT INTO products (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	if err := r.db.QueryRowxContext(ctx, query, name).StructScan(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Product) Update(ctx context.Context, id int, name string) (*domain.Product, error) {
	var product domain.Product

	query := `
		UPDATE products
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, name, created_at, updated_at
	`

	if err := r.db.QueryRowxContext(ctx, query, name, id).StructScan(&product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &product, nil
}

func (r *Product) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE products
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
