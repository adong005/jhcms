<script setup lang="ts">
import type { VxeGridListeners, VxeGridProps } from '#/adapter/vxe-table';
import { computed, onMounted, ref } from 'vue';
import { Button, message, Modal } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getSiteGroupListApi, createSiteGroupApi, updateSiteGroupApi, deleteSiteGroupApi, batchDeleteSiteGroupApi, getSiteGroupAdminsApi } from '#/api/site-group';
import type { SiteGroupApi } from '#/api/site-group';
import { getSiteConfigApi } from '#/api/site-config';
import { useUserStore } from '@vben/stores';

type SiteGroupRecord = SiteGroupApi.SiteGroup;
const userStore = useUserStore();

const isSuperAdmin = computed(() => {
  const u = userStore.userInfo as {
    roles?: string[];
    isPlatformSuperAdmin?: boolean;
    isAdmin?: number | boolean;
  } | null;
  if (!u) return false;
  if (u.isPlatformSuperAdmin) return true;
  if (u.isAdmin === true || u.isAdmin === 1) return true;
  return (u.roles ?? []).includes('super_admin');
});

const siteConfig = ref({ domain: '' });
const adminOptions = ref<Array<{ label: string; value: string }>>([]);

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentSiteGroupId.value = '';
    }
  },
});

const currentSiteGroupId = ref('');
const isEdit = computed(() => !!currentSiteGroupId.value);
const drawerTitle = computed(() => isEdit.value ? '编辑站群' : '新增站群');
const submitLoading = ref(false);

const [SiteGroupForm, formApi] = useVbenForm({
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
        placeholder: '请输入关键词',
      },
      fieldName: 'keyword',
      label: '关键词',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入二级域名',
      },
      fieldName: 'subdomain',
      label: '二级域名',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入网站标题',
      },
      fieldName: 'title',
      label: '标题',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入网站关键词',
      },
      fieldName: 'keywords',
      label: '关键词',
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入网站描述',
        rows: 4,
      },
      fieldName: 'description',
      label: '描述',
      rules: 'required',
    },
  ],
  showDefaultActions: false,
});

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入关键词',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入二级域名',
      },
      fieldName: 'subdomain',
      label: '二级域名',
    },
    ...(isSuperAdmin.value ? [{
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择管理员',
        options: adminOptions,
      },
      fieldName: 'adminId',
      label: '管理员',
    }] : []),
  ],
  showCollapseButton: true,
  submitButtonOptions: {
    content: '查询',
  },
};

const gridOptions: VxeGridProps<SiteGroupRecord> = {
  border: true,
  stripe: true,
  rowConfig: {
    keyField: 'id',
  },
  editConfig: {
    trigger: 'click',
    mode: 'cell',
    showStatus: true,
  },
  checkboxConfig: {
    highlight: true,
    reserve: true,
  },
  keepSource: true,
  columns: [
    { type: 'checkbox', width: 50, fixed: 'left' },
    { field: 'id', title: 'ID', width: 100 },
    ...(isSuperAdmin.value ? [{ field: 'adminName', title: '管理员', width: 140 }] : []),
    {
      field: 'keyword',
      title: '关键词',
      width: 150,
      editRender: {
        name: 'input',
      },
    },
    { 
      field: 'subdomain', 
      title: '二级域名', 
      width: 200,
      slots: { default: 'subdomain' },
    },
    {
      field: 'title',
      title: '标题',
      minWidth: 200,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'keywords',
      title: '关键词',
      minWidth: 200,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'description',
      title: '描述',
      minWidth: 250,
      editRender: {
        name: 'input',
      },
    },
    { field: 'createTime', title: '创建时间', width: 180 },
    {
      title: '操作',
      width: 150,
      slots: { default: 'action' },
      fixed: 'right',
    },
  ],
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const resp = await getSiteGroupListApi({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
        return {
          ...resp,
          items: (resp.items || []).map((item: any) => ({
            ...item,
            id: String(item.id),
          })),
        };
      },
    },
  },
  toolbarConfig: {
    custom: true,
    export: true,
    refresh: true,
    zoom: true,
    // @ts-ignore - 正式环境时有完整的类型声明
    search: true,
  },
};

