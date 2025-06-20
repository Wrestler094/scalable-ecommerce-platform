.PHONY: user-up user-down catalog-up catalog-down cart-up cart-down \
        payment-up payment-down kafka-up kafka-down \
        monitoring-up monitoring-down network-create \
        all-up all-down infra-up infra-down

# === DEPLOY PATHS ===
DEPLOY_DIR := ./deploy
USER_DEPLOY := $(DEPLOY_DIR)/user
CATALOG_DEPLOY := $(DEPLOY_DIR)/catalog
CART_DEPLOY := $(DEPLOY_DIR)/cart
PAYMENT_DEPLOY := $(DEPLOY_DIR)/payment
KAFKA_DEPLOY := $(DEPLOY_DIR)/kafka
MONITORING_DEPLOY := $(DEPLOY_DIR)/monitoring

# === NETWORK ===
NETWORK_NAME := backend

network-create:
	@if ! docker network ls --format '{{.Name}}' | grep -q "^$(NETWORK_NAME)$$"; then \
		echo "Creating network $(NETWORK_NAME)..."; \
		docker network create $(NETWORK_NAME); \
	else \
		echo "Network $(NETWORK_NAME) already exists."; \
	fi

# === USER SERVICE ===
USER_COMPOSE := $(USER_DEPLOY)/docker-compose.yml
USER_ENV := $(USER_DEPLOY)/user.env

user-up: network-create
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) up -d --build

user-down:
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) down

# === CATALOG SERVICE ===
CATALOG_COMPOSE := $(CATALOG_DEPLOY)/docker-compose.yml
CATALOG_ENV := $(CATALOG_DEPLOY)/catalog.env

catalog-up: network-create
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) up -d --build

catalog-down:
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) down

# === CART SERVICE ===
CART_COMPOSE := $(CART_DEPLOY)/docker-compose.yml
CART_ENV := $(CART_DEPLOY)/cart.env

cart-up: network-create
	docker compose -f $(CART_COMPOSE) --env-file $(CART_ENV) up -d --build

cart-down:
	docker compose -f $(CART_COMPOSE) --env-file $(CART_ENV) down

# === PAYMENT SERVICE ===
PAYMENT_COMPOSE := $(PAYMENT_DEPLOY)/docker-compose.yml
PAYMENT_ENV := $(PAYMENT_DEPLOY)/payment.env

payment-up: network-create kafka-up
	docker compose -f $(PAYMENT_COMPOSE) --env-file $(PAYMENT_ENV) up -d --build

payment-down:
	docker compose -f $(PAYMENT_COMPOSE) --env-file $(PAYMENT_ENV) down

# === KAFKA STACK ===
KAFKA_COMPOSE := $(KAFKA_DEPLOY)/docker-compose.yml
KAFKA_ENV := $(KAFKA_DEPLOY)/kafka.env

kafka-up: network-create
	docker compose -f $(KAFKA_COMPOSE) --env-file $(KAFKA_ENV) up -d --build

kafka-down:
	docker compose -f $(KAFKA_COMPOSE) --env-file $(KAFKA_ENV) down

# === MONITORING ===
MONITORING_COMPOSE := $(MONITORING_DEPLOY)/docker-compose.yml

monitoring-up: network-create
	docker compose -f $(MONITORING_COMPOSE) up -d --build

monitoring-down:
	docker compose -f $(MONITORING_COMPOSE) down

# === INFRASTRUCTURE ===
infra-up: monitoring-up kafka-up

infra-down: kafka-down monitoring-down

# === ALL SERVICES ===
all-up: infra-up user-up catalog-up cart-up payment-up

all-down: payment-down cart-down catalog-down user-down infra-down
