package domain

import (
	"context"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, order Order) (string, error)
}
