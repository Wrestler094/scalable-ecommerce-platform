package dto

import (
	"order-service/internal/domain"
	"time"
)

type OrderItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	UUID        string      `json:"uuid"`
	UserID      int64       `json:"user_id"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
}

// ====== CreateOrder ======

type CreateOrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItem `json:"items"`
}

type CreateOrderResponse struct {
	Order      Order  `json:"order"`
	PaymentURL string `json:"payment_url"`
}

// ====== ListOrders ======

type GetOrdersListResponse struct {
	Orders []Order `json:"orders"`
}

// ====== GetOrderByID ======

type GetOrderByIDResponse struct {
	Order Order `json:"order"`
}

// ====== Convertors ======

func (r CreateOrderRequest) ToDomainItems() []domain.OrderItemInput {
	items := make([]domain.OrderItemInput, len(r.Items))
	for i, item := range r.Items {
		items[i] = domain.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return items
}

func FromOrders(orders []domain.Order) []Order {
	result := make([]Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, FromOrder(o))
	}
	return result
}

func FromOrder(o domain.Order) Order {
	return Order{
		UUID:        o.UUID,
		UserID:      o.UserID,
		Status:      o.Status,
		Items:       fromOrderItems(o.Items),
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
	}
}

func fromOrderItems(items []domain.OrderItem) []OrderItem {
	res := make([]OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	return res
}
