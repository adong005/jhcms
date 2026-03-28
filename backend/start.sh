#!/bin/bash

# ADCMS Backend 启动脚本
# 自动停止旧进程并启动新进程

set -e

echo "======================================"
echo "  ADCMS Backend 启动脚本"
echo "======================================"
echo ""

# 检查可执行文件是否存在
if [ ! -f "./adcms-server" ]; then
    echo "❌ 错误: adcms-server 可执行文件不存在"
    echo "请先运行: ./build.sh"
    exit 1
fi

# 停止旧进程
echo "1. 检查并停止旧进程..."
OLD_PIDS=$(lsof -ti:8081 2>/dev/null || true)

if [ -n "$OLD_PIDS" ]; then
    echo "   发现占用 8081 端口的进程: $OLD_PIDS"
    echo "   正在停止..."
    echo "$OLD_PIDS" | xargs kill -9 2>/dev/null || true
    sleep 1
    echo "   ✅ 旧进程已停止"
else
    echo "   ℹ️  没有发现旧进程"
fi

# 检查进程是否还在运行
if lsof -ti:8081 >/dev/null 2>&1; then
    echo "   ⚠️  警告: 端口 8081 仍被占用，强制清理..."
    lsof -ti:8081 | xargs kill -9 2>/dev/null || true
    sleep 2
fi

echo ""
echo "2. 启动新进程..."

# 启动服务（后台运行）
nohup ./adcms-server > logs/server.log 2>&1 &
NEW_PID=$!

echo "   进程 ID: $NEW_PID"
echo "   日志文件: logs/server.log"

# 等待服务启动
echo ""
echo "3. 等待服务启动..."
sleep 3

# 检查服务是否正常运行
if ps -p $NEW_PID > /dev/null 2>&1; then
    # 测试健康检查接口
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo ""
        echo "✅ 服务启动成功！"
        echo ""
        echo "服务信息:"
        echo "  - 地址: http://localhost:8081"
        echo "  - 进程 ID: $NEW_PID"
        echo "  - 健康检查: http://localhost:8081/health"
        echo ""
        echo "查看日志: tail -f logs/server.log"
        echo "停止服务: kill $NEW_PID 或 lsof -ti:8081 | xargs kill -9"
    else
        echo ""
        echo "⚠️  警告: 服务已启动但健康检查失败"
        echo "请检查日志: tail -f logs/server.log"
    fi
else
    echo ""
    echo "❌ 错误: 服务启动失败"
    echo "请检查日志: tail -f logs/server.log"
    exit 1
fi

echo ""
echo "======================================"
