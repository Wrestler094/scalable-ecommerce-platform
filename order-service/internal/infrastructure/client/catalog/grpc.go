package catalog

import (
	"context"
	"fmt"

	pb "github.com/Wrestler094/scalable-ecommerce-platform/gen/go/catalog/v1"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.CatalogServiceClient
}

var _ domain.ProductProvider = (*Client)(nil)

func NewClient(ctx context.Context, grpcURL string) (*Client, error) {
	conn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to catalog service: %w", err)
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	return &Client{
		CatalogServiceClient: pb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *Client) GetProductsByIDs(ctx context.Context, ids []int64) ([]domain.Product, error) {
	// Имена в запросе и RPC-вызове соответствуют вашему .proto файлу
	resp, err := c.CatalogServiceClient.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{ProductIds: ids})
	if err != nil {
		// Улучшение: добавляем контекст к ошибке
		return nil, fmt.Errorf("failed to get products from catalog service: %w", err)
	}

	products := make([]domain.Product, len(resp.GetProducts()))
	for i, p := range resp.GetProducts() {
		products[i] = domain.Product{
			ID:          p.Id,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
			CategoryID:  p.CategoryId,
		}
	}

	return products, nil
}
