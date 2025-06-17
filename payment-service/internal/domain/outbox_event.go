package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	EventID   uuid.UUID
	EventType string
	Timestamp time.Time
	Payload   map[string]any
}

type OutboxWriter interface {
	Write(ctx context.Context, evt OutboxEvent) error
}

type OutboxReader interface {
	FetchUnpublished(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
}

type EventProducer interface {
	Produce(ctx context.Context, topic string, eventType string, key uuid.UUID, timestamp time.Time, payload any) error
}
