package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"pkg/events"
	"pkg/logger"

	"order-service/internal/delivery/kafka/dto"
	"order-service/internal/domain"
)

type Consumer struct {
	reader  *kafka.Reader
	logger  logger.Logger
	usecase domain.OrderPaymentUseCase
}

func NewConsumer(
	brokerAddresses []string,
	topic string,
	groupID string,
	uc domain.OrderPaymentUseCase,
	logger logger.Logger,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddresses,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.LastOffset,
	})

	return &Consumer{
		reader:  reader,
		logger:  logger,
		usecase: uc,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	const op = "kafka.Consumer.Start"

	log := c.logger.WithOp(op)
	log.Info("Kafka consumer started", "topic", c.reader.Config().Topic)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info("Kafka consumer stopped by context")
				return
			}

			log.WithError(err).Error("Failed to read message")
			continue
		}

		var envelope events.Envelope[dto.PaymentPayload]
		if err := json.Unmarshal(m.Value, &envelope); err != nil {
			log.WithError(err).Error("Failed to unmarshal envelope", "message_key", m.Key)
			continue
		}

		if envelope.EventType != events.EventPaymentSuccessful {
			log.Warn("Skipping unsupported event type", "type", envelope.EventType)
			continue
		}

		if err := c.usecase.MarkOrderAsPaid(ctx, envelope.Payload.OrderUUID); err != nil {
			log.WithError(err).Error("Failed to mark order as paid", "message_key", m.Key, "event_id", envelope.EventID)
			continue
		}

		log.Info("Order marked as paid", "order_id", envelope.Payload.OrderUUID, "event_id", envelope.EventID)
	}
}

func (c *Consumer) Close() error {
	const op = "kafka.Consumer.Close"

	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("%s: failed to close Kafka reader: %w", op, err)
	}

	return nil
}
