package main

import (
	"log"

	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/app"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
