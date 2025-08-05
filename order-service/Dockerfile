# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только файлы зависимостей
COPY order-service/go.mod order-service/go.sum ./order-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/order-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY order-service/ ./order-service/

# Сборка бинарников
WORKDIR /app/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем два собранных бинарника из builder-этапа
COPY --from=builder /app/order-service/app .
COPY --from=builder /app/order-service/migrator .
COPY --from=builder /app/order-service/migrations ./migrations

RUN chmod +x ./app ./migrator

# Команда запуска (миграции + приложение)
CMD sh -c "./migrator && ./app"
