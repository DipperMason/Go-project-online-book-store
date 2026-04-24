#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PIDS_FILE="$ROOT_DIR/.demo.pids"

echo "[demo] stopping app processes"
if [[ -f "$PIDS_FILE" ]]; then
  while read -r pid; do
    if [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done < "$PIDS_FILE"
  rm -f "$PIDS_FILE"
fi

for port in 8088 8003 8001; do
  if fuser "${port}/tcp" >/dev/null 2>&1; then
    echo "[demo] killing listeners on :${port}"
    fuser -k "${port}/tcp" >/dev/null 2>&1 || true
  fi
done

# Fallback for orphaned go-run children.
pkill -f '/home/maks/Litsee/web' >/dev/null 2>&1 || true
pkill -f '/home/maks/Litsee/services/profile' >/dev/null 2>&1 || true
pkill -f '/home/maks/Litsee/services/auth' >/dev/null 2>&1 || true

echo "[demo] stopping postgres container"
docker rm -f litsee-profile-db >/dev/null 2>&1 || true

echo "[demo] done"
