package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/usecase/dto"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, input dto.CreateProductInput) (*dto.CreateProductOutput, error)
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProductsByID(ctx context.Context, ids []int64) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, input dto.UpdateProductInput) (*dto.UpdateProductOutput, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context) ([]domain.Product, error)
	ListByCategory(ctx context.Context, categoryID int64) ([]domain.Product, error)
}

type productUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(r domain.ProductRepository) ProductUseCase {
	return &productUseCase{repo: r}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, input dto.CreateProductInput) (*dto.CreateProductOutput, error) {
	p := domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       roundTo2DecimalPlaces(input.Price),
		CategoryID:  input.CategoryID,
	}
	id, err := uc.repo.Save(ctx, p)
	if err != nil {
		return nil, err
	}
	return &dto.CreateProductOutput{ID: id}, nil
}

func (uc *productUseCase) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *productUseCase) GetProductsByID(ctx context.Context, ids []int64) ([]domain.Product, error) {
	return uc.repo.FindByIDs(ctx, ids)
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id int64, input dto.UpdateProductInput) (*dto.UpdateProductOutput, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Price != nil {
		existing.Price = roundTo2DecimalPlaces(*input.Price)
	}
	if input.CategoryID != nil {
		existing.CategoryID = *input.CategoryID
	}

	if err := uc.repo.Update(ctx, *existing); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &dto.UpdateProductOutput{
		ID:          existing.ID,
		Name:        existing.Name,
		Description: existing.Description,
		Price:       existing.Price,
		CategoryID:  existing.CategoryID,
	}, nil
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *productUseCase) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *productUseCase) ListByCategory(ctx context.Context, categoryID int64) ([]domain.Product, error) {
	return uc.repo.FindByCategoryID(ctx, categoryID)
}

func roundTo2DecimalPlaces(val float64) float64 {
	return math.Round(val*100) / 100
}
