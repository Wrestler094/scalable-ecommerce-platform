# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.work
COPY catalog-service/go.work ./catalog-service/

# Копируем только файлы зависимостей
COPY catalog-service/go.mod catalog-service/go.sum ./catalog-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/catalog-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY catalog-service/ ./catalog-service/

# Сборка бинарников
WORKDIR /app/catalog-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser/

COPY --from=builder /app/catalog-service/app .
COPY --from=builder /app/catalog-service/migrator .
COPY --from=builder /app/catalog-service/migrations ./migrations

RUN chmod +x ./app ./migrator

CMD sh -c "./migrator && ./app"
