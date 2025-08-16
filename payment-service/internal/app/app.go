package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/adapters"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/events"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/validator"

	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/delivery/http/infra"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/delivery/http/v1"
	kafkainfra "github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/infrastructure/kafka"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/infrastructure/postgres"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/infrastructure/redis"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/infrastructure/txmanager"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/usecase"
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

	// Connect Redis
	rdb, err := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		runLogger.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	runLogger.Info("Redis connected", "addr", cfg.Redis.Addr)

	// Helpers/Deps
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	txManager := txmanager.NewTxManager(pg.DB, baseLogger)
	healthManager := healthcheck.NewManager()

	// Repositories
	paymentRepo := postgres.NewPaymentRepository(pg.DB)
	outboxRepo := postgres.NewOutboxRepository(pg.DB)
	idempRepo := redis.NewIdempotencyRepository(rdb.Client)

	// Kafka
	producer := kafkainfra.NewProducer[events.PaymentSuccessfulPayload](cfg.Kafka.Brokers)
	defer producer.Close()

	runLogger.Info("Kafka producer initialized")

	poller := kafkainfra.NewPoller(outboxRepo, producer, baseLogger, events.TopicPayments, 5*time.Second, 100)
	ctx, pollerCancel := context.WithCancel(context.Background())
	go poller.Run(ctx)

	runLogger.Info("Kafka poller initialized")

	// Use-Cases
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, outboxRepo, idempRepo, txManager)

	// Handlers
	paymentHandler := v1.NewPaymentHandler(paymentUseCase, httpValidator, baseLogger)
	monitoringHandler := infra.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		V1Handlers: v1.Handlers{
			PaymentHandler: paymentHandler,
		},
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

	pollerCancel()
}
