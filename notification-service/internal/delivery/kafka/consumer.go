package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"pkg/events"
	"pkg/logger"

	"notification-service/internal/domain"
)

type Consumer struct {
	reader  *kafka.Reader
	logger  logger.Logger
	usecase domain.NotificationUseCase
}

func NewConsumer(
	brokerAddresses []string,
	topic string,
	groupID string,
	uc domain.NotificationUseCase,
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

		var envelope events.Envelope[events.PaymentSuccessfulPayload]
		if err := json.Unmarshal(m.Value, &envelope); err != nil {
			log.WithError(err).Error("Failed to unmarshal envelope", "message_key", m.Key)
			continue
		}

		if envelope.EventType != events.EventPaymentSuccessful {
			log.Warn("Skipping unsupported event type", "event_type", envelope.EventType)
			continue
		}

		payload := envelope.Payload
		notif := domain.Notification{
			UserID:  payload.UserID,
			To:      "", // TODO: Implement email retrieval from Kafka or User Service
			Type:    domain.EmailNotification,
			Subject: "Ваш платёж прошёл успешно",
			Message: fmt.Sprintf("Спасибо за оплату заказа %s на сумму %.2f₽", payload.OrderID, payload.Amount),
		}

		if err := c.usecase.Send(notif); err != nil {
			log.WithError(err).Error("Failed to send notification", "message_key", m.Key, "event_id", envelope.EventID)
			continue
		}

		log.Info("Notification sent", "order_id", payload.OrderID)
	}
}

func (c *Consumer) Close() error {
	const op = "kafka.Consumer.Close"

	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("%s: failed to close Kafka reader: %w", op, err)
	}

	return nil
}
