# 多租户权限管理方案设计

> 基于 Vben Admin 5.7.0 的多租户 SaaS 系统权限管理方案

## 一、方案概述

### 1.1 核心特性

- ✅ **多租户隔离** - 租户间数据完全隔离
- ✅ **简化设计** - 用户ID = 租户ID，一个租户一个管理员
- ✅ **层级权限** - 系统超管 → 租户管理员 → 自定义角色 → 普通会员
- ✅ **混合控制** - 前端 + 后端混合访问控制模式
- ✅ **动态分配** - 支持动态创建角色和分配权限

### 1.2 设计理念

本方案采用**简化的租户模型**：
- 每个租户只有一个管理员（所有者）
- 管理员的用户ID即为租户ID
- 不需要单独的租户表，降低系统复杂度
- 适合中小型 SaaS 应用

## 二、权限层级

### 2.1 角色体系

```
系统超管 (super_admin)
└─ 创建和管理租户（管理员账号）

租户管理员 (admin) [用户ID = 租户ID]
├─ 租户内最高权限
├─ 创建自定义角色
├─ 创建普通会员
└─ 管理租户所有数据

自定义角色 (custom)
├─ 由管理员创建
├─ 权限不超过管理员
└─ 只能访问本租户数据

普通会员 (member)
└─ 基础权限，访问自己的数据
```

### 2.2 角色定义

```typescript
// 系统角色
const SYSTEM_ROLES = {
  SUPER_ADMIN: 'super_admin',  // 平台超管
} as const;

// 租户角色
const TENANT_ROLES = {
  ADMIN: 'admin',              // 租户管理员
  CUSTOM: 'custom',            // 自定义角色
  MEMBER: 'member',            // 普通会员
} as const;
```

## 三、数据库设计

### 3.1 核心设计思路

**简化原则**：
- 不需要独立的租户表
- 管理员的 `id` 即为租户ID
- 会员的 `parent_id` 指向管理员ID（租户ID）
- 通过 `parent_id` 字段实现租户隔离

### 3.2 用户表 (users)

```sql
CREATE TABLE users (
  id VARCHAR(50) PRIMARY KEY COMMENT '用户ID',
  parent_id VARCHAR(50) COMMENT '父用户ID，NULL=管理员，有值=会员',
  username VARCHAR(50) NOT NULL COMMENT '用户名',
  password VARCHAR(255) NOT NULL COMMENT '密码',
  real_name VARCHAR(50) COMMENT '真实姓名',
  email VARCHAR(100) COMMENT '邮箱',
  phone VARCHAR(20) COMMENT '手机号',
  role VARCHAR(50) NOT NULL COMMENT '角色：super_admin/admin/custom/member',
  role_id VARCHAR(50) COMMENT '自定义角色ID',
  tenant_name VARCHAR(100) COMMENT '租户名称（仅管理员有）',
  status INT DEFAULT 1 COMMENT '状态：1-启用 0-禁用',
  created_by VARCHAR(50) COMMENT '创建者ID',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  INDEX idx_parent (parent_id),
  INDEX idx_role (role),
  INDEX idx_status (status),
  UNIQUE KEY uk_parent_username (parent_id, username)
) COMMENT='用户表（包含租户信息）';
```

**字段说明**：
- `parent_id = NULL`：表示是管理员（租户所有者）
- `parent_id = 管理员ID`：表示是该租户下的会员
- `tenant_name`：只有管理员才有租户名称

### 3.3 角色表 (roles)

```sql
CREATE TABLE roles (
  id VARCHAR(50) PRIMARY KEY COMMENT '角色ID',
  admin_id VARCHAR(50) NOT NULL COMMENT '管理员ID（即租户ID）',
  name VARCHAR(50) NOT NULL COMMENT '角色名称',
  code VARCHAR(50) NOT NULL COMMENT '角色编码',
  description VARCHAR(255) COMMENT '角色描述',
  permissions JSON COMMENT '权限配置',
  status INT DEFAULT 1 COMMENT '状态：1-启用 0-禁用',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  INDEX idx_admin (admin_id),
  UNIQUE KEY uk_admin_code (admin_id, code)
) COMMENT='角色表';
```

