<script setup lang="ts">
import type { BasicOption } from '@vben/types';

import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';

import { message } from 'ant-design-vue';
import { ProfileBaseSetting } from '@vben/common-ui';
import { useUserStore } from '@vben/stores';

import { getUserInfoApi } from '#/api';
import { updateProfileApi } from '#/api/user';

const profileBaseSettingRef = ref();

const MOCK_ROLES_OPTIONS: BasicOption[] = [
  {
    label: '超级管理员',
    value: 'super_admin',
  },
  {
    label: '超级管理员',
    value: 'super',
  },
  {
    label: '管理员',
    value: 'admin',
  },
  {
    label: '用户',
    value: 'user',
  },
  {
    label: '测试',
    value: 'test',
  },
];

// 将角色代码转换为中文标签
function formatRoles(roles: string[] | string): string {
  if (!roles) return '';

  const rolesArray = Array.isArray(roles) ? roles : [roles];
  const roleLabels = rolesArray.map(role => {
    const option = MOCK_ROLES_OPTIONS.find(opt => opt.value === role);
    // 如果找不到映射，将下划线转为空格并首字母大写
    if (!option) {
      return role.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
    }
    return option.label;
  });

  return roleLabels.join('、');
}

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      fieldName: 'nickName',
      component: 'Input',
      label: '昵称',
    },
    {
      fieldName: 'username',
      component: 'Input',
      componentProps: {
        disabled: true,
      },
      label: '用户名',
    },
    {
      fieldName: 'rolesDisplay',
      component: 'Input',
      componentProps: {
        disabled: true,
        placeholder: '暂无角色',
      },
      label: '角色',
    },
  ];
});

const userStore = useUserStore();

onMounted(async () => {
  const data = await getUserInfoApi();
  // 格式化角色显示
  const formattedData = {
    ...data,
    rolesDisplay: formatRoles(data.roles || []),
  };
  profileBaseSettingRef.value.getFormApi().setValues(formattedData);
});

// 处理表单提交
async function handleSubmit(values: any) {
  try {
    // 只提交可编辑的字段
    const { nickName } = values;

    await updateProfileApi({
      nickName,
    });

    // 重新获取用户信息并更新到 store
    const updatedUserInfo = await getUserInfoApi();
    userStore.setUserInfo(updatedUserInfo);

    message.success('个人信息更新成功');
  } catch (error) {
    message.error('个人信息更新失败');
    console.error('更新个人信息失败:', error);
  }
}
</script>
<template>
  <ProfileBaseSetting ref="profileBaseSettingRef" :form-schema="formSchema" @submit="handleSubmit" />
</template>
