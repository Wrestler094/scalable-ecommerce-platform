package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/events"

	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/domain"
)

type PaymentUseCase struct {
	paymentRepo     domain.PaymentRepository
	outboxWriter    domain.OutboxWriter[events.PaymentSuccessfulPayload]
	idempotencyRepo domain.IdempotencyRepository
	txManager       domain.TxManager
}

func NewPaymentUseCase(
	paymentRepo domain.PaymentRepository,
	outbox domain.OutboxWriter[events.PaymentSuccessfulPayload],
	idempotencyRepo domain.IdempotencyRepository,
	txManager domain.TxManager,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:     paymentRepo,
		outboxWriter:    outbox,
		idempotencyRepo: idempotencyRepo,
		txManager:       txManager,
	}
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, cmd domain.PayCommand) error {
	const op = "paymentUseCase.ProcessPayment"

	// 1. Проверка идемпотентности ДО транзакции
	exists, err := uc.idempotencyRepo.Exists(ctx, cmd.IdempotencyKey)
	if err != nil {
		return fmt.Errorf("%s: failed to check idempotency: %w", op, err)
	}
	if exists {
		return fmt.Errorf("%s: idempotency key already used, %w", op, domain.ErrDuplicatePayment)
	}

	// 2. Всё бизнес-действие — внутри транзакции
	err = uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		payment := domain.Payment{
			OrderUUID: cmd.OrderUUID,
			UserID:    cmd.UserID,
			Amount:    cmd.Amount,
			CreatedAt: time.Now().UTC(),
		}

		if err = uc.paymentRepo.Create(txCtx, payment); err != nil {
			return fmt.Errorf("%s: failed to create payment: %w", op, err)
		}

		// TODO: Подумать над тем, что передавать
		event := domain.OutboxEvent[events.PaymentSuccessfulPayload]{
			EventID:   uuid.New(),
			EventType: events.EventPaymentSuccessful,
			Timestamp: time.Now(),
			Payload: events.PaymentSuccessfulPayload{
				OrderUUID: cmd.OrderUUID.String(),
				UserID:    cmd.UserID,
				Amount:    cmd.Amount,
			},
		}

		if err = uc.outboxWriter.Write(txCtx, event); err != nil {
			return fmt.Errorf("%s: failed to write payment event to outbox: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// 3. Регистрируем идемпотентность в Redis (вне транзакции)
	if err := uc.idempotencyRepo.Register(ctx, cmd.IdempotencyKey); err != nil {
		return fmt.Errorf("%s: failed to register idempotency key: %w", op, domain.ErrIdempotencyRegistrationFailed)
	}

	return nil
}
