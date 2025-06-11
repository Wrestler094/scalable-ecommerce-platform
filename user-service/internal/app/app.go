package app

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pkg/cache"
	"pkg/healthcheck"
	"pkg/httpserver"
	"pkg/logger"
	"pkg/validator"

	"user-service/internal/config"
	"user-service/internal/delivery/http"
	"user-service/internal/infrastructure/hasher"
	"user-service/internal/infrastructure/jwt"
	"user-service/internal/infrastructure/postgres"
	redisinfra "user-service/internal/infrastructure/redis"
	"user-service/internal/usecase"

	"pkg/adapters"
)

// Run creates objects via constructors and starts the application.
func Run(cfg *config.Config) {
	start := time.Now()

	l, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	l = l.WithOp("app.Run")
	l.Info("Logger initialized", "level", cfg.Log.Level)

	// Connect Postgres
	pg, err := postgres.NewConnect(cfg.PG.URL)
	if err != nil {
		l.Fatal("DB initialization failed", "error", err)
	}
	defer pg.Close()

	u, _ := url.Parse(cfg.PG.URL)
	l.Info("PostgreSQL connected", "host", u.Hostname(), "port", u.Port())

	// Connect Redis
	rdb, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	l.Info("Redis connected", "addr", cfg.Redis.Addr)

	// Helpers/Deps
	tokenManager := jwt.NewManager(cfg.JWT.AccessSecret, time.Duration(cfg.JWT.TokenTTL)*time.Second)
	passwordHasher := hasher.NewBcryptHasher()
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	redisCache := cache.NewRedisCache(rdb.Client)
	healthManager := healthcheck.NewManager()

	// Repositories
	userRepo := postgres.NewUserRepository(pg.DB)
	cachedUserRepo := postgres.NewCachedUserRepository(userRepo, redisCache, l)
	refreshRepo := redisinfra.NewRefreshTokenRepository(rdb.Client)

	// Use-Cases
	userUseCase := usecase.NewUserUseCase(cachedUserRepo, refreshRepo, tokenManager, passwordHasher)

	// Handlers
	userHandler := http.NewUserHandler(userUseCase, httpValidator, l)
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

	l.Info("HTTP server is starting", "port", cfg.HTTP.Port)

	// Start server
	if err := httpServer.Start(); err != nil {
		l.WithError(err).Fatal("HTTP server failed to start")
	}

	healthManager.SetReady(true)

	l.Info("Startup complete", logger.LogKeyDurationMS, time.Since(start).String())

	// Graceful shutdown handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("Received shutdown signal", "signal", sig)
	case err := <-httpServer.Notify():
		l.WithError(err).Error("HTTP server reported error")
	}

	healthManager.SetReady(false)
	healthManager.SetAlive(false)

	if err := httpServer.Shutdown(); err != nil {
		l.WithError(err).Error("HTTP server shutdown failed")
	} else {
		l.Info("HTTP server gracefully stopped")
	}
}
