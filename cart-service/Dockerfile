# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.work
COPY cart-service/go.work ./cart-service/

# Копируем только файлы зависимостей
COPY cart-service/go.mod cart-service/go.sum ./cart-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/cart-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY cart-service/ ./cart-service/

# Сборка бинарников
WORKDIR /app/cart-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# 2. Финальный минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser/

COPY --from=builder /app/cart-service/app .

RUN chmod +x ./app

CMD sh -c "./app"