**说明**：`admin_id` 即为租户ID（管理员的用户ID）

### 3.4 数据结构示例

```typescript
// 管理员（租户所有者）
{
  id: 'user_001',           // 用户ID = 租户ID
  parentId: null,           // NULL 表示是管理员
  username: 'admin',
  role: 'admin',
  tenantName: '某某公司',
}

// 普通会员
{
  id: 'user_002',
  parentId: 'user_001',     // 父用户ID = 租户ID
  username: 'member1',
  role: 'member',
}

// 自定义角色的会员
{
  id: 'user_003',
  parentId: 'user_001',     // 同一个租户
  username: 'staff1',
  role: 'custom',
  roleId: 'role_001',
}
```

## 四、权限码设计

### 4.1 权限码格式

```typescript
// 系统级权限
'system:admin:create'      // 创建管理员（租户）
'system:admin:delete'      // 删除管理员（租户）

// 租户级权限（不需要租户ID前缀，通过用户关系判断）
'user:view'                // 查看用户
'user:create'              // 创建用户
'user:edit'                // 编辑用户
'user:delete'              // 删除用户
'role:manage'              // 管理角色
'member:view'              // 查看会员
```

### 4.2 权限码命名规范

```
格式：{module}:{action}

module: user | role | member | order | product 等
action: view | create | edit | delete | manage | export 等
```

### 4.3 预设权限

```typescript
// 租户管理员权限（拥有全部）
const ADMIN_PERMISSIONS = [
  'user:view', 'user:create', 'user:edit', 'user:delete',
  'role:view', 'role:create', 'role:edit', 'role:delete',
  'member:view', 'member:create', 'member:edit', 'member:delete',
];

// 普通会员基础权限
const MEMBER_PERMISSIONS = [
  'profile:view',
  'profile:edit',
];
```

## 五、核心业务流程

### 5.1 创建管理员（租户）

```typescript
/**
 * 系统超管创建管理员（即创建租户）
 */
async function createAdmin(data: {
  username: string;
  password: string;
  tenantName: string;
  email: string;
}) {
  return await db.users.create({
    id: generateId('user'),      // 用户ID = 租户ID
    parentId: null,               // NULL 表示是管理员
    username: data.username,
    password: hashPassword(data.password),
    email: data.email,
    role: 'admin',
    tenantName: data.tenantName,  // 租户名称
    status: 1,
  });
}
```

### 5.2 创建自定义角色

```typescript
/**
 * 管理员创建自定义角色
 */
async function createRole(adminId: string, data: {
  name: string;
  code: string;
  permissions: string[];
}) {
  const admin = await db.users.findById(adminId);
  
  if (admin.role !== 'admin') {
    throw new Error('只有管理员可以创建角色');
  }
  
  // 验证权限不超过管理员
  const invalidPerms = data.permissions.filter(
    p => !ADMIN_PERMISSIONS.includes(p)
  );
  if (invalidPerms.length > 0) {
    throw new Error('角色权限不能超过管理员权限');
  }
  
  return await db.roles.create({
    id: generateId('role'),
    adminId: adminId,             // 管理员ID = 租户ID
    name: data.name,
    code: data.code,
    permissions: data.permissions,
    status: 1,
  });
}
```

### 5.3 创建会员

```typescript
/**
 * 管理员创建会员
 */
async function createMember(adminId: string, data: {
  username: string;
  password: string;
  roleId?: string;
}) {
  const admin = await db.users.findById(adminId);
  
  if (admin.role !== 'admin') {
    throw new Error('只有管理员可以创建会员');
  }
  
  return await db.users.create({
    id: generateId('user'),
    parentId: adminId,            // 父用户ID = 管理员ID = 租户ID
    username: data.username,
    password: hashPassword(data.password),
    role: data.roleId ? 'custom' : 'member',
    roleId: data.roleId,
    status: 1,
  });
}
```

## 六、权限控制实现

### 6.1 获取租户ID工具函数

```typescript
/**
 * 获取用户的租户ID
 */
function getTenantId(user: User): string | null {
  if (user.role === 'super_admin') {
    return null;  // 系统超管无租户
  }
  if (user.role === 'admin') {
    return user.id;  // 管理员的ID即租户ID
  }
  return user.parentId;  // 会员的父用户ID即租户ID
}
```

