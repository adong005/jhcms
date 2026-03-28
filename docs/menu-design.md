# ADCMS 系统菜单设计

## 一、菜单结构概览

### 1.1 一级菜单

| 菜单名称 | 图标 | 路由 | 说明 |
|---------|------|------|------|
| 工作台 | `mdi:view-dashboard` | `/dashboard` | 系统首页，展示概览信息 |
| 系统管理 | `mdi:cog` | `/system` | 系统配置和管理功能 |
| 信息管理 | `mdi:information` | `/info` | 信息内容管理 |
| 个人中心 | `mdi:account-circle` | `/profile` | 个人设置和偏好 |

### 1.2 菜单层级关系

```
工作台
├── 概览
├── 快捷入口
└── 最新动态

系统管理
├── 用户管理
├── 权限管理
│   ├── 角色管理
│   └── 菜单管理
└── 系统配置
    ├── 基础设置
    ├── 字典管理
    └── 操作日志

信息管理
├── 信息分类
│   ├── 分类列表
│   └── 分类配置
├── 信息列表
│   ├── 信息发布
│   ├── 信息审核
│   └── 信息归档
└── 内容管理
    ├── 文章管理
    ├── 资源管理
    └── 评论管理

个人中心
├── 个人信息
├── 修改密码
├── 消息通知
└── 操作记录
```

## 二、详细菜单配置

### 2.1 工作台

```json
{
  "id": 1,
  "name": "Dashboard",
  "path": "/dashboard",
  "type": "menu",
  "component": "/dashboard/index",
  "meta": {
    "title": "工作台",
    "icon": "mdi:view-dashboard",
    "order": 1,
    "keepAlive": true
  }
}
```

### 2.2 系统管理

```json
{
  "id": 2,
  "name": "System",
  "path": "/system",
  "type": "catalog",
  "meta": {
    "title": "系统管理",
    "icon": "mdi:cog",
    "order": 2
  },
  "children": [
    {
      "id": 21,
      "pid": 2,
      "name": "UserManagement",
      "path": "/system/user",
      "type": "menu",
      "component": "/system/user/list",
      "meta": {
        "title": "用户管理",
        "icon": "mdi:account-group",
        "authority": ["system:user:view"]
      },
      "buttons": [
        { "code": "system:user:create", "text": "新增" },
        { "code": "system:user:edit", "text": "编辑" },
        { "code": "system:user:delete", "text": "删除" },
        { "code": "system:user:export", "text": "导出" },
        { "code": "system:user:reset", "text": "重置密码" },
        { "code": "system:user:role:assign", "text": "分配角色" }
      ]
    },
    {
      "id": 22,
      "pid": 2,
      "name": "PermissionManagement",
      "path": "/system/permission",
      "type": "catalog",
      "meta": {
        "title": "权限管理",
        "icon": "mdi:shield-account"
      },
      "children": [
        {
          "id": 221,
          "pid": 22,
          "name": "RoleManagement",
          "path": "/system/permission/role",
          "type": "menu",
          "component": "/system/permission/role",
          "meta": {
            "title": "角色管理",
            "icon": "mdi:shield-account",
            "authority": ["system:role:view"]
          },
          "buttons": [
            { "code": "system:role:create", "text": "新增" },
            { "code": "system:role:edit", "text": "编辑" },
            { "code": "system:role:delete", "text": "删除" },
            { "code": "system:role:copy", "text": "复制" },
            { "code": "system:role:menu:assign", "text": "配置菜单" }
          ]
        },
        {
          "id": 222,
          "pid": 22,
          "name": "MenuManagement",
          "path": "/system/permission/menu",
          "type": "menu",
          "component": "/system/permission/menu",
          "meta": {
            "title": "菜单管理",
            "icon": "mdi:menu",
            "authority": ["system:menu:view"]
          },
          "buttons": [
            { "code": "system:menu:create", "text": "新增" },
            { "code": "system:menu:edit", "text": "编辑" },
            { "code": "system:menu:delete", "text": "删除" },
            { "code": "system:menu:sort", "text": "排序" }
          ]
        }
      ]
    },
    {
      "id": 24,
      "pid": 2,
      "name": "SystemConfig",
      "path": "/system/config",
      "type": "catalog",
      "meta": {
        "title": "系统配置",
        "icon": "mdi:cog-outline"
      },
      "children": [
        {
          "id": 241,
          "pid": 24,
          "name": "BasicSettings",
          "path": "/system/config/basic",
          "type": "menu",
          "component": "/system/config/basic",
          "meta": {
            "title": "基础设置",
            "icon": "mdi:application-settings",
            "authority": ["system:config:view"]
          },
          "buttons": [
            { "code": "system:config:edit", "text": "修改" },
            { "code": "system:config:reset", "text": "重置" }
          ]
        },
        {
          "id": 242,
          "pid": 24,
          "name": "DictManagement",
          "path": "/system/config/dict",
          "type": "menu",
          "component": "/system/config/dict",
          "meta": {
            "title": "字典管理",
            "icon": "mdi:book-alphabet",
            "authority": ["system:dict:view"]
          },
          "buttons": [
            { "code": "system:dict:create", "text": "新增" },
            { "code": "system:dict:edit", "text": "编辑" },
            { "code": "system:dict:delete", "text": "删除" }
          ]
        },
        {
          "id": 243,
          "pid": 24,
          "name": "OperationLog",
          "path": "/system/config/log",
          "type": "menu",
          "component": "/system/config/log",
          "meta": {
            "title": "操作日志",
            "icon": "mdi:history",
            "authority": ["system:log:view"]
          },
          "buttons": [
            { "code": "system:log:export", "text": "导出" },
            { "code": "system:log:clear", "text": "清理" }
          ]
        }
      ]
    }
  ]
}
```

