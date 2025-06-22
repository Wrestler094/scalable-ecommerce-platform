package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"payment-service/internal/domain"
	"pkg/events"
)

var _ domain.EventProducer = (*Producer)(nil)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddresses []string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddresses...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) Produce(
	ctx context.Context,
	topic string,
	eventType string,
	key uuid.UUID,
	timestamp time.Time,
	payload any,
) error {
	const op = "kafka.Produce"

	envelope := events.Envelope[any]{
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

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("%s: failed to write message: %w", op, err)
	}

	return nil
}

func (p *Producer) Close() error {
	const op = "kafka.Close"

	err := p.writer.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close writer: %w", op, err)
	}

	return nil
}
