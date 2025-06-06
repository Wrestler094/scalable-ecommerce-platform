# 🧱 Scalable E-Commerce Platform

Микросервисная e-commerce платформа на Go, Docker и Redis/PostgreSQL. Разделена на отдельные сервисы с возможностью масштабирования и независимого деплоя.

## 📦 Сервисы

- `user-service` — регистрация, авторизация
- `catalog-service` — продукты, категории
- `cart-service` — корзина пользователя
- `pkg/` — общие модули: `authenticator`, `roles`, `httphelper`, `validator` и др.

## 📁 Структура проекта

```text
scalable-ecommerce-platform/
├── Makefile                         # Общие команды запуска
├── pkg/                             # Общие переиспользуемые модули
│   ├── authenticator/               # JWT, контекст, middleware
│   ├── httphelper/                  # Обработка ошибок, JSON-ответов
│   └── ...
├── deploy/                          # Docker Compose и env-файлы
│   ├── docker-compose.user.yml
│   ├── docker-compose.catalog.yml
│   ├── docker-compose.cart.yml
│   └── envs/
│       ├── user.env
│       ├── catalog.env
│       └── cart.env
├── user-service/                    # Сервис пользователей (JWT, Redis)
├── catalog-service/                 # Сервис каталога (PostgreSQL)
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

## ⚙️ Makefile команды

Упрощают сборку и запуск сервисов.

```bash
# Запустить только user-service
make user-up

# Остановить user-service
make user-down

# Запустить только catalog-service
make catalog-up

# Остановить catalog-service
make catalog-down

# Запустить оба сервиса
make all-up

# Остановить оба сервиса
make all-down
```

## ⚙️ Зависимости
- Go 1.23.0
- Docker, Docker Compose
- Redis, PostgreSQL
- go-chi, bcrypt, go-redis, jwt-go

## 🔒 Авторизация

Для защиты эндпоинтов используется middleware:

`RequireRoles(...)`: проверка роли (user, admin)

## 📝 Примечания

- Каждый сервис имеет свой go.work и Dockerfile. 
- Все сервисы используют общий pkg/, подключённый через go.work.
