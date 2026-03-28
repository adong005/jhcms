<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Button, Card, Image, Upload, message } from 'ant-design-vue';
import { useVbenForm } from '#/adapter/form';
import { useUserStore } from '@vben/stores';
import { getSiteConfigApi, updateSiteConfigApi, uploadSiteLogoApi } from '#/api/site-config';

const userStore = useUserStore();
const submitLoading = ref(false);
const logoUploading = ref(false);
const logoUrl = ref('');

function resolveLogoUrl(url: string) {
  if (!url) return '';
  if (/^https?:\/\//i.test(url)) return url;
  const apiBase = (import.meta.env.VITE_GLOB_API_URL || '').replace(/\/+$/, '');
  const serverBase = apiBase.endsWith('/api') ? apiBase.slice(0, -4) : apiBase;
  if (!serverBase) return url;
  return `${serverBase}${url.startsWith('/') ? '' : '/'}${url}`;
}

const [Form, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  layout: 'vertical',
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入网站标题',
      },
      fieldName: 'title',
      label: '网站标题',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入网站关键词，多个关键词用逗号分隔',
      },
      fieldName: 'keywords',
      label: '网站关键词',
      rules: 'required',
      help: '多个关键词请用英文逗号分隔',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入网站描述',
        rows: 4,
      },
      fieldName: 'description',
      label: '网站描述',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入网站域名，例如：www.example.com',
      },
      fieldName: 'domain',
      label: '网站域名',
      rules: 'required',
      help: '请输入完整的域名，不包含 http:// 或 https://',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入备案码，例如：京ICP备12345678号',
      },
      fieldName: 'icpCode',
      label: '备案码',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入联系方式',
      },
      fieldName: 'contactPhone',
      label: '联系方式',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入联系地址',
      },
      fieldName: 'contactAddress',
      label: '联系地址',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入联系邮箱',
      },
      fieldName: 'contactEmail',
      label: '联系邮箱',
    },
  ],
  showDefaultActions: false,
});

// 加载网站配置
onMounted(async () => {
  try {
    const data = await getSiteConfigApi();
    logoUrl.value = data.logo || '';
    formApi.setValues({
      title: data.title,
      keywords: data.keywords,
      description: data.description,
      domain: data.domain,
      logo: data.logo,
      icpCode: data.icpCode,
      contactPhone: data.contactPhone,
      contactAddress: data.contactAddress,
      contactEmail: data.contactEmail,
    });
  } catch (error) {
    console.error('获取网站配置失败:', error);
  }
});

// 提交表单
async function handleSubmit() {
  try {
    await formApi.validate();
    submitLoading.value = true;

    const formValues = await formApi.getValues();

    // 提交数据时带上用户 ID
    const submitData = {
      userId: userStore.userInfo?.userId || '',
      title: formValues.title,
      keywords: formValues.keywords,
      description: formValues.description,
      domain: formValues.domain,
      logo: logoUrl.value || '',
      icpCode: formValues.icpCode || '',
      contactPhone: formValues.contactPhone || '',
      contactAddress: formValues.contactAddress || '',
      contactEmail: formValues.contactEmail || '',
    };

    await updateSiteConfigApi(submitData);

    message.success('网站配置更新成功');
  } catch (error) {
    message.error('网站配置更新失败');
    console.error('更新失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

async function handleLogoUpload(file: File) {
  try {
    logoUploading.value = true;
    const resp = await uploadSiteLogoApi(file);
    logoUrl.value = resp.logo || '';
    formApi.setValues({ logo: logoUrl.value });
    message.success('Logo上传成功');
  } catch (error) {
    message.error('Logo上传失败');
  } finally {
    logoUploading.value = false;
  }
  return false;
}

// 重置表单
function handleReset() {
  formApi.resetForm();
}
</script>

<template>
  <div class="p-4">
    <Card title="网站配置" :bordered="false">
      <template #extra>
        <span class="text-gray-500">配置网站的基本信息</span>
      </template>

      <div class="max-w-3xl">
        <div class="mb-4">
          <div class="mb-2 text-sm text-gray-600">网站Logo</div>
          <div class="flex items-center gap-3">
            <Upload
              :show-upload-list="false"
              accept="image/*"
              :before-upload="handleLogoUpload"
            >
              <Button :loading="logoUploading">上传Logo</Button>
            </Upload>
            <Image
              v-if="logoUrl"
              :src="resolveLogoUrl(logoUrl)"
              :width="72"
              :height="72"
              style="object-fit: contain; border: 1px solid #f0f0f0; border-radius: 8px; padding: 4px;"
            />
            <span v-else class="text-gray-400 text-sm">未上传</span>
          </div>
        </div>

        <Form />

        <div class="mt-6 flex gap-2">
          <Button type="primary" :loading="submitLoading" @click="handleSubmit">
            <template #icon><span class="i-ant-design:save-outlined" /></template>
            保存配置
          </Button>
          <Button @click="handleReset">
            <template #icon><span class="i-ant-design:redo-outlined" /></template>
            重置
          </Button>
        </div>
      </div>
    </Card>
  </div>
</template>
