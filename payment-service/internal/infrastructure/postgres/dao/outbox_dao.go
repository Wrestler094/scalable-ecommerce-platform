package dao

import (
	"encoding/json"
	"time"

	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID        uuid.UUID       `db:"id"`
	EventType string          `db:"event_type"`
	CreatedAt time.Time       `db:"created_at"`
	Payload   json.RawMessage `db:"payload"`
}

func FromDomainEvent[T any](e domain.OutboxEvent[T]) (OutboxEvent, error) {
	payloadJSON, err := json.Marshal(e.Payload)
	if err != nil {
		return OutboxEvent{}, err
	}

	return OutboxEvent{
		ID:        e.EventID,
		EventType: e.EventType,
		CreatedAt: e.Timestamp,
		Payload:   payloadJSON,
	}, nil
}

func ToDomainEvent[T any](e OutboxEvent) (domain.OutboxEvent[T], error) {
	var payload T
	err := json.Unmarshal(e.Payload, &payload)
	if err != nil {
		return domain.OutboxEvent[T]{}, err
	}

	return domain.OutboxEvent[T]{
		EventID:   e.ID,
		EventType: e.EventType,
		Timestamp: e.CreatedAt,
		Payload:   payload,
	}, nil
}

func ToDomainEventList[T any](events []OutboxEvent) ([]domain.OutboxEvent[T], error) {
	result := make([]domain.OutboxEvent[T], len(events))
	for i, e := range events {
		converted, err := ToDomainEvent[T](e)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}
