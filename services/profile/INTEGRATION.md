# Profile Service - Интеграция с go.work

## Интеграция в основной проект

### Шаг 1: Проверка go.work

Убедитесь что в главном `go.work` файле содержится:

```
go 1.21

use (
    .
    ./pkg/config
    ./pkg/jwt
    ./pkg/logger
    ./pkg/middlewares
    ./pkg/serde
    ./services/auth
    ./services/account
    ./services/catalog
    ./services/gateway
    ./services/notifier
    ./services/order
    ./services/profile  // <- Добавьте эту строку
)
```

### Шаг 2: Обновление go.work

Выполните:
```bash
cd /home/maks/Litsee
go work sync
```

### Шаг 3: Запуск Profile Service

```bash
cd services/profile
chmod +x run.sh check.sh
./check.sh  # Проверка структуры
./run.sh dev  # Запуск в development режиме
```

## Интеграция с API Gateway

### Конфигурация маршрутов

В файле конфигурации Gateway добавьте:

```yaml
routes:
  - path: /api/v1/profile
    service: profile
    port: 8003
    methods: [GET, PUT, POST, DELETE]
    authentication: jwt
```

Или через environment переменные:

```bash
PROFILE_SERVICE_URL=http://localhost:8003
```

## Интеграция с другими сервисами

### Auth Service

Profile Service использует JWT токены от Auth Service. После успешной аутентификации, клиент получает токен:

```bash
# 1. Регистрация
curl -X POST http://localhost:8001/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# 2. Аутентификация
curl -X POST http://localhost:8001/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Получите token из ответа

# 3. Используйте токен для Profile Service
curl -H "Authorization: Bearer $token" http://localhost:8003/api/v1/profile
```

### Notification Service (будущая интеграция)

Profile Service может отправлять события на Notification Service:

```go
// Пример:
client := &http.Client{}
request, _ := http.NewRequest("POST", 
    "http://localhost:8007/api/v1/send-notification",
    bytes.NewBuffer([]byte(`{
        "user_id": 1,
        "type": "profile_update",
        "message": "Your profile was updated"
    }`)))
client.Do(request)
```

## Мониторинг и логирование

### Включение центрального логирования

```bash
# Настройте логирование через Loki/ELK
export LOG_LEVEL=info
export LOG_FORMAT=json
./run.sh dev
```

### Метрики Prometheus (опционально)

Добавьте в `main.go`:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

func init() {
    http.Handle("/metrics", promhttp.Handler())
}
```

## Переменные окружения для production

Создайте файл `.env.production`:

```env
PORT=8003
LOG_LEVEL=warn
DATABASE_URL=postgres://user:password@db.example.com:5432/litsee_profile
JWT_SECRET=your-super-secret-key-here
ALLOWED_ORIGINS=https://example.com,https://app.example.com
```

## Развертывание в Docker

### Использование Docker Compose для всех сервисов

Обновите корневой `docker-compose.yml`:

```yaml
version: '3.8'

services:
  # ... другие сервисы ...
  
  profile:
    build:
      context: .
      dockerfile: services/profile/Dockerfile
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
      - LOG_LEVEL=info
      - DATABASE_URL=postgres://postgres:password@db:5432/litsee_profile
    depends_on:
      - db
    networks:
      - litsee
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=litsee_profile
    volumes:
      - db_data:/var/lib/postgresql/data
      - ./services/profile/migrations:/docker-entrypoint-initdb.d
    networks:
      - litsee

volumes:
  db_data:

networks:
  litsee:
    driver: bridge
```

Затем запустите:

```bash
docker-compose up -d profile
```

## Kubernetes развертывание

### ConfigMap для конфигурации

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: profile-service-config
data:
  PORT: "8003"
  LOG_LEVEL: "info"
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: profile-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: profile-service
  template:
    metadata:
      labels:
        app: profile-service
    spec:
      containers:
      - name: profile
        image: profile-service:latest
        ports:
        - containerPort: 8003
          name: http
        envFrom:
        - configMapRef:
            name: profile-service-config
        livenessProbe:
          httpGet:
            path: /api/v1/profile
            port: 8003
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/profile
            port: 8003
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: profile-service
spec:
  selector:
    app: profile-service
  ports:
  - port: 8003
    targetPort: 8003
  type: ClusterIP
```

Применить:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

## Тестирование интеграции

### E2E тест

```bash
#!/bin/bash

# 1. Запустить все сервисы
docker-compose up -d

# 2. Подождать 5 секунд
sleep 5

# 3. Пройти регистрацию
TOKEN=$(curl -s -X POST http://localhost:8001/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq -r '.token')

# 4. Получить профиль
curl -H "Authorization: Bearer $TOKEN" http://localhost:8003/api/v1/profile

# 5. Остановить сервисы
docker-compose down
```

## Troubleshooting

### Profile Service не отвечает

```bash
# Проверить логи
docker logs profile-service

# Проверить что сервис слушает на порту
lsof -i :8003

# Проверить конфигурацию БД
psql -h localhost -U postgres -d litsee_profile -c "SELECT 1"
```

### Ошибки JWT аутентификации

```bash
# Убедитесь что токен есть в Authorization заголовке
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8003/api/v1/profile

# Проверить что токен не истекший
echo $TOKEN | jq -R 'split(".") | .[1] | @base64d'
```

### Проблемы с БД

```bash
# Проверить что все миграции применены
psql -d litsee_profile -c "\dt"

# Применить миграции заново
psql -d litsee_profile < migrations/001_init.sql
```

## Дополнительные ресурсы

- [README.md](README.md) - Полная документация API
- [EXAMPLES.md](EXAMPLES.md) - Примеры использования
- [STRUCTURE.md](STRUCTURE.md) - Описание структуры проекта
