import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:settings',
      order: 5,
      title: '网站配置',
    },
    name: 'SiteConfig',
    path: '/site-config',
    children: [
      {
        name: 'SiteConfigIndex',
        path: '/site-config',
        component: () => import('#/views/site-config/index.vue'),
        meta: {
          icon: 'lucide:settings',
          title: '网站配置',
          authority: ['system:config:view'],
        },
      },
    ],
  },
];

export default routes;
