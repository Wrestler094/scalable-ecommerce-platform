package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"payment-service/internal/infrastructure/txmanager"
	"pkg/events"

	"pkg/adapters"
	"pkg/authenticator"
	"pkg/healthcheck"
	"pkg/httpserver"
	"pkg/logger"
	"pkg/validator"

	"payment-service/internal/config"
	"payment-service/internal/delivery/http"
	kafkainfra "payment-service/internal/infrastructure/kafka"
	"payment-service/internal/infrastructure/postgres"
	"payment-service/internal/infrastructure/redis"
	"payment-service/internal/usecase"
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
	pg, err := postgres.NewConnect(cfg.PG.DSN())
	if err != nil {
		l.Fatal("DB initialization failed", "error", err)
	}
	defer pg.Close()

	l.Info("PostgreSQL connected", "host", cfg.PG.Host, "port", cfg.PG.Port, "db", cfg.PG.DBName)

	// Connect Redis
	rdb, err := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("Redis initialization failed", "error", err)
	}
	defer rdb.Close()

	l.Info("Redis connected", "addr", cfg.Redis.Addr)

	// Connect Kafka
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	producer := kafkainfra.NewProducer(writer)
	defer writer.Close() // TODO: Сделать через обертку

	l.Info("Kafka producer initialized")

	// Helpers/Deps
	rawValidator := validator.NewPlaygroundValidator()
	httpValidator := adapters.NewHttpValidatorAdapter(rawValidator)
	JWTAuthenticator := authenticator.NewJWTAuthenticator(cfg.JWT.AccessSecret)
	txManager := txmanager.NewTxManager(pg.DB, l)
	healthManager := healthcheck.NewManager()

	// Repositories
	paymentRepo := postgres.NewPaymentRepository(pg.DB)
	outboxRepo := postgres.NewOutboxRepository(pg.DB)
	idempRepo := redis.NewIdempotencyRepository(rdb.Client)

	// Kafka poller
	poller := kafkainfra.NewPoller(outboxRepo, producer, l, events.TopicPayments, 5*time.Second, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go poller.Run(ctx) // TODO: Подумать как лучше закрывать
	l.Info("Kafka poller initialized")

	// Use-Cases
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, outboxRepo, idempRepo, txManager)

	// Handlers
	paymentHandler := http.NewPaymentHandler(paymentUseCase, httpValidator, l)
	monitoringHandler := http.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		PaymentHandler:    paymentHandler,
		MonitoringHandler: monitoringHandler,
	}, JWTAuthenticator)

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
