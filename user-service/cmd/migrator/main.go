package main

import (
	"log"
	"os"

	"pkg/migrator"
)

const (
	defaultPath = "./migrations"
	op          = "cmd.migrator"
)

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		log.Fatalf("%s: missing environment variable: POSTGRES_URL", op)
	}

	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = defaultPath
	}

	log.Printf("%s: running DB migrations from %q ...", op, path)

	if err := migrator.Run(dsn, path); err != nil {
		log.Fatalf("%s: migration error: %v", op, err)
	}

	log.Println("migrations applied successfully")
}
