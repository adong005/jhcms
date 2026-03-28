<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Button, message } from 'ant-design-vue';

import { bindGoogleAuthApi, getGoogleAuthSettingApi, unbindGoogleAuthApi } from '#/api/user';

const qrCodeUrl = ref('');
const secretKey = ref('');
const verifyCode = ref('');
const isBound = ref(false);
const loading = ref(false);

onMounted(async () => {
  try {
    const data = await getGoogleAuthSettingApi();
    isBound.value = data.isBound;
    qrCodeUrl.value = data.qrCodeUrl || '';
    secretKey.value = data.secretKey || '';
  } catch (error) {
    console.error('获取谷歌验证器设置失败:', error);
  }
});

async function handleBind() {
  if (!verifyCode.value) {
    message.error('请输入验证码');
    return;
  }

  loading.value = true;
  try {
    await bindGoogleAuthApi({ verifyCode: verifyCode.value });
    isBound.value = true;
    verifyCode.value = '';
    message.success('谷歌验证器绑定成功');
  } catch (error) {
    message.error('绑定失败，请检查验证码是否正确');
    console.error('绑定谷歌验证器失败:', error);
  } finally {
    loading.value = false;
  }
}

async function handleUnbind() {
  loading.value = true;
  try {
    await unbindGoogleAuthApi();
    isBound.value = false;
    verifyCode.value = '';
    message.success('谷歌验证器解绑成功');
  } catch (error) {
    message.error('解绑失败');
    console.error('解绑谷歌验证器失败:', error);
  } finally {
    loading.value = false;
  }
}

function copySecretKey() {
  if (window.navigator && window.navigator.clipboard) {
    window.navigator.clipboard.writeText(secretKey.value);
    message.success('密钥已复制到剪贴板');
  }
}
</script>

<template>
  <div class="w-2/3">
    <div v-if="!isBound" class="space-y-6">
      <div class="rounded-lg border p-4">
        <h3 class="mb-4 text-base font-medium">第一步：扫描二维码</h3>
        <div class="flex flex-col items-center space-y-4">
          <img :src="qrCodeUrl" alt="QR Code" class="h-48 w-48" />
          <p class="text-sm text-muted-foreground">
            使用 Google Authenticator 或其他验证器应用扫描此二维码
          </p>
        </div>
      </div>

      <div class="rounded-lg border p-4">
        <h3 class="mb-4 text-base font-medium">第二步：手动输入密钥（可选）</h3>
        <div class="space-y-2">
          <p class="text-sm text-muted-foreground">
            如果无法扫描二维码，请手动输入以下密钥：
          </p>
          <div class="flex items-center space-x-2">
            <code class="rounded bg-muted px-3 py-2 text-sm">{{ secretKey }}</code>
            <Button size="small" @click="copySecretKey">
              复制
            </Button>
          </div>
        </div>
      </div>

      <div class="rounded-lg border p-4">
        <h3 class="mb-4 text-base font-medium">第三步：输入验证码</h3>
        <div class="space-y-4">
          <div>
            <label class="mb-2 block text-sm font-medium">验证码</label>
            <input v-model="verifyCode" type="text" placeholder="请输入6位验证码" maxlength="6"
              class="w-full rounded-md border px-3 py-2" />
          </div>
          <Button type="primary" :loading="loading" @click="handleBind">
            绑定谷歌验证器
          </Button>
        </div>
      </div>
    </div>

    <div v-else class="space-y-4">
      <div class="rounded-lg border border-green-200 bg-green-50 p-4">
        <div class="flex items-center space-x-2">
          <span class="text-green-600">✓</span>
          <span class="font-medium text-green-800">已绑定谷歌验证器</span>
        </div>
        <p class="mt-2 text-sm text-green-700">
          您的账户已启用谷歌验证器二次验证，登录时需要输入验证码
        </p>
      </div>
      <Button danger :loading="loading" @click="handleUnbind">
        解绑谷歌验证器
      </Button>
    </div>
  </div>
</template>
