package dao

import (
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
)

type DBOrder struct {
	ID          int64     `db:"id"`
	UUID        string    `db:"uuid"`
	UserID      int64     `db:"user_id"`
	Status      string    `db:"status"`
	TotalAmount float64   `db:"total_amount"`
	CreatedAt   time.Time `db:"created_at"`
}

type DBOrderItem struct {
	OrderID   int64   `db:"order_id"`
	ProductID int64   `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

// ======= Converots ========

func ToDBOrderItems(orderID int64, items []domain.OrderItem) []DBOrderItem {
	dbItems := make([]DBOrderItem, len(items))
	for i, item := range items {
		dbItems[i] = DBOrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	return dbItems
}

func ToDomainOrderItem(item DBOrderItem) domain.OrderItem {
	return domain.OrderItem{
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
	}
}
