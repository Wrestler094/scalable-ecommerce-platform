.PHONY: user-up user-down catalog-up catalog-down cart-up cart-down all-up all-down

# === DEPLOY PATHS ===
DEPLOY_DIR := ./deploy
USER_DEPLOY := $(DEPLOY_DIR)/user
CATALOG_DEPLOY := $(DEPLOY_DIR)/catalog
CART_DEPLOY := $(DEPLOY_DIR)/cart

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

# === ALL SERVICES ===
all-up: user-up catalog-up cart-up

all-down: user-down catalog-down cart-down
