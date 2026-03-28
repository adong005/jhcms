<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, reactive } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';

import { getQuestionSettingApi, updateQuestionSettingApi } from '#/api/user';

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      fieldName: 'question1',
      label: '密保问题1',
      component: 'Select',
      componentProps: {
        placeholder: '请选择密保问题',
        options: [
          { label: '您的出生地是？', value: '1' },
          { label: '您的母亲姓名是？', value: '2' },
          { label: '您的小学名称是？', value: '3' },
        ],
      },
      rules: 'required',
    },
    {
      fieldName: 'answer1',
      label: '答案1',
      component: 'Input',
      componentProps: {
        placeholder: '请输入答案',
      },
      rules: 'required',
    },
    {
      fieldName: 'question2',
      label: '密保问题2',
      component: 'Select',
      componentProps: {
        placeholder: '请选择密保问题',
        options: [
          { label: '您的父亲姓名是？', value: '4' },
          { label: '您的宠物名字是？', value: '5' },
          { label: '您最喜欢的颜色是？', value: '6' },
        ],
      },
      rules: 'required',
    },
    {
      fieldName: 'answer2',
      label: '答案2',
      component: 'Input',
      componentProps: {
        placeholder: '请输入答案',
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
    const data = await getQuestionSettingApi();
    if (data.question1) {
      formApi.setValues(data);
    }
  } catch (error) {
    console.error('获取密保问题失败:', error);
  }
});

async function handleSubmit() {
  try {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues();
    await updateQuestionSettingApi({
      question1: values.question1,
      answer1: values.answer1,
      question2: values.question2,
      answer2: values.answer2,
    });

    message.success('密保问题设置成功');
  } catch (error) {
    message.error('密保问题设置失败');
    console.error('更新密保问题失败:', error);
  }
}
</script>

<template>
  <div class="w-2/3">
    <Form />
    <button type="button"
      class="mt-4 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      @click="handleSubmit">
      设置密保问题
    </button>
  </div>
</template>
