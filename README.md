# 🧱 Scalable E-Commerce Platform

Микросервисная e-commerce платформа на Go, Docker, Redis и PostgreSQL. Разделена на изолированные сервисы с возможностью масштабирования, независимого деплоя и мониторинга. Использует Kafka для обмена событиями между сервисами. Поддерживает централизованную авторизацию, единый подход к валидации, логированию и метрикам.

## 📦 Сервисы

- `user-service` — регистрация, авторизация (PostgreSQL, Redis)
- `catalog-service` — продукты, категории (PostgreSQL)
- `cart-service` — корзина пользователя (Redis)
- `payment-service` — заглушка для обработки платежей (PostgreSQL, Redis, Kafka)
- `notification-service` — отправка уведомлений по email и/или SMS (Kafka)
- `pkg/` — общие модули: `authenticator`, `roles`, `httphelper`, `validator` и др. 
- `deploy/monitoring/` — мониторинг на базе Prometheus + Grafana, c auto-provisioning дашбордов

## 📁 Структура проекта

```text
scalable-ecommerce-platform/
├── Makefile                         # Общие команды запуска
├── pkg/                             # Общие переиспользуемые модули
│   ├── authenticator/               # JWT, контекст, middleware
│   ├── httphelper/                  # Обработка ошибок, JSON-ответов
│   ├── events/                      # Общие топики события Kafka
│   └── ...
├── deploy/                          # Docker Compose и .env файлы по сервисам
│   ├── user/
│   │   ├── docker-compose.yml
│   │   ├── user.env.example
│   │   └── user.env
│   ├── catalog/
│   ├── cart/
│   ├── payment/
│   ├── kafka/                       # Kafka stack: Kafka + ZooKeeper + Kafka UI
│   └── monitoring/                  # Monitoring stack: Prometheus + Grafana
│       ├── docker-compose.yml       # Сборка мониторинга
│       ├── prometheus/
│       │   └── prometheus.yml       # Конфигурация Prometheus (scrape configs)
│       └── grafana/
│           ├── dashboards/          # JSON-файлы дашбордов (provisioning)
│           ├── dashboards.yml       # Описание провайдеров дашбордов
│           └── datasources/         # Prometheus datasource provisioning
├── user-service/                    # Сервис пользователей (PostgreSQL, Redis)
├── catalog-service/                 # Сервис каталога (PostgreSQL)
├── payment-service/                 # Сервис обработки платежей (PostgreSQL, Redis, Kafka-producer)
├── notification-service/            # Сервис уведомлений (Kafka-consumer)
└── cart-service/                    # Сервис корзины (Redis)
```
## 📁 Структура сервиса

```text
user-service/
├── cmd/                  # Точки входа (app и мигратор)
│   ├── app/              # Основной бинарник сервиса
│   └── migrator/         # Миграции базы данных
├── internal/             # Внутренние слои (недоступны извне)
│   ├── app/              # Инициализация слоёв: DI, роутер, конфиг
│   ├── delivery/         # HTTP-хендлеры и DTO
│   ├── domain/           # Сущности и интерфейсы (UseCase, Repository)
│   ├── infrastructure/   # Внешние зависимости: БД, Redis, JWT и т.д.
│   ├── usecase/          # Бизнес-логика
│   └── config/           # Загрузка конфигурации из env/config файла
├── migrations/           # SQL-миграции PostgreSQL
├── Dockerfile            # Docker-сборка сервиса
├── go.mod                # Go-модуль сервиса
└── go.work               # Рабочее пространство Go
```

## 🚀 Запуск проекта
Для удобства запуска и управления сервисами используется Makefile. Он позволяет поднимать отдельные сервисы или все сразу.

### 🔧 Подготовка
Сначала создай .env файлы на основе шаблонов:

```bash
cp deploy/user/user.env.example deploy/user/user.env
cp deploy/catalog/catalog.env.example deploy/catalog/catalog.env
cp deploy/cart/cart.env.example deploy/cart/cart.env
cp deploy/payment/payment.env.example deploy/payment/payment.env
cp deploy/notification/notification.env.example deploy/notification/notification.env
cp deploy/kafka/kafka.env.example deploy/kafka/kafka.env
```

### ⚙️ Makefile команды

Упрощают сборку и запуск сервисов.

```bash
# Отдельные сервисы
make user-up          # или catalog-up, cart-up, payment-up, notification-up
make user-down        # или catalog-down, cart-down, payment-down, notification-down

# Все сервисы
make all-up
make all-down

# Мониторинг
make monitoring-up
make monitoring-down

# Kafka стек (Kafka + ZooKeeper + Kafka UI)
make kafka-up
make kafka-down

# Вся инфраструктура сразу (Kafka + Monitoring)
make infra-up
make infra-down

# Docker-сеть (создаётся один раз)
make network-create
```
💡 Все команды определены в корневом Makefile.

## ⚙️ Зависимости
- Go 1.23.0
- Docker, Docker Compose 
- Kafka (взаимодействие между сервисами)
- Prometheus + Grafana (мониторинг Go-сервисов)
- go-chi, sqlx, bcrypt, go-redis, jwt-go, validator.v10

## 🔒 Авторизация

Для защиты эндпоинтов используется middleware:

`RequireRoles(...)`: проверка роли (user, admin)

JWT-токены валидируются через модуль authenticator, реализующий интерфейс Authenticator.

## 📝 Примечания

- Каждый сервис имеет свой go.work и Dockerfile и docker-compose.yml. 
- Все сервисы используют общий pkg/, подключённый через go.work. 
- Для мониторинга используется deploy/monitoring/:
  - Prometheus собирает метрики с /metrics каждого сервиса 
  - Grafana автоматически импортирует дашборды через provisioning 
  - Поддерживаются базовые runtime-метрики (go_*, process_*, promhttp_*)
- Топики Kafka создаются вручную. Они не создаются автоматически при запуске consumer'ов, чтобы избежать неявных ошибок и обеспечить контроль конфигурации.
- Kafka UI (kafka-ui от provectuslabs) разворачивается в docker-compose и доступен по адресу http://localhost:8090 (порт может отличаться в зависимости от конфигурации).
