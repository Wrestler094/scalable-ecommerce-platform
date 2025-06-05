package postgres

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"user-service/internal/domain"

	"pkg/cache"
)

var _ domain.UserRepository = (*cachedUserRepository)(nil)

type cachedUserRepository struct {
	dbRepo      domain.UserRepository
	cache       cache.Cache
	cacheTTL    time.Duration
	cachePrefix string
}

func NewCachedUserRepository(dbRepo domain.UserRepository, c cache.Cache) domain.UserRepository {
	return &cachedUserRepository{
		dbRepo:      dbRepo,
		cache:       c,
		cacheTTL:    60 * time.Minute,
		cachePrefix: "USER:",
	}
}

func (r *cachedUserRepository) CreateUser(ctx context.Context, user domain.UserWithPassword) (int64, error) {
	return r.dbRepo.CreateUser(ctx, user)
}

func (r *cachedUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	return r.dbRepo.GetUserByEmail(ctx, email)
}

func (r *cachedUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	cacheKey := r.cachePrefix + strconv.FormatInt(id, 10)

	// Попробуем взять из кэша
	data, err := r.cache.Get(ctx, cacheKey)
	if err == nil && data != nil {
		var user domain.User
		if err := json.Unmarshal(data, &user); err == nil {
			return &user, nil
		}
		// Если unmarshal не удался — инвалидируем кэш и пробуем из базы
		// TODO: log warning
		_ = r.invalidate(ctx, id)
	}

	// Из базы
	user, err := r.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Кэшируем
	raw, err := json.Marshal(user)
	if err == nil {
		_ = r.cache.Set(ctx, cacheKey, raw, r.cacheTTL)
	}

	return user, nil
}

func (r *cachedUserRepository) invalidate(ctx context.Context, id int64) error {
	return r.cache.Delete(ctx, r.cachePrefix+strconv.FormatInt(id, 10))
}
