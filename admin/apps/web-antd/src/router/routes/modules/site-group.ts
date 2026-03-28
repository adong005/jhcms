import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:network',
      order: 7,
      title: '站群管理',
    },
    name: 'SiteGroup',
    path: '/site-group',
    children: [
      {
        name: 'SiteGroupList',
        path: '/site-group/list',
        component: () => import('#/views/site-group/list.vue'),
        meta: {
          icon: 'lucide:list',
          title: '站群列表',
          authority: ['sitegroup:view'],
        },
      },
    ],
  },
];

export default routes;
