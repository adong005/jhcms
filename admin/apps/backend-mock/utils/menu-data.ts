// 定义菜单组件类型
interface MenuComponent {
  id: number;
  pid?: number;
  name: string;
  path: string;
  type: 'menu' | 'catalog' | 'button';
  component?: string;
  meta: {
    title: string;
    icon?: string;
    order?: number;
    authority?: string[];
    keepAlive?: boolean;
    [key: string]: any;
  };
  buttons?: Array<{
    code: string;
    text: string;
  }>;
  children?: MenuComponent[];
  status?: number;
  authCode?: string;
}

// 工作台菜单
const dashboardMenus: MenuComponent[] = [
  {
    id: 1,
    name: 'Analytics',
    path: '/analytics',
    type: 'menu',
    component: '/dashboard/analytics/index',
    meta: {
      title: '工作台',
      icon: 'mdi:view-dashboard',
      order: 1,
      keepAlive: true,
      affixTab: true,
    },
  },
];

// 系统管理菜单 - 一级目录，包含用户、角色、菜单、日志管理
const systemMenus: MenuComponent[] = [
  {
    id: 2,
    name: 'System',
    path: '/system',
    type: 'catalog',
    meta: {
      title: '系统管理',
      icon: 'lucide:settings',
      order: 2,
    },
    children: [
      {
        id: 21,
        pid: 2,
        name: 'Users',
        path: '/users/list',
        type: 'menu',
        component: '/users/list',
        meta: {
          title: '用户管理',
          icon: 'lucide:users',
          authority: ['system:user:view'],
        },
        buttons: [
          { code: 'system:user:create', text: '新增' },
          { code: 'system:user:edit', text: '编辑' },
          { code: 'system:user:delete', text: '删除' },
          { code: 'system:user:export', text: '导出' },
          { code: 'system:user:reset', text: '重置密码' },
          { code: 'system:user:role:assign', text: '分配角色' },
        ],
      },
      {
        id: 22,
        pid: 2,
        name: 'Roles',
        path: '/roles/list',
        type: 'menu',
        component: '/roles/list',
        meta: {
          title: '角色管理',
          icon: 'lucide:shield',
          authority: ['system:role:view'],
        },
        buttons: [
          { code: 'system:role:create', text: '新增' },
          { code: 'system:role:edit', text: '编辑' },
          { code: 'system:role:delete', text: '删除' },
          { code: 'system:role:copy', text: '复制' },
          { code: 'system:role:menu:assign', text: '配置菜单' },
        ],
      },
      {
        id: 23,
        pid: 2,
        name: 'Menus',
        path: '/menus/list',
        type: 'menu',
        component: '/menus/list',
        meta: {
          title: '菜单管理',
          icon: 'lucide:menu',
          authority: ['system:menu:view'],
        },
        buttons: [
          { code: 'system:menu:create', text: '新增' },
          { code: 'system:menu:edit', text: '编辑' },
          { code: 'system:menu:delete', text: '删除' },
          { code: 'system:menu:sort', text: '排序' },
        ],
      },
      {
        id: 24,
        pid: 2,
        name: 'Logs',
        path: '/logs/list',
        type: 'menu',
        component: '/system-logs/list',
        meta: {
          title: '日志管理',
          icon: 'lucide:file-text',
          authority: ['system:log:view'],
        },
        buttons: [
          { code: 'system:log:view', text: '查看' },
          { code: 'system:log:export', text: '导出' },
          { code: 'system:log:delete', text: '删除' },
          { code: 'system:log:clear', text: '清空' },
        ],
      },
      {
        id: 25,
        pid: 2,
        name: 'SiteConfig',
        path: '/site-config',
        type: 'menu',
        component: '/site-config/index',
        meta: {
          title: '网站配置',
          icon: 'lucide:settings',
          authority: ['system:config:view'],
        },
        buttons: [
          { code: 'system:config:view', text: '查看' },
          { code: 'system:config:edit', text: '编辑' },
        ],
      },
    ],
  },
];

