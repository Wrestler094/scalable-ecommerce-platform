package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pkg/adapters"
	"pkg/authenticator"
	"pkg/events"
	"pkg/healthcheck"
	"pkg/httpserver"
	"pkg/logger"
	"pkg/validator"

	"order-service/internal/config"
	"order-service/internal/delivery/http"
	"order-service/internal/delivery/kafka"
	paymentmock "order-service/internal/infrastructure/payment/mock"
	"order-service/internal/infrastructure/postgres"
	productmock "order-service/internal/infrastructure/product/mock"
	"order-service/internal/usecase"
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
	pg, err := postgres.NewConnect(cfg.PG.DSN())
	if err != nil {
		runLogger.Fatal("DB initialization failed", "error", err)
	}
	defer pg.Close()

	runLogger.Info("PostgreSQL connected", "host", cfg.PG.Host, "port", cfg.PG.Port, "db", cfg.PG.DBName)

	// Helpers/Deps
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	healthManager := healthcheck.NewManager()
	authenticatorImpl := authenticator.NewJWTAuthenticator(cfg.JWT.AccessSecret)

	// Services
	productService := productmock.NewMockProductService()
	paymentService := paymentmock.NewMockPaymentService()

	// Repositories
	orderRepo := postgres.NewOrderRepository(pg.DB)

	// Use-Cases
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productService, paymentService)
	paymentUseCase := usecase.NewPaymentUseCase(orderRepo)

	// Handlers
	orderHandler := http.NewOrderHandler(orderUseCase, httpValidator, baseLogger)
	monitoringHandler := http.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		OrderHandler:      orderHandler,
		MonitoringHandler: monitoringHandler,
	}, authenticatorImpl)

	// Kafka Consumer
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		events.TopicPayments,
		events.OrderGroup,
		paymentUseCase,
		baseLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

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

	cancel()
	if err := consumer.Close(); err != nil {
		runLogger.WithError(err).Error("Failed to close consumer")
	} else {
		runLogger.Info("Kafka consumer closed successfully")
	}
}
