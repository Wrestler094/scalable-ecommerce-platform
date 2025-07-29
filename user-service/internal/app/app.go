package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pkg/adapters"
	"pkg/cache"
	"pkg/healthcheck"
	"pkg/httpserver"
	"pkg/logger"
	"pkg/validator"
	"user-service/internal/infrastructure/idgenerator"

	"user-service/internal/config"
	"user-service/internal/delivery/http"
	"user-service/internal/infrastructure/hasher"
	"user-service/internal/infrastructure/jwt"
	"user-service/internal/infrastructure/postgres"
	redisinfra "user-service/internal/infrastructure/redis"
	"user-service/internal/usecase"
)

// Run creates objects via constructors and starts the application.
func Run(cfg *config.Config) {
	start := time.Now()

	baseLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("app.Run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	// Connect Postgres
	shardRouter, err := postgres.NewShardRouter(cfg.PGShards)
	if err != nil {
		runLogger.Fatal("DB initialization failed", "error", err)
	}
	defer shardRouter.Close()

	runLogger.Info("PostgreSQL connected", "shards", len(cfg.PGShards))

	// Connect Redis
	rdb, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		runLogger.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	runLogger.Info("Redis connected", "addr", cfg.Redis.Addr)

	// Helpers/Deps
	idGenerator, err := idgenerator.NewSnowflakeGenerator(cfg.Snowflake.NodeID, cfg.Snowflake.Epoch)
	if err != nil {
		runLogger.Fatal("Failed to init Snowflake ID generator", "error", err)
	}
	tokenManager := jwt.NewManager(cfg.JWT.AccessSecret, time.Duration(cfg.JWT.TokenTTL)*time.Second)
	passwordHasher := hasher.NewBcryptHasher()
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	redisCache := cache.NewRedisCache(rdb.Client)
	healthManager := healthcheck.NewManager()

	// Repositories
	userRepo := postgres.NewUserRepository(shardRouter, idGenerator)

	cachedUserRepo := postgres.NewCachedUserRepository(userRepo, redisCache, baseLogger)
	refreshRepo := redisinfra.NewRefreshTokenRepository(rdb.Client)

	// Use-Cases
	userUseCase := usecase.NewUserUseCase(cachedUserRepo, refreshRepo, tokenManager, passwordHasher)

	// Handlers
	userHandler := http.NewUserHandler(userUseCase, httpValidator, baseLogger)
	monitoringHandler := http.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		UserHandler:       userHandler,
		MonitoringHandler: monitoringHandler,
	})

	// HTTP Server
	httpServer := httpserver.NewServer(
		httpserver.Port(fmt.Sprintf(":%d", cfg.HTTP.Port)),
		httpserver.Handler(router),
	)

	runLogger.Info("HTTP server is starting", "port", cfg.HTTP.Port)

	// Start server
	if err := httpServer.Start(); err != nil {
		runLogger.WithError(err).Fatal("HTTP server failed to start")
	}

	healthManager.SetReady(true)

	runLogger.Info("Startup complete", logger.LogKeyDurationMS, time.Since(start).String())

	// Graceful shutdown handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		runLogger.Info("Received shutdown signal", "signal", sig)
	case err := <-httpServer.Notify():
		runLogger.WithError(err).Error("HTTP server reported error")
	}

	healthManager.SetReady(false)
	healthManager.SetAlive(false)

	if err := httpServer.Shutdown(); err != nil {
		runLogger.WithError(err).Error("HTTP server shutdown failed")
	} else {
		runLogger.Info("HTTP server gracefully stopped")
	}
}
