# 1. Сборочный контейнер
FROM golang:1.23.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY catalog-service/go.mod catalog-service/go.sum ./catalog-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Загружаем зависимости
WORKDIR /app/catalog-service
RUN go mod download

# Копируем весь исходный код (и зависимости)
WORKDIR /app
COPY pkg/ ./pkg/
COPY catalog-service/ ./catalog-service/

# Статическая сборка бинарников
WORKDIR /app/catalog-service
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

# 2. Финальный минимальный образ
FROM alpine:latest

# Добавляем корневые сертификаты (для HTTPS-запросов из приложения)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /home/appuser/

# Копируем собранные бинарники из builder-этапа
COPY --from=builder /app/catalog-service/app .
COPY --from=builder /app/catalog-service/migrator .
COPY --from=builder /app/catalog-service/migrations ./migrations

# Делаем бинарники исполняемыми внутри финального контейнера
RUN chmod +x ./app ./migrator

# Команда запуска (миграции + приложение)
CMD sh -c "./migrator && ./app"
