# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.work
COPY notification-service/go.work ./notification-service/

# Копируем только файлы зависимостей
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/notification-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY notification-service/ ./notification-service/

# Сборка бинарников
WORKDIR /app/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# 2. Финальный минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser/

COPY --from=builder /app/notification-service/app .

RUN chmod +x ./app

CMD sh -c "./app"
