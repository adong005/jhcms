<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Button, message } from 'ant-design-vue';

import { getSecuritySettingsApi } from '#/api/user';

interface SecuritySettings {
  accountPassword: boolean;
  securityPhone: boolean;
  securityPhoneNumber?: string;
  securityQuestion: boolean;
  securityEmail: boolean;
  securityEmailAddress?: string;
  securityMfa: boolean;
  passwordStrength?: string;
}

const securitySettings = ref<SecuritySettings>({
  accountPassword: true,
  securityPhone: false,
  securityQuestion: false,
  securityEmail: false,
  securityMfa: false,
});

const securityItems = [
  {
    key: 'accountPassword',
    label: '账户密码',
    getDescription: (settings: SecuritySettings) =>
      `当前密码强度：${settings.passwordStrength || '中'}`,
  },
  {
    key: 'securityPhone',
    label: '密保手机',
    getDescription: (settings: SecuritySettings) =>
      settings.securityPhoneNumber
        ? `已绑定手机：${settings.securityPhoneNumber}`
        : '未绑定手机，绑定后可用于账户找回',
  },
  {
    key: 'securityQuestion',
    label: '密保问题',
    getDescription: (settings: SecuritySettings) =>
      settings.securityQuestion
        ? '已设置密保问题'
        : '未设置密保问题，密保问题可有效保护账户安全',
  },
  {
    key: 'securityEmail',
    label: '备用邮箱',
    getDescription: (settings: SecuritySettings) =>
      settings.securityEmailAddress
        ? `已绑定邮箱：${settings.securityEmailAddress}`
        : '未绑定邮箱，绑定后可用于接收安全通知',
  },
  {
    key: 'securityMfa',
    label: 'MFA 设备',
    getDescription: (settings: SecuritySettings) =>
      settings.securityMfa
        ? '已绑定 MFA 设备，登录时需要二次验证'
        : '未绑定 MFA 设备，绑定后可以进行二次确认',
  },
];

onMounted(async () => {
  try {
    const data = await getSecuritySettingsApi();
    securitySettings.value = data;
  } catch (error) {
    console.error('获取安全设置失败:', error);
  }
});

async function handleModify(key: string) {
  message.info(`点击了 ${securityItems.find(item => item.key === key)?.label} 的更改按钮`);
  // 这里可以打开对应的修改弹窗或跳转到修改页面
}
</script>
<template>
  <div class="space-y-4">
    <div v-for="item in securityItems" :key="item.key" class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex-1 space-y-0.5">
        <div class="text-base font-medium">{{ item.label }}</div>
        <div class="text-sm text-muted-foreground">
          {{ item.getDescription(securitySettings) }}
        </div>
      </div>
      <Button type="link" @click="handleModify(item.key)">
        更改
      </Button>
    </div>
  </div>
</template>
