# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY payment-service/go.mod payment-service/go.sum ./payment-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/payment-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY payment-service/ ./payment-service/

# Статическая сборка бинарников
WORKDIR /app/payment-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем собранные бинарники из builder-этапа
COPY --from=builder /app/payment-service/app .
COPY --from=builder /app/payment-service/migrator .
COPY --from=builder /app/payment-service/migrations ./migrations

RUN chmod +x ./app ./migrator

# Команда запуска (миграции + приложение)
CMD sh -c "./migrator && ./app"
