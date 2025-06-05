package redis

import (
	"context"
	"fmt"
	"time"

	"user-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

var _ domain.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

type RefreshTokenRepository struct {
	client     *redis.Client
	prefix     string
	defaultTTL time.Duration
}

func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		client:     client,
		prefix:     "REFRESH",
		defaultTTL: 30 * 24 * time.Hour,
	}
}

func (r *RefreshTokenRepository) key(userID int64, token string) string {
	return fmt.Sprintf("%s:%d:%s", r.prefix, userID, token)
}

func (r *RefreshTokenRepository) Store(ctx context.Context, userID int64, token string) error {
	key := r.key(userID, token)
	return r.client.Set(ctx, key, userID, r.defaultTTL).Err()
}

func (r *RefreshTokenRepository) GetUserID(ctx context.Context, token string) (int64, error) {
	iter := r.client.Scan(ctx, 0, fmt.Sprintf("%s:*:%s", r.prefix, token), 1).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		return r.client.Get(ctx, key).Int64()
	}
	if err := iter.Err(); err != nil {
		return 0, err
	}
	return 0, redis.Nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	iter := r.client.Scan(ctx, 0, fmt.Sprintf("%s:*:%s", r.prefix, token), 1).Iterator()
	for iter.Next(ctx) {
		return r.client.Del(ctx, iter.Val()).Err()
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) Replace(
	ctx context.Context,
	oldToken string,
	newToken string,
	userID int64,
) error {
	if err := r.Delete(ctx, oldToken); err != nil {
		return err
	}
	return r.Store(ctx, userID, newToken)
}
