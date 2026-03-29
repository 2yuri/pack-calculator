package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/internal/repository"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := sqlx.ConnectContext(context.Background(), "postgres", dbURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.ExecContext(context.Background(), "DELETE FROM packs")
		db.ExecContext(context.Background(), "DELETE FROM products")
		db.Close()
	})

	db.ExecContext(context.Background(), "DELETE FROM packs")
	db.ExecContext(context.Background(), "DELETE FROM products")

	return db
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	product, err := repo.Create(context.Background(), "Test Product")
	require.NoError(t, err)
	assert.NotZero(t, product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.NotZero(t, product.CreatedAt)
	assert.NotZero(t, product.UpdatedAt)
}

func TestProductRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	repo.Create(context.Background(), "Product A")
	repo.Create(context.Background(), "Product B")

	products, err := repo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, products, 2)
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	created, _ := repo.Create(context.Background(), "Test Product")

	product, err := repo.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, product.ID)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	_, err := repo.GetByID(context.Background(), 99999)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestProductRepository_GetByIDs(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	p1, _ := repo.Create(context.Background(), "Product A")
	p2, _ := repo.Create(context.Background(), "Product B")
	repo.Create(context.Background(), "Product C")

	products, err := repo.GetByIDs(context.Background(), []int{p1.ID, p2.ID})
	require.NoError(t, err)
	assert.Len(t, products, 2)
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	created, _ := repo.Create(context.Background(), "Old Name")

	updated, err := repo.Update(context.Background(), created.ID, "New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	created, _ := repo.Create(context.Background(), "To Delete")

	err := repo.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), created.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProduct(db)

	err := repo.Delete(context.Background(), 99999)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
