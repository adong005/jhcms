import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:shield',
      order: 3,
      title: $t('page.roles.title'),
    },
    name: 'Roles',
    path: '/roles',
    redirect: '/roles/list',
    children: [
      {
        name: 'RoleList',
        path: '/roles/list',
        component: () => import('#/views/roles/list.vue'),
        meta: {
          title: $t('page.roles.title'),
        },
      },
    ],
  },
];

export default routes;
