package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App       App
		HTTP      HTTP
		JWT       JWT
		Log       Log
		Snowflake Snowflake
		PGShards  []PGShard
		Redis     Redis
		Metrics   Metrics
		Swagger   Swagger
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port int `env:"HTTP_PORT,required"`
	}

	JWT struct {
		AccessSecret string `env:"ACCESS_SECRET,required"`
		TokenTTL     int    `env:"TOKEN_TTL" env-default:"900"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	Snowflake struct {
		NodeID int64 `env:"SNOWFLAKE_NODE_ID,required"`
		Epoch  int64 `env:"SNOWFLAKE_EPOCH,required"`
	}

	PGShard struct {
		Name     string
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	Redis struct {
		Addr     string `env:"REDIS_ADDR,required"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB" env-default:"0"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" env-default:"false"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" env-default:"false"`
	}
)

// NewConfig returns parsed app config.
func NewConfig() (*Config, error) {
	const op = "config.NewConfig"

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to parse env: %w", op, err)
	}

	host0 := getEnvOrDefault("DB_HOST_0", "user-db-shard-0")
	host1 := getEnvOrDefault("DB_HOST_1", "user-db-shard-1")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "user")
	pass := getEnvOrDefault("DB_PASSWORD", "password")
	db0 := getEnvOrDefault("DB_NAME_0", "users_shard_0")
	db1 := getEnvOrDefault("DB_NAME_1", "users_shard_1")

	cfg.PGShards = []PGShard{
		{Name: "shard-0", Host: host0, Port: port, User: user, Password: pass, DBName: db0},
		{Name: "shard-1", Host: host1, Port: port, User: user, Password: pass, DBName: db1},
	}

	return cfg, nil
}

// DSN builds a PostgreSQL DSN.
func (pg PGShard) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User, pg.Password, pg.Host, pg.Port, pg.DBName,
	)
}

func getEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	log.Printf("env %s not set, using default: %s", key, fallback)
	return fallback
}
