package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/domain"
)

var _ domain.CartRepository = (*redisCartRepo)(nil)

type redisCartRepo struct {
	rdb    *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisCartRepo(rdb *redis.Client) domain.CartRepository {
	return &redisCartRepo{
		rdb:    rdb,
		prefix: "CART",
		ttl:    30 * 24 * time.Hour,
	}
}

func (r *redisCartRepo) key(userID int64) string {
	return fmt.Sprintf("%s:%d", r.prefix, userID)
}

func (r *redisCartRepo) Get(ctx context.Context, userID int64) ([]domain.CartItem, error) {
	data, err := r.rdb.HGetAll(ctx, r.key(userID)).Result()
	if err != nil {
		return nil, err
	}

	var items []domain.CartItem
	for pid, qty := range data {
		pidInt, _ := strconv.ParseInt(pid, 10, 64)
		qtyInt, _ := strconv.Atoi(qty)
		items = append(items, domain.CartItem{ProductID: pidInt, Quantity: qtyInt})
	}

	return items, nil
}

func (r *redisCartRepo) Add(ctx context.Context, userID, productID int64, quantity int) error {
	pipe := r.rdb.TxPipeline()
	pipe.HIncrBy(ctx, r.key(userID), strconv.FormatInt(productID, 10), int64(quantity))
	pipe.Expire(ctx, r.key(userID), r.ttl)
	_, err := pipe.Exec(ctx)

	return err
}

func (r *redisCartRepo) Update(ctx context.Context, userID, productID int64, quantity int) error {
	pipe := r.rdb.TxPipeline()
	if quantity > 0 {
		pipe.HSet(ctx, r.key(userID), strconv.FormatInt(productID, 10), quantity)
	} else {
		pipe.HDel(ctx, r.key(userID), strconv.FormatInt(productID, 10))
	}
	pipe.Expire(ctx, r.key(userID), r.ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisCartRepo) Remove(ctx context.Context, userID, productID int64) error {
	return r.rdb.HDel(ctx, r.key(userID), strconv.FormatInt(productID, 10)).Err()
}

func (r *redisCartRepo) Clear(ctx context.Context, userID int64) error {
	return r.rdb.Del(ctx, r.key(userID)).Err()
}
