export interface UserInfo {
  id: number;
  password: string;
  realName: string;
  roles: string[];
  username: string;
  homePath?: string;
}

export interface TimezoneOption {
  offset: number;
  timezone: string;
}

import { faker } from '@faker-js/faker';
import {
  dashboardMenus,
  systemMenus,
  infoMenus,
  siteGroupMenus,
  formManageMenus,
  createMenusByRole
} from './menu-data';

export const MOCK_USERS = [
  {
    accessToken: 'accessToken-vben',
    refreshToken: 'refreshToken-vben',
    username: 'vben',
    password: '123456',
    roles: ['super_admin'],
    userId: '100000',
    id: 'user_super_001',
    parentId: null,
    role: 'super_admin',
    realName: '系统超级管理员',
  },
  {
    accessToken: 'accessToken-admin',
    refreshToken: 'refreshToken-admin',
    username: 'admin',
    password: '123456',
    roles: ['admin'],
    userId: '100001',
    id: 'user_admin_001',
    parentId: null,
    role: 'admin',
    realName: '租户管理员',
    tenantName: '测试租户A',
  },
  {
    accessToken: 'accessToken-member1',
    refreshToken: 'refreshToken-member1',
    username: 'member1',
    password: '123456',
    roles: ['member'],
    userId: '100002',
    id: 'user_member_001',
    parentId: 'user_admin_001',
    role: 'member',
    realName: '普通会员',
  },
];

export const MOCK_CODES = [
  // super_admin - 系统级所有权限
  {
    codes: [
      'system:admin:create',
      'system:admin:delete',
      'system:tenant:manage',
      'system:*:*',
      'AC_100100',
      'AC_100110',
      'AC_100120',
      'AC_100010',
    ],
    username: 'vben',
  },
  // admin - 租户管理权限
  {
    codes: [
      'user:view',
      'user:create',
      'user:edit',
      'user:delete',
      'role:view',
      'role:create',
      'role:manage',
      'member:view',
      'member:create',
      'member:manage',
      'menu:view',
      'menu:manage',
      'AC_100010',
      'AC_100020',
      'AC_100030',
    ],
    username: 'admin',
  },
  // member - 基础权限
  {
    codes: [
      'profile:view',
      'profile:edit',
      'AC_1000001',
      'AC_1000002',
    ],
    username: 'member1',
  },
];

// 定义并导出菜单数据
export const MOCK_MENUS = [
  {
    username: 'vben',
    roles: ['super_admin'],
    menus: [...dashboardMenus, ...createMenusByRole('super_admin')],
  },
  {
    username: 'admin',
    roles: ['admin'],
    menus: [...dashboardMenus, ...createMenusByRole('admin')],
  },
  {
    username: 'member1',
    roles: ['member'],
    menus: [...dashboardMenus, ...createMenusByRole('member')],
  },
];

// 导出菜单列表（用于菜单管理）
export const MOCK_MENU_LIST = [
  ...dashboardMenus,
  ...systemMenus,
  ...infoMenus,
  ...siteGroupMenus,
  ...formManageMenus,
];

export function getMenuIds(menus: any[]) {
  const ids: number[] = [];
  menus.forEach((item) => {
    ids.push(item.id);
    if (item.children && item.children.length > 0) {
      ids.push(...getMenuIds(item.children));
    }
  });
  return ids;
}

/**
 * 时区选项
 */
export const TIME_ZONE_OPTIONS: TimezoneOption[] = [
  {
    offset: -5,
    timezone: 'America/New_York',
  },
  {
    offset: 0,
    timezone: 'Europe/London',
  },
  {
    offset: 8,
    timezone: 'Asia/Shanghai',
  },
  {
    offset: 9,
    timezone: 'Asia/Tokyo',
  },
  {
    offset: 9,
    timezone: 'Asia/Seoul',
  },
];
