package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
)

var _ domain.OrderUseCase = (*OrderUseCase)(nil)

type OrderUseCase struct {
	orderRepo       domain.OrderRepository
	productProvider domain.ProductProvider
	paymentService  domain.PaymentService
}

func NewOrderUseCase(orderRepo domain.OrderRepository, productProvider domain.ProductProvider, paymentService domain.PaymentService) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		productProvider: productProvider,
		paymentService:  paymentService,
	}
}

func (s *OrderUseCase) CreateOrder(
	ctx context.Context,
	userID int64,
	items []domain.OrderItemInput,
) (domain.Order, string, error) {
	const op = "orderUseCase.CreateOrder"

	if len(items) == 0 {
		return domain.Order{}, "", fmt.Errorf("%s: no items", op)
	}

	orderItems, total, err := calculateOrderItems(ctx, items, s.productProvider)
	if err != nil {
		return domain.Order{}, "", fmt.Errorf("%s: failed to calculate order items: %w", op, err)
	}

	order := domain.Order{
		UUID:        uuid.New().String(),
		UserID:      userID,
		Status:      "pending",
		Items:       orderItems,
		TotalAmount: total,
		CreatedAt:   time.Now().UTC(),
	}

	if err = s.orderRepo.Create(ctx, order); err != nil {
		return domain.Order{}, "", fmt.Errorf("%s: failed to save order: %w", op, err)
	}

	// TODO: добавить ретрай/очередь на случай падения
	paymentURL, err := s.paymentService.CreatePayment(ctx, order)
	if err != nil {
		return domain.Order{}, "", fmt.Errorf("%s: failed to create payment: %w", op, err)
	}

	return order, paymentURL, nil
}

func (s *OrderUseCase) ListOrdersByUser(ctx context.Context, userID int64) ([]domain.Order, error) {
	return s.orderRepo.FindByUserID(ctx, userID)
}

func (s *OrderUseCase) GetOrderByUUID(ctx context.Context, orderID string) (domain.Order, error) {
	return s.orderRepo.FindByUUID(ctx, orderID)
}

func calculateOrderItems(
	ctx context.Context,
	items []domain.OrderItemInput,
	ps domain.ProductProvider,
) ([]domain.OrderItem, float64, error) {
	const op = "orderUseCase.calculateOrderItems"

	productIDs := make([]int64, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}

	products, err := ps.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: failed to get products info: %w", op, err)
	}

	// Проверка, что каталог-сервис вернул все запрошенные продукты
	if len(products) != len(productIDs) {
		return nil, 0, fmt.Errorf("%s: could not find all requested products", op)
	}

	productMap := make(map[int64]domain.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	orderItems := make([]domain.OrderItem, len(items))
	var total float64
	for i, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			// Эта проверка на случай, если каталог-сервис вернул не все товары,
			// что мы уже проверили выше, но это добавляет надежности.
			return nil, 0, fmt.Errorf("%s: product with id %d not found after batch call", op, item.ProductID)
		}

		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		total += product.Price * float64(item.Quantity)
	}

	return orderItems, total, nil
}
