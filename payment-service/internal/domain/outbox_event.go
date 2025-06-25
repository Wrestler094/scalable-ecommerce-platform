package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent[T any] struct {
	EventID   uuid.UUID
	EventType string
	Timestamp time.Time
	Payload   T
}

type OutboxWriter[T any] interface {
	Write(ctx context.Context, evt OutboxEvent[T]) error
}

type OutboxReader[T any] interface {
	FetchUnpublished(ctx context.Context, limit int) ([]OutboxEvent[T], error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
}

type EventProducer[T any] interface {
	Produce(ctx context.Context, topic string, eventType string, key uuid.UUID, timestamp time.Time, payload T) error
}
