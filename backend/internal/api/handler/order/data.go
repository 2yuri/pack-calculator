package order

import "github.com/2yuri/pack-calculator/internal/domain"

type ItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CalculateRequest struct {
	Items []ItemRequest `json:"items"`
}

type CalculateResponse struct {
	Results []*domain.CalculationResult `json:"results"`
}
