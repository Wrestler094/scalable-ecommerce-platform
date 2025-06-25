package usecase

import (
	"context"
	"fmt"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
)

var _ domain.OrderUseCase = (*OrderUseCase)(nil)

type OrderUseCase struct {
	orderRepo      domain.OrderRepository
	productService domain.ProductService
	paymentService domain.PaymentService
}

func NewOrderUseCase(orderRepo domain.OrderRepository, productService domain.ProductService, paymentService domain.PaymentService) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		productService: productService,
		paymentService: paymentService,
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

	orderItems, total, err := calculateOrderItems(ctx, items, s.productService)
	if err != nil {
		return domain.Order{}, "", fmt.Errorf("%s: failed to calculate order items: %w", op, err)
	}

	order := domain.Order{
		UUID:        uuid.New().String(),
		UserID:      userID,
		Status:      "pending",
		Items:       orderItems,
		TotalAmount: total,
		CreatedAt:   time.Now(),
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
	ps domain.ProductService,
) ([]domain.OrderItem, float64, error) {
	const op = "usecase.calculateOrderItems"

	orderItems := make([]domain.OrderItem, 0, len(items))
	var total float64

	// TODO: Сделать получение батчем
	for _, item := range items {
		price, err := ps.GetPrice(ctx, item.ProductID)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: failed to get price for product %d: %w", op, item.ProductID, err)
		}

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})

		total += price * float64(item.Quantity)
	}

	return orderItems, total, nil
}
