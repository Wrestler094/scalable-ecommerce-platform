package usecase

import (
	"context"

	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/domain"
)

var _ domain.CartUseCase = (*cartUseCase)(nil)

type cartUseCase struct {
	repo domain.CartRepository
}

func NewCartUseCase(repo domain.CartRepository) domain.CartUseCase {
	return &cartUseCase{repo: repo}
}

func (uc *cartUseCase) GetCart(ctx context.Context, userID int64) ([]domain.CartItem, error) {
	return uc.repo.Get(ctx, userID)
}

func (uc *cartUseCase) AddItem(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repo.Add(ctx, userID, productID, quantity)
}

func (uc *cartUseCase) UpdateItem(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repo.Update(ctx, userID, productID, quantity)
}

func (uc *cartUseCase) RemoveItem(ctx context.Context, userID, productID int64) error {
	return uc.repo.Remove(ctx, userID, productID)
}

func (uc *cartUseCase) ClearCart(ctx context.Context, userID int64) error {
	return uc.repo.Clear(ctx, userID)
}
