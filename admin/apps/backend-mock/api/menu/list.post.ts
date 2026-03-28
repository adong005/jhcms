import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  // 获取请求参数
  const body = await readBody(event);
  const { page = 1, pageSize = 100, name, type, status } = body;

  // 生成菜单数据（树形结构，使用 null 表示根节点）
  const generateMenus = () => {
    const menus = [
      // ========== 工作台 ==========
      {
        id: '1',
        name: '工作台',
        path: '/analytics',
        type: 'menu',
        icon: 'mdi:view-dashboard',
        component: '/dashboard/analytics/index',
        parentId: null,
        order: 1,
        status: 1,
        createTime: '2024-01-01 10:00:00',
        updateTime: '2024-03-20 14:00:00',
      },

      // ========== 用户管理（一级菜单）==========
      {
        id: '2',
        name: '用户管理',
        path: '/users/list',
        type: 'menu',
        icon: 'lucide:users',
        component: '/users/list',
        parentId: null,
        order: 2,
        status: 1,
        createTime: '2024-01-02 10:00:00',
        updateTime: '2024-03-21 14:00:00',
      },

      // ========== 角色管理（一级菜单）==========
      {
        id: '3',
        name: '角色管理',
        path: '/roles/list',
        type: 'menu',
        icon: 'lucide:shield',
        component: '/roles/list',
        parentId: null,
        order: 3,
        status: 1,
        createTime: '2024-01-03 10:00:00',
        updateTime: '2024-03-22 14:00:00',
      },

      // ========== 菜单管理（一级菜单）==========
      {
        id: '4',
        name: '菜单管理',
        path: '/menus/list',
        type: 'menu',
        icon: 'lucide:menu',
        component: '/menus/list',
        parentId: null,
        order: 4,
        status: 1,
        createTime: '2024-01-04 10:00:00',
        updateTime: '2024-03-23 14:00:00',
      },

      // ========== 网站配置（一级菜单）==========
      {
        id: '25',
        name: '网站配置',
        path: '/site-config',
        type: 'menu',
        icon: 'lucide:settings',
        component: '/site-config/index',
        parentId: null,
        order: 5,
        status: 1,
        createTime: '2024-03-25 12:00:00',
        updateTime: '2024-03-25 12:00:00',
      },

      // ========== 信息管理（一级菜单）==========
      {
        id: '5',
        name: '信息管理',
        path: '/info',
        type: 'catalog',
        icon: 'mdi:information',
        component: '',
        parentId: null,
        order: 5,
        status: 1,
        createTime: '2024-01-05 10:00:00',
        updateTime: '2024-03-24 14:00:00',
      },
      {
        id: '51',
        name: '信息分类',
        path: '/info/category/list',
        type: 'menu',
        icon: 'mdi:folder-multiple',
        component: '/info/category/list',
        parentId: '5',
        order: 1,
        status: 1,
        createTime: '2024-01-05 10:10:00',
        updateTime: '2024-03-24 14:10:00',
      },
      {
        id: '52',
        name: '信息列表',
        path: '/info/list',
        type: 'menu',
        icon: 'mdi:file-document-multiple',
        component: '/info/list',
        parentId: '5',
        order: 2,
        status: 1,
        createTime: '2024-01-05 11:00:00',
        updateTime: '2024-03-24 15:00:00',
      },
      {
        id: '53',
        name: '发布信息',
        path: '/info/publish',
        type: 'menu',
        icon: 'mdi:plus-circle',
        component: '/info/publish',
        parentId: '5',
        order: 3,
        status: 1,
        createTime: '2024-01-05 12:00:00',
        updateTime: '2024-03-24 16:00:00',
      },

      // ========== 站群管理（一级菜单）==========
      {
        id: '7',
        name: '站群管理',
        path: '/site-group',
        type: 'catalog',
        icon: 'lucide:network',
        component: '',
        parentId: null,
        order: 7,
        status: 1,
        createTime: '2024-03-25 12:00:00',
        updateTime: '2024-03-25 12:00:00',
      },
      {
        id: '71',
        name: '站群列表',
        path: '/site-group/list',
        type: 'menu',
        icon: 'lucide:list',
        component: '/site-group/list',
        parentId: '7',
        order: 1,
        status: 1,
        createTime: '2024-03-25 12:10:00',
        updateTime: '2024-03-25 12:10:00',
      },

      // ========== 表单管理（一级菜单）==========
      {
        id: '8',
        name: '表单管理',
        path: '/form-manage',
        type: 'catalog',
        icon: 'lucide:file-text',
        component: '',
        parentId: null,
        order: 8,
        status: 1,
        createTime: '2024-03-25 12:20:00',
        updateTime: '2024-03-25 12:20:00',
      },
      {
        id: '81',
        name: '表单列表',
        path: '/form-manage/list',
        type: 'menu',
        icon: 'lucide:list',
        component: '/form-manage/list',
        parentId: '8',
        order: 1,
        status: 1,
        createTime: '2024-03-25 12:30:00',
        updateTime: '2024-03-25 12:30:00',
      },
    ];

    return menus;
  };

  const allMenus = generateMenus();

  // 过滤数据
  let filteredData = allMenus;

  if (name) {
    filteredData = filteredData.filter(item =>
      item.name.includes(name)
    );
  }

  if (type) {
    filteredData = filteredData.filter(item =>
      item.type === type
    );
  }

  if (status !== undefined && status !== '') {
    filteredData = filteredData.filter(item =>
      item.status === Number(status)
    );
  }

  // 分页
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = filteredData.slice(start, end);

  return useResponseSuccess({
    items,
    total: filteredData.length,
  });
});
