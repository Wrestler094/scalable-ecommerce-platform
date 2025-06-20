# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.work
COPY payment-service/go.work ./payment-service/

# Копируем только файлы зависимостей для кэширования
COPY payment-service/go.mod payment-service/go.sum ./payment-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Предварительная загрузка зависимостей
WORKDIR /app/payment-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY payment-service/ ./payment-service/

# Статическая сборка двух независимых бинарников
WORKDIR /app/payment-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ для продакшена
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем два собранных бинарника из builder-этапа
COPY --from=builder /app/payment-service/app .
COPY --from=builder /app/payment-service/migrator .
COPY --from=builder /app/payment-service/migrations ./migrations

RUN chmod +x ./app ./migrator

CMD sh -c "./migrator && ./app"
