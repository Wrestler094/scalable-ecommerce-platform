package grpc

import (
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"google.golang.org/grpc"

	grpcV1 "github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/delivery/grpc/v1"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/usecase"
	catalogv1 "github.com/Wrestler094/scalable-ecommerce-platform/gen/go/catalog/v1"
)

// RegisterServices регистрирует все gRPC обработчики на сервере.
func RegisterServices(
	gRPCServer *grpc.Server,
	productUC usecase.ProductUseCase,
	logger logger.Logger,
) {
	// Создаем и регистрируем обработчик для продуктов
	productHandler := grpcV1.NewProductHandler(productUC, logger)
	catalogv1.RegisterCatalogServiceServer(gRPCServer, productHandler)
}
