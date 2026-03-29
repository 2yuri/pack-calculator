package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ domain.PackRepository = (*Pack)(nil)

type Pack struct {
	db *sqlx.DB
}

func NewPack(db *sqlx.DB) *Pack {
	return &Pack{db: db}
}

func (r *Pack) GetAll(ctx context.Context) ([]*domain.Pack, error) {
	var packs []*domain.Pack

	query := `
		SELECT id, product_id, size, created_at, updated_at 
		FROM packs WHERE deleted_at IS NULL 
		ORDER BY id
	`

	if err := r.db.SelectContext(ctx, &packs, query); err != nil {
		return nil, err
	}

	return packs, nil
}

func (r *Pack) GetByID(ctx context.Context, id int) (*domain.Pack, error) {
	var pack domain.Pack

	query := `
		SELECT id, product_id, size, created_at, updated_at 
		FROM packs 
		WHERE id = $1 AND deleted_at IS NULL
	`

	if err := r.db.GetContext(ctx, &pack, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &pack, nil
}

func (r *Pack) GetByProductID(ctx context.Context, productID int) ([]*domain.Pack, error) {
	var packs []*domain.Pack

	query := `
		SELECT id, product_id, size, created_at, updated_at 
		FROM packs 
		WHERE product_id = $1 AND deleted_at IS NULL 
		ORDER BY size
	`

	if err := r.db.SelectContext(ctx, &packs, query, productID); err != nil {
		return nil, err
	}

	return packs, nil
}

func (r *Pack) Create(ctx context.Context, productID int, size int) (*domain.Pack, error) {
	var pack domain.Pack

	query := `
		INSERT INTO packs (product_id, size)
		VALUES ($1, $2)
		RETURNING id, product_id, size, created_at, updated_at
	`

	if err := r.db.QueryRowxContext(ctx, query, productID, size).StructScan(&pack); err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, domain.ErrDuplicatePackSize
		}

		return nil, err
	}

	return &pack, nil
}

func (r *Pack) CreateBatch(ctx context.Context, productID int, sizes []int) ([]*domain.Pack, error) {
	var packs []*domain.Pack

	query := `
		INSERT INTO packs (product_id, size)
		SELECT $1, unnest($2::int[])
		RETURNING id, product_id, size, created_at, updated_at
	`

	if err := r.db.SelectContext(ctx, &packs, query, productID, pq.Array(sizes)); err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, domain.ErrDuplicatePackSize
		}
		return nil, err
	}

	return packs, nil
}

func (r *Pack) Update(ctx context.Context, id int, size int) (*domain.Pack, error) {
	var pack domain.Pack

	query := `
		UPDATE packs 
		SET size = $1, updated_at = NOW() 
		WHERE id = $2 AND deleted_at IS NULL 
		RETURNING id, product_id, size, created_at, updated_at
	`

	if err := r.db.QueryRowxContext(ctx, query, size, id).StructScan(&pack); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		if strings.Contains(err.Error(), "unique") {
			return nil, domain.ErrDuplicatePackSize
		}

		return nil, err
	}

	return &pack, nil
}

func (r *Pack) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE packs 
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
