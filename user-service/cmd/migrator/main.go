package main

import (
	"fmt"
	"log"
	"os"

	"pkg/migrator"
)

const (
	defaultPath = "./migrations"
	op          = "cmd.migrator"
)

func main() {
	port := getEnvOrFail("DB_PORT")
	user := getEnvOrFail("DB_USER")
	password := getEnvOrFail("DB_PASSWORD")

	shards := []struct {
		Name string
		Host string
		DB   string
	}{
		{
			Name: "shard-0",
			Host: getEnvOrFail("DB_HOST_0"),
			DB:   getEnvOrFail("DB_NAME_0"),
		},
		{
			Name: "shard-1",
			Host: getEnvOrFail("DB_HOST_1"),
			DB:   getEnvOrFail("DB_NAME_1"),
		},
	}

	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = defaultPath
	}

	log.Printf("%s: running DB migrations from %q ...", op, path)

	for _, shard := range shards {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password, shard.Host, port, shard.DB,
		)

		log.Printf("%s: applying migrations to %s (%s)", op, shard.Name, shard.DB)

		if err := migrator.Run(dsn, path); err != nil {
			log.Fatalf("%s: migration error on %s: %v", op, shard.Name, err)
		}

		log.Printf("%s: ✅ migrations applied to %s", op, shard.Name)
	}

	log.Printf("%s: ✅ all migrations completed successfully", op)
}

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s: missing required environment variable: %s", op, key)
	}
	return val
}
