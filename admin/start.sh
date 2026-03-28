#!/usr/bin/env bash

# ADCMS Admin startup script
# - stop old vben admin dev processes
# - ensure only one web-antd dev server
# - force fixed port 5666

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$ROOT_DIR/apps/web-antd"
LOG_DIR="$ROOT_DIR/logs"
LOG_FILE="$LOG_DIR/web-antd-dev.log"
TARGET_PORT=5666

echo "======================================"
echo "  ADCMS Admin startup script"
echo "======================================"
echo

if [ ! -d "$APP_DIR" ]; then
  echo "ERROR: app directory not found: $APP_DIR"
  exit 1
fi

mkdir -p "$LOG_DIR"

echo "1) Stop old admin dev processes..."

# Collect potential old dev pids under this admin workspace.
mapfile -t CANDIDATE_PIDS < <(pgrep -f "/backup/projects/adcms/admin" || true)
STOPPED_ANY=0

for pid in "${CANDIDATE_PIDS[@]:-}"; do
  if [ -z "$pid" ]; then
    continue
  fi

  if [ "$pid" = "$$" ]; then
    continue
  fi

  cmdline="$(ps -p "$pid" -o args= 2>/dev/null || true)"
  if [ -z "$cmdline" ]; then
    continue
  fi

  if [[ "$cmdline" == *"turbo-run dev"* ]] || [[ "$cmdline" == *"pnpm run dev"* ]] || [[ "$cmdline" == *"vite --mode development"* ]]; then
    echo "   stopping pid=$pid : $cmdline"
    kill -9 "$pid" 2>/dev/null || true
    STOPPED_ANY=1
  fi
done

# Also clear target port if still occupied.
mapfile -t PORT_PIDS < <(lsof -ti:"$TARGET_PORT" 2>/dev/null || true)
for pid in "${PORT_PIDS[@]:-}"; do
  if [ -z "$pid" ]; then
    continue
  fi

  echo "   freeing port $TARGET_PORT (pid=$pid)"
  kill -9 "$pid" 2>/dev/null || true
  STOPPED_ANY=1
done

if [ "$STOPPED_ANY" -eq 0 ]; then
  echo "   no old dev process found"
else
  sleep 1
  echo "   old processes cleared"
fi

echo
echo "2) Start unique web-antd dev server on :$TARGET_PORT ..."

cd "$APP_DIR"
nohup pnpm vite --mode development --host 0.0.0.0 --port "$TARGET_PORT" --strictPort > "$LOG_FILE" 2>&1 &
NEW_PID=$!

echo "   pid: $NEW_PID"
echo "   log: $LOG_FILE"

sleep 3

if ps -p "$NEW_PID" >/dev/null 2>&1 && lsof -ti:"$TARGET_PORT" >/dev/null 2>&1; then
  echo
  echo "OK: admin started successfully"
  echo "URL: http://localhost:$TARGET_PORT/"
  echo "Tail log: tail -f $LOG_FILE"
  echo "Stop: kill $NEW_PID"
else
  echo
  echo "ERROR: admin failed to start"
  echo "Check log: $LOG_FILE"
  exit 1
fi

echo
echo "======================================"
