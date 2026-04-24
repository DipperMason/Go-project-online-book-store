# Litsee Web Demo

Минимальная веб-часть для демонстрации работоспособности платформы.

## Что показывает

- регистрация пользователя (`auth`)
- вход и получение JWT (`auth`)
- чтение/редактирование профиля (`profile`)
- просмотр купленных книг (`profile/books`)
- просмотр активности (`profile/activity`)

## Запуск

```bash
cd web
GOWORK=off go run .
```

По умолчанию UI доступен на `http://localhost:8088`.

## Переменные окружения

- `WEB_ADDR` (default `:8088`)
- `AUTH_URL` (default `http://localhost:8001`)
- `PROFILE_URL` (default `http://localhost:8003`)

## Важно

Web-сервер проксирует браузерные запросы к backend-сервисам, чтобы не упираться в CORS.
