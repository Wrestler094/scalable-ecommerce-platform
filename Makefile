.PHONY: user-up user-down catalog-up catalog-down cart-up cart-down \
        monitoring-up monitoring-down network-create \
        all-up all-down infra-up infra-down

# === DEPLOY PATHS ===
DEPLOY_DIR := ./deploy
USER_DEPLOY := $(DEPLOY_DIR)/user
CATALOG_DEPLOY := $(DEPLOY_DIR)/catalog
CART_DEPLOY := $(DEPLOY_DIR)/cart
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

user-up:
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) up -d --build

user-down:
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) down

# === CATALOG SERVICE ===
CATALOG_COMPOSE := $(CATALOG_DEPLOY)/docker-compose.yml
CATALOG_ENV := $(CATALOG_DEPLOY)/catalog.env

catalog-up:
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) up -d --build

catalog-down:
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) down

# === CART SERVICE ===
CART_COMPOSE := $(CART_DEPLOY)/docker-compose.yml
CART_ENV := $(CART_DEPLOY)/cart.env

cart-up:
	docker compose -f $(CART_COMPOSE) --env-file $(CART_ENV) up -d --build

cart-down:
	docker compose -f $(CART_COMPOSE) --env-file $(CART_ENV) down

# === MONITORING ===
MONITORING_COMPOSE := $(MONITORING_DEPLOY)/docker-compose.yml

monitoring-up: network-create
	docker compose -f $(MONITORING_COMPOSE) up -d --build

monitoring-down:
	docker compose -f $(MONITORING_COMPOSE) down

# === INFRASTRUCTURE ===
infra-up: monitoring-up

infra-down: monitoring-down

# === ALL SERVICES ===
all-up: infra-up user-up catalog-up cart-up

all-down: cart-down catalog-down user-down infra-down
