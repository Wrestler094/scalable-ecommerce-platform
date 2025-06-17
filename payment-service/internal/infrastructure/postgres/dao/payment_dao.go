package dao

import (
	"time"

	"payment-service/internal/domain"
)

type Payment struct {
	OrderID   int64     `db:"order_id"`
	UserID    int64     `db:"user_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

func FromDomainPayment(p domain.Payment) Payment {
	return Payment{
		OrderID:   p.OrderID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt,
	}
}

func (p Payment) ToDomainPayment() domain.Payment {
	return domain.Payment{
		OrderID:   p.OrderID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt,
	}
}
