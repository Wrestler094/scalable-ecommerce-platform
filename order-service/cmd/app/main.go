package main

import (
	"log"

	"order-service/internal/app"
	"order-service/internal/config"
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
