package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/events"

	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/domain"
)

var _ domain.EventProducer[events.PaymentSuccessfulPayload] = (*Producer[events.PaymentSuccessfulPayload])(nil)

type Producer[T any] struct {
	writer *kafka.Writer
}

func NewProducer[T any](brokerAddresses []string) *Producer[T] {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddresses...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer[T]{writer: writer}
}

func (p *Producer[T]) Produce(
	ctx context.Context,
	topic string,
	eventType string,
	key uuid.UUID,
	timestamp time.Time,
	payload T,
) error {
	const op = "kafka.Produce"

	envelope := events.Envelope[T]{
		EventID:   key,
		EventType: eventType,
		Timestamp: timestamp.Format(time.RFC3339),
		Payload:   payload,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal envelope: %w", op, err)
	}

	msg := kafka.Message{
		Key:   []byte(key.String()),
		Value: data,
		Topic: topic,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer[T]) Close() error {
	const op = "kafka.Close"

	err := p.writer.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close writer: %w", op, err)
	}

	return nil
}
