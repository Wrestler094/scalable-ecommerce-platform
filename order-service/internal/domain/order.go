package domain

import (
	"context"
	"time"
)

// Domain entities

type OrderItem struct {
	ProductID int64
	Quantity  int
	Price     float64
}

type Order struct {
	ID          int64
	UUID        string
	UserID      int64
	Status      string
	Items       []OrderItem
	TotalAmount float64
	CreatedAt   time.Time
}

// Services

type ProductService interface {
	GetPrice(ctx context.Context, productID int64) (float64, error)
}

type PaymentService interface {
	CreatePayment(ctx context.Context, order Order) (string, error)
}

// UseCases

type OrderItemInput struct {
	ProductID int64
	Quantity  int
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID int64, items []OrderItemInput) (Order, string, error)
	ListOrdersByUser(ctx context.Context, userID int64) ([]Order, error)
	GetOrderByUUID(ctx context.Context, uuid string) (Order, error)
}

type OrderPaymentUseCase interface {
	MarkOrderAsPaid(ctx context.Context, uuid string) error
}

// Repositories

type OrderRepository interface {
	Create(ctx context.Context, order Order) error
	FindByUUID(ctx context.Context, uuid string) (Order, error)
	FindByUserID(ctx context.Context, userID int64) ([]Order, error)
}

type OrderPaymentRepository interface {
	MarkAsPaid(ctx context.Context, orderUUID string) error
}
