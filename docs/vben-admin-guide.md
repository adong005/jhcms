# Vben Admin 5.7.0 快速开始指南

## 项目简介

Vue Vben Admin 是一个基于 Vue3、Vite、TypeScript 的免费开源中后台管理系统模板。

- **版本**: 5.7.0
- **官方文档**: https://doc.vben.pro
- **GitHub**: https://github.com/vbenjs/vue-vben-admin
- **在线预览**: https://vben.pro

## 核心特点

### 技术栈
- **前端框架**: Vue 3.5.30
- **构建工具**: Vite 8.0.1
- **语言**: TypeScript 5.9.3
- **包管理**: pnpm 10.32.1
- **架构**: Pnpm Monorepo + Turborepo

### 主要特性
- ✅ 最新技术栈：Vue3、Vite、TypeScript
- ✅ 国际化：内置完善的国际化方案
- ✅ 权限验证：按钮级别权限控制
- ✅ 多主题：支持多种主题配置和黑暗模式
- ✅ 动态菜单：根据权限配置显示菜单
- ✅ Mock 数据：基于 Nitro 的本地高性能 Mock 方案
- ✅ 代码规范：Oxfmt、Oxlint、ESLint、Stylelint
- ✅ 多 UI 库支持：Ant Design Vue、Element Plus、Naive UI、TDesign

## 环境要求

### 必需环境
- **Node.js**: 20.15.0+ (推荐 20.19.0 / 22.18.0 / 24.0.0)
- **Git**: 任意版本
- **pnpm**: 10.0.0+ (通过 corepack 自动安装)

### 版本管理工具（可选）
- fnm
- nvm
- pnpm env

### 验证环境
```bash
# 查看 Node 版本
node -v

# 查看 Git 版本
git -v
```

## 快速开始

### 1. 获取源码

```bash
# GitHub（推荐）
git clone https://github.com/vbenjs/vue-vben-admin.git

# Gitee（可能不是最新）
git clone https://gitee.com/annsion/vue-vben-admin.git
```

⚠️ **注意**: 代码目录及所有父级目录不能包含中文、韩文、日文及空格

### 2. 安装依赖

```bash
# 进入项目目录
cd vue-vben-admin

# 启用 corepack（自动安装指定版本 pnpm）
npm i -g corepack

# 安装依赖
pnpm install
```

#### 国内网络优化
如果无法访问 npm 源，设置环境变量：
```bash
export COREPACK_NPM_REGISTRY=https://registry.npmmirror.com
pnpm install
```

### 3. 运行项目

#### 方式一：交互式选择
```bash
pnpm dev
```
会出现选择菜单：
```
│ ◆ Select the app you need to run [dev]:
│ ● @vben/web-antd          # Ant Design Vue
│ ○ @vben/web-antdv-next    # Ant Design Vue Next
│ ○ @vben/web-ele           # Element Plus
│ ○ @vben/web-naive         # Naive UI
│ ○ @vben/web-tdesign       # TDesign
│ ○ @vben/docs              # 文档
│ ○ @vben/playground        # 演练场
```

#### 方式二：直接运行指定项目
```bash
pnpm run dev:antd      # Ant Design Vue 版本
pnpm run dev:antdv-next # Ant Design Vue Next 版本
pnpm run dev:ele       # Element Plus 版本
pnpm run dev:naive     # Naive UI 版本
pnpm run dev:tdesign   # TDesign 版本
pnpm run dev:docs      # 文档
pnpm run dev:play      # 演练场
```

### 4. 访问项目

默认地址：http://localhost:5555 （端口可能不同）

Mock API 地址：http://localhost:5320/api

## 项目结构（Monorepo）

### 基础概念

#### 大仓（Monorepo）
整个项目仓库，包含所有代码、包、应用、规范、文档、配置

#### 应用（Apps）
位于 `apps/` 目录，每个应用独立运行、构建、测试、部署
- `apps/web-antd` - Ant Design Vue 版本
- `apps/web-antdv-next` - Ant Design Vue Next 版本
- `apps/web-ele` - Element Plus 版本
- `apps/web-naive` - Naive UI 版本
- `apps/web-tdesign` - TDesign 版本
- `apps/backend-mock` - Mock 后端服务（Nitro）
- `apps/docs` - 文档站点

#### 包（Packages）
位于 `packages/` 目录，可被多个应用引用的独立模块
- 组件库
- 工具库
- 类型定义
- 配置文件等

### 目录结构
```
vue-vben-admin/
├── apps/                 # 应用目录
│   ├── web-antd/        # Ant Design Vue 应用
│   ├── backend-mock/    # Mock 服务
│   └── ...
├── packages/            # 共享包目录
│   ├── @vben/          # 核心包
│   └── ...
├── internal/            # 内部工具和配置
├── docs/                # 文档
├── scripts/             # 脚本
└── pnpm-workspace.yaml  # pnpm 工作空间配置
```

## 浏览器支持

### 开发环境
- Chrome 最新版（不支持 Chrome 80 以下）

### 生产环境
- 支持现代浏览器
- ❌ 不支持 IE

## 常用命令

```bash
# 开发
pnpm dev              # 交互式选择应用
pnpm dev:antd         # 运行 Ant Design Vue 版本

# 构建
pnpm build            # 构建所有应用
pnpm build:antd       # 构建 Ant Design Vue 版本

# 代码检查
pnpm lint             # 代码检查
pnpm format           # 代码格式化
pnpm check            # 完整检查（循环依赖、类型、拼写等）

# 测试
pnpm test:unit        # 单元测试
pnpm test:e2e         # E2E 测试

# 其他
pnpm clean            # 清理
pnpm reinstall        # 重新安装依赖
```

## Mock 数据

### Mock 服务
- 基于 **Nitro** 框架
- 位置：`apps/backend-mock/`
- 地址：http://localhost:5320/api
- 特点：高性能、支持 JWT 认证、权限控制

### Mock 数据文件
- 核心数据：`apps/backend-mock/utils/mock-data.ts`
- API 路由：`apps/backend-mock/api/`

### 默认测试账号
```
用户名: vben      密码: 123456  角色: super
用户名: admin     密码: 123456  角色: admin
用户名: jack      密码: 123456  角色: user
```

## 配置文件

### 应用配置
- `apps/web-antd/src/preferences.ts` - 应用偏好设置
- `apps/web-antd/.env` - 环境变量

### 重要提示
⚠️ 更改配置后请清空缓存，否则可能不生效

## 注意事项

1. **只支持 pnpm** 进行依赖安装
2. **路径不能包含中文、韩文、日文、空格**
3. **Node.js 版本必须 20.15.0+**
4. **更改配置后需清空缓存**
5. **使用 corepack 管理 pnpm 版本**

## 贡献方式

1. 长期提交 PR
2. 提供有价值的建议
3. 参与讨论，帮助解决 issue
4. 共同维护文档

## 相关链接

- 官方文档：https://doc.vben.pro
- GitHub：https://github.com/vbenjs/vue-vben-admin
- 在线预览：https://vben.pro
- 讨论区：https://github.com/vbenjs/vue-vben-admin/discussions

## 许可证

MIT © 2020-2026 Vben
