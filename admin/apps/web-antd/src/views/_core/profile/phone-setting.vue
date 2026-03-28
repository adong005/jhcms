<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, reactive, ref } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';

import { getPhoneSettingApi, updatePhoneSettingApi } from '#/api/user';

const currentPhone = ref('');

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      fieldName: 'currentPhone',
      label: '当前手机号',
      component: 'Input',
      componentProps: {
        disabled: true,
        placeholder: currentPhone.value || '未绑定',
      },
    },
    {
      fieldName: 'newPhone',
      label: '新手机号',
      component: 'Input',
      componentProps: {
        placeholder: '请输入新手机号',
      },
      rules: 'required',
    },
    {
      fieldName: 'verifyCode',
      label: '验证码',
      component: 'Input',
      componentProps: {
        placeholder: '请输入验证码',
      },
      rules: 'required',
    },
  ];
});

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      labelWidth: 130,
      componentProps: {
        class: 'w-full',
      },
    },
    layout: 'horizontal',
    schema: formSchema,
    showDefaultActions: false,
  }),
);

onMounted(async () => {
  try {
    const data = await getPhoneSettingApi();
    currentPhone.value = data.phone || '';
  } catch (error) {
    console.error('获取手机号失败:', error);
  }
});

async function handleSubmit() {
  try {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues();
    await updatePhoneSettingApi({
      newPhone: values.newPhone,
      verifyCode: values.verifyCode,
    });

    currentPhone.value = values.newPhone;
    formApi.resetForm();
    message.success('手机号更新成功');
  } catch (error) {
    message.error('手机号更新失败');
    console.error('更新手机号失败:', error);
  }
}
</script>

<template>
  <div class="w-2/3">
    <Form />
    <button type="button"
      class="mt-4 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      @click="handleSubmit">
      更新手机号
    </button>
  </div>
</template>
