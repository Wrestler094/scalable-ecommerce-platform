# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/api-gateway
RUN go mod download

# Копируем исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY api-gateway/ ./api-gateway/

# Сборка бинарника
WORKDIR /app/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# 2. Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser/

COPY --from=builder /app/api-gateway/app .

RUN chmod +x ./app

CMD ["./app"]
