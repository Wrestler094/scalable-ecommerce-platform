package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		JWT  JWT
		Log  Log
		//Redis     Redis
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

	JWT struct {
		AccessSecret string `env:"ACCESS_SECRET,required"`
		TokenTTL     int    `env:"TOKEN_TTL" env-default:"900"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	//Redis struct {
	//	Addr     string `env:"REDIS_ADDR,required"`
	//	Password string `env:"REDIS_PASSWORD"`
	//	DB       int    `env:"REDIS_DB" env-default:"0"`
	//}

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

	return cfg, nil
}