const gridEvents: VxeGridListeners<SiteGroupRecord> = {
  // 单元格编辑完成后自动提交
  editClosed: async ({ row, column }) => {
    if (String(row.id || '').startsWith('virtual-')) {
      return;
    }
    // 只处理可编辑的字段
    const editableFields = ['keyword', 'title', 'keywords', 'description'];
    if (!editableFields.includes(column.field)) {
      return;
    }

    // 检查该行是否有更新（使用 keepSource 时可用）
    const $grid = gridApi.grid;
    if (!$grid.isUpdateByRow(row)) {
      return; // 数据未修改，不提交
    }

    try {
      const updateData: any = {};
      updateData[column.field] = row[column.field as keyof SiteGroupRecord];

      await updateSiteGroupApi({
        id: row.id,
        keyword: row.keyword,
        subdomain: row.subdomain,
        title: row.title,
        keywords: row.keywords,
        description: row.description,
      });
      message.success('更新成功');
      // 刷新列表
      gridApi.reload();
    } catch (error) {
      message.error('更新失败');
      // 刷新列表恢复原数据
      gridApi.reload();
    }
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents,
  gridOptions,
});

async function handleAdd() {
  currentSiteGroupId.value = '';
  drawerApi.open();
}

async function handleEdit(record: SiteGroupRecord) {
  if (String(record.id || '').startsWith('virtual-')) {
    message.warning('默认城市站群为运行时数据，请先新增后再编辑');
    return;
  }
  currentSiteGroupId.value = record.id;
  formApi.setValues({
    keyword: record.keyword,
    subdomain: record.subdomain,
    title: record.title,
    keywords: record.keywords,
    description: record.description,
  });
  drawerApi.open();
}

async function handleDelete(record: SiteGroupRecord) {
  if (String(record.id || '').startsWith('virtual-')) {
    message.warning('默认城市站群为运行时数据，无法删除');
    return;
  }
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除站群"${record.title}"吗？`,
    onOk: async () => {
      await deleteSiteGroupApi(record.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

async function handleBatchDelete() {
  const selectedRows = gridApi.grid.getCheckboxRecords();
  if (selectedRows.length === 0) {
    message.warning('请选择要删除的站群');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRows.length} 个站群吗？`,
    onOk: async () => {
      const ids = selectedRows
        .filter((row) => !String(row.id || '').startsWith('virtual-'))
        .map(row => row.id);
      if (ids.length === 0) {
        message.warning('默认城市站群为运行时数据，无法批量删除');
        return;
      }
      await batchDeleteSiteGroupApi(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
  });
}

async function handleSubmit() {
  try {
    await formApi.validate();
    submitLoading.value = true;

    const formValues = await formApi.getValues();

    if (isEdit.value) {
      await updateSiteGroupApi({
        id: currentSiteGroupId.value,
        keyword: formValues.keyword,
        subdomain: formValues.subdomain,
        title: formValues.title,
        keywords: formValues.keywords,
        description: formValues.description,
      });
      message.success('更新成功');
    } else {
      await createSiteGroupApi({
        keyword: formValues.keyword,
        subdomain: formValues.subdomain,
        title: formValues.title,
        keywords: formValues.keywords,
        description: formValues.description,
      });
      message.success('创建成功');
    }

    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    console.error('提交失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

// 加载网站配置
onMounted(async () => {
  try {
    const config = await getSiteConfigApi();
    siteConfig.value = config;
    if (isSuperAdmin.value) {
      const admins = await getSiteGroupAdminsApi();
      adminOptions.value = (admins || []).map((a) => ({
        label: a.nickName || a.username,
        value: a.userId,
      }));
    }
  } catch (error) {
    console.error('获取网站配置失败:', error);
  }
});

// 处理二级域名点击
function handleSubdomainClick(subdomain: string) {
  const raw = String(subdomain || '').trim();
  if (!raw) {
    return;
  }

  const domain = String(siteConfig.value.domain || '').trim().replace(/^https?:\/\//, '').replace(/\/+$/, '');

  // 如果 subdomain 已经是完整域名（包含配置的主域名），直接使用
  let host = raw;
  if (domain && !raw.endsWith(domain)) {
    host = `${raw}.${domain}`;
  }

  const fullUrl = /^https?:\/\//i.test(host) ? host : `http://${host}`;
  window.open(fullUrl, '_blank');
}
</script>

<template>
  <div class="p-4">
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          新增站群
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <template #subdomain="{ row }">
        <Button 
          type="link" 
          size="small" 
          @click="handleSubdomainClick(row.subdomain)"
          class="p-0"
        >
          {{ row.subdomain }}
        </Button>
      </template>

      <template #action="{ row }">
        <Button type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button type="link" size="small" danger @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>

    <Drawer :title="drawerTitle" :loading="submitLoading" class="w-full md:w-3/4 lg:w-1/2">
      <SiteGroupForm />
      <template #footer>
        <Button class="mr-2" @click="drawerApi.close">取消</Button>
        <Button type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </Button>
      </template>
    </Drawer>
  </div>
</template>
