package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/delivery/http/v1/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/domain"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (int64, error)
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, req dto.UpdateProductRequest) (domain.Product, error)
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

func (uc *productUseCase) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (int64, error) {
	p := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       roundTo2DecimalPlaces(req.Price),
		CategoryID:  req.CategoryID,
	}
	return uc.repo.Save(ctx, p)
}

func (uc *productUseCase) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id int64, req dto.UpdateProductRequest) (domain.Product, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("find product: %w", err)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Price != nil {
		existing.Price = roundTo2DecimalPlaces(*req.Price)
	}
	if req.CategoryID != nil {
		existing.CategoryID = *req.CategoryID
	}

	if err := uc.repo.Update(ctx, *existing); err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}

	return *existing, nil
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
