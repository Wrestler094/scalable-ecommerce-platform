package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/events"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	infralogger "github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"

	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/delivery/http/infra"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/delivery/http/v1"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/delivery/kafka"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/infrastructure/sender"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/usecase"
)

// Run creates objects via constructors and starts the application.
func Run(cfg *config.Config) {
	start := time.Now()

	baseLogger, err := infralogger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("app.Run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	// Helpers/Deps
	healthManager := healthcheck.NewManager()
	emailSender := sender.NewEmailSender(cfg.ElasticEmail.APIKey, cfg.ElasticEmail.FromEmail, cfg.ElasticEmail.FromName)

	// Use-Cases
	notificationUseCase := usecase.NewNotificationUseCase(emailSender)

	// Handlers
	monitoringHandler := infra.NewMonitoringHandler(healthManager)

	// Router
	router := http.NewRouter(http.Handlers{
		V1Handlers: v1.Handlers{
			// No handlers yet, but structure ready for future expansion
		},
		MonitoringHandler: monitoringHandler,
	})

	// Kafka Consumer
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		events.TopicPayments,
		events.NotificationGroup,
		notificationUseCase,
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

	runLogger.Info("Startup complete", infralogger.LogKeyDurationMS, time.Since(start).String())

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
