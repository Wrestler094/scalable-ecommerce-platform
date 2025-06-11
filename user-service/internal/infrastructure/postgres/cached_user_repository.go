package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"pkg/logger"
	"user-service/internal/domain"

	"pkg/cache"
)

var _ domain.UserRepository = (*cachedUserRepository)(nil)

type cachedUserRepository struct {
	dbRepo      domain.UserRepository
	cache       cache.Cache
	cacheTTL    time.Duration
	cachePrefix string
	logger      logger.Logger
}

func NewCachedUserRepository(dbRepo domain.UserRepository, c cache.Cache, logger logger.Logger) domain.UserRepository {
	return &cachedUserRepository{
		dbRepo:      dbRepo,
		cache:       c,
		cacheTTL:    60 * time.Minute,
		cachePrefix: "USER:",
		logger:      logger,
	}
}

func (r *cachedUserRepository) CreateUser(ctx context.Context, user domain.UserWithPassword) (int64, error) {
	return r.dbRepo.CreateUser(ctx, user)
}

func (r *cachedUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	return r.dbRepo.GetUserByEmail(ctx, email)
}

func (r *cachedUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	const op = "cachedUserRepository.GetUserByID"

	log := r.logger.WithOp(op).WithUserID(id).WithRequestID(middleware.GetReqID(ctx))

	cacheKey := r.cachePrefix + strconv.FormatInt(id, 10)

	// Попробуем взять из кэша
	data, err := r.cache.Get(ctx, cacheKey)
	if err == nil && data != nil {
		var user domain.User
		if err := json.Unmarshal(data, &user); err == nil {
			return &user, nil
		}

		// Если unmarshal не удался — инвалидируем кэш и пробуем из базы
		log.WithError(err).Warn("failed to unmarshal user from cache", "cache_key", cacheKey)

		if err := r.invalidate(ctx, id); err != nil {
			log.WithError(err).Warn("failed to invalidate user cache", "cache_key", cacheKey)
		}
	}

	// Из базы
	user, err := r.dbRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user from DB: %w", op, err)
	}

	raw, err := json.Marshal(user)
	if err != nil {
		log.WithError(err).Error("failed to marshal user for cache")
	}

	if err == nil {
		if err := r.cache.Set(ctx, cacheKey, raw, r.cacheTTL); err != nil {
			log.WithError(err).Error("failed to set user to cache", "cache_key", cacheKey)
		}
	}

	return user, nil
}

func (r *cachedUserRepository) invalidate(ctx context.Context, id int64) error {
	const op = "cachedUserRepository.invalidate"

	cacheKey := r.cachePrefix + strconv.FormatInt(id, 10)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		return fmt.Errorf("%s: failed to delete cache: %w", op, err)
	}

	return nil
}
