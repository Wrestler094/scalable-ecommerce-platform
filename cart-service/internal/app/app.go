package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/adapters"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/validator"

	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/delivery/http/infra"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/delivery/http/v1"
	redisinfra "github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/infrastructure/redis"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/usecase"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	l.Info("Logger initialized", "level", cfg.Log.Level)

	// Connect Redis
	rdb, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	l.Info("Redis initialized", "url", cfg.Redis.Addr)

	// Helpers/Deps
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	healthManager := healthcheck.NewManager()

	// Repository
	cartRepository := redisinfra.NewRedisCartRepo(rdb.Client)

	// Use-Case
	cartUseCase := usecase.NewCartUseCase(cartRepository)

	// Handlers
	cartHandler := v1.NewCartHandler(cartUseCase, httpValidator)
	monitoringHandler := infra.NewMonitoringHandler(healthManager)

	handlers := http.Handlers{
		V1Handlers: v1.Handlers{
			CartHandler: cartHandler,
		},
		MonitoringHandler: monitoringHandler,
	}

	// Router
	router := http.NewRouter(handlers)

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

	healthManager.SetReady(true)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("Received signal", "signal", s)
	case err = <-httpServer.Notify():
		l.Error("httpServer.Notify", "error", err)
	}

	healthManager.SetReady(false)
	healthManager.SetAlive(false)

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("httpServer.Shutdown", "error", err)
	}
}