### 2.3 信息管理

```json
{
  "id": 3,
  "name": "InfoManagement",
  "path": "/info",
  "type": "catalog",
  "meta": {
    "title": "信息管理",
    "icon": "mdi:information",
    "order": 3
  },
  "children": [
    {
      "id": 31,
      "pid": 3,
      "name": "InfoCategory",
      "path": "/info/category",
      "type": "catalog",
      "meta": {
        "title": "信息分类",
        "icon": "mdi:folder-multiple"
      },
      "children": [
        {
          "id": 311,
          "pid": 31,
          "name": "CategoryList",
          "path": "/info/category/list",
          "type": "menu",
          "component": "/info/category/list",
          "meta": {
            "title": "分类列表",
            "icon": "mdi:format-list-checks",
            "authority": ["info:category:view"]
          },
          "buttons": [
            { "code": "info:category:create", "text": "新增" },
            { "code": "info:category:edit", "text": "编辑" },
            { "code": "info:category:delete", "text": "删除" },
            { "code": "info:category:sort", "text": "排序" }
          ]
        },
        {
          "id": 312,
          "pid": 31,
          "name": "CategoryConfig",
          "path": "/info/category/config",
          "type": "menu",
          "component": "/info/category/config",
          "meta": {
            "title": "分类配置",
            "icon": "mdi:cog-box",
            "authority": ["info:category:config:view"]
          }
        }
      ]
    },
    {
      "id": 32,
      "pid": 3,
      "name": "InfoList",
      "path": "/info/list",
      "type": "catalog",
      "meta": {
        "title": "信息列表",
        "icon": "mdi:file-document-multiple"
      },
      "children": [
        {
          "id": 321,
          "pid": 32,
          "name": "InfoPublish",
          "path": "/info/list/publish",
          "type": "menu",
          "component": "/info/list/publish",
          "meta": {
            "title": "信息发布",
            "icon": "mdi:plus-circle",
            "authority": ["info:publish:view"]
          },
          "buttons": [
            { "code": "info:publish:create", "text": "发布" },
            { "code": "info:publish:draft", "text": "草稿" },
            { "code": "info:publish:schedule", "text": "定时发布" }
          ]
        },
        {
          "id": 322,
          "pid": 32,
          "name": "InfoReview",
          "path": "/info/list/review",
          "type": "menu",
          "component": "/info/list/review",
          "meta": {
            "title": "信息审核",
            "icon": "mdi:clipboard-check",
            "authority": ["info:review:view"]
          },
          "buttons": [
            { "code": "info:review:approve", "text": "通过" },
            { "code": "info:review:reject", "text": "拒绝" },
            { "code": "info:review:batch", "text": "批量审核" }
          ]
        },
        {
          "id": 323,
          "pid": 32,
          "name": "InfoArchive",
          "path": "/info/list/archive",
          "type": "menu",
          "component": "/info/list/archive",
          "meta": {
            "title": "信息归档",
            "icon": "mdi:archive",
            "authority": ["info:archive:view"]
          },
          "buttons": [
            { "code": "info:archive:restore", "text": "恢复" },
            { "code": "info:archive:delete", "text": "彻底删除" }
          ]
        }
      ]
    },
    {
      "id": 33,
      "pid": 3,
      "name": "ContentManagement",
      "path": "/info/content",
      "type": "catalog",
      "meta": {
        "title": "内容管理",
        "icon": "mdi:file-multiple"
      },
      "children": [
        {
          "id": 331,
          "pid": 33,
          "name": "ArticleManagement",
          "path": "/info/content/article",
          "type": "menu",
          "component": "/info/content/article",
          "meta": {
            "title": "文章管理",
            "icon": "mdi:file-document",
            "authority": ["info:article:view"]
          },
          "buttons": [
            { "code": "info:article:create", "text": "新增" },
            { "code": "info:article:edit", "text": "编辑" },
            { "code": "info:article:delete", "text": "删除" },
            { "code": "info:article:copy", "text": "复制" }
          ]
        },
        {
          "id": 332,
          "pid": 33,
          "name": "ResourceManagement",
          "path": "/info/content/resource",
          "type": "menu",
          "component": "/info/content/resource",
          "meta": {
            "title": "资源管理",
            "icon": "mdi:folder-multiple-image",
            "authority": ["info:resource:view"]
          },
          "buttons": [
            { "code": "info:resource:upload", "text": "上传" },
            { "code": "info:resource:delete", "text": "删除" },
            { "code": "info:resource:move", "text": "移动" }
          ]
        },
        {
          "id": 333,
          "pid": 33,
          "name": "CommentManagement",
          "path": "/info/content/comment",
          "type": "menu",
          "component": "/info/content/comment",
          "meta": {
            "title": "评论管理",
            "icon": "mdi:comment-multiple",
            "authority": ["info:comment:view"]
          },
          "buttons": [
            { "code": "info:comment:approve", "text": "审核通过" },
            { "code": "info:comment:reject", "text": "拒绝" },
            { "code": "info:comment:delete", "text": "删除" }
          ]
        }
      ]
    }
  ]
}
```

