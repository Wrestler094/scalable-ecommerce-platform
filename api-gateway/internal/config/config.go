package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App    App
		HTTP   HTTP
		JWT    JWT
		Routes Routes
		Log    Log
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
		// TODO: На будущее, когда access токен будет лежать в куках
		// AccessCookieName  string `env:"ACCESS_COOKIE_NAME" env-default:"access_token"`
		// TODO: На будущее, для внутренних вызовов через гейтвей
		// ServiceKeyHeader string `env:"SERVICE_KEY_HEADER" env-default:"X-Service-Key"`
		// ServiceKey       string `env:"SERVICE_KEY" env-default:""`
	}

	Routes struct {
		UserService         string `env:"USER_SERVICE_URL" env-default:"http://user-service:4000"`
		CatalogService      string `env:"CATALOG_SERVICE_URL" env-default:"http://catalog-service:4000"`
		CartService         string `env:"CART_SERVICE_URL" env-default:"http://cart-service:4000"`
		OrderService        string `env:"ORDER_SERVICE_URL" env-default:"http://order-service:4000"`
		PaymentService      string `env:"PAYMENT_SERVICE_URL" env-default:"http://payment-service:4000"`
		NotificationService string `env:"NOTIFICATION_SERVICE_URL" env-default:"http://notification-service:4000"`
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
