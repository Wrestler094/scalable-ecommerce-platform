package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/infrastructure/client/catalog"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/adapters"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/events"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/validator"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/delivery/http/infra"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/delivery/http/v1"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/delivery/kafka"
	paymentmock "github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/infrastructure/client/payment/mock"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/infrastructure/postgres"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/usecase"
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

	// Services
	paymentService := paymentmock.NewMockPaymentService()
	productProvider, err := catalog.NewClient(context.Background(), cfg.Clients.Catalog)
	if err != nil {
		runLogger.WithError(err).Fatal("failed to create catalog service client")
	}

	// Repositories
	orderRepo := postgres.NewOrderRepository(pg.DB)

	// Use-Cases
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productProvider, paymentService)
	paymentUseCase := usecase.NewPaymentUseCase(orderRepo)

	// Handlers
	orderHandler := v1.NewOrderHandler(orderUseCase, httpValidator, baseLogger)
	monitoringHandler := infra.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		V1Handlers: v1.Handlers{
			OrderHandler: orderHandler,
		},
		MonitoringHandler: monitoringHandler,
	})

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
