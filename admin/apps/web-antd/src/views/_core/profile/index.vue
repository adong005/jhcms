<script setup lang="ts">
import { ref } from 'vue';

import { Profile } from '@vben/common-ui';
import { useUserStore } from '@vben/stores';

import ProfileBase from './base-setting.vue';
import ProfileEmailSetting from './email-setting.vue';
import ProfileGoogleAuthSetting from './google-auth-setting.vue';
import ProfileNotificationSetting from './notification-setting.vue';
import ProfilePasswordSetting from './password-setting.vue';
import ProfilePhoneSetting from './phone-setting.vue';
import ProfileQuestionSetting from './question-setting.vue';

const userStore = useUserStore();

const tabsValue = ref<string>('basic');

const tabs = ref([
  {
    label: '基本设置',
    value: 'basic',
  },
  {
    label: '修改密码',
    value: 'password',
  },
  {
    label: '密保手机',
    value: 'phone',
  },
  {
    label: '密保问题',
    value: 'question',
  },
  {
    label: '联系邮箱',
    value: 'email',
  },
  {
    label: '绑定谷歌验证器',
    value: 'google-auth',
  },
  {
    label: '新消息提醒',
    value: 'notice',
  },
]);
</script>
<template>
  <Profile v-model:model-value="tabsValue" title="个人中心" :user-info="userStore.userInfo" :tabs="tabs">
    <template #content>
      <ProfileBase v-if="tabsValue === 'basic'" />
      <ProfilePasswordSetting v-if="tabsValue === 'password'" />
      <ProfilePhoneSetting v-if="tabsValue === 'phone'" />
      <ProfileQuestionSetting v-if="tabsValue === 'question'" />
      <ProfileEmailSetting v-if="tabsValue === 'email'" />
      <ProfileGoogleAuthSetting v-if="tabsValue === 'google-auth'" />
      <ProfileNotificationSetting v-if="tabsValue === 'notice'" />
    </template>
  </Profile>
</template>
