<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';

import { ProfilePasswordSetting } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { getEmailSettingApi, updateEmailSettingApi } from '#/api/user';

const currentEmail = ref('');
const profileRef = ref();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      fieldName: 'currentEmail',
      label: '当前邮箱',
      component: 'Input',
      componentProps: {
        disabled: true,
        placeholder: currentEmail.value || '未绑定',
      },
    },
    {
      fieldName: 'newEmail',
      label: '新邮箱',
      component: 'Input',
      componentProps: {
        placeholder: '请输入新邮箱地址',
      },
      rules: 'required|email',
    },
    {
      fieldName: 'verifyCode',
      label: '验证码',
      component: 'Input',
      componentProps: {
        placeholder: '请输入邮箱验证码',
      },
      rules: 'required',
    },
  ];
});

onMounted(async () => {
  try {
    const data = await getEmailSettingApi();
    currentEmail.value = data.email || '';
  } catch (error) {
    console.error('获取邮箱失败:', error);
  }
});

async function handleSubmit(values: any) {
  try {
    await updateEmailSettingApi({
      newEmail: values.newEmail,
      verifyCode: values.verifyCode,
    });

    currentEmail.value = values.newEmail;
    profileRef.value?.getFormApi?.()?.resetForm();
    message.success('邮箱更新成功');
  } catch (error) {
    message.error('邮箱更新失败');
    console.error('更新邮箱失败:', error);
  }
}
</script>

<template>
  <ProfilePasswordSetting ref="profileRef" class="w-2/3" :form-schema="formSchema" @submit="handleSubmit" />
</template>
