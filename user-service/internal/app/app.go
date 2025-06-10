package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"pkg/cache"
	"user-service/internal/config"
	"user-service/internal/delivery/http"
	"user-service/internal/infrastructure/hasher"
	"user-service/internal/infrastructure/jwt"
	"user-service/internal/infrastructure/postgres"
	redisinfra "user-service/internal/infrastructure/redis"
	"user-service/internal/usecase"

	"pkg/adapters"
	"pkg/httpserver"
	"pkg/logger"
	"pkg/validator"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	l.Info("Logger initialized", "level", cfg.Log.Level)

	// Connect Postgres
	pg, err := postgres.NewConnect(cfg.PG.URL)
	if err != nil {
		l.Fatal("DB initialization failed", "error", err)
	}
	defer pg.Close()

	l.Info("DB initialized", "url", strings.Split(cfg.PG.URL, "@")[1])

	// Connect Redis
	rdb, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	l.Info("Redis initialized", "url", cfg.Redis.Addr)

	// Helpers/Deps
	tokenManager := jwt.NewManager(cfg.JWT.AccessSecret, time.Duration(cfg.JWT.TokenTTL)*time.Second)
	passwordHasher := hasher.NewBcryptHasher()
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	redisCache := cache.NewRedisCache(rdb.Client)

	// Repository
	userRepo := postgres.NewUserRepository(pg.DB)
	cashedUserRepo := postgres.NewCachedUserRepository(userRepo, redisCache, l)
	refreshRepo := redisinfra.NewRefreshTokenRepository(rdb.Client)

	// Use-Case
	userUseCase := usecase.NewUserUseCase(cashedUserRepo, refreshRepo, tokenManager, passwordHasher)

	// Handlers
	userHandler := http.NewUserHandler(userUseCase, httpValidator)

	// Router
	router := http.NewRouter(http.Handlers{
		UserHandler: userHandler,
	})

	// HTTP Server
	httpServer := httpserver.NewServer(
		httpserver.Port(fmt.Sprintf(":%d", cfg.HTTP.Port)),
		httpserver.Handler(router),
	)

	l.Info("HTTP Server running", "port", cfg.HTTP.Port)

	// Start servers
	err = httpServer.Start()
	if err != nil {
		l.Fatal("Failed to start server", "error", err)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("Received signal", "signal", s)
	case err = <-httpServer.Notify():
		l.Error("httpServer.Notify", "error", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("httpServer.Shutdown", "error", err)
	}
}
