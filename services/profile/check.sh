#!/usr/bin/env bash

# Скрипт для проверки что всё работает

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔍 Проверка структуры проекта Profile Service..."
echo ""

# Проверка основных файлов
files=(
    "main.go"
    "go.mod"
    "README.md"
    "EXAMPLES.md"
    "STRUCTURE.md"
    "Dockerfile"
    "docker-compose.yml"
    "postman_collection.json"
    "run.sh"
)

missing_files=0

for file in "${files[@]}"; do
    if [ -f "$SCRIPT_DIR/$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file - НЕ НАЙДЕН"
        missing_files=$((missing_files + 1))
    fi
done

echo ""
echo "🗂️  Проверка директорий..."

dirs=(
    "internal/app"
    "internal/config"
    "internal/domain"
    "internal/repo"
    "internal/services"
    "internal/transport/http"
    "migrations"
)

missing_dirs=0

for dir in "${dirs[@]}"; do
    if [ -d "$SCRIPT_DIR/$dir" ]; then
        echo "✅ $dir"
    else
        echo "❌ $dir - НЕ НАЙДЕНА"
        missing_dirs=$((missing_dirs + 1))
    fi
done

echo ""
echo "📝 Проверка исходных файлов..."

src_files=(
    "internal/app/app.go"
    "internal/config/config.go"
    "internal/domain/user_profile.go"
    "internal/domain/session.go"
    "internal/domain/activity_log.go"
    "internal/domain/errors.go"
    "internal/repo/profile_repo.go"
    "internal/repo/profile_postgres_repo.go"
    "internal/repo/session_repo.go"
    "internal/repo/session_postgres_repo.go"
    "internal/repo/activity_repo.go"
    "internal/repo/activity_postgres_repo.go"
    "internal/services/profile_service.go"
    "internal/services/profile_service_test.go"
    "internal/transport/http/handler.go"
    "internal/transport/http/responses.go"
    "internal/transport/http/auth_middleware.go"
)

missing_src=0

for file in "${src_files[@]}"; do
    if [ -f "$SCRIPT_DIR/$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file - НЕ НАЙДЕН"
        missing_src=$((missing_src + 1))
    fi
done

echo ""
echo "=== ИТОГИ ПРОВЕРКИ ==="
total_missing=$((missing_files + missing_dirs + missing_src))

if [ $total_missing -eq 0 ]; then
    echo "✅ Все файлы на месте! Проект готов к использованию."
    echo ""
    echo "💡 Для запуска выполните:"
    echo "   cd $SCRIPT_DIR"
    echo "   ./run.sh dev"
    exit 0
else
    echo "❌ Найдено проблем: $total_missing"
    echo ""
    echo "⚠️  Пожалуйста, проверьте проект."
    exit 1
fi
