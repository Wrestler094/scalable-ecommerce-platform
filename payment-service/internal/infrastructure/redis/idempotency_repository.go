package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"payment-service/internal/domain"
)

var _ domain.IdempotencyRepository = (*IdempotencyRepository)(nil)

type IdempotencyRepository struct {
	client     *redis.Client
	prefix     string
	defaultTTL time.Duration
}

func NewIdempotencyRepository(client *redis.Client) *IdempotencyRepository {
	return &IdempotencyRepository{
		client:     client,
		prefix:     "IDEMP_KEY",
		defaultTTL: 30 * 24 * time.Hour,
	}
}

func (r *IdempotencyRepository) key(raw string) string {
	return fmt.Sprintf("%s:%s", r.prefix, raw)
}

func (r *IdempotencyRepository) Exists(ctx context.Context, key string) (bool, error) {
	const op = "idempotencyRepository.Exists"

	val, err := r.client.Exists(ctx, r.key(key)).Result()
	if err != nil {
		return false, fmt.Errorf("%s: failed to check existence in Redis: %w", op, err)
	}

	return val == 1, nil
}

func (r *IdempotencyRepository) Register(ctx context.Context, key string) error {
	const op = "idempotencyRepository.Register"

	ok, err := r.client.SetNX(ctx, r.key(key), "1", r.defaultTTL).Result()
	if err != nil {
		return fmt.Errorf("%s: failed to set idempotency key in Redis: %w", op, err)
	}

	// Если ключ уже существовал — значит кто-то успел записать
	if !ok {
		return fmt.Errorf("%s: failed to register idempotency key — already exists", op)
	}

	return nil
}
