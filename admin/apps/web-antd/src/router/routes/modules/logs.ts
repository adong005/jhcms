import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:file-text',
      order: 4,
      title: '日志管理',
    },
    name: 'Logs',
    path: '/logs',
    children: [
      {
        name: 'LogList',
        path: '/logs/list',
        component: () => import('#/views/system-logs/list.vue'),
        meta: {
          icon: 'lucide:file-text',
          title: '日志管理',
          authority: ['system:log:view'],
        },
      },
    ],
  },
];

export default routes;
