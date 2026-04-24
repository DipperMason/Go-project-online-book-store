# Примеры использования Profile Service API

## Локальное тестирование

### 1. Запуск сервиса

```bash
cd services/profile
go run main.go
```

Сервис запустится на `http://localhost:8003`

### 2. Установка зависимостей

```bash
go mod download
go mod tidy
```

### 3. Сборка Docker образа

```bash
docker build -t profile-service:latest .
docker run -p 8003:8003 profile-service:latest
```

## Примеры API вызовов

### Переменные для тестирования

```bash
# Установите переменные окружения
export TOKEN="your_jwt_token_here"
export BASE_URL="http://localhost:8003"
export USER_ID="1"
```

### 1. Получить профиль

```bash
curl -X GET "$BASE_URL/api/v1/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**Response:**
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

### 2. Обновить профиль

```bash
curl -X PUT "$BASE_URL/api/v1/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "avatar": "https://example.com/new-avatar.jpg",
    "bio": "Updated bio"
  }'
```

### 3. Изменить пароль

```bash
curl -X POST "$BASE_URL/api/v1/profile/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldPassword123",
    "new_password": "newPassword456"
  }'
```

### 4. Получить активные сессии

```bash
curl -X GET "$BASE_URL/api/v1/profile/sessions" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
[
  {
    "id": "session-uuid-1",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "created_at": "2026-04-20T10:00:00Z",
    "expires_at": "2026-04-21T10:00:00Z",
    "is_active": true
  }
]
```

### 5. Завершить одну сессию

```bash
curl -X DELETE "$BASE_URL/api/v1/profile/sessions/session-uuid-1" \
  -H "Authorization: Bearer $TOKEN"
```

### 6. Завершить все сессии

```bash
curl -X POST "$BASE_URL/api/v1/profile/sessions/terminate-all" \
  -H "Authorization: Bearer $TOKEN"
```

### 7. Получить историю активности

```bash
# Получить первые 50 записей
curl -X GET "$BASE_URL/api/v1/profile/activity?limit=50&offset=0" \
  -H "Authorization: Bearer $TOKEN"

# Получить следующие 50 записей
curl -X GET "$BASE_URL/api/v1/profile/activity?limit=50&offset=50" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
[
  {
    "id": 1,
    "action": "login",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "status": "success",
    "details": "Login from web",
    "created_at": "2026-04-20T10:00:00Z"
  },
  {
    "id": 2,
    "action": "profile_update",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "status": "success",
    "details": "Updated profile",
    "created_at": "2026-04-20T11:00:00Z"
  }
]
```

## Тестирование с использованием Postman

1. Откройте Postman
2. Импортируйте файл `postman_collection.json`
3. Установите переменные в Collection:
   - `token` - ваш JWT токен
   - `session_id` - ID сессии для тестирования
4. Запустите нужные запросы

## Интеграция с фронтенд

### JavaScript/TypeScript пример

```javascript
const API_BASE = 'http://localhost:8003';
const token = localStorage.getItem('auth_token');

// Получить профиль
async function getProfile() {
  const response = await fetch(`${API_BASE}/api/v1/profile`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch profile');
  }
  
  return response.json();
}

// Обновить профиль
async function updateProfile(data) {
  const response = await fetch(`${API_BASE}/api/v1/profile`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });
  
  if (!response.ok) {
    throw new Error('Failed to update profile');
  }
  
  return response.json();
}

// Получить активные сессии
async function getActiveSessions() {
  const response = await fetch(`${API_BASE}/api/v1/profile/sessions`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch sessions');
  }
  
  return response.json();
}

// Получить историю активности
async function getActivityHistory(limit = 50, offset = 0) {
  const response = await fetch(
    `${API_BASE}/api/v1/profile/activity?limit=${limit}&offset=${offset}`,
    {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to fetch activity history');
  }
  
  return response.json();
}

// Изменить пароль
async function changePassword(oldPassword, newPassword) {
  const response = await fetch(`${API_BASE}/api/v1/profile/password`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      old_password: oldPassword,
      new_password: newPassword
    })
  });
  
  if (!response.ok) {
    throw new Error('Failed to change password');
  }
  
  return response.json();
}
```

## Обработка ошибок

Все ошибки возвращаются в стандартном формате:

```json
{
  "error": "Описание ошибки",
  "status": 400
}
```

Коды ошибок:
- `400` - Invalid request (неверные параметры)
- `401` - Unauthorized (отсутствует или неверный токен)
- `404` - Not found (ресурс не найден)
- `500` - Internal server error

## Тестирование с помощью Go

```go
package main

import (
    "testing"
    "profile/internal/domain"
    "profile/internal/repo"
    "profile/internal/services"
    "context"
)

func TestProfileService(t *testing.T) {
    // Создание моков
    profileRepo := repo.NewMockProfileRepository()
    sessionRepo := repo.NewMockSessionRepository()
    activityRepo := repo.NewMockActivityRepository()
    
    // Создание сервиса
    service := services.NewProfileService(profileRepo, sessionRepo, activityRepo)
    
    // Тестирование получения профиля
    profile, err := service.GetProfile(context.Background(), 1)
    if err != nil {
        t.Fatalf("Failed to get profile: %v", err)
    }
    
    if profile.ID != 1 {
        t.Errorf("Expected ID 1, got %d", profile.ID)
    }
}
```

## Запуск тестов

```bash
cd services/profile
go test ./...
go test -v ./...
go test -cover ./...
```

## Отладка

### Включить подробное логирование

```bash
export LOG_LEVEL=debug
./profile
```

### Проверить порт

```bash
# Linux/Mac
lsof -i :8003

# Windows
netstat -ano | findstr :8003
```

### Прямое подключение к БД (PostgreSQL)

```bash
psql -h localhost -U postgres -d litsee_profile
\dt  -- показать все таблицы
SELECT * FROM user_profiles;
SELECT * FROM sessions;
SELECT * FROM activity_logs;
```
