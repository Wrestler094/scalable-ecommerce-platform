package usecase

import (
	"context"

	"catalog-service/internal/domain"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, name string) (int64, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
}

type categoryUseCase struct {
	repo domain.CategoryRepository
}

func NewCategoryUseCase(r domain.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{repo: r}
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, name string) (int64, error) {
	return uc.repo.Save(ctx, domain.Category{Name: name})
}

func (uc *categoryUseCase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.repo.FindAll(ctx)
}
