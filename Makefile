.PHONY: user-up user-up-prod user-down \
        catalog-up catalog-up-prod catalog-down \
        cart-up cart-up-prod cart-down \
        order-up order-up-prod order-down \
        payment-up payment-up-prod payment-down \
        notification-up notification-up-prod notification-down \
        gateway-up gateway-up-prod gateway-down \
        traefik-up traefik-up-prod traefik-down \
        kafka-up kafka-down \
        monitoring-up monitoring-down \
        network-create \
        all-up all-up-prod all-down \
        clean-containers clean-images clean-hard \
        infra-up infra-down

# === GENERAL VARS ===
DEPLOY_DIR := ./deploy

# === NETWORK ===
NETWORK_NAME := backend

network-create:
	@if ! docker network ls --format '{{.Name}}' | grep -q "^$(NETWORK_NAME)$$"; then \
		echo "🔧 Creating network $(NETWORK_NAME)..."; \
		docker network create $(NETWORK_NAME); \
	else \
		echo "🔧 Network $(NETWORK_NAME) already exists."; \
	fi

# === USER SERVICE ===
USER_DEPLOY := $(DEPLOY_DIR)/user
USER_COMPOSE := $(USER_DEPLOY)/docker-compose.yml
USER_COMPOSE_DEV := $(USER_DEPLOY)/docker-compose.override.yml
USER_ENV := $(USER_DEPLOY)/user.env

user-up: network-create
	@echo "🔧 Starting User Service (dev config)..."
	docker compose -f $(USER_COMPOSE) -f $(USER_COMPOSE_DEV) --env-file $(USER_ENV) up -d --build

user-up-prod: network-create
	@echo "🔧 Starting User Service (prod config)..."
	docker compose -f $(USER_COMPOSE) --env-file $(USER_ENV) up -d --build

user-down:
	@echo "🛑 Shutdown User Service..."
	docker compose -f $(USER_COMPOSE) -f $(USER_COMPOSE_DEV) --env-file $(USER_ENV) down

# === CATALOG SERVICE ===
CATALOG_DEPLOY := $(DEPLOY_DIR)/catalog
CATALOG_COMPOSE := $(CATALOG_DEPLOY)/docker-compose.yml
CATALOG_COMPOSE_DEV := $(CATALOG_DEPLOY)/docker-compose.override.yml
CATALOG_ENV := $(CATALOG_DEPLOY)/catalog.env

catalog-up: network-create
	@echo "🔧 Starting Catalog Service (dev config)..."
	docker compose -f $(CATALOG_COMPOSE) -f $(CATALOG_COMPOSE_DEV) --env-file $(CATALOG_ENV) up -d --build

catalog-up-prod: network-create
	@echo "🔧 Starting Catalog Service (prod config)..."
	docker compose -f $(CATALOG_COMPOSE) --env-file $(CATALOG_ENV) up -d --build

catalog-down:
	@echo "🛑 Shutting down Catalog Service..."
	docker compose -f $(CATALOG_COMPOSE) -f $(CATALOG_COMPOSE_DEV) --env-file $(CATALOG_ENV) down

# === CART SERVICE ===
CART_DEPLOY := $(DEPLOY_DIR)/cart
CART_COMPOSE := $(CART_DEPLOY)/docker-compose.yml
CART_COMPOSE_DEV := $(CART_DEPLOY)/docker-compose.override.yml
CART_ENV := $(CART_DEPLOY)/cart.env

cart-up: network-create
	@echo "🔧 Starting Cart Service (dev config)..."
	docker compose -f $(CART_COMPOSE) -f $(CART_COMPOSE_DEV) --env-file $(CART_ENV) up -d --build

cart-up-prod: network-create
	@echo "🔧 Starting Cart Service (prod config)..."
	docker compose -f $(CART_COMPOSE) --env-file $(CART_ENV) up -d --build

cart-down:
	@echo "🛑 Shutting down Cart Service..."
	docker compose -f $(CART_COMPOSE) -f $(CART_COMPOSE_DEV) --env-file $(CART_ENV) down

# === ORDER SERVICE ===
ORDER_DEPLOY := $(DEPLOY_DIR)/order
ORDER_COMPOSE := $(ORDER_DEPLOY)/docker-compose.yml
ORDER_COMPOSE_DEV := $(ORDER_DEPLOY)/docker-compose.override.yml
ORDER_ENV := $(ORDER_DEPLOY)/order.env

