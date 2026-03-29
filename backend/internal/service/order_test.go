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

func TestOrderService_Calculate(t *testing.T) {
	tests := []struct {
		name           string
		items          []domain.OrderItem
		cachedPacks    []*domain.ProductPacks
		expectedResult []*domain.CalculationResult
		expectedErr    error
	}{
		{
			name: "1 item orders 1x250 pack",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 1},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
					{ID: 3, ProductID: 1, Size: 1000},
					{ID: 4, ProductID: 1, Size: 2000},
					{ID: 5, ProductID: 1, Size: 5000},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   1,
					ProductID:  1,
					TotalItems: 250,
					TotalPacks: 1,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 1, ProductID: 1, Size: 250}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "250 items orders 1x250 pack",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 250},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
					{ID: 3, ProductID: 1, Size: 1000},
					{ID: 4, ProductID: 1, Size: 2000},
					{ID: 5, ProductID: 1, Size: 5000},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   250,
					ProductID:  1,
					TotalItems: 250,
					TotalPacks: 1,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 1, ProductID: 1, Size: 250}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "251 items orders 1x500 pack",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 251},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
					{ID: 3, ProductID: 1, Size: 1000},
					{ID: 4, ProductID: 1, Size: 2000},
					{ID: 5, ProductID: 1, Size: 5000},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   251,
					ProductID:  1,
					TotalItems: 500,
					TotalPacks: 1,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 2, ProductID: 1, Size: 500}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "501 items orders 1x500 + 1x250",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 501},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
					{ID: 3, ProductID: 1, Size: 1000},
					{ID: 4, ProductID: 1, Size: 2000},
					{ID: 5, ProductID: 1, Size: 5000},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   501,
					ProductID:  1,
					TotalItems: 750,
					TotalPacks: 2,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 2, ProductID: 1, Size: 500}, Quantity: 1},
						{Pack: &domain.Pack{ID: 1, ProductID: 1, Size: 250}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "12001 items orders 2x5000 + 1x2000 + 1x250",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 12001},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
					{ID: 3, ProductID: 1, Size: 1000},
					{ID: 4, ProductID: 1, Size: 2000},
					{ID: 5, ProductID: 1, Size: 5000},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   12001,
					ProductID:  1,
					TotalItems: 12250,
					TotalPacks: 4,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 5, ProductID: 1, Size: 5000}, Quantity: 2},
						{Pack: &domain.Pack{ID: 4, ProductID: 1, Size: 2000}, Quantity: 1},
						{Pack: &domain.Pack{ID: 1, ProductID: 1, Size: 250}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "multi-product order",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 250},
				{ProductID: 2, Quantity: 501},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
				}},
				{ProductID: 2, Packs: []*domain.Pack{
					{ID: 3, ProductID: 2, Size: 250},
					{ID: 4, ProductID: 2, Size: 500},
				}},
			},
			expectedResult: []*domain.CalculationResult{
				{
					Quantity:   250,
					ProductID:  1,
					TotalItems: 250,
					TotalPacks: 1,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 1, ProductID: 1, Size: 250}, Quantity: 1},
					},
				},
				{
					Quantity:   501,
					ProductID:  2,
					TotalItems: 750,
					TotalPacks: 2,
					Packs: []domain.CalculationResultPack{
						{Pack: &domain.Pack{ID: 4, ProductID: 2, Size: 500}, Quantity: 1},
						{Pack: &domain.Pack{ID: 3, ProductID: 2, Size: 250}, Quantity: 1},
					},
				},
			},
		},
		{
			name: "invalid quantity returns error",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 0},
			},
			expectedErr: domain.ErrInvalidQuantity,
		},
		{
			name: "negative quantity returns error",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: -5},
			},
			expectedErr: domain.ErrInvalidQuantity,
		},
		{
			name: "duplicate product_id returns error",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 100},
				{ProductID: 1, Quantity: 200},
			},
			expectedErr: domain.ErrDuplicateProductID,
		},
		{
			name: "no packs available returns error",
			items: []domain.OrderItem{
				{ProductID: 1, Quantity: 100},
			},
			cachedPacks: []*domain.ProductPacks{
				{ProductID: 1, Packs: []*domain.Pack{}},
			},
			expectedErr: domain.ErrNoPacksAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheRepo := mocks.NewMockCacheRepository(ctrl)

			if tt.cachedPacks != nil {
				var productIDs []int
				for _, item := range tt.items {
					productIDs = append(productIDs, item.ProductID)
				}
				cacheRepo.EXPECT().
					GetPacksByProductIDs(gomock.Any(), productIDs).
					Return(tt.cachedPacks, nil)
			}

			svc := NewOrder(OrderDeps{
				CacheRepo: cacheRepo,
			})

			results, err := svc.Calculate(context.Background(), tt.items)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Len(t, results, len(tt.expectedResult))

			for i, expected := range tt.expectedResult {
				actual := results[i]
				assert.Equal(t, expected.Quantity, actual.Quantity)
				assert.Equal(t, expected.ProductID, actual.ProductID)
				assert.Equal(t, expected.TotalItems, actual.TotalItems)
				assert.Equal(t, expected.TotalPacks, actual.TotalPacks)
				require.Len(t, actual.Packs, len(expected.Packs))
				for j, ep := range expected.Packs {
					assert.Equal(t, ep.Pack.Size, actual.Packs[j].Pack.Size)
					assert.Equal(t, ep.Quantity, actual.Packs[j].Quantity)
				}
			}
		})
	}
}
