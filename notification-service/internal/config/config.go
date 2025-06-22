package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App          App
		HTTP         HTTP
		Log          Log
		Kafka        Kafka
		ElasticEmail ElasticEmail
		Metrics      Metrics
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port int `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	Kafka struct {
		Brokers []string `env:"KAFKA_BROKERS,required"`
	}

	ElasticEmail struct {
		APIKey    string `env:"ELASTIC_API_KEY"`
		FromEmail string `env:"ELASTIC_FROM_EMAIL"`
		FromName  string `env:"ELASTIC_FROM_NAME"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"false"`
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