### 6.2 权限检查中间件

```typescript
/**
 * 权限检查中间件
 */
function checkPermission(requiredPermission: string) {
  return async (req: Request, res: Response, next: NextFunction) => {
    const user = req.user;
    
    if (!user) {
      return res.status(401).json({ message: '未登录' });
    }
    
    // 系统超管拥有所有权限
    if (user.role === 'super_admin') {
      return next();
    }
    
    // 管理员拥有租户内所有权限
    if (user.role === 'admin') {
      return next();
    }
    
    // 自定义角色检查具体权限
    if (user.role === 'custom' && user.roleId) {
      const role = await db.roles.findById(user.roleId);
      if (role.permissions?.includes(requiredPermission)) {
        return next();
      }
    }
    
    // 普通会员检查基础权限
    if (user.role === 'member') {
      if (MEMBER_PERMISSIONS.includes(requiredPermission)) {
        return next();
      }
    }
    
    return res.status(403).json({ message: '无权限' });
  };
}
```

### 6.3 数据隔离中间件

```typescript
/**
 * 租户数据隔离
 */
function tenantIsolation() {
  return async (req: Request, res: Response, next: NextFunction) => {
    const user = req.user;
    const tenantId = getTenantId(user);
    
    // 注入租户ID到请求
    req.tenantId = tenantId;
    
    next();
  };
}
```

### 6.4 查询租户数据

```typescript
/**
 * 查询租户下的用户列表
 */
async function getUserList(currentUser: User) {
  const tenantId = getTenantId(currentUser);
  
  if (!tenantId) {
    // 系统超管查看所有用户
    return await db.users.find();
  }
  
  // 查询租户下的所有用户（管理员 + 会员）
  return await db.users.find({
    $or: [
      { id: tenantId },        // 管理员自己
      { parentId: tenantId },  // 所有会员
    ],
  });
}
```

## 七、前端权限配置

### 7.1 访问模式配置

```typescript
// src/preferences.ts
import { defineOverridesPreferences } from '@vben/preferences';

export const overridesPreferences = defineOverridesPreferences({
  app: {
    accessMode: 'mixed', // 混合访问控制模式
  },
});
```

### 7.2 路由权限配置

```typescript
// src/router/routes/modules/system.ts
import type { RouteRecordRaw } from 'vue-router';
import { BasicLayout } from '#/layouts';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'mdi:cog',
      title: '系统管理',
      authority: ['admin'], // 只有租户管理员可见
    },
    name: 'System',
    path: '/system',
    children: [
      {
        name: 'UserManagement',
        path: '/users/list',
        component: () => import('#/views/users/list.vue'),
        meta: {
          title: '用户管理',
          icon: 'mdi:account-group',
          authority: ['admin'],
        },
      },
      {
        name: 'RoleManagement',
        path: '/roles/list',
        component: () => import('#/views/roles/list.vue'),
        meta: {
          title: '角色管理',
          icon: 'mdi:shield-account',
          authority: ['admin'],
        },
      },
      {
        name: 'MemberManagement',
        path: '/members/list',
        component: () => import('#/views/members/list.vue'),
        meta: {
          title: '会员管理',
          icon: 'mdi:account-multiple',
          authority: ['admin'],
        },
      },
    ],
  },
];

export default routes;
```

### 7.3 按钮权限控制

```vue
<script setup lang="ts">
import { useAccess } from '@vben/access';
import { AccessControl } from '@vben/access';

const { hasAccessByCodes, hasAccessByRoles } = useAccess();
const authStore = useAuthStore();
const tenantId = authStore.userInfo?.tenantId;
</script>

<template>
  <div class="p-4">
    <!-- 方式1：使用角色控制 -->
    <Button 
      v-if="hasAccessByRoles(['admin'])" 
      type="primary"
      @click="handleCreateUser"
    >
      新增用户
    </Button>
    
    <!-- 方式2：使用权限码控制 -->
    <Button 
      v-if="hasAccessByCodes([`${tenantId}:user:create`])" 
      type="primary"
      @click="handleCreateUser"
    >
      新增用户
    </Button>
    
    <!-- 方式3：使用组件控制 -->
    <AccessControl :codes="['admin']">
      <Button danger @click="handleBatchDelete">
        批量删除
      </Button>
    </AccessControl>
    
    <!-- 方式4：使用指令控制 -->
    <Button v-access:role="'admin'" @click="handleExport">
      导出数据
    </Button>
  </div>
</template>
```