### 2.4 个人中心

```json
{
  "id": 4,
  "name": "Profile",
  "path": "/profile",
  "type": "catalog",
  "meta": {
    "title": "个人中心",
    "icon": "mdi:account-circle",
    "order": 4
  },
  "children": [
    {
      "id": 41,
      "pid": 4,
      "name": "PersonalInfo",
      "path": "/profile/info",
      "type": "menu",
      "component": "/profile/info",
      "meta": {
        "title": "个人信息",
        "icon": "mdi:account-details"
      },
      "buttons": [
        { "code": "profile:info:edit", "text": "编辑" },
        { "code": "profile:avatar:upload", "text": "上传头像" }
      ]
    },
    {
      "id": 42,
      "pid": 4,
      "name": "ChangePassword",
      "path": "/profile/password",
      "type": "menu",
      "component": "/profile/password",
      "meta": {
        "title": "修改密码",
        "icon": "mdi:lock-reset"
      },
      "buttons": [
        { "code": "profile:password:change", "text": "修改密码" }
      ]
    },
    {
      "id": 43,
      "pid": 4,
      "name": "MessageNotification",
      "path": "/profile/message",
      "type": "menu",
      "component": "/profile/message",
      "meta": {
        "title": "消息通知",
        "icon": "mdi:bell"
      },
      "buttons": [
        { "code": "profile:message:read", "text": "标记已读" },
        { "code": "profile:message:delete", "text": "删除" },
        { "code": "profile:message:setting", "text": "通知设置" }
      ]
    },
    {
      "id": 44,
      "pid": 4,
      "name": "OperationHistory",
      "path": "/profile/history",
      "type": "menu",
      "component": "/profile/history",
      "meta": {
        "title": "操作记录",
        "icon": "mdi:history"
      },
      "buttons": [
        { "code": "profile:history:export", "text": "导出" },
        { "code": "profile:history:clear", "text": "清理" }
      ]
    }
  ]
}
```

## 三、权限设计

### 3.1 角色定义

| 角色 | 说明 | 权限范围 |
|------|------|----------|
| super | 超级管理员 | 所有权限 |
| admin | 系统管理员 | 系统管理、信息管理 |
| editor | 内容编辑 | 信息管理、个人中心 |
| viewer | 只读用户 | 查看权限、个人中心 |

### 3.2 权限码规范

格式：`模块:功能:操作`

示例：
- `system:user:view` - 系统管理-用户管理-查看
- `info:article:create` - 信息管理-文章管理-新增
- `profile:info:edit` - 个人中心-个人信息-编辑

### 3.3 按钮权限

每个菜单可配置按钮权限：
- **查看** - view
- **新增** - create
- **编辑** - edit
- **删除** - delete
- **导出** - export
- **导入** - import
- **审核** - approve
- **发布** - publish
- **重置** - reset

## 四、Mock 数据配置

将以上菜单配置更新到 `apps/backend-mock/utils/mock-data.ts` 的 `MOCK_MENUS` 数组中。

### 4.1 用户菜单映射

```typescript
// super 用户 - 所有菜单
super: {
  username: 'vben',
  roles: ['super'],
  menus: [/* 工作台、系统管理、信息管理、个人中心 */]
},

// admin 用户 - 系统管理、信息管理
admin: {
  username: 'admin',
  roles: ['admin'],
  menus: [/* 工作台、系统管理、信息管理、个人中心 */]
},

// editor 用户 - 信息管理
editor: {
  username: 'editor',
  roles: ['editor'],
  menus: [/* 工作台、信息管理、个人中心 */]
},

// viewer 用户 - 只读权限
viewer: {
  username: 'viewer',
  roles: ['viewer'],
  menus: [/* 工作台、个人中心 */]
}
```

## 五、实施步骤

1. **更新 Mock 数据** - 将菜单配置写入 `mock-data.ts`
2. **清空缓存** - 浏览器清空缓存
3. **重新登录** - 使用不同账号测试
4. **验证权限** - 检查不同角色看到的菜单
5. **测试按钮** - 验证按钮权限控制

## 六、扩展建议

1. **菜单图标** - 使用统一的图标库（Material Design Icons）
2. **面包屑** - 自动根据路由生成面包屑导航
3. **标签页** - 支持多标签页切换
4. **搜索菜单** - 添加菜单搜索功能
5. **收藏菜单** - 支持收藏常用菜单
