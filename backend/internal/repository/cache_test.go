package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/internal/repository"
)

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set, skipping integration test")
	}

	client := redis.NewClient(&redis.Options{Addr: addr})

	t.Cleanup(func() {
		client.FlushDB(context.Background())
		client.Close()
	})

	client.FlushDB(context.Background())
	return client
}

func TestCacheRepository_SetAndGetPacks(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	packs := []*domain.Pack{
		{ID: 1, ProductID: 1, Size: 250},
		{ID: 2, ProductID: 1, Size: 500},
		{ID: 3, ProductID: 2, Size: 1000},
	}

	err := repo.SetPacks(context.Background(), packs)
	require.NoError(t, err)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{1, 2})
	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, 1, result[0].ProductID)
	assert.Len(t, result[0].Packs, 2)
	assert.Equal(t, 2, result[1].ProductID)
	assert.Len(t, result[1].Packs, 1)
}

func TestCacheRepository_SetPack(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	pack := &domain.Pack{ID: 1, ProductID: 1, Size: 250}
	err := repo.SetPack(context.Background(), pack)
	require.NoError(t, err)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{1})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Len(t, result[0].Packs, 1)
	assert.Equal(t, 250, result[0].Packs[0].Size)

	pack2 := &domain.Pack{ID: 2, ProductID: 1, Size: 500}
	err = repo.SetPack(context.Background(), pack2)
	require.NoError(t, err)

	result, err = repo.GetPacksByProductIDs(context.Background(), []int{1})
	require.NoError(t, err)
	assert.Len(t, result[0].Packs, 2)
}

func TestCacheRepository_SetPack_UpdateExisting(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	pack := &domain.Pack{ID: 1, ProductID: 1, Size: 250}
	repo.SetPack(context.Background(), pack)

	updated := &domain.Pack{ID: 1, ProductID: 1, Size: 300}
	err := repo.SetPack(context.Background(), updated)
	require.NoError(t, err)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{1})
	require.NoError(t, err)
	assert.Len(t, result[0].Packs, 1)
	assert.Equal(t, 300, result[0].Packs[0].Size)
}

func TestCacheRepository_RemovePack(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	packs := []*domain.Pack{
		{ID: 1, ProductID: 1, Size: 250},
		{ID: 2, ProductID: 1, Size: 500},
	}
	repo.SetPacks(context.Background(), packs)

	err := repo.RemovePack(context.Background(), 1)
	require.NoError(t, err)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{1})
	require.NoError(t, err)
	assert.Len(t, result[0].Packs, 1)
	assert.Equal(t, 500, result[0].Packs[0].Size)
}

func TestCacheRepository_Invalidate(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	packs := []*domain.Pack{
		{ID: 1, ProductID: 1, Size: 250},
		{ID: 2, ProductID: 2, Size: 500},
	}
	repo.SetPacks(context.Background(), packs)

	err := repo.Invalidate(context.Background())
	require.NoError(t, err)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{1, 2})
	require.NoError(t, err)
	assert.Len(t, result[0].Packs, 0)
	assert.Len(t, result[1].Packs, 0)
}

func TestCacheRepository_GetPacksByProductIDs_MissingProduct(t *testing.T) {
	client := setupTestRedis(t)
	repo := repository.NewCache(client)

	result, err := repo.GetPacksByProductIDs(context.Background(), []int{999})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 999, result[0].ProductID)
	assert.Empty(t, result[0].Packs)
}
