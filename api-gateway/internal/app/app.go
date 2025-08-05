package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/lo"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"

	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/config"
	gatewayHTTP "github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/delivery/http"
)

var proxyRoutes = map[string]string{
	"user":         "http://user-service:4000",
	"catalog":      "http://catalog-service:4000",
	"cart":         "http://cart-service:4000",
	"order":        "http://order-service:4000",
	"payment":      "http://payment-service:4000",
	"notification": "http://notification-service:4000",
}

func Run(cfg *config.Config) {
	start := time.Now()

	baseLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("main.run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	healthManager := healthcheck.NewManager()

	monitoringHandler := gatewayHTTP.NewMonitoringHandler(healthManager)
	proxyHandler := gatewayHTTP.NewStaticProxyHandler(proxyRoutes, baseLogger)

	handlers := gatewayHTTP.Handlers{
		ProxyHandler:      proxyHandler,
		MonitoringHandler: monitoringHandler,
	}

	router := gatewayHTTP.NewRouter(handlers, lo.Keys(proxyRoutes))

	// HTTP Server
	httpServer := httpserver.NewServer(
		httpserver.Port(fmt.Sprintf(":%d", cfg.HTTP.Port)),
		httpserver.Handler(router),
	)

	runLogger.Info("HTTP server is starting", "app", cfg.App.Name, "port", cfg.HTTP.Port)

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
