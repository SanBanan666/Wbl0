# Используем официальный образ Go 1.23
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Настраиваем GOPROXY для обхода блокировки
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./Wbl0/cmd/api

# Используем минимальный образ для финального контейнера
FROM alpine:3.19

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из builder stage
COPY --from=builder /app/main .

# Меняем владельца файла
RUN chown appuser:appgroup main

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./main"]