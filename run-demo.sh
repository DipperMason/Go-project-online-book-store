#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PIDS_FILE="$ROOT_DIR/.demo.pids"

log() {
  echo "[demo] $*"
}

cleanup_previous() {
  if [[ -f "$PIDS_FILE" ]]; then
    while read -r pid; do
      if [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null; then
        kill "$pid" 2>/dev/null || true
      fi
    done < "$PIDS_FILE"
    rm -f "$PIDS_FILE"
  fi

  # Ensure previous listeners are gone even if go run spawned child processes.
  for port in 8088 8003 8001; do
    if fuser "${port}/tcp" >/dev/null 2>&1; then
      fuser -k "${port}/tcp" >/dev/null 2>&1 || true
    fi
  done
}

start_db() {
  if ! docker ps --format '{{.Names}}' | grep -q '^litsee-profile-db$'; then
    log "starting postgres container litsee-profile-db"
    docker run -d \
      --name litsee-profile-db \
      -e POSTGRES_USER=litsee \
      -e POSTGRES_PASSWORD=litsee \
      -e POSTGRES_DB=litsee_profile \
      -p 55432:5432 \
      postgres:16-alpine >/dev/null
  else
    log "postgres container already running"
  fi

  for _ in {1..30}; do
    if docker exec litsee-profile-db pg_isready -U litsee -d litsee_profile >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done

  log "applying profile migrations"
  docker exec -i litsee-profile-db psql -U litsee -d litsee_profile < "$ROOT_DIR/services/profile/migrations/001_init.sql" >/dev/null
}

start_auth() {
  log "starting auth on :8001"
  (
    cd "$ROOT_DIR/services/auth"
    HTTP_ADDR=':8001' JWT_SECRET='hihihaha' go run . > /tmp/litsee-auth.log 2>&1
  ) &
  echo "$!" >> "$PIDS_FILE"
}

start_profile() {
  log "starting profile on :8003"
  (
    cd "$ROOT_DIR/services/profile"
    PORT='8003' \
    DATABASE_URL='postgres://litsee:litsee@localhost:55432/litsee_profile?sslmode=disable' \
    JWT_SECRET='hihihaha' \
    REDPANDA_ORDER_PAID_TOPIC='' \
    GOWORK=off go run . > /tmp/litsee-profile.log 2>&1
  ) &
  echo "$!" >> "$PIDS_FILE"
}

start_web() {
  log "starting web on :8088"
  (
    cd "$ROOT_DIR/web"
    WEB_ADDR=':8088' \
    AUTH_URL='http://localhost:8001' \
    PROFILE_URL='http://localhost:8003' \
    GOWORK=off go run . > /tmp/litsee-web.log 2>&1
  ) &
  echo "$!" >> "$PIDS_FILE"
}

wait_ports() {
  for _ in {1..40}; do
    if ss -ltn | grep -q ':8001 ' && ss -ltn | grep -q ':8003 ' && ss -ltn | grep -q ':8088 '; then
      return 0
    fi
    sleep 1
  done
  return 1
}

print_status() {
  echo
  log "services status"
  ss -ltnp | grep -E ':8001 |:8003 |:8088 |:55432 ' || true

  echo
  log "quick checks"
  echo -n "web: "
  curl --noproxy '*' -s -o /dev/null -w '%{http_code}' http://127.0.0.1:8088/
  echo

  echo -n "auth (probe): "
  curl --noproxy '*' -s -o /dev/null -w '%{http_code}' \
    -X POST http://127.0.0.1:8001/api/v1/auth \
    -H 'Content-Type: application/json' \
    -d '{"email":"probe@local","password":"wrong"}'
  echo

  echo -n "profile (unauth): "
  curl --noproxy '*' -s -o /dev/null -w '%{http_code}' http://127.0.0.1:8003/api/v1/profile
  echo

  echo
  log "open http://localhost:8088"
  log "logs: /tmp/litsee-auth.log /tmp/litsee-profile.log /tmp/litsee-web.log"
}

main() {
  cleanup_previous
  : > "$PIDS_FILE"
  start_db
  start_auth
  start_profile
  start_web

  if ! wait_ports; then
    log "failed to start all services"
    log "auth log:"; tail -n 60 /tmp/litsee-auth.log || true
    log "profile log:"; tail -n 60 /tmp/litsee-profile.log || true
    log "web log:"; tail -n 60 /tmp/litsee-web.log || true
    exit 1
  fi

  print_status
}

main "$@"
