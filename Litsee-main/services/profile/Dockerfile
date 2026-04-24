# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY services/profile/go.mod services/profile/go.mod
COPY pkg/ pkg/

# Копируем исходный код
COPY services/profile/ .

# Скачиваем зависимости
RUN go mod download

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o profile main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/profile .

# Открываем порт
EXPOSE 8003

# Запускаем приложение
CMD ["./profile"]
