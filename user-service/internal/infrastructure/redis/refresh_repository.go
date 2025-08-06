package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/domain"
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
	const op = "RefreshTokenRepository.Store"

	key := r.key(userID, token)
	if err := r.client.Set(ctx, key, userID, r.defaultTTL).Err(); err != nil {
		return fmt.Errorf("%s: failed to set refresh token: %w", op, err)
	}

	return nil
}

func (r *RefreshTokenRepository) GetUserID(ctx context.Context, token string) (int64, error) {
	const op = "RefreshTokenRepository.GetUserID"

	iter := r.client.Scan(ctx, 0, fmt.Sprintf("%s:*:%s", r.prefix, token), 1).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		userID, err := r.client.Get(ctx, key).Int64()
		if err != nil {
			return 0, fmt.Errorf("%s: failed to get user ID: %w", op, err)
		}
		return userID, nil
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("%s: failed to scan keys: %w", op, err)
	}

	return 0, redis.Nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	const op = "RefreshTokenRepository.Delete"

	iter := r.client.Scan(ctx, 0, fmt.Sprintf("%s:*:%s", r.prefix, token), 1).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("%s: failed to delete token: %w", op, err)
		}
		return nil
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("%s: failed to scan keys: %w", op, err)
	}

	return nil
}

func (r *RefreshTokenRepository) Replace(ctx context.Context, oldToken string, newToken string, userID int64) error {
	const op = "RefreshTokenRepository.Replace"

	if err := r.Delete(ctx, oldToken); err != nil {
		return fmt.Errorf("%s: failed to delete old token: %w", op, err)
	}

	if err := r.Store(ctx, userID, newToken); err != nil {
		return fmt.Errorf("%s: failed to store new token: %w", op, err)
	}

	return nil
}
