package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"catalog-service/internal/config"
	"catalog-service/internal/delivery/http"
	"catalog-service/internal/infrastructure/postgres"
	"catalog-service/internal/usecase"
	"pkg/healthcheck"

	"pkg/adapters"
	"pkg/authenticator"
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
	pg, err := postgres.NewConnect(cfg.PG.DSN())
	if err != nil {
		l.Fatal("DB initialization failed", "error", err)
	}
	defer pg.Close()

	l.Info("PostgreSQL connected", "host", cfg.PG.Host, "port", cfg.PG.Port, "db", cfg.PG.DBName)

	// Helpers/Deps
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	JWTAuthenticator := authenticator.NewJWTAuthenticator(cfg.JWT.AccessSecret)
	healthManager := healthcheck.NewManager()

	// Repository
	categoryRepository := postgres.NewCategoryRepository(pg.DB)
	productRepository := postgres.NewProductRepository(pg.DB)

	// Use-Case
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	productUseCase := usecase.NewProductUseCase(productRepository)

	// Handlers
	categoryHandler := http.NewCategoryHandler(categoryUseCase, productUseCase, httpValidator)
	productHandler := http.NewProductHandler(productUseCase, httpValidator)
	monitoringHandler := http.NewMonitoringHandler(healthManager)

	handlers := http.Handlers{
		ProductHandler:    productHandler,
		CategoryHandler:   categoryHandler,
		MonitoringHandler: monitoringHandler,
	}

	// Router
	router := http.NewRouter(handlers, JWTAuthenticator)

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
