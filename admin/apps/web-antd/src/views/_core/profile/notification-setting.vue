<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { ProfileNotificationSetting } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import {
  getNotificationSettingsApi,
  updateNotificationSettingsApi,
} from '#/api/user';

const formSchema = ref([
  {
    value: true,
    fieldName: 'accountPassword',
    label: '账户密码',
    description: '其他用户的消息将以站内信的形式通知',
  },
  {
    value: true,
    fieldName: 'systemMessage',
    label: '系统消息',
    description: '系统消息将以站内信的形式通知',
  },
  {
    value: true,
    fieldName: 'todoTask',
    label: '待办任务',
    description: '待办任务将以站内信的形式通知',
  },
]);

onMounted(async () => {
  try {
    const settings = await getNotificationSettingsApi();
    formSchema.value = formSchema.value.map((item) => ({
      ...item,
      value: settings[item.fieldName as keyof typeof settings] ?? true,
    }));
  } catch (error) {
    console.error('获取通知设置失败:', error);
  }
});

async function handleSubmit(values: Record<string, boolean>) {
  await updateNotificationSettingsApi(values);
  message.success('通知设置已保存');
}
</script>
<template>
  <ProfileNotificationSetting :form-schema="formSchema" @submit="handleSubmit" />
</template>
