package usecase

import (
	"context"
	"fmt"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
)

var _ domain.OrderPaymentUseCase = (*PaymentUseCase)(nil)

type PaymentUseCase struct {
	orderPaymentRepo domain.OrderPaymentRepository
}

func NewPaymentUseCase(orderPaymentRepo domain.OrderPaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{
		orderPaymentRepo: orderPaymentRepo,
	}
}

func (u *PaymentUseCase) MarkOrderAsPaid(ctx context.Context, orderUUID string) error {
	const op = "paymentUseCase.MarkOrderAsPaid"

	err := u.orderPaymentRepo.MarkAsPaid(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("%s: failed to mark order as paid: %w", op, err)
	}

	return nil
}
