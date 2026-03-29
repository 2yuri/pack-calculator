package domain

import "context"

type CalculationResultPack struct {
	Pack     *Pack `json:"pack"`
	Quantity int   `json:"quantity"`
}

type CalculationResult struct {
	Quantity   int                     `json:"quantity"`
	ProductID  int                     `json:"product_id"`
	TotalItems int                     `json:"total_items"`
	TotalPacks int                     `json:"total_packs"`
	Packs      []CalculationResultPack `json:"packs"`
}

type OrderItem struct {
	ProductID int
	Quantity  int
}

type OrderService interface {
	Calculate(ctx context.Context, items []OrderItem) ([]*CalculationResult, error)
}
