package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/migrator"
)

const (
	defaultPath = "./migrations"
	op          = "cmd.migrator"
)

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("%s: one of the DB env variables is missing", op)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

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
