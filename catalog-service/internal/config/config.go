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
		Metrics Metrics
		Swagger Swagger
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

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"false"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func (pg PG) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User, pg.Password, pg.Host, pg.Port, pg.DBName,
	)
}
