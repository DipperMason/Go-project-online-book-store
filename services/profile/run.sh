#!/bin/bash

# Скрипт для запуска Profile Service в разных окружениях

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для вывода информации
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Получение текущей директории
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Проверка Go установлена
if ! command -v go &> /dev/null; then
    log_error "Go не установлен. Пожалуйста, установите Go 1.21 или выше."
    exit 1
fi

log_info "Go версия: $(go version)"

# Проверка зависимостей
log_info "Проверка зависимостей..."
cd "$SCRIPT_DIR"

if [ ! -f "go.mod" ]; then
    log_error "go.mod не найден в директории $SCRIPT_DIR"
    exit 1
fi

log_info "Загрузка зависимостей..."
go mod download
go mod tidy

# Выбор режима запуска
if [ -z "$1" ]; then
    log_warn "Режим запуска не указан. Использую режим development."
    MODE="dev"
else
    MODE="$1"
fi

case "$MODE" in
    dev|development)
        log_info "Запуск в режиме development..."
        export PORT="${PORT:-8003}"
        export LOG_LEVEL="${LOG_LEVEL:-debug}"
        go run main.go
        ;;
    build)
        log_info "Сборка приложения..."
        go build -o profile main.go
        log_info "Приложение собрано успешно: ./profile"
        ;;
    build-release)
        log_info "Сборка release версии..."
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o profile_linux_amd64 main.go
        log_info "Release версия собрана: ./profile_linux_amd64"
        ;;
    test)
        log_info "Запуск тестов..."
        go test -v -cover ./...
        ;;
    lint)
        log_info "Проверка кода с помощью golangci-lint..."
        if ! command -v golangci-lint &> /dev/null; then
            log_warn "golangci-lint не установлен. Пропускаю проверку."
        else
            golangci-lint run ./...
        fi
        ;;
    docker)
        log_info "Сборка Docker образа..."
        docker build -t profile-service:latest -f Dockerfile "$PROJECT_ROOT"
        log_info "Docker образ собран: profile-service:latest"
        ;;
    docker-run)
        log_info "Запуск Docker контейнера..."
        docker run -p 8003:8003 -e PORT=8003 -e LOG_LEVEL=info profile-service:latest
        ;;
    docker-compose)
        log_info "Запуск Docker Compose..."
        docker-compose up
        ;;
    *)
        log_error "Неизвестный режим: $MODE"
        echo ""
        echo "Доступные режимы:"
        echo "  dev              - Запуск в режиме development"
        echo "  build            - Сборка приложения"
        echo "  build-release    - Сборка release версии для Linux"
        echo "  test             - Запуск тестов"
        echo "  lint             - Проверка кода"
        echo "  docker           - Сборка Docker образа"
        echo "  docker-run       - Запуск Docker контейнера"
        echo "  docker-compose   - Запуск Docker Compose"
        exit 1
        ;;
esac

log_info "Готово!"
