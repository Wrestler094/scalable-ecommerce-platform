# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/api-gateway
RUN go mod download

# Копируем весь исходный код (и зависимости)
WORKDIR /app
COPY pkg/ ./pkg/
COPY api-gateway/ ./api-gateway/

# Статическая сборка бинарника
WORKDIR /app/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# 2. Финальный минимальный образ
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем собранный бинарник из builder-этапа
COPY --from=builder /app/api-gateway/app .

# Делаем бинарник исполняемым внутри финального контейнера
RUN chmod +x ./app

# Команда запуска
CMD sh -c "./app"
