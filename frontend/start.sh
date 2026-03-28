#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

if [ ! -d "node_modules" ]; then
  echo "[frontend] node_modules 不存在，开始安装依赖..."
  npm install
fi

echo "[frontend] 启动开发服务: http://localhost:5174"
npm run dev -- --host 0.0.0.0 --port 5174
