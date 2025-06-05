package main

import (
	"log"
	"os"

	"pkg/migrator"
)

const defaultPath = "./migrations"

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		log.Fatal("missing POSTGRES_URL")
	}

	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = defaultPath
	}

	log.Printf("running migrations from %q", path)

	if err := migrator.Run(dsn, path); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	log.Println("migrations applied successfully")
}
