package domain

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicatePackSize  = errors.New("pack size already exists for this product")
	ErrInvalidQuantity    = errors.New("quantity must be greater than 0")
	ErrDuplicateProductID = errors.New("duplicate product_id in order")
	ErrNoPacksAvailable   = errors.New("no packs available for product")
)
