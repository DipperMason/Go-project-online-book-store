FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod .

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

# Второй этап: минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из первого этапа
COPY --from=builder /app/gateway .

EXPOSE 8000

CMD ["./gateway"]