# ADCMS Backend - Go Gin 多租户后端

基于 Go Gin 框架的 B2B2C 多租户内容管理系统后端。

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **认证**: JWT
- **日志**: Zap
- **配置**: godotenv

## 项目结构

```
backend/
├── cmd/
│   ├── server/          # 主程序入口
│   └── migrate/         # 数据库迁移工具
├── internal/
│   ├── config/          # 配置管理
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器
│   ├── router/          # 路由配置
│   └── pkg/             # 工具包
├── .env                 # 环境变量配置
├── Makefile             # 构建脚本
└── go.mod               # Go 模块依赖
```

## 快速开始

### 1. 配置环境变量

复制 `.env.example` 并修改为 `.env`：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接信息。

### 2. 运行数据库迁移

```bash
make migrate
```

这将创建所有必要的数据表，并创建默认超级管理员账号：
- 用户名: `admin`
- 密码: `admin123`

### 3. 启动开发服务器

```bash
make run
```

服务器将在 `http://localhost:8080` 启动。

### 4. 健康检查

```bash
curl http://localhost:8080/health
```

## 编译部署

### 编译 Linux 可执行程序

```bash
make build-linux
```

这将生成 `adcms-server` 可执行文件。

### 部署到服务器

1. 上传文件到服务器：
```bash
scp adcms-server root@192.168.11.3:/www/wwwroot/adcms/
scp .env root@192.168.11.3:/www/wwwroot/adcms/
```

2. 设置权限：
```bash
chmod +x /www/wwwroot/adcms/adcms-server
chmod 600 /www/wwwroot/adcms/.env
```

3. 运行数据库迁移：
```bash
cd /www/wwwroot/adcms
./adcms-server migrate
```

4. 启动服务（使用 systemd 或宝塔面板进程守护）

详细部署步骤请参考项目文档。

## B2B2C 多租户架构

### 角色层级

```
超级管理员 (super_admin)
    ├── 管理员A (admin) - 租户A
    │   ├── 用户A1 (user)
    │   └── 用户A2 (user)
    └── 管理员B (admin) - 租户B
        └── 用户B1 (user)
```

### 数据隔离规则

- **超级管理员**: 可以访问所有租户数据
- **管理员**: 只能访问自己租户的数据（tenant_id = 管理员的 user_id）
- **用户**: 只能访问自己创建的数据（parent_id = 租户ID, created_by = 用户ID）

### 表单分表设计

每个租户拥有独立的表单表：`forms_{tenant_id}`

创建管理员时自动创建对应的表单表。

## API 接口

### 认证接口

- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `POST /api/auth/refresh` - 刷新 Token

### 用户管理

- `POST /api/user/list` - 用户列表
- `POST /api/user/create` - 创建用户
- `GET /api/user/info` - 获取当前用户信息

### 网站配置

- `GET /api/site-config` - 获取网站配置
- `POST /api/site-config` - 更新网站配置

更多接口请参考 API 文档。

## 开发命令

```bash
# 安装依赖
make deps

# 运行服务
make run

# 运行迁移
make migrate

# 编译（本地）
make build-local

# 编译（Linux）
make build-linux

# 清理构建文件
make clean

# 运行测试
make test
```

## 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务器端口 | 8080 |
| SERVER_MODE | 运行模式 (debug/release) | debug |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_DATABASE | 数据库名称 | adcms |
| JWT_SECRET | JWT 密钥 | - |
| JWT_ACCESS_TOKEN_EXPIRE | Access Token 过期时间（秒） | 7200 |
| JWT_REFRESH_TOKEN_EXPIRE | Refresh Token 过期时间（秒） | 604800 |
| CORS_ALLOW_ORIGINS | 允许的跨域来源 | http://localhost:5666 |

## 注意事项

1. 生产环境务必修改 `JWT_SECRET` 为强密码
2. 数据库用户需要有创建表的权限（用于表单分表）
3. 建议使用 Nginx 反向代理并配置 SSL 证书
4. 定期备份数据库，特别是动态创建的表单表

## License

MIT
