package service

import (
	"context"
	"maps"
	"sort"

	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/internal/domain"
)

type OrderDeps struct {
	CacheRepo domain.CacheRepository
}

var _ domain.OrderService = (*Order)(nil)

type Order struct {
	cacheRepo domain.CacheRepository
}

func NewOrder(deps OrderDeps) *Order {
	return &Order{
		cacheRepo: deps.CacheRepo,
	}
}

func (s *Order) Calculate(ctx context.Context, items []domain.OrderItem) ([]*domain.CalculationResult, error) {
	productsMap := make(map[int]bool)
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, domain.ErrInvalidQuantity
		}

		if productsMap[item.ProductID] {
			return nil, domain.ErrDuplicateProductID
		}

		productsMap[item.ProductID] = true
	}

	ids := make([]int, 0, len(productsMap))
	for productID := range productsMap {
		ids = append(ids, productID)
	}

	productPacks, err := s.cacheRepo.GetPacksByProductIDs(ctx, ids)
	if err != nil {
		zap.L().Error("failed to fetch packs from cache", zap.Error(err))
		return nil, err
	}

	packsMap := make(map[int][]*domain.Pack)
	for _, pp := range productPacks {
		packsMap[pp.ProductID] = pp.Packs
	}

	results := make([]*domain.CalculationResult, 0, len(items))
	for _, item := range items {
		packs := packsMap[item.ProductID]
		if len(packs) == 0 {
			return nil, domain.ErrNoPacksAvailable
		}

		result, err := s.calculateForItem(item, packs)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *Order) calculateForItem(item domain.OrderItem, packs []*domain.Pack) (*domain.CalculationResult, error) {
	smallest := packs[0]
	for _, p := range packs[1:] {
		if p.Size < smallest.Size {
			smallest = p
		}
	}

	if item.Quantity <= smallest.Size {
		return &domain.CalculationResult{
			Quantity:   item.Quantity,
			ProductID:  item.ProductID,
			TotalItems: smallest.Size,
			TotalPacks: 1,
			Packs: []domain.CalculationResultPack{
				{
					Pack:     smallest,
					Quantity: 1,
				},
			},
		}, nil
	}

	result := s.dpSolve(item.Quantity, packs)

	totalItems := 0
	totalPacks := 0
	for _, c := range result {
		totalItems += c.Pack.Size * c.Quantity
		totalPacks += c.Quantity
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Pack.Size > result[j].Pack.Size
	})

	return &domain.CalculationResult{
		Quantity:   item.Quantity,
		ProductID:  item.ProductID,
		TotalItems: totalItems,
		TotalPacks: totalPacks,
		Packs:      result,
	}, nil
}

type dpState struct {
	valid      bool
	totalItems int
	totalPacks int
	packs      map[int]int
}

func (s *Order) dpSolve(quantity int, packs []*domain.Pack) []domain.CalculationResultPack {
	sizes := make([]int, len(packs))
	packsMap := make(map[int]*domain.Pack)

	for i, p := range packs {
		sizes[i] = p.Size
		packsMap[p.Size] = p
	}

	maxItems := quantity + sizes[0]
	for _, size := range sizes {
		if quantity+size > maxItems {
			maxItems = quantity + size
		}
	}

	dp := make([]dpState, maxItems+1)
	dp[0] = dpState{valid: true, totalItems: 0, totalPacks: 0, packs: make(map[int]int)}

	for i := 1; i <= maxItems; i++ {
		for _, size := range sizes {
			current := dp[i]

			prevIdx := i - size
			if prevIdx < 0 {
				continue
			}

			prev := dp[prevIdx]
			if !prev.valid {
				continue
			}

			newItems := prev.totalItems + size
			newPacks := prev.totalPacks + 1

			if !current.valid || newItems < current.totalItems ||
				(newItems == current.totalItems && newPacks < current.totalPacks) {
				packs := make(map[int]int)
				maps.Copy(packs, prev.packs)
				packs[size]++
				dp[i] = dpState{valid: true, totalItems: newItems, totalPacks: newPacks, packs: packs}
			}
		}
	}

	var best dpState
	for i := quantity; i <= maxItems; i++ {
		if !dp[i].valid {
			continue
		}

		if !best.valid || dp[i].totalItems < best.totalItems ||
			(dp[i].totalItems == best.totalItems && dp[i].totalPacks < best.totalPacks) {
			best = dp[i]
		}
	}

	return s.formatResult(best, packsMap)
}

func (s *Order) formatResult(state dpState, packsMap map[int]*domain.Pack) []domain.CalculationResultPack {
	var result []domain.CalculationResultPack

	for size, count := range state.packs {
		result = append(result, domain.CalculationResultPack{
			Pack:     packsMap[size],
			Quantity: count,
		})
	}

	return result
}
