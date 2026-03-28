# ADCMS Backend 部署指南

## 当前完成状态

✅ **已完成的功能**：
1. 项目结构初始化
2. 配置管理（.env 文件）
3. JWT 认证系统
4. 中间件（认证、租户隔离、CORS、日志）
5. 数据模型（用户、角色、菜单、信息、站群、表单、系统日志）
6. 认证接口（登录、登出、刷新令牌、获取用户信息、获取权限码）
7. 网站配置接口（获取、更新）
8. 数据库迁移工具

## 部署步骤

### 1. 配置数据库连接

编辑 `.env` 文件，确保数据库连接信息正确：

```bash
# 数据库配置
DB_HOST=192.168.11.3
DB_PORT=3306
DB_USERNAME=jhcms
DB_PASSWORD=4XFiGi8simGJAzK4
DB_DATABASE=adcms
```

**注意**：如果数据库连接失败，可能需要：
1. 检查数据库用户权限
2. 确保数据库允许远程连接
3. 检查防火墙设置

### 2. 运行数据库迁移

```bash
# 方式1：使用 Makefile
make migrate

# 方式2：直接运行
go run cmd/migrate/main.go
```

这将：
- 创建所有数据表
- 创建默认超级管理员账号（username: admin, password: admin123）

### 3. 启动服务

```bash
# 开发环境
make run

# 或直接运行
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

### 4. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 预期返回
{
  "status": "ok",
  "message": "ADCMS Backend is running"
}
```

## API 接口测试

### 1. 登录接口

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**成功响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGc...",
    "refreshToken": "eyJhbGc...",
    "user": {
      "userId": "1",
      "username": "admin",
      "realName": "超级管理员",
      "role": "super_admin",
      "roles": ["super_admin"]
    }
  }
}
```

### 2. 获取用户信息

```bash
curl -X GET http://localhost:8080/api/user/info \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 3. 获取权限码

```bash
curl -X GET http://localhost:8080/api/auth/codes \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. 获取网站配置

```bash
curl -X GET http://localhost:8080/api/site-config \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. 更新网站配置

```bash
curl -X POST http://localhost:8080/api/site-config \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "1",
    "title": "ADCMS 管理系统",
    "keywords": "内容管理,CMS,多租户",
    "description": "基于 Go Gin 的多租户内容管理系统",
    "domain": "www.adcms.com"
  }'
```

## 生产环境部署

### 1. 编译 Linux 可执行程序

```bash
make build-linux
```

这将生成 `adcms-server` 可执行文件。

### 2. 上传到服务器

```bash
# 上传程序
scp adcms-server root@192.168.11.3:/www/wwwroot/adcms/

# 上传配置文件
scp .env root@192.168.11.3:/www/wwwroot/adcms/
```

### 3. 服务器上设置权限

```bash
ssh root@192.168.11.3

cd /www/wwwroot/adcms

# 设置可执行权限
chmod +x adcms-server

# 设置 .env 文件权限
chmod 600 .env
```

### 4. 运行数据库迁移

```bash
# 首次部署需要运行迁移
cd /www/wwwroot/adcms
./adcms-server migrate
```

**注意**：如果出现 "command not found" 错误，说明迁移功能还未实现。
可以手动执行 SQL 创建表，或者编译迁移工具：

```bash
# 在本地编译迁移工具
GOOS=linux GOARCH=amd64 go build -o adcms-migrate cmd/migrate/main.go

# 上传并运行
scp adcms-migrate root@192.168.11.3:/www/wwwroot/adcms/
ssh root@192.168.11.3 "cd /www/wwwroot/adcms && ./adcms-migrate"
```

### 5. 配置 Systemd 服务

创建服务文件：

```bash
sudo nano /etc/systemd/system/adcms.service
```

内容：

```ini
[Unit]
Description=ADCMS Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www
WorkingDirectory=/www/wwwroot/adcms
ExecStart=/www/wwwroot/adcms/adcms-server
Restart=always
RestartSec=5
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start adcms

# 设置开机自启
sudo systemctl enable adcms

# 查看服务状态
sudo systemctl status adcms

# 查看日志
sudo journalctl -u adcms -f
```

### 6. 配置 Nginx 反向代理

在宝塔面板中：
1. 创建网站
2. 网站设置 → 反向代理
3. 添加代理：
   - 代理名称: adcms-api
   - 目标URL: http://127.0.0.1:8080
   - 发送域名: $host

或手动配置 Nginx：

```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

## 常见问题

### 1. 数据库连接失败

**错误**：`Error 1045 (28000): Access denied for user`

**解决方案**：
1. 检查数据库用户权限
2. 确保数据库允许远程连接
3. 检查 .env 文件中的密码是否正确

### 2. 端口被占用

**错误**：`bind: address already in use`

**解决方案**：
1. 修改 .env 中的 SERVER_PORT
2. 或停止占用端口的进程

### 3. CORS 跨域问题

**解决方案**：
在 .env 中配置允许的来源：

```bash
CORS_ALLOW_ORIGINS=http://localhost:5666,https://yourdomain.com
```

## 下一步开发计划

🔲 **待实现功能**：
1. 用户管理（创建、编辑、删除、列表）
2. 角色管理
3. 菜单管理
4. 信息管理
5. 站群管理
6. 表单管理（分表实现）
7. 系统日志

## 技术支持

如有问题，请检查：
1. 服务日志：`journalctl -u adcms -f`
2. Nginx 日志：`/www/wwwlogs/`
3. 数据库连接状态

## 默认账号

- **用户名**: admin
- **密码**: admin123
- **角色**: super_admin

**重要**：生产环境请立即修改默认密码！
