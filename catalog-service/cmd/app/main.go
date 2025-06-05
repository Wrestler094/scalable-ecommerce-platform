package main

import (
	"log"

	"catalog-service/internal/app"
	"catalog-service/internal/config"
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
