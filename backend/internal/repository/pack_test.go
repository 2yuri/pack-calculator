package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/internal/repository"
)

func createTestProduct(t *testing.T, repo *repository.Product) *domain.Product {
	t.Helper()
	p, err := repo.Create(context.Background(), "Test Product")
	require.NoError(t, err)
	return p
}

func TestPackRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)

	pack, err := packRepo.Create(context.Background(), product.ID, 250)
	require.NoError(t, err)
	assert.NotZero(t, pack.ID)
	assert.Equal(t, product.ID, pack.ProductID)
	assert.Equal(t, 250, pack.Size)
}

func TestPackRepository_Create_DuplicateSize(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)

	_, err := packRepo.Create(context.Background(), product.ID, 250)
	require.NoError(t, err)

	_, err = packRepo.Create(context.Background(), product.ID, 250)
	require.ErrorIs(t, err, domain.ErrDuplicatePackSize)
}

func TestPackRepository_CreateBatch(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)

	packs, err := packRepo.CreateBatch(context.Background(), product.ID, []int{250, 500, 1000})
	require.NoError(t, err)
	assert.Len(t, packs, 3)
}

func TestPackRepository_GetByProductID(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	packRepo.CreateBatch(context.Background(), product.ID, []int{250, 500})

	packs, err := packRepo.GetByProductID(context.Background(), product.ID)
	require.NoError(t, err)
	assert.Len(t, packs, 2)
	assert.Equal(t, 250, packs[0].Size)
	assert.Equal(t, 500, packs[1].Size)
}

func TestPackRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	created, _ := packRepo.Create(context.Background(), product.ID, 250)

	updated, err := packRepo.Update(context.Background(), created.ID, 300)
	require.NoError(t, err)
	assert.Equal(t, 300, updated.Size)
}

func TestPackRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	created, _ := packRepo.Create(context.Background(), product.ID, 250)

	err := packRepo.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = packRepo.GetByID(context.Background(), created.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPackRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	p1 := createTestProduct(t, productRepo)
	p2 := createTestProduct(t, productRepo)

	packRepo.Create(context.Background(), p1.ID, 250)
	packRepo.Create(context.Background(), p1.ID, 500)
	packRepo.Create(context.Background(), p2.ID, 1000)

	packs, err := packRepo.GetAll(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(packs), 3)
}

func TestPackRepository_GetAll_ExcludesSoftDeleted(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	created, _ := packRepo.Create(context.Background(), product.ID, 250)
	packRepo.Create(context.Background(), product.ID, 500)

	packRepo.Delete(context.Background(), created.ID)

	packs, err := packRepo.GetAll(context.Background())
	require.NoError(t, err)
	for _, p := range packs {
		assert.NotEqual(t, created.ID, p.ID)
	}
}

func TestPackRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	created, _ := packRepo.Create(context.Background(), product.ID, 250)

	pack, err := packRepo.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, pack.ID)
	assert.Equal(t, product.ID, pack.ProductID)
	assert.Equal(t, 250, pack.Size)
}

func TestPackRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	packRepo := repository.NewPack(db)

	_, err := packRepo.GetByID(context.Background(), 999999)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPackRepository_SoftDelete_ExcludedFromGetByProductID(t *testing.T) {
	db := setupTestDB(t)
	productRepo := repository.NewProduct(db)
	packRepo := repository.NewPack(db)

	product := createTestProduct(t, productRepo)
	p1, _ := packRepo.Create(context.Background(), product.ID, 250)
	packRepo.Create(context.Background(), product.ID, 500)

	packRepo.Delete(context.Background(), p1.ID)

	packs, err := packRepo.GetByProductID(context.Background(), product.ID)
	require.NoError(t, err)
	assert.Len(t, packs, 1)
	assert.Equal(t, 500, packs[0].Size)
}
