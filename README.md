# 🧱 Scalable E-Commerce Platform

Микросервисная e-commerce платформа на Go, Docker и Redis/PostgreSQL. Разделена на отдельные сервисы с возможностью масштабирования и независимого деплоя.

## 📦 Сервисы

- `user-service` — регистрация, авторизация
- `catalog-service` — продукты, категории
- `pkg/` — общие модули: `authenticator`, `roles`, `httphelper`, `validator` и др.

## 📁 Структура проекта

```text
scalable-ecommerce-platform/
├── Makefile
├── pkg/
│   ├── authenticator/
│   ├── httphelper/
│   ├── ...
│   └── go.mod
├── deploy/
│   ├── docker-compose.catalog.yml
│   ├── docker-compose.user.yml
│   └── envs/
│       ├── catalog.env
│       └── user.env
├── user-service/
│   ├── cmd/
│   ├── internal/
│   ├── migrations
│   ├── Dockerfile
│   ├── go.mod
│   └── go.work
└── catalog-service/
│   ├── cmd/
│   ├── internal/
│   ├── migrations
│   ├── Dockerfile
│   ├── go.mod
│   └── go.work
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

- Каждый сервис имеет свой go.work, содержащий ., ../pkg
- Каждый сервис имеет свой Dockerfile. Все сервисы используют общий pkg/, подключённый через go.work.
