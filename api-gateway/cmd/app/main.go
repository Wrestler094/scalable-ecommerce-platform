package main

import (
	"log"

	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/app"
	"github.com/Wrestler094/scalable-ecommerce-platform/api-gateway/internal/config"
)

func main() {
	const op = "cmd.app"

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("%s: config error: %s", op, err)
	}

	// Run
	app.Run(cfg)
}
