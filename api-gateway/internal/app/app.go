package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/delivery/http/middleware"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
	"github.com/samber/lo"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httpserver"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"

	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/config"
	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/delivery/http"
	gatewayHandlers "github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/delivery/http/handlers"
)

func Run(cfg *config.Config) {
	start := time.Now()

	var proxyRoutes = map[string]string{
		"user":         cfg.Routes.UserService,
		"catalog":      cfg.Routes.CatalogService,
		"cart":         cfg.Routes.CartService,
		"order":        cfg.Routes.OrderService,
		"payment":      cfg.Routes.PaymentService,
		"notification": cfg.Routes.NotificationService,
	}

	baseLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("main.run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	healthManager := healthcheck.NewManager()
	jwtAuthenticator := authenticator.NewJWTAuthenticator(cfg.JWT.AccessSecret)
	tokenMiddleware := middleware.NewTokenMiddleware(jwtAuthenticator, baseLogger)

	monitoringHandler := gatewayHandlers.NewMonitoringHandler(healthManager)
	proxyHandler := gatewayHandlers.NewProxyHandler(proxyRoutes, baseLogger)

	handlers := http.Handlers{
		ProxyHandler:      proxyHandler,
		MonitoringHandler: monitoringHandler,
	}

	router := http.NewRouter(handlers, tokenMiddleware, lo.Keys(proxyRoutes))

	// HTTP Server
	httpServer := httpserver.NewServer(
		httpserver.Port(fmt.Sprintf(":%d", cfg.HTTP.Port)),
		httpserver.Handler(router),
	)

	runLogger.Info("HTTP server is starting", "app", cfg.App.Name, "port", cfg.HTTP.Port)

	// Start
	if err := httpServer.Start(); err != nil {
		runLogger.WithError(err).Fatal("HTTP server failed to start")
	}

	healthManager.SetReady(true)

	runLogger.Info("Startup complete", logger.LogKeyDurationMS, time.Since(start).String())

	// Graceful shutdown
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
