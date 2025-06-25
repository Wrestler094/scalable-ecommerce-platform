package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"payment-service/internal/domain"
	"payment-service/internal/infrastructure/postgres/dao"
	"payment-service/internal/infrastructure/txmanager"
)

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p domain.Payment) error {
	const op = "paymentRepository.Create"
	const query = `
		INSERT INTO payments (order_uuid, user_id, amount, created_at)
		VALUES ($1, $2, $3, $4)
	`

	daoPayment := dao.FromDomainPayment(p)

	// Пытаемся получить транзакцию из контекста
	if tx, ok := txmanager.ExtractTx(ctx); ok {
		_, err := tx.ExecContext(ctx, query, daoPayment.OrderUUID, daoPayment.UserID, daoPayment.Amount, daoPayment.CreatedAt)
		if err != nil {
			return fmt.Errorf("%s: failed to create payment in transaction: %w", op, err)
		}
		return nil
	}

	// Если транзакции нет - используем обычное соединение
	_, err := r.db.ExecContext(ctx, query, daoPayment.OrderUUID, daoPayment.UserID, daoPayment.Amount, daoPayment.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: failed to create payment: %w", op, err)
	}
	return nil
}
