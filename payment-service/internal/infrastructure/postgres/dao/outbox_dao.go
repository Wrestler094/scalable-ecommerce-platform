package dao

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"payment-service/internal/domain"
)

type OutboxEvent struct {
	ID        uuid.UUID       `db:"id"`
	EventType string          `db:"event_type"`
	CreatedAt time.Time       `db:"created_at"`
	Payload   json.RawMessage `db:"payload"`
}

func FromDomainEvent(e domain.OutboxEvent) (OutboxEvent, error) {
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

func (e OutboxEvent) ToDomainEvent() (domain.OutboxEvent, error) {
	var payload map[string]any
	err := json.Unmarshal(e.Payload, &payload)
	if err != nil {
		return domain.OutboxEvent{}, err
	}

	return domain.OutboxEvent{
		EventID:   e.ID,
		EventType: e.EventType,
		Timestamp: e.CreatedAt,
		Payload:   payload,
	}, nil
}

func FromDomainEventList(events []domain.OutboxEvent) ([]OutboxEvent, error) {
	result := make([]OutboxEvent, len(events))
	for i, e := range events {
		converted, err := FromDomainEvent(e)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func ToDomainEventList(events []OutboxEvent) ([]domain.OutboxEvent, error) {
	result := make([]domain.OutboxEvent, len(events))
	for i, e := range events {
		converted, err := e.ToDomainEvent()
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}
