import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:users',
      order: 2,
      title: $t('page.users.title'),
    },
    name: 'Users',
    path: '/users',
    redirect: '/users/list',
    children: [
      {
        name: 'UserList',
        path: '/users/list',
        component: () => import('#/views/users/list.vue'),
        meta: {
          title: $t('page.users.title'),
        },
      },
    ],
  },
];

export default routes;