order-up: network-create kafka-up
	@echo "🔧 Starting Order Service (dev config)..."
	docker compose -f $(ORDER_COMPOSE) -f $(ORDER_COMPOSE_DEV) --env-file $(ORDER_ENV) up -d --build

order-up-prod: network-create kafka-up
	@echo "🔧 Starting Order Service (prod config)..."
	docker compose -f $(ORDER_COMPOSE) --env-file $(ORDER_ENV) up -d --build

order-down:
	@echo "🛑 Shutting down Order Service..."
	docker compose -f $(ORDER_COMPOSE) -f $(ORDER_COMPOSE_DEV) --env-file $(ORDER_ENV) down

# === PAYMENT SERVICE ===
PAYMENT_DEPLOY := $(DEPLOY_DIR)/payment
PAYMENT_COMPOSE := $(PAYMENT_DEPLOY)/docker-compose.yml
PAYMENT_COMPOSE_DEV := $(PAYMENT_DEPLOY)/docker-compose.override.yml
PAYMENT_ENV := $(PAYMENT_DEPLOY)/payment.env

payment-up: network-create kafka-up
	@echo "🔧 Starting Payment Service (dev config)..."
	docker compose -f $(PAYMENT_COMPOSE) -f $(PAYMENT_COMPOSE_DEV) --env-file $(PAYMENT_ENV) up -d --build

payment-up-prod: network-create kafka-up
	@echo "🔧 Starting Payment Service (prod config)..."
	docker compose -f $(PAYMENT_COMPOSE) --env-file $(PAYMENT_ENV) up -d --build

payment-down:
	@echo "🛑 Shutting down Payment Service..."
	docker compose -f $(PAYMENT_COMPOSE) -f $(PAYMENT_COMPOSE_DEV) --env-file $(PAYMENT_ENV) down

# === NOTIFICATION SERVICE ===
NOTIFICATION_DEPLOY := $(DEPLOY_DIR)/notification
NOTIFICATION_COMPOSE := $(NOTIFICATION_DEPLOY)/docker-compose.yml
NOTIFICATION_COMPOSE_DEV := $(NOTIFICATION_DEPLOY)/docker-compose.override.yml
NOTIFICATION_ENV := $(NOTIFICATION_DEPLOY)/notification.env

notification-up: network-create kafka-up
	@echo "🔧 Starting Notification Service (dev config)..."
	docker compose -f $(NOTIFICATION_COMPOSE) -f $(NOTIFICATION_COMPOSE_DEV) --env-file $(NOTIFICATION_ENV) up -d --build

notification-up-prod: network-create kafka-up
	@echo "🔧 Starting Notification Service (prod config)..."
	docker compose -f $(NOTIFICATION_COMPOSE) --env-file $(NOTIFICATION_ENV) up -d --build

notification-down:
	@echo "🛑 Shutting down Notification Service..."
	docker compose -f $(NOTIFICATION_COMPOSE) -f $(NOTIFICATION_COMPOSE_DEV) --env-file $(NOTIFICATION_ENV) down

# === TRAEFIK (REVERSE PROXY) ===
TRAEFIK_DEPLOY := $(DEPLOY_DIR)/traefik
TRAEFIK_COMPOSE := $(TRAEFIK_DEPLOY)/docker-compose.yml
TRAEFIK_COMPOSE_DEV := $(TRAEFIK_DEPLOY)/docker-compose.override.yml

traefik-up: network-create
	@echo "🔧 Starting Traefik (dev config)..."
	docker compose -f $(TRAEFIK_COMPOSE) -f $(TRAEFIK_COMPOSE_DEV) up -d --build

traefik-up-prod: network-create
	@echo "🔧 Starting Traefik (prod config)..."
	docker compose -f $(TRAEFIK_COMPOSE) up -d --build

traefik-down:
	@echo "🛑 Shutting down Traefik..."
	docker compose -f $(TRAEFIK_COMPOSE) -f $(TRAEFIK_COMPOSE_DEV) down

# === API GATEWAY ===
GATEWAY_DEPLOY := $(DEPLOY_DIR)/gateway
GATEWAY_COMPOSE := $(GATEWAY_DEPLOY)/docker-compose.yml
GATEWAY_ENV := $(GATEWAY_DEPLOY)/gateway.env

gateway-up: network-create traefik-up
	@echo "🔧 Starting API Gateway..."
	docker compose -f $(GATEWAY_COMPOSE) --env-file $(GATEWAY_ENV) up -d --build

gateway-up-prod: network-create traefik-up-prod
	@echo "🔧 Starting API Gateway (prod config)..."
	docker compose -f $(GATEWAY_COMPOSE) --env-file $(GATEWAY_ENV) up -d --build

