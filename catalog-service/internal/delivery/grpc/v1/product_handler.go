package v1

import (
	"context"
	"errors"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/usecase"
	catalogv1 "github.com/Wrestler094/scalable-ecommerce-platform/gen/go/catalog/v1"
)

type ProductHandler struct {
	catalogv1.UnimplementedCatalogServiceServer
	productUC usecase.ProductUseCase
	logger    logger.Logger
}

func NewProductHandler(productUC usecase.ProductUseCase, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		productUC: productUC,
		logger:    logger,
	}
}

func (h *ProductHandler) GetProductsByIDs(ctx context.Context, req *catalogv1.GetProductsByIDsRequest) (*catalogv1.GetProductsByIDsResponse, error) {
	products, err := h.productUC.GetProductsByID(ctx, req.GetProductIds())
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			h.logger.WithError(err).Warn("products not found")
			return nil, status.Errorf(codes.NotFound, "products not found: %v", err)
		}

		h.logger.WithError(err).Error("failed to get products")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &catalogv1.GetProductsByIDsResponse{Products: convertProductsToProto(products)}, nil
}

func convertProductsToProto(products []domain.Product) []*catalogv1.Product {
	pbProducts := make([]*catalogv1.Product, 0, len(products))
	for _, p := range products {
		pbProducts = append(pbProducts, &catalogv1.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CategoryId:  p.CategoryID,
		})
	}
	return pbProducts
}