// 信息管理菜单
const infoMenus: MenuComponent[] = [
  {
    id: 5,
    name: 'InfoManagement',
    path: '/info',
    type: 'catalog',
    meta: {
      title: '信息管理',
      icon: 'mdi:information',
      order: 5,
    },
    children: [
      {
        id: 51,
        pid: 5,
        name: 'InfoCategory',
        path: '/info/category/list',
        type: 'menu',
        component: '/info/category/list',
        meta: {
          title: '信息分类',
          icon: 'mdi:folder-multiple',
          authority: ['info:category:view'],
        },
        buttons: [
          { code: 'info:category:create', text: '新增' },
          { code: 'info:category:edit', text: '编辑' },
          { code: 'info:category:delete', text: '删除' },
          { code: 'info:category:sort', text: '排序' },
        ],
      },
      {
        id: 52,
        pid: 5,
        name: 'InfoList',
        path: '/info/list',
        type: 'menu',
        component: '/info/list',
        meta: {
          title: '信息列表',
          icon: 'mdi:file-document-multiple',
          authority: ['info:list:view'],
        },
        buttons: [
          { code: 'info:list:view', text: '查看' },
          { code: 'info:list:edit', text: '编辑' },
          { code: 'info:list:delete', text: '删除' },
          { code: 'info:list:export', text: '导出' },
        ],
      },
      {
        id: 53,
        pid: 5,
        name: 'InfoPublish',
        path: '/info/publish',
        type: 'menu',
        component: '/info/publish',
        meta: {
          title: '发布信息',
          icon: 'mdi:plus-circle',
          authority: ['info:publish:view'],
        },
        buttons: [
          { code: 'info:publish:create', text: '发布' },
          { code: 'info:publish:draft', text: '保存草稿' },
          { code: 'info:publish:schedule', text: '定时发布' },
        ],
      },
    ],
  },
];

// 站群管理菜单
const siteGroupMenus: MenuComponent[] = [
  {
    id: 7,
    name: 'SiteGroup',
    path: '/site-group',
    type: 'catalog',
    meta: {
      title: '站群管理',
      icon: 'lucide:network',
      order: 7,
    },
    children: [
      {
        id: 71,
        pid: 7,
        name: 'SiteGroupList',
        path: '/site-group/list',
        type: 'menu',
        component: '/site-group/list',
        meta: {
          title: '站群列表',
          icon: 'lucide:list',
          authority: ['sitegroup:view'],
        },
        buttons: [
          { code: 'sitegroup:view', text: '查看' },
          { code: 'sitegroup:create', text: '新增' },
          { code: 'sitegroup:edit', text: '编辑' },
          { code: 'sitegroup:delete', text: '删除' },
        ],
      },
    ],
  },
];

// 表单管理菜单
const formManageMenus: MenuComponent[] = [
  {
    id: 8,
    name: 'FormManage',
    path: '/form-manage',
    type: 'catalog',
    meta: {
      title: '表单管理',
      icon: 'lucide:file-text',
      order: 8,
    },
    children: [
      {
        id: 81,
        pid: 8,
        name: 'FormManageList',
        path: '/form-manage/list',
        type: 'menu',
        component: '/form-manage/list',
        meta: {
          title: '表单列表',
          icon: 'lucide:list',
          authority: ['form:view'],
        },
        buttons: [
          { code: 'form:view', text: '查看' },
          { code: 'form:delete', text: '删除' },
          { code: 'form:export', text: '导出' },
        ],
      },
    ],
  },
];

// 创建不同角色的菜单
function createMenusByRole(role: string): MenuComponent[] {
  switch (role) {
    case 'super_admin':
      // 超级管理员 - 所有权限（系统管理、信息管理、站群管理、表单管理）
      return [...systemMenus, ...infoMenus, ...siteGroupMenus, ...formManageMenus];

    case 'admin':
      // 租户管理员 - 系统管理、信息管理、站群管理、表单管理
      return [...systemMenus, ...infoMenus, ...siteGroupMenus, ...formManageMenus];

    case 'member':
      // 普通会员 - 仅信息管理
      return [...infoMenus];

    default:
      // 默认显示信息管理
      return [...infoMenus];
  }
}

// 导出菜单基础数据
export { dashboardMenus, systemMenus, infoMenus, siteGroupMenus, formManageMenus };
export { createMenusByRole };
