import {
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from 'vue-router';

import { resetStaticRoutes } from '@vben/utils';
import { useAccessStore } from '@vben/stores';
import { LOGIN_PATH } from '@vben/constants';

import { createRouterGuard } from './guard';
import { routes } from './routes';

/**
 *  @zh_CN 创建vue-router实例
 */
const router = createRouter({
  history:
    import.meta.env.VITE_ROUTER_HISTORY === 'hash'
      ? createWebHashHistory(import.meta.env.VITE_BASE)
      : createWebHistory(import.meta.env.VITE_BASE),
  // 应该添加到路由的初始路由列表。
  routes,
  scrollBehavior: (to, _from, savedPosition) => {
    if (savedPosition) {
      return savedPosition;
    }
    return to.hash ? { behavior: 'smooth', el: to.hash } : { left: 0, top: 0 };
  },
  // 是否应该禁止尾部斜杠。
  // strict: true,
});

const resetRoutes = () => resetStaticRoutes(router, routes);

// 创建路由守卫
createRouterGuard(router);

// 路由运行时兜底：出现未捕获异常时回到登录页，避免启动阶段卡死
router.onError((error) => {
  console.error('[router] unexpected navigation error', error);
  const accessStore = useAccessStore();
  accessStore.setAccessToken(null);
  accessStore.setIsAccessChecked(false);
  if (router.currentRoute.value.path !== LOGIN_PATH) {
    router.replace({ path: LOGIN_PATH, replace: true }).catch(() => {});
  }
});

export { resetRoutes, router };
