import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:file-text',
      order: 8,
      title: '表单管理',
    },
    name: 'FormManage',
    path: '/form-manage',
    children: [
      {
        name: 'FormManageList',
        path: '/form-manage/list',
        component: () => import('#/views/form-manage/list.vue'),
        meta: {
          icon: 'lucide:list',
          title: '表单列表',
          authority: ['form:view'],
        },
      },
    ],
  },
];

export default routes;
