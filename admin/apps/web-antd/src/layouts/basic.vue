<script lang="ts" setup>
import type { NotificationItem } from '@vben/layouts';

import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationLoginExpiredModal } from '@vben/common-ui';
import { useWatermark } from '@vben/hooks';
import {
  BasicLayout,
  LockScreen,
  Notification,
  UserDropdown,
} from '@vben/layouts';
import { preferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';

import { $t } from '#/locales';
import { useAuthStore } from '#/store';
import LoginForm from '#/views/_core/authentication/login.vue';

const notifications = ref<NotificationItem[]>([
  {
    id: 1,
    avatar: 'https://avatar.vercel.sh/vercel.svg?text=VB',
    date: '3小时前',
    isRead: true,
    message: '描述信息描述信息描述信息',
    title: '收到了 14 份新周报',
  },
  {
    id: 2,
    avatar: 'https://avatar.vercel.sh/1',
    date: '刚刚',
    isRead: false,
    message: '描述信息描述信息描述信息',
    title: '朱偏右 回复了你',
  },
  {
    id: 3,
    avatar: 'https://avatar.vercel.sh/1',
    date: '2024-01-01',
    isRead: false,
    message: '描述信息描述信息描述信息',
    title: '曲丽丽 评论了你',
  },
  {
    id: 4,
    avatar: 'https://avatar.vercel.sh/satori',
    date: '1天前',
    isRead: false,
    message: '描述信息描述信息描述信息',
    title: '代办提醒',
  },
  {
    id: 5,
    avatar: 'https://avatar.vercel.sh/satori',
    date: '1天前',
    isRead: false,
    message: '描述信息描述信息描述信息',
    title: '跳转Workspace示例',
    link: '/workspace',
  },
  {
    id: 6,
    avatar: 'https://avatar.vercel.sh/satori',
    date: '1天前',
    isRead: false,
    message: '描述信息描述信息描述信息',
    title: '跳转外部链接示例',
    link: 'https://doc.vben.pro',
  },
]);

const router = useRouter();
const userStore = useUserStore();
const authStore = useAuthStore();
const accessStore = useAccessStore();
const { destroyWatermark, updateWatermark } = useWatermark();
const showDot = computed(() =>
  notifications.value.some((item) => !item.isRead),
);

const menus = computed(() => [
  {
    handler: () => {
      router.push('/profile');
    },
    icon: 'lucide:user',
    text: $t('page.auth.profile'),
  },
]);

/** 后端 /user/info 扩展字段（BasicUserInfo 未声明） */
type SessionUserInfo = {
  username?: string;
  realName?: string;
  nickName?: string;
  email?: string;
  role?: string;
  roles?: string[];
};

const sessionUser = computed(
  () => (userStore.userInfo ?? null) as SessionUserInfo | null,
);

const userDropdownText = computed(() => {
  const u = sessionUser.value;
  if (!u) return '';
  const name = (u.nickName || u.realName || u.username || '').trim();
  return name || u.username || '';
});

const userDropdownDescription = computed(() => {
  const u = sessionUser.value;
  if (!u) return '';
  const email = (u.email || '').trim();
  if (email) return email;
  const un = (u.username || '').trim();
  return un ? `@${un}` : '';
});

const roleLabelMap: Record<string, string> = {
  super_admin: '超级管理员',
  admin: '管理员',
  user: '用户',
};

const userDropdownTag = computed(() => {
  const u = sessionUser.value;
  const role = (u?.role || u?.roles?.[0] || '').trim();
  if (!role) return '';
  return roleLabelMap[role] || role;
});

const avatar = computed(() => {
  return userStore.userInfo?.avatar ?? preferences.app.defaultAvatar;
});

async function handleLogout() {
  await authStore.logout(false);
}

function handleNoticeClear() {
  notifications.value = [];
}

function markRead(id: number | string) {
  const item = notifications.value.find((item) => item.id === id);
  if (item) {
    item.isRead = true;
  }
}

function remove(id: number | string) {
  notifications.value = notifications.value.filter((item) => item.id !== id);
}

function handleMakeAll() {
  notifications.value.forEach((item) => (item.isRead = true));
}
watch(
  () => ({
    enable: preferences.app.watermark,
    content: preferences.app.watermarkContent,
  }),
  async ({ enable, content }) => {
    if (enable) {
      await updateWatermark({
        content:
          content ||
          (() => {
            const u = userStore.userInfo as SessionUserInfo | null;
            const label = (u?.nickName || u?.realName || u?.username || '').trim() || u?.username || '';
            return `${userStore.userInfo?.username ?? ''} - ${label}`;
          })(),
      });
    } else {
      destroyWatermark();
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <BasicLayout @clear-preferences-and-logout="handleLogout">
    <template #user-dropdown>
      <UserDropdown
        :avatar
        :menus
        :text="userDropdownText"
        :description="userDropdownDescription"
        :tag-text="userDropdownTag"
        @logout="handleLogout"
      />
    </template>
    <template #notification>
      <Notification :dot="showDot" :notifications="notifications" @clear="handleNoticeClear"
        @read="(item) => item.id && markRead(item.id)" @remove="(item) => item.id && remove(item.id)"
        @make-all="handleMakeAll" />
    </template>
    <template #extra>
      <AuthenticationLoginExpiredModal v-model:open="accessStore.loginExpired" :avatar>
        <LoginForm />
      </AuthenticationLoginExpiredModal>
    </template>
    <template #lock-screen>
      <LockScreen :avatar @to-login="handleLogout" />
    </template>
  </BasicLayout>
</template>
