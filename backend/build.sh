#!/bin/bash

# ADCMS Backend 编译脚本
# 用于编译 Linux 平台的可执行文件

set -e

echo "======================================"
echo "  ADCMS Backend 编译脚本"
echo "======================================"
echo ""

# 设置编译参数
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 获取版本信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "版本信息:"
echo "  Version: $VERSION"
echo "  Build Time: $BUILD_TIME"
echo "  Git Commit: $GIT_COMMIT"
echo ""

# 编译
echo "开始编译..."
go build -v \
  -ldflags "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
  -o adcms-server \
  cmd/server/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 编译成功！"
    echo ""
    echo "可执行文件: ./adcms-server"
    ls -lh adcms-server
    echo ""
    echo "运行命令: ./start.sh"
else
    echo ""
    echo "❌ 编译失败！"
    exit 1
fi
