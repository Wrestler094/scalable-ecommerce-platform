package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewClient(addr, password string, db int) (*Redis, error) {
	const op = "redis.NewClient"

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to connect to redis: %w", op, err)
	}

	return &Redis{Client: rdb}, nil
}

func (r *Redis) Close() error {
	const op = "redis.Close"

	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("%s: failed to close redis client: %w", op, err)
	}

	return nil
}
