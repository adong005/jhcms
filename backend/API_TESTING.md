# ADCMS Backend API 测试文档

## 服务信息

- **服务地址**: http://localhost:8081
- **数据库**: jhcms @ 192.168.11.3
- **默认管理员**: admin / admin123

## ✅ 已测试通过的接口

### 1. 健康检查

```bash
curl http://localhost:8081/health
```

**响应**:
```json
{
  "status": "ok",
  "message": "ADCMS Backend is running"
}
```

---

### 2. 用户登录

**接口**: `POST /api/auth/login`

```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGc...",
    "refreshToken": "eyJhbGc...",
    "user": {
      "userId": "2",
      "username": "admin",
      "realName": "超级管理员",
      "role": "super_admin",
      "roles": ["super_admin"]
    }
  }
}
```

---

### 3. 获取用户信息

**接口**: `GET /api/user/info`

**需要认证**: ✅

```bash
TOKEN="YOUR_ACCESS_TOKEN"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/user/info
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userId": "2",
    "username": "admin",
    "realName": "超级管理员",
    "role": "super_admin",
    "roles": ["super_admin"]
  }
}
```

---

### 4. 获取权限码

**接口**: `GET /api/auth/codes`

**需要认证**: ✅

```bash
TOKEN="YOUR_ACCESS_TOKEN"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/auth/codes
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    "AC_100100",
    "AC_100110",
    "AC_100120",
    "AC_100010",
    "AC_100020",
    "AC_100030",
    "AC_1000031",
    "AC_1000032"
  ]
}
```

---

### 5. 获取网站配置

**接口**: `GET /api/site-config`

**需要认证**: ✅

```bash
TOKEN="YOUR_ACCESS_TOKEN"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "",
    "keywords": "",
    "description": "",
    "domain": ""
  }
}
```

---

### 6. 更新网站配置

**接口**: `POST /api/site-config`

**需要认证**: ✅

```bash
TOKEN="YOUR_ACCESS_TOKEN"
curl -X POST http://localhost:8081/api/site-config \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "2",
    "title": "ADCMS 管理系统",
    "keywords": "内容管理,CMS,多租户",
    "description": "基于 Go Gin 的 B2B2C 多租户内容管理系统",
    "domain": "www.adcms.com"
  }'
```

**响应**:
```json
{
  "code": 0,
  "message": "网站配置更新成功"
}
```

---

### 7. 用户登出

**接口**: `POST /api/auth/logout`

**需要认证**: ✅

```bash
TOKEN="YOUR_ACCESS_TOKEN"
curl -X POST http://localhost:8081/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

**响应**:
```json
{
  "code": 0,
  "message": "登出成功"
}
```

---

### 8. 刷新令牌

**接口**: `POST /api/auth/refresh`

```bash
curl -X POST http://localhost:8081/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGc..."
  }
}
```

---

## 前端对接配置

### 修改前端 API 地址

编辑 `/backup/projects/adcms/admin/apps/web-antd/.env.development`:

```bash
# API 地址
VITE_GLOB_API_URL=http://localhost:8081/api
```

### 响应格式说明

后端统一响应格式：

```typescript
interface Response<T = any> {
  code: number;      // 0=成功, 其他=失败
  message: string;   // 提示信息
  data?: T;          // 响应数据
}
```

### 分页响应格式

```typescript
interface PageResponse<T = any> {
  code: 0;
  message: "success";
  data: {
    items: T[];      // 数据列表
    total: number;   // 总数
  }
}
```

---

## 权限码说明

### 超级管理员 (super_admin)
```
AC_100100, AC_100110, AC_100120, AC_100010,
AC_100020, AC_100030, AC_1000031, AC_1000032
```

### 管理员 (admin)
```
AC_100030, AC_1000031, AC_1000032
```

### 普通用户 (user)
```
AC_1000032
```

---

## 错误码说明

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 1 | 业务错误 |
| 401 | 未授权（未登录或令牌过期） |
| 403 | 禁止访问（无权限） |

---

## 测试脚本

### 完整登录流程测试

```bash
#!/bin/bash

# 1. 登录获取 Token
echo "=== 1. 登录 ==="
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo $LOGIN_RESPONSE | python3 -m json.tool

# 提取 Token
TOKEN=$(echo $LOGIN_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['accessToken'])")

echo -e "\n=== 2. 获取用户信息 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/user/info | python3 -m json.tool

echo -e "\n=== 3. 获取权限码 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/auth/codes | python3 -m json.tool

echo -e "\n=== 4. 获取网站配置 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config | python3 -m json.tool

echo -e "\n=== 5. 更新网站配置 ==="
curl -s -X POST http://localhost:8081/api/site-config \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "2",
    "title": "ADCMS 管理系统",
    "keywords": "内容管理,CMS,多租户",
    "description": "基于 Go Gin 的 B2B2C 多租户内容管理系统",
    "domain": "www.adcms.com"
  }' | python3 -m json.tool

echo -e "\n=== 6. 验证更新后的配置 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config | python3 -m json.tool

echo -e "\n=== 测试完成 ==="
```

保存为 `test-api.sh` 并运行：
```bash
chmod +x test-api.sh
./test-api.sh
```

---

## 下一步开发计划

### 待实现的接口

1. **用户管理**
   - POST /api/user/list - 用户列表
   - POST /api/user/create - 创建用户
   - PUT /api/user/:id - 更新用户
   - DELETE /api/user/:id - 删除用户

2. **站群管理**
   - POST /api/site-group/list - 站群列表
   - POST /api/site-group/create - 创建站群
   - PUT /api/site-group/:id - 更新站群
   - DELETE /api/site-group/:id - 删除站群

3. **表单管理**
   - POST /api/form-manage/list - 表单列表
   - POST /api/form-manage/export - 导出表单

4. **系统日志**
   - POST /api/system-logs/list - 日志列表

---

## 注意事项

1. **CORS 配置**: 前端地址 `http://localhost:5666` 已添加到 CORS 白名单
2. **Token 过期时间**: 
   - Access Token: 7200 秒（2小时）
   - Refresh Token: 604800 秒（7天）
3. **数据库连接**: 确保数据库服务正常运行
4. **端口占用**: 后端使用 8081 端口（8080 被占用）

---

## 生产环境部署

参考 `DEPLOYMENT.md` 文档进行生产环境部署。
