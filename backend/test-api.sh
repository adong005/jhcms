#!/bin/bash

# ADCMS Backend API 完整测试脚本

echo "=========================================="
echo "ADCMS Backend API 测试"
echo "=========================================="

# 1. 登录获取 Token
echo -e "\n=== 1. 登录测试 ==="
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo $LOGIN_RESPONSE | python3 -m json.tool

# 提取 Token
TOKEN=$(echo $LOGIN_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['accessToken'])")

if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败，无法获取 Token"
  exit 1
fi

echo "✅ 登录成功，Token: ${TOKEN:0:50}..."

# 2. 获取用户信息
echo -e "\n=== 2. 获取用户信息 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/user/info | python3 -m json.tool

# 3. 获取权限码
echo -e "\n=== 3. 获取权限码 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/auth/codes | python3 -m json.tool

# 4. 获取网站配置
echo -e "\n=== 4. 获取网站配置（更新前） ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config | python3 -m json.tool

# 5. 更新网站配置
echo -e "\n=== 5. 更新网站配置 ==="
curl -s -X POST http://localhost:8081/api/site-config \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "2",
    "title": "ADCMS 管理系统",
    "keywords": "内容管理,CMS,多租户,Go,Gin",
    "description": "基于 Go Gin 的 B2B2C 多租户内容管理系统，支持租户隔离和权限管理",
    "domain": "www.adcms.com"
  }' | python3 -m json.tool

# 6. 验证更新后的配置
echo -e "\n=== 6. 验证更新后的配置 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config | python3 -m json.tool

echo -e "\n=========================================="
echo "✅ 所有 API 测试完成！"
echo "=========================================="
