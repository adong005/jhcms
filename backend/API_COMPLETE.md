# ADCMS Backend API 完整接口文档

## 服务信息

- **服务地址**: http://localhost:8081
- **数据库**: jhcms @ 192.168.11.3
- **默认管理员**: admin / admin123

---

## 一、认证接口

### 1.1 用户登录
**接口**: `POST /api/auth/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "admin123"
}
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

### 1.2 用户登出
**接口**: `POST /api/auth/logout`

**需要认证**: ✅

**响应**:
```json
{
  "code": 0,
  "message": "登出成功"
}
```

### 1.3 刷新令牌
**接口**: `POST /api/auth/refresh`

**请求参数**:
```json
{
  "refreshToken": "YOUR_REFRESH_TOKEN"
}
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

### 1.4 获取用户信息
**接口**: `GET /api/user/info`

**需要认证**: ✅

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

### 1.5 获取权限码
**接口**: `GET /api/auth/codes`

**需要认证**: ✅

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

## 二、菜单接口

### 2.1 获取系统导航菜单
**接口**: `GET /api/menu/all`

**需要认证**: ✅

**用途**: 前端路由系统使用

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "name": "Dashboard",
      "path": "/",
      "meta": {
        "title": "首页",
        "icon": "lucide:layout-dashboard",
        "order": -1
      },
      "redirect": "/analytics",
      "children": [
        {
          "name": "Analytics",
          "path": "/analytics",
          "component": "/dashboard/analytics/index",
          "meta": {
            "title": "分析页",
            "icon": "lucide:area-chart",
            "affixTab": true
          }
        }
      ]
    }
  ]
}
```

### 2.2 获取菜单管理列表
**接口**: `POST /api/menu/list`

**需要认证**: ✅

**用途**: 菜单管理页面使用

**请求参数**:
```json
{
  "page": 1,
  "pageSize": 100,
  "name": "",
  "type": "",
  "status": null
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "1",
        "name": "工作台",
        "path": "/analytics",
        "type": "menu",
        "icon": "mdi:view-dashboard",
        "component": "/dashboard/analytics/index",
        "parentId": null,
        "order": 1,
        "status": 1,
        "createTime": "2024-01-01 10:00:00",
        "updateTime": "2024-03-20 14:00:00"
      }
    ],
    "total": 13
  }
}
```

---

## 三、网站配置接口

### 3.1 获取网站配置
**接口**: `GET /api/site-config`

**需要认证**: ✅

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "ADCMS 管理系统",
    "keywords": "内容管理,CMS,多租户,Go,Gin",
    "description": "基于 Go Gin 的 B2B2C 多租户内容管理系统",
    "domain": "www.adcms.com"
  }
}
```

### 3.2 更新网站配置
**接口**: `POST /api/site-config`

**需要认证**: ✅

**请求参数**:
```json
{
  "userId": "2",
  "title": "ADCMS 管理系统",
  "keywords": "内容管理,CMS,多租户,Go,Gin",
  "description": "基于 Go Gin 的 B2B2C 多租户内容管理系统",
  "domain": "www.adcms.com"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "网站配置更新成功"
}
```

---

## 四、信息管理接口

### 4.1 获取信息列表
**接口**: `POST /api/info/list`

**需要认证**: ✅

**请求参数**:
```json
{
  "page": 1,
  "pageSize": 10,
  "title": "",
  "status": null
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "测试信息标题",
        "content": "这是一条测试信息的内容",
        "status": 1,
        "createdBy": 2,
        "createdAt": "2026-03-25T10:40:57.659Z",
        "updatedAt": "2026-03-25T10:40:57.659Z"
      }
    ],
    "total": 1
  }
}
```

**租户隔离规则**:
- **超级管理员**: 可以看到所有信息
- **管理员**: 只能看到自己租户的信息
- **普通用户**: 只能看到自己创建的信息

### 4.2 获取信息详情
**接口**: `GET /api/info/:id`

