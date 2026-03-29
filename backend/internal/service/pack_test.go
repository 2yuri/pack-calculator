package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/internal/mocks"
)

func TestPackService_GetByProductID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	expected := []*domain.Pack{
		{ID: 1, ProductID: 1, Size: 250},
		{ID: 2, ProductID: 1, Size: 500},
	}
	packRepo.EXPECT().GetByProductID(gomock.Any(), 1).Return(expected, nil)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	result, err := svc.GetByProductID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestPackService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	created := &domain.Pack{ID: 1, ProductID: 1, Size: 250}
	packRepo.EXPECT().Create(gomock.Any(), 1, 250).Return(created, nil)
	cacheRepo.EXPECT().SetPack(gomock.Any(), created).Return(nil)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	result, err := svc.Create(context.Background(), 1, 250)

	require.NoError(t, err)
	assert.Equal(t, created, result)
}

func TestPackService_CreateBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	created := []*domain.Pack{
		{ID: 1, ProductID: 1, Size: 250},
		{ID: 2, ProductID: 1, Size: 500},
	}
	packRepo.EXPECT().CreateBatch(gomock.Any(), 1, []int{250, 500}).Return(created, nil)
	cacheRepo.EXPECT().SetPacks(gomock.Any(), created).Return(nil)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	result, err := svc.CreateBatch(context.Background(), 1, []int{250, 500})

	require.NoError(t, err)
	assert.Equal(t, created, result)
}

func TestPackService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	updated := &domain.Pack{ID: 1, ProductID: 1, Size: 300}
	packRepo.EXPECT().Update(gomock.Any(), 1, 300).Return(updated, nil)
	cacheRepo.EXPECT().SetPack(gomock.Any(), updated).Return(nil)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	result, err := svc.Update(context.Background(), 1, 300)

	require.NoError(t, err)
	assert.Equal(t, updated, result)
}

func TestPackService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	packRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)
	cacheRepo.EXPECT().RemovePack(gomock.Any(), 1).Return(nil)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
}

func TestPackService_Create_DuplicateSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packRepo := mocks.NewMockPackRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	packRepo.EXPECT().Create(gomock.Any(), 1, 250).Return(nil, domain.ErrDuplicatePackSize)

	svc := NewPack(PackDeps{
		Repo:      packRepo,
		CacheRepo: cacheRepo,
	})
	_, err := svc.Create(context.Background(), 1, 250)

	require.ErrorIs(t, err, domain.ErrDuplicatePackSize)
}
