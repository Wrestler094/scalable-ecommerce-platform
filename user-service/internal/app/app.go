package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/config"
)

// Run creates objects via constructors and starts the application.
func Run(cfg *config.Config) {
	start := time.Now()

	// Init base logger
	baseLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("app.Run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	healthManager := healthcheck.NewManager()

	// DI
	httpServer, cleanup, err := InitDI(cfg, baseLogger, healthManager)
	if err != nil {
		runLogger.WithError(err).Fatal("Initialization failed")
	}
	defer func() {
		if cerr := cleanup(); cerr != nil {
			runLogger.WithError(cerr).Error("Cleanup finished with errors")
		}
	}()

	runLogger.Info("HTTP server is starting", "port", cfg.HTTP.Port)

	// Start HTTP server
	if err := httpServer.Start(); err != nil {
		runLogger.WithError(err).Fatal("HTTP server failed to start")
	}

	healthManager.SetReady(true)

	runLogger.Info("Startup complete", logger.LogKeyDurationMS, time.Since(start).String())

	// Graceful shutdown on signal or server error
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		runLogger.Info("Received shutdown signal", "signal", sig)
	case err := <-httpServer.Notify():
		runLogger.WithError(err).Error("HTTP server reported error")
	}

	// Change healt status and stop server
	healthManager.SetReady(false)
	healthManager.SetAlive(false)

	if err := httpServer.Shutdown(); err != nil {
		runLogger.WithError(err).Error("HTTP server shutdown failed")
	} else {
		runLogger.Info("HTTP server gracefully stopped")
	}
}
