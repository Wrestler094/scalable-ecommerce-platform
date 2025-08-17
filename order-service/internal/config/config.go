package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		JWT     JWT
		Log     Log
		PG      PG
		Kafka   Kafka
		Metrics Metrics
		Swagger Swagger
		Clients Clients
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port int `env:"HTTP_PORT,required"`
	}

	JWT struct{}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		Host     string `env:"DB_HOST,required"`
		Port     string `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		DBName   string `env:"DB_NAME,required"`
	}

	Kafka struct {
		Brokers []string `env:"KAFKA_BROKERS,required"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"false"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	Clients struct {
		Catalog string `env:"CATALOG_SERVICE_URL,required"`
		// TODO: Uncomment when payment client is implemented
		// Payment string `env:"PAYMENT_SERVICE_URL,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	const op = "config.NewConfig"

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to parse env: %w", op, err)
	}

	return cfg, nil
}

func (pg PG) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User, pg.Password, pg.Host, pg.Port, pg.DBName,
	)
}
