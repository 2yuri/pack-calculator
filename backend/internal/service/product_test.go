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

func TestProductService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	expected := []*domain.Product{
		{ID: 1, Name: "Product A"},
		{ID: 2, Name: "Product B"},
	}
	repo.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	result, err := svc.GetAll(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestProductService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	expected := &domain.Product{ID: 1, Name: "Product A"}
	repo.EXPECT().GetByID(gomock.Any(), 1).Return(expected, nil)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	result, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	repo.EXPECT().GetByID(gomock.Any(), 99).Return(nil, domain.ErrNotFound)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	_, err := svc.GetByID(context.Background(), 99)

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestProductService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	expected := &domain.Product{ID: 1, Name: "New Product"}
	repo.EXPECT().Create(gomock.Any(), "New Product").Return(expected, nil)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	result, err := svc.Create(context.Background(), "New Product")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestProductService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	expected := &domain.Product{ID: 1, Name: "Updated"}
	repo.EXPECT().Update(gomock.Any(), 1, "Updated").Return(expected, nil)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	result, err := svc.Update(context.Background(), 1, "Updated")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestProductService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)

	repo.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	svc := NewProduct(ProductDeps{
		Repo: repo,
	})
	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
}
