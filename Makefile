.PHONY: user-up user-down catalog-up catalog-down all-up all-down

# === ENV PATHS ===
ENV_DIR := ./deploy/envs
COMPOSE_DIR := ./deploy

# === USER SERVICE ===
USER_COMPOSE := $(COMPOSE_DIR)/docker-compose.user.yml
USER_ENV := $(ENV_DIR)/user.env

user-up:
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) up -d --build

user-down:
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) down

# === CATALOG SERVICE ===
CATALOG_COMPOSE := $(COMPOSE_DIR)/docker-compose.catalog.yml
CATALOG_ENV := $(ENV_DIR)/catalog.env

catalog-up:
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) up -d --build

catalog-down:
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) down

# === ALL SERVICES ===
all-up: user-up catalog-up

all-down: user-down catalog-down