**需要认证**: ✅

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "测试信息标题",
    "content": "这是一条测试信息的内容",
    "status": 1,
    "createdBy": 2,
    "createdAt": "2026-03-25T10:40:57.659Z",
    "updatedAt": "2026-03-25T10:40:57.659Z"
  }
}
```

### 4.3 创建信息
**接口**: `POST /api/info/create`

**需要认证**: ✅

**请求参数**:
```json
{
  "title": "测试信息标题",
  "content": "这是一条测试信息的内容",
  "categoryId": null,
  "status": 1,
  "publishDate": ""
}
```

**响应**:
```json
{
  "code": 0,
  "message": "创建信息成功"
}
```

### 4.4 更新信息
**接口**: `PUT /api/info/:id`

**需要认证**: ✅

**请求参数**:
```json
{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "categoryId": null,
  "status": 1,
  "publishDate": ""
}
```

**响应**:
```json
{
  "code": 0,
  "message": "更新信息成功"
}
```

### 4.5 删除信息
**接口**: `DELETE /api/info/:id`

**需要认证**: ✅

**响应**:
```json
{
  "code": 0,
  "message": "删除信息成功"
}
```

---

## 五、健康检查

### 5.1 健康检查
**接口**: `GET /health`

**响应**:
```json
{
  "status": "ok",
  "message": "ADCMS Backend is running"
}
```

---

## 六、错误码说明

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 1 | 业务错误 |
| 401 | 未授权（未登录或令牌过期） |
| 403 | 禁止访问（无权限） |

---

## 七、认证说明

### Token 使用

所有需要认证的接口都需要在请求头中携带 Token：

```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Token 过期时间

- **Access Token**: 7200 秒（2小时）
- **Refresh Token**: 604800 秒（7天）

---

## 八、租户隔离规则

### 角色说明

1. **super_admin（超级管理员）**
   - tenant_id = 0
   - 可以访问所有数据
   - 拥有所有权限码

2. **admin（管理员）**
   - tenant_id = user_id
   - 只能访问自己租户的数据
   - 可以管理自己租户下的用户

3. **user（普通用户）**
   - tenant_id = parent_id
   - 只能访问自己创建的数据
   - 权限受限

---

## 九、测试脚本

### 完整功能测试

```bash
#!/bin/bash

# 1. 登录
TOKEN=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['data']['accessToken'])")

echo "Token: ${TOKEN:0:50}..."

# 2. 获取用户信息
echo -e "\n=== 用户信息 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/user/info | python3 -m json.tool

# 3. 获取权限码
echo -e "\n=== 权限码 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/auth/codes | python3 -m json.tool

# 4. 获取系统菜单
echo -e "\n=== 系统菜单 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/menu/all | python3 -m json.tool | head -50

# 5. 获取网站配置
echo -e "\n=== 网站配置 ==="
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/site-config | python3 -m json.tool

# 6. 创建信息
echo -e "\n=== 创建信息 ==="
curl -s -X POST http://localhost:8081/api/info/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试信息","content":"测试内容","status":1}' | \
  python3 -m json.tool

# 7. 获取信息列表
echo -e "\n=== 信息列表 ==="
curl -s -X POST http://localhost:8081/api/info/list \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"page":1,"pageSize":10}' | python3 -m json.tool

echo -e "\n=== 测试完成 ==="
```

---

## 十、前端对接配置

### 环境变量配置

编辑 `/backup/projects/adcms/admin/apps/web-antd/.env.development`:

```bash
# API 地址
VITE_GLOB_API_URL=http://localhost:8081/api

# 关闭 Mock 服务
VITE_NITRO_MOCK=false
```

### 权限模式配置

编辑 `/backup/projects/adcms/admin/apps/web-antd/src/preferences.ts`:

```typescript
export const overridesPreferences = defineOverridesPreferences({
  app: {
    accessMode: 'backend', // 后端控制模式
  },
});
```

---

## 十一、已实现的功能模块

✅ **认证模块**
- 登录、登出、刷新令牌
- 获取用户信息和权限码

✅ **菜单模块**
- 系统导航菜单（用于前端路由）
- 菜单管理列表（树形结构）

✅ **网站配置模块**
- 获取和更新网站配置
- 支持超级管理员和管理员

✅ **信息管理模块**
- 信息列表（分页、搜索、租户隔离）
- 信息详情、创建、更新、删除

---

## 十二、待实现的功能模块

⏳ **用户管理**
- 用户列表、创建、更新、删除

⏳ **角色管理**
- 角色列表、创建、更新、删除

⏳ **站群管理**
- 站群列表、创建、更新、删除

⏳ **表单管理**
- 表单列表、导出

⏳ **系统日志**
- 日志列表、查询

---

## 十三、注意事项

1. **CORS 配置**: 前端地址 `http://localhost:5666` 已添加到 CORS 白名单
2. **数据库连接**: 确保数据库服务正常运行
3. **端口占用**: 后端使用 8081 端口
4. **字符编码**: 数据库使用 UTF-8 编码
5. **软删除**: 删除操作使用软删除，数据不会真正删除

---

## 十四、生产环境部署

参考 `DEPLOYMENT.md` 文档进行生产环境部署。

主要步骤：
1. 编译：`make build-linux`
2. 配置环境变量
3. 启动服务：`./adcms-server`
4. 配置 Nginx 反向代理
5. 配置 Systemd 服务
