package domain

import (
	"context"
)

type CartItem struct {
	ProductID int64
	Quantity  int
}

type CartRepository interface {
	Get(ctx context.Context, userID int64) ([]CartItem, error)
	Add(ctx context.Context, userID, productID int64, quantity int) error
	Update(ctx context.Context, userID, productID int64, quantity int) error
	Remove(ctx context.Context, userID, productID int64) error
	Clear(ctx context.Context, userID int64) error
}

type CartUseCase interface {
	GetCart(ctx context.Context, userID int64) ([]CartItem, error)
	AddItem(ctx context.Context, userID, productID int64, quantity int) error
	UpdateItem(ctx context.Context, userID, productID int64, quantity int) error
	RemoveItem(ctx context.Context, userID, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}
