# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.work
COPY user-service/go.work ./user-service/

# Копируем только файлы зависимостей для кэширования
COPY user-service/go.mod user-service/go.sum ./user-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Предварительная загрузка зависимостей
WORKDIR /app/user-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY user-service/ ./user-service/

# Статическая сборка двух независимых бинарников
WORKDIR /app/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ для продакшена
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем два собранных бинарника из builder-этапа
COPY --from=builder /app/user-service/app .
COPY --from=builder /app/user-service/migrator .
COPY --from=builder /app/user-service/migrations ./migrations

RUN chmod +x ./app ./migrator

CMD sh -c "./migrator && ./app"
