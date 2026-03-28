# Vben Admin 完整开发指南

> 基于 Vben Admin 5.7.0 官方文档整理

## 目录

- [一、基础概念](#一基础概念)
- [二、本地开发](#二本地开发)
- [三、路由和菜单](#三路由和菜单)
- [四、配置系统](#四配置系统)
- [五、图标使用](#五图标使用)
- [六、样式开发](#六样式开发)
- [七、外部模块](#七外部模块)
- [八、构建与部署](#八构建与部署)
- [九、服务端交互](#九服务端交互)
- [十、登录认证](#十登录认证)
- [十一、主题定制](#十一主题定制)
- [十二、权限控制](#十二权限控制)
- [十三、国际化](#十三国际化)
- [十四、常用功能](#十四常用功能)

---

## 一、基础概念

### 1.1 Monorepo 架构

#### 大仓（Monorepo）
整个项目仓库，包含所有代码、包、应用、规范、文档、配置等。

#### 应用（Apps）
位于 `apps/` 目录，每个应用独立运行、构建、测试、部署：
- `apps/web-antd` - Ant Design Vue 版本
- `apps/web-antdv-next` - Ant Design Vue Next 版本
- `apps/web-ele` - Element Plus 版本
- `apps/web-naive` - Naive UI 版本
- `apps/web-tdesign` - TDesign 版本
- `apps/backend-mock` - Mock 后端服务（Nitro）

⚠️ **注意**: 应用不限于前端，`apps/backend-mock` 就是一个后端服务。

#### 包（Packages）
位于 `packages/` 目录，可被多个应用引用的独立模块，类似独立的 npm 包。

### 1.2 包的使用

#### 包引入
```typescript
// 从 workspace 引入
import { useVbenForm } from '@vben/common-ui';

// 从适配器引入（推荐）
import { useVbenForm } from '#/adapter/form';
```

#### 包使用
与使用 npm 包完全一致，支持 tree-shaking。

### 1.3 路径别名

| 别名 | 说明 | 示例 |
|------|------|------|
| `#/` | 应用内 `src` 目录 | `#/api/user` |
| `@vben/` | workspace 包 | `@vben/common-ui` |
| `~/` | 当前包根目录 | `~/utils` |

---

## 二、本地开发

### 2.1 前置准备

#### 需要掌握的基础知识
- Vue 3 基础
- TypeScript 基础
- ES6+ 语法
- Vite 基本使用
- 包管理工具（pnpm）

#### 工具配置
推荐使用 **VSCode** + 以下插件：
- Vue - Official
- TypeScript Vue Plugin (Volar)
- ESLint
- Prettier
- Tailwind CSS IntelliSense

### 2.2 Npm Scripts

```bash
# 开发
pnpm dev              # 交互式选择应用
pnpm dev:antd         # 运行 Ant Design Vue 版本
pnpm dev:ele          # 运行 Element Plus 版本
pnpm dev:naive        # 运行 Naive UI 版本
pnpm dev:docs         # 运行文档

# 构建
pnpm build            # 构建所有应用
pnpm build:antd       # 构建 Ant Design Vue 版本
pnpm build:analyze    # 构建并分析

# 代码质量
pnpm lint             # 代码检查
pnpm format           # 代码格式化
pnpm check            # 完整检查（循环依赖、类型、拼写）
pnpm check:type       # 类型检查

# 测试
pnpm test:unit        # 单元测试
pnpm test:e2e         # E2E 测试

# 其他
pnpm clean            # 清理
pnpm reinstall        # 重新安装依赖
pnpm preview          # 预览构建结果
```

### 2.3 本地运行

```bash
# 1. 进入项目目录
cd vue-vben-admin

# 2. 安装依赖
pnpm install

# 3. 运行项目
pnpm dev:antd

# 4. 访问
# http://localhost:5666
```

### 2.4 区分构建环境

通过 `.env` 文件区分环境：
- `.env` - 所有环境加载
- `.env.development` - 开发环境
- `.env.production` - 生产环境
- `.env.analyze` - 分析环境

### 2.5 公共静态资源

放在 `public/` 目录下的资源会被直接复制到构建目录。

```html
<!-- 引用 public 目录资源 -->
<img src="/logo.png" />
```

### 2.6 DevTools

项目内置 Vue DevTools，开发环境自动启用。

---

## 三、路由和菜单

### 3.1 路由类型

#### 核心路由
框架内置的基础路由，如登录页、404 页等，不可删除。

#### 静态路由
在代码中直接定义的路由，适合权限简单的应用。

#### 动态路由
从后端获取的路由，适合复杂权限场景。

### 3.2 路由定义

#### 一级路由
```typescript
// src/router/routes/modules/demo.ts
import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/demo',
    name: 'Demo',
    component: () => import('#/views/demo/index.vue'),
    meta: {
      title: '演示页面',
      icon: 'mdi:home',
    },
  },
];

export default routes;
```

#### 二级路由
```typescript
const routes: RouteRecordRaw[] = [
  {
    path: '/system',
    name: 'System',
    redirect: '/system/user',
    meta: {
      title: '系统管理',
      icon: 'carbon:settings',
    },
    children: [
      {
        path: 'user',
        name: 'SystemUser',
        component: () => import('#/views/system/user/index.vue'),
        meta: {
          title: '用户管理',
        },
      },
    ],
  },
];
```

#### 多级路由
支持无限层级嵌套，但建议不超过 3 级。

### 3.3 新增页面

#### 步骤 1: 添加路由
在 `src/router/routes/modules/` 下创建路由文件。

#### 步骤 2: 添加页面组件
在 `src/views/` 下创建对应的 Vue 组件。

#### 步骤 3: 验证
访问对应路径，检查页面是否正常显示。

### 3.4 路由配置（Meta）

| 属性 | 类型 | 说明 |
|------|------|------|
| `title` | `string` | 路由标题（菜单名称） |
| `icon` | `string` | 菜单图标 |
| `activeIcon` | `string` | 激活状态图标 |
| `keepAlive` | `boolean` | 是否缓存页面 |
| `hideInMenu` | `boolean` | 是否在菜单中隐藏 |
| `hideInTab` | `boolean` | 是否在标签页中隐藏 |
| `hideInBreadcrumb` | `boolean` | 是否在面包屑中隐藏 |
| `hideChildrenInMenu` | `boolean` | 隐藏子菜单 |
| `authority` | `string[]` | 权限码 |
| `badge` | `string \| number` | 徽标内容 |
| `badgeType` | `'dot' \| 'normal'` | 徽标类型 |
| `badgeVariants` | `string` | 徽标样式 |
| `affixTab` | `boolean` | 是否固定标签页 |
| `affixTabOrder` | `number` | 固定标签页排序 |
| `iframeSrc` | `string` | iframe 地址 |
| `ignoreAccess` | `boolean` | 忽略权限 |
| `link` | `string` | 外部链接 |
| `openInNewWindow` | `boolean` | 新窗口打开 |
| `order` | `number` | 菜单排序 |
| `activePath` | `string` | 激活路径 |
| `maxNumOfOpenTab` | `number` | 最大打开标签数 |

### 3.5 路由刷新

```typescript
import { useRouter } from 'vue-router';

const router = useRouter();

// 刷新当前路由
router.replace({
  path: '/redirect' + router.currentRoute.value.fullPath,
});
```

### 3.6 标签页与路由控制

```typescript
import { useTabs } from '@vben/hooks';

const { close, closeAll, closeLeft, closeRight } = useTabs();

// 关闭当前标签页
close();

// 关闭所有标签页
closeAll();
```

---

## 四、配置系统

### 4.1 环境变量配置

#### 环境配置文件
```bash
.env                  # 所有环境
.env.development      # 开发环境
.env.production       # 生产环境
```

#### 环境变量示例
```bash
# .env.development
VITE_APP_TITLE=Vben Admin
VITE_API_URL=http://localhost:5320/api
VITE_GLOB_APP_SHORT_NAME=vben_admin
```

#### 使用环境变量
```typescript
// 在代码中使用
const apiUrl = import.meta.env.VITE_API_URL;
const title = import.meta.env.VITE_APP_TITLE;
```

### 4.2 生产环境动态配置

#### 作用
无需重新构建，即可修改生产环境配置。

#### 使用
配置文件位于 `public/_app.config.js`：

```javascript
window._VBEN_ADMIN_PRO_APP_CONF_ = {
  VITE_API_URL: 'https://api.example.com',
};
```

#### 新增配置
1. 在 `.env` 中定义变量（以 `VITE_GLOB_` 开头）
2. 在 `_app.config.js` 中添加对应配置
3. 在代码中通过 `import.meta.env` 使用

### 4.3 偏好设置

#### 配置文件
`src/preferences.ts` - 应用偏好设置

```typescript
import { defineOverridesPreferences } from '@vben/preferences';

export const overridesPreferences = defineOverridesPreferences({
  // 应用配置
  app: {
    name: 'Vben Admin',
    enableCheckUpdates: true,
    enableRefreshToken: true,
  },
  
  // 主题配置
  theme: {
    mode: 'light', // 'light' | 'dark' | 'auto'
    colorPrimary: '#1677ff',
    radius: '0.5rem',
  },
  
  // 布局配置
  layout: {
    mode: 'sidebar-nav', // 布局模式
    sidebarCollapsed: false, // 侧边栏是否折叠
    headerHeight: 50, // 顶栏高度
    sidebarWidth: 230, // 侧边栏宽度
  },
  
  // 标签页配置
  tabbar: {
    enable: true, // 是否启用
    height: 40, // 高度
    keepAlive: true, // 是否缓存
  },
  
  // 面包屑配置
  breadcrumb: {
    enable: true,
    showHome: true,
    showIcon: true,
  },
  
  // 页脚配置
  footer: {
    enable: true,
    fixed: false,
  },
  
  // 过渡动画
  transition: {
    enable: true,
    name: 'fade-slide',
  },
});
```

⚠️ **重要**: 更改配置后需清空浏览器缓存！

---

## 五、图标使用

### 5.1 Iconify 图标（推荐）

#### 特点
- 10万+ 图标
- 按需加载
- 支持在线/离线使用

#### 查找图标
访问 [Iconify](https://icon-sets.iconify.design/) 或 [Icônes](https://icones.js.org/)

#### 使用方式

**在路由中使用：**
```typescript
{
  meta: {
    icon: 'mdi:home', // 格式：图标集:图标名
  },
}
```

**在组件中使用：**
```vue
<template>
  <Icon icon="mdi:home" />
</template>

<script setup>
import { Icon } from '@vben/icons';
</script>
```

#### 新增图标集
```typescript
// 在 src/main.ts 中
import { addCollection } from '@iconify/vue';
import myIcons from './my-icons.json';

addCollection(myIcons);
```

### 5.2 SVG 图标（推荐）

#### 新增 SVG
将 SVG 文件放入 `src/assets/icons/` 目录。

#### 使用
```vue
<template>
  <!-- 自动引入 src/assets/icons/logo.svg -->
  <SvgIcon name="logo" />
</template>

<script setup>
import { SvgIcon } from '@vben/icons';
</script>
```

### 5.3 Tailwind CSS 图标

使用 `@iconify/tailwind` 插件：

```html
<span class="i-mdi-home text-2xl"></span>
```

---

## 六、样式开发

### 6.1 项目结构

```
src/
  ├── styles/
  │   ├── index.scss        # 全局样式入口
  │   ├── variables.scss    # 变量定义
  │   └── mixins.scss       # 混入
```

### 6.2 技术栈

- **Tailwind CSS** - 原子化 CSS 框架（主要）
- **SCSS** - CSS 预处理器
- **PostCSS** - CSS 后处理器

### 6.3 Tailwind CSS

#### 使用方式
```vue
<template>
  <div class="flex items-center justify-between p-4 bg-white rounded-lg shadow">
    <h1 class="text-2xl font-bold text-gray-800">标题</h1>
    <button class="px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600">
      按钮
    </button>
  </div>
</template>
```

#### 自定义配置
在 `tailwind.config.js` 中扩展：

```javascript
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: '#1677ff',
      },
    },
  },
};
```

### 6.4 SCSS

#### 使用变量
```scss
// variables.scss
$primary-color: #1677ff;
$border-radius: 4px;

// 组件中使用
.my-component {
  color: $primary-color;
  border-radius: $border-radius;
}
```

### 6.5 BEM 规范

推荐使用 BEM 命名规范：

```scss
.block {
  &__element {
    // 元素样式
  }
  
  &--modifier {
    // 修饰符样式
  }
}
```

### 6.6 CSS Modules

```vue
<template>
  <div :class="$style.container">
    <h1 :class="$style.title">标题</h1>
  </div>
</template>

<style module>
.container {
  padding: 20px;
}

.title {
  font-size: 24px;
}
</style>
```

---

## 七、外部模块

### 7.1 安装依赖

```bash
# 安装到指定应用
pnpm --filter @vben/web-antd add lodash-es

# 安装到所有应用
pnpm add lodash-es -w

# 安装开发依赖
pnpm add -D @types/lodash-es
```

### 7.2 使用方式

#### 全局引入
```typescript
// src/main.ts
import 'element-plus/dist/index.css';
```

#### 局部引入（推荐）
```vue
<script setup>
import { debounce } from 'lodash-es';

const handleInput = debounce(() => {
  // 处理逻辑
}, 300);
</script>
```

---

## 八、构建与部署

### 8.1 构建

```bash
# 构建生产环境
pnpm build

# 构建指定应用
pnpm build:antd

# 构建并分析
pnpm build:analyze
```

### 8.2 预览

```bash
# 预览构建结果
pnpm preview
```

### 8.3 压缩

#### 开启 gzip 压缩
```typescript
// vite.config.ts
import viteCompression from 'vite-plugin-compression';

export default {
  plugins: [
    viteCompression({
      algorithm: 'gzip',
    }),
  ],
};
```

#### 开启 brotli 压缩
```typescript
viteCompression({
  algorithm: 'brotliCompress',
});
```

### 8.4 构建分析

```bash
pnpm build:analyze
```

会自动打开分析页面，查看各模块大小。

### 8.5 部署

#### 前端路由模式

**Hash 模式：**
```typescript
// router/index.ts
createRouter({
  history: createWebHashHistory(),
});
```

**History 模式：**
```typescript
createRouter({
  history: createWebHistory(),
});
```

#### Nginx 配置（History 模式）

```nginx
server {
  listen 80;
  server_name example.com;
  
  location / {
    root /usr/share/nginx/html;
    try_files $uri $uri/ /index.html;
  }
  
  # API 代理
  location /api {
    proxy_pass http://backend:3000;
  }
}
```

### 8.6 跨域处理

#### 开发环境
在 `vite.config.ts` 中配置代理：

```typescript
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
};
```

#### 生产环境
- 使用 Nginx 反向代理
- 后端配置 CORS
- 使用同域部署

---

## 九、服务端交互

### 9.1 开发环境交互

#### 本地开发跨域配置
```typescript
// vite.config.ts
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:5320',
        changeOrigin: true,
      },
    },
  },
};
```

#### 没有跨域时的配置
直接配置 API 地址：

```bash
# .env.development
VITE_API_URL=http://localhost:5320/api
```

### 9.2 生产环境交互

#### 接口地址配置
```bash
# .env.production
VITE_GLOB_API_URL=https://api.example.com
```

#### 跨域处理
1. Nginx 反向代理（推荐）
2. 后端 CORS 配置
3. 同域部署

### 9.3 接口请求配置

#### 基础配置
```typescript
// src/api/request.ts
import { requestClient } from '#/api/request';

export async function getUserInfo() {
  return requestClient.get('/user/info');
}
```

#### 扩展配置
```typescript
requestClient.get('/user/info', {
  headers: {
    'Custom-Header': 'value',
  },
  params: {
    id: 1,
  },
});
```

#### 请求示例
```typescript
// GET 请求
const data = await requestClient.get('/api/user');

// POST 请求
const result = await requestClient.post('/api/user', {
  name: '张三',
  age: 18,
});

// PUT 请求
await requestClient.put('/api/user/1', {
  name: '李四',
});

// DELETE 请求
await requestClient.delete('/api/user/1');
```

### 9.4 多个接口地址

```typescript
// 创建新的请求实例
import { createHttpClient } from '@vben/request';

const otherClient = createHttpClient({
  baseURL: 'https://other-api.com',
});
```

### 9.5 刷新 Token

框架内置 Token 刷新机制，在 `src/api/core/auth.ts` 中配置。

### 9.6 数据 Mock

#### Nitro Mock 服务
位于 `apps/backend-mock/`，基于 Nitro 框架。

#### 使用方式
```bash
# 自动启动（运行 pnpm dev 时）
pnpm dev:antd

# Mock API 地址
http://localhost:5320/api
```

#### 添加 Mock 接口
在 `apps/backend-mock/api/` 下创建文件：

```typescript
// apps/backend-mock/api/user/list.ts
import { eventHandler } from 'h3';

export default eventHandler(() => {
  return {
    code: 0,
    data: [
      { id: 1, name: '张三' },
      { id: 2, name: '李四' },
    ],
  };
});
```

#### 关闭 Mock 服务
修改 `package.json` 中的启动脚本，移除 Mock 相关配置。

---

## 十、登录认证

### 10.1 登录页面调整

登录页面位于：`src/views/_core/authentication/login.vue`

### 10.2 登录表单调整

修改表单配置：

```typescript
const [Form, formApi] = useVbenForm({
  schema: [
    {
      component: 'Input',
      fieldName: 'username',
      label: '用户名',
      rules: 'required',
    },
    {
      component: 'InputPassword',
      fieldName: 'password',
      label: '密码',
      rules: 'required',
    },
  ],
});
```

### 10.3 接口对接流程

#### 前置条件
确保后端接口返回格式：

```json
{
  "code": 0,
  "data": {
    "accessToken": "xxx",
    "refreshToken": "xxx",
    "user": {
      "id": 1,
      "username": "admin",
      "realName": "管理员",
      "roles": ["admin"]
    }
  },
  "message": "success"
}
```

#### 登录接口
修改 `src/api/core/auth.ts`：

```typescript
export async function loginApi(data: LoginParams) {
  return requestClient.post<UserInfo>('/auth/login', data);
}
```

#### 获取用户信息
```typescript
export async function getUserInfoApi() {
  return requestClient.get<UserInfo>('/user/info');
}
```

---

## 十一、主题定制

### 11.1 CSS 变量

框架使用 CSS 变量实现主题切换，主要变量：

```css
:root {
  --primary: #1677ff;
  --background: #ffffff;
  --foreground: #000000;
  --border: #e5e7eb;
  --radius: 0.5rem;
}
```

### 11.2 覆盖默认 CSS 变量

#### 默认主题下
```css
/* src/styles/theme.css */
:root {
  --primary: #00b96b;
  --radius: 0.75rem;
}
```

#### 黑暗模式下
```css
.dark {
  --primary: #00b96b;
  --background: #1f1f1f;
}
```

### 11.3 更改品牌主色

在 `src/preferences.ts` 中：

```typescript
export const overridesPreferences = defineOverridesPreferences({
  theme: {
    colorPrimary: '#00b96b',
  },
});
```

### 11.4 内置主题

框架内置多套主题：
- **默认主题** - Default
- **深蓝主题** - Deep Blue
- **深绿主题** - Deep Green
- **玫瑰主题** - Rose
- **天空蓝主题** - Sky Blue
- **紫色主题** - Violet
- **黄色主题** - Yellow
- **青色主题** - Cyan
- **粉色主题** - Pink
- **绿色主题** - Green
- **灰色主题** - Gray

### 11.5 新增主题

```typescript
// src/preferences.ts
export const overridesPreferences = defineOverridesPreferences({
  theme: {
    builtinType: 'custom',
    colorPrimary: '#ff6b6b',
  },
});
```

### 11.6 黑暗模式

```typescript
export const overridesPreferences = defineOverridesPreferences({
  theme: {
    mode: 'dark', // 'light' | 'dark' | 'auto'
  },
});
```

### 11.7 自定义侧边栏颜色

#### 默认主题
```css
:root {
  --sidebar-background: #001529;
  --sidebar-foreground: #ffffff;
}
```

#### 黑暗模式
```css
.dark {
  --sidebar-background: #1f1f1f;
}
```

### 11.8 色弱模式

```typescript
export const overridesPreferences = defineOverridesPreferences({
  theme: {
    colorWeakMode: true,
  },
});
```

### 11.9 灰色模式

```typescript
export const overridesPreferences = defineOverridesPreferences({
  theme: {
    grayMode: true,
  },
});
```

---

## 十二、权限控制

### 12.1 前端访问控制

#### 步骤
1. 在路由 meta 中定义权限码
2. 登录后获取用户权限
3. 框架自动过滤无权限路由

#### 示例
```typescript
{
  path: '/system/user',
  meta: {
    authority: ['system:user:view'],
  },
}
```

#### 菜单可见但禁止访问
```typescript
{
  meta: {
    authority: ['admin'],
    menuVisibleWithForbidden: true,
  },
}
```

### 12.2 后端访问控制

#### 步骤
1. 登录后从后端获取路由数据
2. 框架动态生成路由和菜单
3. 用户只能看到有权限的菜单

#### 配置
```typescript
// src/router/access.ts
export const accessRoutes = {
  mode: 'backend', // 'frontend' | 'backend' | 'mixed'
};
```

### 12.3 混合访问控制

前端定义路由结构，后端返回权限码，两者结合。

```typescript
export const accessRoutes = {
  mode: 'mixed',
};
```

### 12.4 按钮细粒度控制

#### 使用权限码
```vue
<template>
  <VbenButton v-if="hasAuthority('system:user:create')">
    新增用户
  </VbenButton>
</template>

<script setup>
import { useAccess } from '@vben/access';

const { hasAuthority } = useAccess();
</script>
```

#### 使用角色
```vue
<template>
  <VbenButton v-if="hasRole('admin')">
    管理员操作
  </VbenButton>
</template>

<script setup>
import { useAccess } from '@vben/access';

const { hasRole } = useAccess();
</script>
```

---

## 十三、国际化

### 13.1 IDE 插件

推荐安装 **i18n Ally** 插件，提供可视化翻译管理。

### 13.2 配置默认语言

```typescript
// src/preferences.ts
export const overridesPreferences = defineOverridesPreferences({
  app: {
    locale: 'zh-CN', // 'zh-CN' | 'en-US'
  },
});
```

### 13.3 动态切换语言

```typescript
import { useLocale } from '@vben/locales';

const { setLocale } = useLocale();

// 切换到英文
setLocale('en-US');
```

### 13.4 新增翻译文本

在 `src/locales/langs/` 下添加翻译：

```typescript
// src/locales/langs/zh-CN/common.json
{
  "save": "保存",
  "cancel": "取消",
  "confirm": "确认"
}

// src/locales/langs/en-US/common.json
{
  "save": "Save",
  "cancel": "Cancel",
  "confirm": "Confirm"
}
```

### 13.5 使用翻译文本

#### 在代码中使用
```typescript
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const message = t('common.save'); // "保存"
```

#### 在模板中使用
```vue
<template>
  <button>{{ $t('common.save') }}</button>
</template>
```

### 13.6 新增语言包

1. 在 `src/locales/langs/` 下创建新语言目录
2. 添加翻译文件
3. 在 `src/locales/index.ts` 中注册

### 13.7 界面切换语言功能

框架内置语言切换组件，在顶栏显示。

### 13.8 远程加载语言包

```typescript
import { loadLocaleMessages } from '@vben/locales';

// 异步加载语言包
await loadLocaleMessages('zh-CN', async () => {
  const messages = await fetch('/api/locale/zh-CN').then(r => r.json());
  return messages;
});
```

### 13.9 移除国际化

如果不需要国际化，可以：
1. 删除 `src/locales/` 目录
2. 移除相关依赖
3. 直接使用中文文本

---

## 十四、常用功能

### 14.1 登录认证过期

#### 跳转登录页面
```typescript
import { useAuthStore } from '@vben/stores';

const authStore = useAuthStore();

// 清除认证信息并跳转登录页
authStore.logout();
```

#### 打开登录弹窗
```typescript
import { useAuthStore } from '@vben/stores';

const authStore = useAuthStore();

// 打开登录弹窗
authStore.openLoginModal();
```

### 14.2 动态标题

```typescript
// src/preferences.ts
export const overridesPreferences = defineOverridesPreferences({
  app: {
    dynamicTitle: true, // 启用动态标题
  },
});
```

标题会自动根据当前路由的 `meta.title` 更新。

### 14.3 页面水印

```typescript
// src/preferences.ts
export const overridesPreferences = defineOverridesPreferences({
  app: {
    watermark: true,
    watermarkContent: '内部系统',
  },
});
```

### 14.4 检查更新

框架内置版本检查功能，会自动检测新版本并提示用户刷新。

#### 配置
```typescript
export const overridesPreferences = defineOverridesPreferences({
  app: {
    enableCheckUpdates: true,
    checkUpdatesInterval: 60000, // 检查间隔（ms）
  },
});
```

### 14.5 全局 Loading

#### 原理
在应用加载时显示 loading 动画，提升用户体验。

#### 关闭
在 `index.html` 中移除相关代码。

#### 自定义
修改 `index.html` 中的 loading 样式。

### 14.6 组件库切换

#### 新增组件库应用
1. 复制现有应用目录（如 `apps/web-antd`）
2. 修改 `package.json` 中的依赖
3. 调整适配器配置
4. 更新组件引用

---

## 附录

### A. 常见问题

#### 1. 依赖安装失败
- 检查 Node.js 版本（需要 20.15.0+）
- 使用 `pnpm clean && pnpm install` 重新安装
- 配置国内镜像源

#### 2. 启动报错
- 清空浏览器缓存
- 检查端口是否被占用
- 查看控制台错误信息

#### 3. 构建失败
- 检查代码是否有语法错误
- 确认所有依赖已安装
- 查看构建日志

### B. 最佳实践

1. **代码规范**: 遵循 ESLint 和 Prettier 配置
2. **组件拆分**: 单个组件不超过 300 行
3. **性能优化**: 使用懒加载、虚拟滚动
4. **类型安全**: 充分利用 TypeScript
5. **注释文档**: 关键逻辑添加注释

### C. 相关链接

- 官方文档：https://doc.vben.pro
- GitHub：https://github.com/vbenjs/vue-vben-admin
- 在线预览：https://vben.pro
- Vue 3 文档：https://cn.vuejs.org
- Vite 文档：https://cn.vitejs.dev
- Tailwind CSS：https://tailwindcss.com

---

**文档版本**: 基于 Vben Admin 5.7.0  
**最后更新**: 2026-03-24
