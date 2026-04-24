# Profile Service - Личный Кабинет

Полнофункциональный микросервис для управления профилем пользователя в системе Litsee.

## 📋 Функциональность

✅ **Управление профилем:**
- Просмотр данных профиля пользователя
- Редактирование информации профиля (имя, фамилия, аватар, биография)
- Изменение пароля с проверкой старого пароля

✅ **Управление сессиями:**
- Список активных сессий пользователя
- Завершение отдельной сессии
- Завершение всех сессий пользователя (выход со всех устройств)

✅ **История активности:**
- Логирование всех действий пользователя (login, logout, profile_update, password_change)
- Просмотр истории с указанием IP адреса, User-Agent, даты и времени
- Пагинация для больших объемов данных

✅ **Безопасность:**
- JWT аутентификация
- Хеширование паролей с bcrypt
- Валидация всех входных данных
- Логирование всех операций

## 🚀 Быстрый старт

### Требования
- Go 1.21+
- PostgreSQL 12+ (опционально, есть in-memory реализация)

### Локальный запуск

```bash
cd services/profile
chmod +x run.sh
./run.sh dev
```

Или напрямую:
```bash
go run main.go
```

Сервис запустится на `http://localhost:8003`

### Docker запуск

```bash
./run.sh docker-build
./run.sh docker-run
```

Или через Docker Compose:
```bash
./run.sh docker-compose
```

## 📁 Структура проекта

```
profile/
├── main.go                                 # Точка входа приложения
├── go.mod                                  # Модульные зависимости
├── run.sh                                  # Скрипт для запуска
├── README.md                               # Этот файл
├── EXAMPLES.md                             # Примеры использования API
├── Dockerfile                              # Docker конфигурация
├── docker-compose.yml                      # Docker Compose конфигурация
├── postman_collection.json                 # Коллекция Postman для тестирования
├── migrations/
│   └── 001_init.sql                       # SQL миграции для PostgreSQL
└── internal/
    ├── app/
    │   └── app.go                         # Инициализация приложения
    ├── config/
    │   └── config.go                      # Конфигурация сервиса
    ├── domain/
    │   ├── user_profile.go                # Модель профиля пользователя
    │   ├── session.go                     # Модель сессии
    │   ├── activity_log.go                # Модель логов активности
    │   └── errors.go                      # Определение ошибок
    ├── repo/
    │   ├── profile_repo.go                # Интерфейс и mock репозитория профилей
    │   ├── profile_postgres_repo.go       # PostgreSQL реализация профилей
    │   ├── session_repo.go                # Интерфейс и mock репозитория сессий
    │   ├── session_postgres_repo.go       # PostgreSQL реализация сессий
    │   ├── activity_repo.go               # Интерфейс и mock репозитория активности
    │   └── activity_postgres_repo.go      # PostgreSQL реализация активности
    ├── services/
    │   ├── profile_service.go             # Бизнес-логика профиля
    │   └── profile_service_test.go        # Unit тесты
    └── transport/
        └── http/
            ├── handler.go                 # HTTP обработчики
            ├── responses.go               # Стандартные HTTP ответы
            └── auth_middleware.go         # JWT middleware
```

## 🔌 API Endpoints

### Профиль
- `GET /api/v1/profile` - Получить профиль пользователя
- `PUT /api/v1/profile` - Обновить профиль пользователя

### Пароль
- `POST /api/v1/profile/password` - Изменить пароль

### Сессии
- `GET /api/v1/profile/sessions` - Получить активные сессии
- `DELETE /api/v1/profile/sessions/{id}` - Завершить одну сессию
- `POST /api/v1/profile/sessions/terminate-all` - Завершить все сессии

### История активности
- `GET /api/v1/profile/activity?limit=50&offset=0` - Получить историю активности

Полная документация API: смотрите [README.md](README.md)

## 🧪 Тестирование

