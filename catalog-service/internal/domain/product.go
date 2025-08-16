package domain

import (
	"context"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	CategoryID  int64
}

type ProductRepository interface {
	Save(ctx context.Context, p Product) (int64, error)
	FindByID(ctx context.Context, id int64) (*Product, error)
	FindByIDs(ctx context.Context, ids []int64) ([]Product, error)
	Update(ctx context.Context, p Product) error
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context) ([]Product, error)
	FindByCategoryID(ctx context.Context, categoryID int64) ([]Product, error)
}
