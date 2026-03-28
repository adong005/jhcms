import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:menu',
      order: 4,
      title: $t('page.menus.title'),
    },
    name: 'Menus',
    path: '/menus',
    redirect: '/menus/list',
    children: [
      {
        name: 'MenuList',
        path: '/menus/list',
        component: () => import('#/views/menus/list.vue'),
        meta: {
          title: $t('page.menus.title'),
        },
      },
    ],
  },
];

export default routes;
