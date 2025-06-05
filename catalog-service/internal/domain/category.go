package domain

import (
	"context"
)

type Category struct {
	ID   int64
	Name string
}

type CategoryRepository interface {
	Save(ctx context.Context, c Category) (int64, error)
	FindAll(ctx context.Context) ([]Category, error)
}
