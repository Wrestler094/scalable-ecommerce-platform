package domain

import (
	"context"
)

type Product struct {
	ID          int64
	Price       float64
	Name        string
	Description string
	CategoryID  int64
}

type ProductProvider interface {
	GetProductsByIDs(ctx context.Context, ids []int64) ([]Product, error)
}
