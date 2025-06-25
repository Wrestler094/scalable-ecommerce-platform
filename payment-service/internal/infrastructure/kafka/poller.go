package kafka

import (
	"context"
	"time"

	"pkg/events"
	"pkg/logger"

	"payment-service/internal/domain"
)

type Poller struct {
	reader   domain.OutboxReader[events.PaymentSuccessfulPayload]
	producer domain.EventProducer[events.PaymentSuccessfulPayload]
	logger   logger.Logger
	topic    string
	interval time.Duration
	batch    int
}

func NewPoller(
	reader domain.OutboxReader[events.PaymentSuccessfulPayload],
	producer domain.EventProducer[events.PaymentSuccessfulPayload],
	logger logger.Logger,
	topic string,
	interval time.Duration,
	batch int,
) *Poller {
	return &Poller{
		reader:   reader,
		producer: producer,
		logger:   logger,
		interval: interval,
		batch:    batch,
		topic:    topic,
	}
}

func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.process(ctx)
		}
	}
}

func (p *Poller) process(ctx context.Context) {
	const op = "kafka.Poller.process"

	unpublishedEvents, err := p.reader.FetchUnpublished(ctx, p.batch)
	if err != nil {
		p.logger.WithOp(op).WithError(err).Error("failed to fetch events")
		return
	}

	for _, evt := range unpublishedEvents {
		if err := p.producer.Produce(ctx, p.topic, evt.EventType, evt.EventID, evt.Timestamp, evt.Payload); err != nil {
			p.logger.WithOp(op).WithError(err).Error("failed to produce events", "event", evt.EventID)
			continue
		}

		if err := p.reader.MarkPublished(ctx, evt.EventID); err != nil {
			p.logger.WithOp(op).WithError(err).Error("failed to mark published", "event", evt.EventID)
		}
	}
}
