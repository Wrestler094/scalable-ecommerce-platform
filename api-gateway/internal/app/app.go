package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/config"
	gatewayHTTP "github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/delivery/http"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
)

func Run(cfg *config.Config) {
	baseLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger initialization failed: %s", err)
	}

	runLogger := baseLogger.WithOp("main.run")
	runLogger.Info("Logger initialized", "level", cfg.Log.Level)

	handlers := gatewayHTTP.NewStaticProxyHandler(map[string]string{
		"user":         "http://user-service:8080",
		"catalog":      "http://catalog-service:8080",
		"cart":         "http://cart-service:8080",
		"order":        "http://order-service:8080",
		"payment":      "http://payment-service:8080",
		"notification": "http://notification-service:8080",
	})

	router := gatewayHTTP.NewRouter(handlers)

	runLogger.Info("🚀 Gateway running", "port", cfg.HTTP.Port, "app", cfg.App.Name)

	srvError := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTP.Port), router)
	if srvError != nil {
		baseLogger.Fatal("Failed to start server", "error", srvError)
	}
}
