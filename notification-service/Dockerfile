# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/notification-service
RUN go mod download

# Копируем весь исходный код
WORKDIR /app
COPY pkg/ ./pkg/
COPY notification-service/ ./notification-service/

# Статическая сборка бинарника
WORKDIR /app/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# 2. Финальный минимальный образ
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем собранный бинарник из builder-этапа
COPY --from=builder /app/notification-service/app .

RUN chmod +x ./app

# Команда запуска
CMD sh -c "./app"
