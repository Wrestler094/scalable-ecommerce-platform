package domain

import (
	"context"
	"time"
)

type Payment struct {
	OrderID int64
	UserID  int64
	Amount  float64
	// TODO: Подумать нужен ли тут CreatedAt
	CreatedAt time.Time
}

type PayCommand struct {
	UserID         int64
	OrderID        int64
	Amount         float64
	IdempotencyKey string
}

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, cmd PayCommand) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment Payment) error
}

type IdempotencyRepository interface {
	Exists(ctx context.Context, key string) (bool, error)
	Register(ctx context.Context, key string) error
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
