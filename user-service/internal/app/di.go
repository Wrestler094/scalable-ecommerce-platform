package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/adapters"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/cache"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/validator"

	httpdelivery "github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/delivery/http/infra"
	v1 "github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/delivery/http/v1"

	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/infrastructure/hasher"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/infrastructure/idgenerator"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/infrastructure/jwt"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/infrastructure/postgres"
	redisinfra "github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/infrastructure/redis"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/usecase"
)

/* ======= Cleanup ======= */
type cleanups struct{ fns []func() error }

func (c *cleanups) add(fn func() error) { c.fns = append(c.fns, fn) }
func (c *cleanups) run() error {
	var all error
	for i := len(c.fns) - 1; i >= 0; i-- {
		if err := c.fns[i](); err != nil {
			all = errors.Join(all, err)
		}
	}
	return all
}

// InitDI собирает все зависимости приложения.
// Возвращает HTTP-сервер, cleanup и ошибку.
func InitDI(cfg *config.Config, baseLogger logger.Logger, healthManager healthcheck.Manager) (
	*httpserver.Server,
	func() error,
	error,
) {
	var cl cleanups
	cleanup := func() error { return cl.run() }

	// Connect Postgres
	shardRouter, err := postgres.NewShardRouter(cfg.PGShards)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to init postgres: %w", err)
	}
	cl.add(func() error {
		shardRouter.Close()
		return nil
	})

	// Connect Redis
	rdb, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to init redis: %w", err)
	}
	cl.add(func() error {
		rdb.Close()
		return nil
	})

	// Helpers/Deps
	idGen, err := idgenerator.NewSnowflakeGenerator(cfg.Snowflake.NodeID, cfg.Snowflake.Epoch)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to init snowflake: %w", err)
	}
	tokenMgr := jwt.NewManager(cfg.JWT.AccessSecret, time.Duration(cfg.JWT.TokenTTL)*time.Second)
	passHasher := hasher.NewBcryptHasher()

	rawVal := validator.NewPlaygroundValidator()
	httpVal := adapters.NewHttpValidatorAdapter(rawVal)

	redisCache := cache.NewRedisCache(rdb.Client)

	// Repositories
	userRepo := postgres.NewUserRepository(shardRouter, idGen)
	cachedUserRepo := postgres.NewCachedUserRepository(userRepo, redisCache, baseLogger)
	refreshRepo := redisinfra.NewRefreshTokenRepository(rdb.Client)

	// Use-Cases
	userUC := usecase.NewUserUseCase(cachedUserRepo, refreshRepo, tokenMgr, passHasher)

	// Handlers
	userH := v1.NewUserHandler(userUC, httpVal, baseLogger)
	monH := infra.NewMonitoringHandler(healthManager)

	// Router
	router := httpdelivery.NewRouter(httpdelivery.Handlers{
		V1Handlers: v1.Handlers{
			UserHandler: userH,
		},
		MonitoringHandler: monH,
	})

	// HTTP Server
	server := httpserver.NewServer(
		httpserver.Port(fmt.Sprintf(":%d", cfg.HTTP.Port)),
		httpserver.Handler(router),
	)

	return server, cleanup, nil
}