### Unit тесты
```bash
cd services/profile
go test -v ./...
```

### С покрытием
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Постман тесты
1. Импортируйте `postman_collection.json` в Postman
2. Установите переменные:
   - `token` - ваш JWT токен
   - `session_id` - ID сессии для тестирования

### cURL примеры
```bash
# Получить профиль
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8003/api/v1/profile

# Обновить профиль
curl -X PUT http://localhost:8003/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe"}'
```

Смотрите [EXAMPLES.md](EXAMPLES.md) для полных примеров

## 🗄️ База данных

### Установка PostgreSQL (опционально)

```bash
# macOS
brew install postgresql@15

# Linux (Ubuntu/Debian)
sudo apt-get install postgresql postgresql-contrib

# Windows
# Скачайте с https://www.postgresql.org/download/windows/
```

### Создание БД
```bash
createdb litsee_profile
psql -d litsee_profile < migrations/001_init.sql
```

### Переключение на PostgreSQL

В `internal/app/app.go`:
```go
// Вместо Mock репозиториев используйте PostgreSQL
db, _ := sql.Open("postgres", "user=postgres password=password dbname=litsee_profile")
profileRepo := repo.NewPostgresProfileRepository(db)
sessionRepo := repo.NewPostgresSessionRepository(db)
activityRepo := repo.NewPostgresActivityRepository(db)
```

## 🔧 Конфигурация

### Переменные окружения
- `PORT` - Порт сервиса (default: 8003)
- `LOG_LEVEL` - Уровень логирования: debug, info, warn, error (default: info)

### Пример .env файла
```env
PORT=8003
LOG_LEVEL=info
DATABASE_URL=postgres://user:password@localhost:5432/litsee_profile
JWT_SECRET=your-secret-key
```

## 📦 Сборка и развертывание

### Локальная сборка
```bash
go build -o profile main.go
./profile
```

### Release сборка
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o profile_linux_amd64 main.go
```

### Docker
```bash
docker build -t profile-service:latest .
docker run -p 8003:8003 profile-service:latest
```

### Kubernetes (пример deployment.yml)
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
        env:
        - name: PORT
          value: "8003"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

## 🔐 Безопасность

- ✅ JWT токены обязательны для всех endpoints (кроме публичных)
- ✅ Пароли хешируются с bcrypt (cost: 10)
- ✅ Все входные данные валидируются
- ✅ Логирование всех операций
- ✅ HTTPS рекомендуется для production
- ✅ Rate limiting рекомендуется добавить в production

## 📊 Интеграция с микросервисами

### API Gateway
Добавьте маршрут в ваш API Gateway:
```yaml
routes:
  - path: /api/v1/profile
    service: profile-service
    port: 8003
```

### Интеграция с Auth Service
Profile Service использует JWT токены от Auth Service. Убедитесь что токены совместимы.

## 📝 Лицензия

Copyright 2026. Все права защищены.

## 🤝 Вклад

Для разработки новых функций:

1. Создайте feature branch: `git checkout -b feature/amazing-feature`
2. Коммитьте изменения: `git commit -m 'Add amazing feature'`
3. Пушьте в branch: `git push origin feature/amazing-feature`
4. Создайте Pull Request

## 📞 Поддержка

При возникновении проблем:
1. Проверьте [EXAMPLES.md](EXAMPLES.md)
2. Смотрите логи: `tail -f /var/log/profile-service.log`
3. Проверьте подключение к БД (если используется PostgreSQL)
4. Убедитесь что JWT токены правильные

## 🎯 TODO

- [ ] Добавить Redis для кэширования активных сессий
- [ ] Двухфакторная аутентификация (2FA)
- [ ] Экспорт данных профиля (JSON, CSV)
- [ ] Мобильное приложение
- [ ] Графики активности
- [ ] Email уведомления о новых логинах
- [ ] Восстановление аккаунта
- [ ] Смена email адреса
