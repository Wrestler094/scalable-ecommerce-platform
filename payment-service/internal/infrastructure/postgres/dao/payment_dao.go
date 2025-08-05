package dao

import (
	"time"

	"github.com/google/uuid"

	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/domain"
)

type Payment struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	UserID    int64     `db:"user_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

func FromDomainPayment(p domain.Payment) Payment {
	return Payment{
		OrderUUID: p.OrderUUID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt,
	}
}

func (p Payment) ToDomainPayment() domain.Payment {
	return domain.Payment{
		OrderUUID: p.OrderUUID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt,
	}
}
