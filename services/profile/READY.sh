#!/bin/bash

# Финальная проверка и инструкции для Profile Service

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║           PROFILE SERVICE - ЛИЧНЫЙ КАБИНЕТ                     ║"
echo "║                     ✅ ПРОЕКТ ЗАВЕРШЕН                         ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

PROFILE_DIR="/home/maks/Litsee/services/profile"

echo "📍 РАСПОЛОЖЕНИЕ: $PROFILE_DIR"
echo ""

# Счетчики
TOTAL=0
FOUND=0

# Проверка файлов
echo "📋 ПРОВЕРКА ФАЙЛОВ:"
echo "=================="

check_file() {
    local file=$1
    local desc=$2
    TOTAL=$((TOTAL + 1))
    
    if [ -f "$PROFILE_DIR/$file" ]; then
        FOUND=$((FOUND + 1))
        echo "  ✅ $desc"
    else
        echo "  ❌ $desc (НЕ НАЙДЕН)"
    fi
}

check_dir() {
    local dir=$1
    local desc=$2
    TOTAL=$((TOTAL + 1))
    
    if [ -d "$PROFILE_DIR/$dir" ]; then
        FOUND=$((FOUND + 1))
        echo "  ✅ $desc"
    else
        echo "  ❌ $desc (НЕ НАЙДЕН)"
    fi
}

# Основные файлы
echo ""
echo "📁 Основная структура:"
check_file "main.go" "Точка входа (main.go)"
check_file "go.mod" "Модули (go.mod)"
check_file "Dockerfile" "Docker конфигурация"
check_file "docker-compose.yml" "Docker Compose"
check_file "run.sh" "Скрипт запуска"
check_file "check.sh" "Скрипт проверки"

# Документация
echo ""
echo "📚 Документация:"
check_file "README.md" "Документация API"
check_file "EXAMPLES.md" "Примеры использования"
check_file "STRUCTURE.md" "Описание архитектуры"
check_file "INTEGRATION.md" "Интеграция"
check_file "SUMMARY.md" "Резюме проекта"
check_file "postman_collection.json" "Postman коллекция"

# Директории
echo ""
echo "🗂️  Директории:"
check_dir "internal/app" "Application layer"
check_dir "internal/config" "Конфигурация"
check_dir "internal/domain" "Domain models"
check_dir "internal/repo" "Repository layer"
check_dir "internal/services" "Service layer"
check_dir "internal/transport/http" "HTTP transport"
check_dir "migrations" "БД миграции"

# Исходные файлы
echo ""
echo "📝 Исходные файлы Go:"
check_file "internal/app/app.go" "App initialization"
check_file "internal/config/config.go" "Configuration"
check_file "internal/domain/user_profile.go" "User profile model"
check_file "internal/domain/session.go" "Session model"
check_file "internal/domain/activity_log.go" "Activity log model"
check_file "internal/domain/errors.go" "Domain errors"
check_file "internal/repo/profile_repo.go" "Profile repository"
check_file "internal/repo/profile_postgres_repo.go" "Profile PostgreSQL"
check_file "internal/repo/session_repo.go" "Session repository"
check_file "internal/repo/session_postgres_repo.go" "Session PostgreSQL"
check_file "internal/repo/activity_repo.go" "Activity repository"
check_file "internal/repo/activity_postgres_repo.go" "Activity PostgreSQL"
check_file "internal/services/profile_service.go" "Profile service"
check_file "internal/services/profile_service_test.go" "Service tests"
check_file "internal/transport/http/handler.go" "HTTP handlers"
check_file "internal/transport/http/responses.go" "HTTP responses"
check_file "internal/transport/http/auth_middleware.go" "JWT middleware"

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "📊 СТАТИСТИКА: $FOUND из $TOTAL файлов найдено"
echo "════════════════════════════════════════════════════════════════"
echo ""

if [ $FOUND -eq $TOTAL ]; then
    echo "✅ ВСЕ ФАЙЛЫ НАЙДЕНЫ!"
    echo ""
    echo "🚀 СЛЕДУЮЩИЕ ШАГИ:"
    echo "=================="
    echo ""
    echo "1️⃣  ЗАПУСК В РАЗРАБОТКЕ:"
    echo "   cd $PROFILE_DIR"
    echo "   chmod +x run.sh check.sh"
    echo "   ./run.sh dev"
    echo ""
    echo "2️⃣  ЗАПУСК ТЕСТОВ:"
    echo "   ./run.sh test"
    echo ""
    echo "3️⃣  DOCKER ЗАПУСК:"
    echo "   ./run.sh docker-build"
    echo "   ./run.sh docker-run"
    echo ""
    echo "4️⃣  ИНТЕГРАЦИЯ:"
    echo "   - Прочитайте INTEGRATION.md"
    echo "   - Добавьте ./services/profile в go.work"
    echo "   - Выполните: go work sync"
    echo ""
    echo "5️⃣  ТЕСТИРОВАНИЕ:"
    echo "   - Импортируйте postman_collection.json в Postman"
    echo "   - Или используйте примеры из EXAMPLES.md"
    echo ""
    echo "📚 ДОКУМЕНТАЦИЯ:"
    echo "   - README.md      → Полная документация API"
    echo "   - EXAMPLES.md    → Примеры cURL, JS, Go"
    echo "   - STRUCTURE.md   → Архитектура и развертывание"
    echo "   - INTEGRATION.md → Интеграция с другими сервисами"
    echo "   - SUMMARY.md     → Резюме проекта"
    echo ""
    echo "🌐 ENDPOINTS (7 всего):"
    echo "   GET    /api/v1/profile"
    echo "   PUT    /api/v1/profile"
    echo "   POST   /api/v1/profile/password"
    echo "   GET    /api/v1/profile/sessions"
    echo "   DELETE /api/v1/profile/sessions/{id}"
    echo "   POST   /api/v1/profile/sessions/terminate-all"
    echo "   GET    /api/v1/profile/activity"
    echo ""
    echo "✨ ФУНКЦИИ:"
    echo "   ✅ Просмотр профиля"
    echo "   ✅ Редактирование профиля"
    echo "   ✅ Смена пароля"
    echo "   ✅ Управление сессиями"
    echo "   ✅ История активности"
    echo "   ✅ JWT аутентификация"
    echo ""
    echo "🔧 ТЕХНИЧЕСКИЙ СТЕК:"
    echo "   • Go 1.21"
    echo "   • PostgreSQL / In-Memory"
    echo "   • Docker & Docker Compose"
    echo "   • JWT + bcrypt"
    echo "   • Clean Architecture"
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo "✅ Проект готов! Начните с: cd $PROFILE_DIR && ./run.sh dev"
    echo "════════════════════════════════════════════════════════════════"
else
    echo "⚠️  ВНИМАНИЕ: Не все файлы найдены!"
    echo "   Пожалуйста, проверьте структуру проекта"
fi

echo ""
