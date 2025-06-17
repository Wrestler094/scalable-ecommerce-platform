package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"payment-service/internal/domain"
	"payment-service/internal/infrastructure/postgres/dao"
	"payment-service/internal/infrastructure/txmanager"
)

var _ domain.OutboxWriter = (*OutboxRepository)(nil)
var _ domain.OutboxReader = (*OutboxRepository)(nil)

type OutboxRepository struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Write(ctx context.Context, evt domain.OutboxEvent) error {
	const op = "outboxRepository.Write"
	const query = `
        INSERT INTO outbox (id, event_type, payload, created_at)
        VALUES ($1, $2, $3, $4)
    `

	daoEvent, err := dao.FromDomainEvent(evt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Пытаемся получить транзакцию из контекста
	if tx, ok := txmanager.ExtractTx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, daoEvent.ID, daoEvent.EventType, daoEvent.Payload, daoEvent.CreatedAt)
		if err != nil {
			return fmt.Errorf("%s: failed to insert outbox event in transaction: %w", op, err)
		}
		return nil
	}

	// Если транзакции нет - используем обычное соединение
	_, err = r.db.ExecContext(ctx, query, daoEvent.ID, daoEvent.EventType, daoEvent.Payload, daoEvent.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: failed to insert outbox event: %w", op, err)
	}

	return nil
}

func (r *OutboxRepository) FetchUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	const op = "outboxRepository.FetchUnpublished"

	const query = `
		SELECT id, event_type, payload, created_at
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
	`

	var daoEvents []dao.OutboxEvent
	if err := r.db.SelectContext(ctx, &daoEvents, query, limit); err != nil {
		return nil, fmt.Errorf("%s: failed to fetch events: %w", op, err)
	}

	domainEvents, err := dao.ToDomainEventList(daoEvents)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to convert events: %w", op, err)
	}

	return domainEvents, nil
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id uuid.UUID) error {
	const op = "outboxRepository.MarkPublished"
	const query = `
		UPDATE outbox
		SET published_at = now()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: failed to mark event as published: %w", op, err)
	}

	return nil
}
