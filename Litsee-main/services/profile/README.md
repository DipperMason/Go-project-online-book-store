# Profile Service - Личный Кабинет

Микросервис для управления профилем пользователя и его активностью.

## Установка и запуск

```bash
cd services/profile
go build -o profile main.go
./profile
```

По умолчанию сервис запускается на `http://localhost:8003`

## API Endpoints

### 1. Профиль пользователя

#### Получить профиль
```
GET /api/v1/profile
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "avatar": "https://example.com/avatar.jpg",
  "bio": "Hello, I'm John!",
  "updated_at": "2026-04-20T12:00:00Z",
  "created_at": "2026-04-01T10:00:00Z"
}
```

#### Обновить профиль
```
PUT /api/v1/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "avatar": "https://example.com/new-avatar.jpg",
  "bio": "Updated bio"
}
```

**Response (200):**
```json
{
  "message": "Profile updated successfully"
}
```

### 2. Смена пароля

#### Изменить пароль
```
POST /api/v1/profile/password
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "old_password": "currentPassword123",
  "new_password": "newPassword456"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Error Responses:**
- `400` - Invalid old password
- `400` - Invalid request body
- `401` - Unauthorized

### 3. Управление сессиями

#### Получить активные сессии
```
GET /api/v1/profile/sessions
Authorization: Bearer <token>
```

**Response (200):**
```json
[
  {
    "id": "session-uuid-1",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "created_at": "2026-04-20T10:00:00Z",
    "expires_at": "2026-04-21T10:00:00Z",
    "is_active": true
  },
  {
    "id": "session-uuid-2",
    "ip_address": "192.168.1.2",
    "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6)...",
    "created_at": "2026-04-19T15:30:00Z",
    "expires_at": "2026-04-20T15:30:00Z",
    "is_active": true
  }
]
```

#### Завершить одну сессию
```
DELETE /api/v1/profile/sessions/{id}
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "message": "Session terminated"
}
```

#### Завершить все сессии
```
POST /api/v1/profile/sessions/terminate-all
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "message": "All sessions terminated"
}
```

### 4. История активности

#### Получить историю активности
```
GET /api/v1/profile/activity?limit=50&offset=0
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit` (int, optional) - Количество записей на странице (default: 50, max: 100)
- `offset` (int, optional) - Смещение от начала (default: 0)

**Response (200):**
```json
[
  {
    "id": 1,
    "action": "login",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "status": "success",
    "details": "Login from web",
    "created_at": "2026-04-20T10:00:00Z"
  },
  {
    "id": 2,
    "action": "profile_update",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "status": "success",
    "details": "Updated first name and bio",
    "created_at": "2026-04-20T11:15:00Z"
  },
  {
    "id": 3,
    "action": "password_change",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "status": "success",
    "details": "",
    "created_at": "2026-04-20T12:30:00Z"
  }
]
```

## Возможные действия (Actions)

- `login` - Вход в систему
- `logout` - Выход из системы
- `profile_update` - Обновление профиля
- `password_change` - Смена пароля
- `session_create` - Создание новой сессии
- `session_terminate` - Завершение сессии

## Коды ошибок

- `400` - Bad Request (invalid parameters or body)
- `401` - Unauthorized (missing or invalid token)
- `404` - Not Found
- `500` - Internal Server Error

## Интеграция с JWT

Сервис профиля использует JWT токены для аутентификации. Токен передается в заголовке `Authorization` в формате:

```
Authorization: Bearer <jwt_token>
```

Токен должен содержать поле `user_id` в payload, которое используется для идентификации пользователя.

## Архитектура

```
profile/
├── main.go                           # Точка входа
├── go.mod                            # Модульные зависимости
└── internal/
    ├── app/
    │   └── app.go                   # Инициализация приложения
    ├── config/
    │   └── config.go                # Конфигурация
    ├── domain/
    │   ├── user_profile.go          # Модель профиля
    │   ├── session.go               # Модель сессии
    │   ├── activity_log.go          # Модель логов активности
    │   └── errors.go                # Определение ошибок
    ├── repo/
    │   ├── profile_repo.go          # Репозиторий профилей
    │   ├── session_repo.go          # Репозиторий сессий
    │   └── activity_repo.go         # Репозиторий активности
    ├── services/
    │   └── profile_service.go       # Бизнес-логика
    └── transport/
        └── http/
            └── handler.go           # HTTP обработчики
```

## Структура данных

### UserProfile
```go
type UserProfile struct {
    ID        int
    Email     string
    FirstName string
    LastName  string
    Avatar    string
    Bio       string
    UpdatedAt time.Time
    CreatedAt time.Time
}
```

### Session
```go
type Session struct {
    ID        string
    UserID    int
    Token     string
    IPAddress string
    UserAgent string
    CreatedAt time.Time
    ExpiresAt time.Time
    IsActive  bool
}
```

### ActivityLog
```go
type ActivityLog struct {
    ID        int
    UserID    int
    Action    string
    IPAddress string
    UserAgent string
    Status    string
    Details   string
    CreatedAt time.Time
}
```

## Развертывание

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o profile main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/profile .
EXPOSE 8003
CMD ["./profile"]
```

## Переменные окружения

- `PORT` - Порт для запуска сервиса (default: 8003)
- `LOG_LEVEL` - Уровень логирования (default: info)

## Примеры использования

### cURL

```bash
# Получить профиль
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8003/api/v1/profile

# Обновить профиль
curl -X PUT http://localhost:8003/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "bio": "New bio"
  }'

# Получить активные сессии
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8003/api/v1/profile/sessions

# Получить историю активности
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8003/api/v1/profile/activity?limit=10&offset=0
```