## 八、工具函数

### 8.1 ID 生成

```typescript
/**
 * 简单的ID生成器
 */
export function generateId(prefix: string = 'id'): string {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 8);
  return `${prefix}_${timestamp}_${random}`;
}

// 使用示例
const userId = generateId('user');    // user_1234567890_abc123
const roleId = generateId('role');    // role_1234567892_ghi789
```

### 8.2 密码加密

```typescript
import bcrypt from 'bcrypt';

/**
 * 密码加密
 */
export async function hashPassword(password: string): Promise<string> {
  return await bcrypt.hash(password, 10);
}

/**
 * 密码验证
 */
export async function verifyPassword(
  password: string, 
  hash: string
): Promise<boolean> {
  return await bcrypt.compare(password, hash);
}
```

## 九、最佳实践

### 9.1 数据隔离原则

- **用户ID = 租户ID** - 管理员的用户ID即为租户ID
- **parent_id 关联** - 会员通过 parent_id 关联到管理员
- **查询过滤** - 所有查询必须过滤租户数据
- **级联删除** - 删除管理员时级联删除所有会员

### 9.2 权限设计原则

- **最小权限** - 用户只拥有必需的权限
- **权限继承** - 自定义角色权限不超过管理员
- **显式授权** - 权限必须显式授予
- **定期审计** - 定期审查权限配置

### 9.3 开发规范

**命名规范**：
- 用户ID：`user_1234567890_abc123`
- 角色ID：`role_1234567890_abc123`
- 权限码：`{module}:{action}`（如 `user:create`）

**HTTP 状态码**：
- `401` - 未登录
- `403` - 无权限
- `404` - 资源不存在
- `500` - 服务器错误

## 十、常见问题

### Q1: 如何删除租户？

删除管理员即删除租户，需级联删除所有数据：

```typescript
async function deleteAdmin(adminId: string) {
  await db.transaction(async (trx) => {
    // 1. 删除所有会员
    await trx('users').where('parent_id', adminId).delete();
    
    // 2. 删除所有角色
    await trx('roles').where('admin_id', adminId).delete();
    
    // 3. 删除管理员
    await trx('users').where('id', adminId).delete();
  });
}
```

### Q2: 如何限制租户配额？

在用户表添加配额字段：

```sql
ALTER TABLE users ADD COLUMN member_quota INT DEFAULT 100;
ALTER TABLE users ADD COLUMN member_count INT DEFAULT 0;
```

创建会员时检查：

```typescript
async function createMember(adminId: string, data: any) {
  const admin = await db.users.findById(adminId);
  
  if (admin.memberCount >= admin.memberQuota) {
    throw new Error('会员数量已达上限');
  }
  
  // 创建并更新计数
  await db.transaction(async (trx) => {
    await trx('users').insert({ ...data, parentId: adminId });
    await trx('users').where('id', adminId).increment('member_count', 1);
  });
}
```

### Q3: 管理员可以转移吗？

不建议转移，因为管理员ID即租户ID。如需转移：
1. 创建新管理员
2. 将所有会员的 parent_id 更新为新管理员ID
3. 删除旧管理员

## 十一、总结

本方案采用**简化的多租户设计**，核心特点：

✅ **简单高效** - 用户ID = 租户ID，无需独立租户表
✅ **数据隔离** - 通过 parent_id 实现租户隔离
✅ **混合控制** - 前端 + 后端混合权限控制
✅ **易于扩展** - 支持动态角色和权限分配
✅ **性能优化** - 减少表关联，查询更快

**适用场景**：
- 中小型 SaaS 应用
- 一个租户一个管理员的业务模型
- 需要快速开发和部署的项目

**技术栈**：
- 前端：Vben Admin 5.7.0 + Vue 3 + TypeScript
- 后端：Node.js + Express/Koa
- 数据库：MySQL/PostgreSQL
- 认证：JWT + bcrypt