gateway-down:
	@echo "🛑 Shutting down API Gateway..."
	docker compose -f $(GATEWAY_COMPOSE) --env-file $(GATEWAY_ENV) down

# === KAFKA STACK ===
KAFKA_DEPLOY := $(DEPLOY_DIR)/kafka
KAFKA_COMPOSE := $(KAFKA_DEPLOY)/docker-compose.yml
KAFKA_COMPOSE_DEV := $(KAFKA_DEPLOY)/docker-compose.override.yml
KAFKA_ENV := $(KAFKA_DEPLOY)/kafka.env

kafka-up: network-create
	@echo "🔧 Starting Kafka Stack (dev config)..."
	docker compose -f $(KAFKA_COMPOSE) -f $(KAFKA_COMPOSE_DEV) --env-file $(KAFKA_ENV) up -d --build

kafka-up-prod: network-create
	@echo "🔧 Starting Kafka Stack (prod config)..."
	docker compose -f $(KAFKA_COMPOSE) --env-file $(KAFKA_ENV) up -d --build

kafka-down:
	@echo "🛑 Shutting down Kafka Stack..."
	docker compose -f $(KAFKA_COMPOSE) -f $(KAFKA_COMPOSE_DEV) --env-file $(KAFKA_ENV) down

# === MONITORING ===
MONITORING_DEPLOY := $(DEPLOY_DIR)/monitoring
MONITORING_COMPOSE := $(MONITORING_DEPLOY)/docker-compose.yml
MONITORING_COMPOSE_DEV := $(MONITORING_DEPLOY)/docker-compose.override.yml

monitoring-up: network-create
	@echo "🔧 Starting Monitoring Stack (dev config)..."
	docker compose -f $(MONITORING_COMPOSE) -f $(MONITORING_COMPOSE_DEV) up -d --build

monitoring-up-prod: network-create
	@echo "🔧 Starting Monitoring Stack (prod config)..."
	docker compose -f $(MONITORING_COMPOSE) up -d --build

monitoring-down:
	@echo "🛑 Shutting down Monitoring Stack..."
	docker compose -f $(MONITORING_COMPOSE) -f $(MONITORING_COMPOSE_DEV) down

# === INFRASTRUCTURE ===
infra-up: network-create
	@echo "🔧 Starting infrastructure (Traefik + Monitoring + Kafka)..."
	$(MAKE) traefik-up
	$(MAKE) monitoring-up
	$(MAKE) kafka-up

infra-up-prod: network-create
	@echo "🔧 Starting infrastructure (Traefik + Monitoring + Kafka) [prod config]..."
	$(MAKE) traefik-up-prod
	$(MAKE) monitoring-up-prod
	$(MAKE) kafka-up-prod

infra-down:
	@echo "🛑 Shutting down infrastructure (Traefik + Kafka + Monitoring)..."
	$(MAKE) kafka-down
	$(MAKE) monitoring-down
	$(MAKE) traefik-down

# === ALL SERVICES ===
all-up:
	$(MAKE) infra-up
	$(MAKE) gateway-up
	$(MAKE) user-up
	$(MAKE) catalog-up
	$(MAKE) cart-up
	$(MAKE) order-up
	$(MAKE) payment-up
	$(MAKE) notification-up
	@echo "🚀 All services started!"

all-up-prod:
	$(MAKE) infra-up-prod
	$(MAKE) gateway-up-prod
	$(MAKE) user-up-prod
	$(MAKE) catalog-up-prod
	$(MAKE) cart-up-prod
	$(MAKE) order-up-prod
	$(MAKE) payment-up-prod
	$(MAKE) notification-up-prod
	@echo "🚀 All services started (prod config)!"

all-down:
	$(MAKE) notification-down
	$(MAKE) payment-down
	$(MAKE) order-down
	$(MAKE) cart-down
	$(MAKE) catalog-down
	$(MAKE) user-down
	$(MAKE) gateway-down
	$(MAKE) infra-down
	@echo "🛑 All services stopped!"

# === CLEAN ===
clean-containers:
	@echo "🧹 Removing all stopped containers..."
	docker container prune -f

clean-images:
	@echo "🧹 Removing all unused images..."
	docker image prune -a -f

clean-hard: all-down
	@echo "🔥 Performing full cleanup: volumes, images, containers, networks..."
	$(MAKE) clean-containers
	$(MAKE) clean-images
	docker network prune -f
