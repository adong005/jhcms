import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { notification } from 'ant-design-vue';
import { defineStore } from 'pinia';

import {
  getAccessCodesApi,
  getUserInfoApi,
  impersonateApi,
  loginApi,
  logoutApi,
} from '#/api';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 异步处理登录操作
   * Asynchronously handle the login process
   * @param params 登录表单数据
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    // 异步处理用户登录操作并获取 accessToken
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const { accessToken } = await loginApi(params);

      // 如果成功获取到 accessToken
      if (accessToken) {
        accessStore.setAccessToken(accessToken);

        // 获取用户信息并存储到 accessStore 中
        const [fetchUserInfoResult, accessCodes] = await Promise.all([
          fetchUserInfo(),
          getAccessCodesApi(),
        ]);

        userInfo = fetchUserInfoResult;

        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);

        if (accessStore.loginExpired) {
          accessStore.setLoginExpired(false);
        } else {
          onSuccess
            ? await onSuccess?.()
            : await router.push(
                userInfo.homePath || preferences.app.defaultHomePath,
              );
        }

        const displayName =
          (userInfo as Record<string, unknown> | null)?.nickName ||
          userInfo?.realName ||
          userInfo?.username;
        if (displayName) {
          notification.success({
            description: `${$t('authentication.loginSuccessDesc')}:${displayName}`,
            duration: 3,
            message: $t('authentication.loginSuccess'),
          });
        }
      }
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  async function logout(redirect: boolean = true) {
    try {
      await logoutApi();
    } catch {
      // 不做任何处理
    }
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登录页带上当前路由地址
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  async function fetchUserInfo() {
    const userInfo = await getUserInfoApi();
    userStore.setUserInfo(userInfo);
    return userInfo;
  }

  /**
   * 刷新当前租户下的用户与权限上下文。
   * 适用于租户切换后重建菜单/按钮权限。
   */
  async function refreshAccessContext() {
    const [userInfo, accessCodes] = await Promise.all([
      fetchUserInfo(),
      getAccessCodesApi(),
    ]);
    userStore.setUserInfo(userInfo);
    accessStore.setAccessCodes(accessCodes);
    accessStore.setIsAccessChecked(false);
    return { accessCodes, userInfo };
  }

  /**
   * 平台超级管理员：切换到目标用户身份并刷新菜单权限（进入对方后台）
   */
  async function enterUserBackend(userId: string) {
    const { accessToken } = await impersonateApi({ userId });
    if (!accessToken) {
      throw new Error('未获取到访问令牌');
    }
    accessStore.setAccessToken(accessToken);
    if (accessStore.loginExpired) {
      accessStore.setLoginExpired(false);
    }
    const [userInfo, accessCodes] = await Promise.all([
      fetchUserInfo(),
      getAccessCodesApi(),
    ]);
    userStore.setUserInfo(userInfo);
    accessStore.setAccessCodes(accessCodes);
    accessStore.setIsAccessChecked(false);
    await router.replace(
      userInfo.homePath || preferences.app.defaultHomePath,
    );
    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    refreshAccessContext,
    enterUserBackend,
    loginLoading,
    logout,
  };
});
