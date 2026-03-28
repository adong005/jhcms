<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, ref } from 'vue';
import { Button, message, Modal, Switch } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getCategoryListApi, updateCategoryStatusApi, updateCategoryInfoApi, createCategoryApi, deleteCategoryApi, batchDeleteCategoryApi } from '#/api/info/category';
import type { CategoryApi } from '#/api/info/category';

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentCategoryId.value = '';
    }
  },
});

const currentCategoryId = ref('');
const isEdit = computed(() => !!currentCategoryId.value);
const drawerTitle = computed(() => isEdit.value ? '编辑分类' : '新增分类');
const submitLoading = ref(false);

const [CategoryForm, formApi] = useVbenForm({
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
        placeholder: '请输入分类名称',
      },
      fieldName: 'name',
      label: '分类名称',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '显示', value: 1 },
          { label: '不显示', value: 0 },
        ],
      },
      defaultValue: 0,
      fieldName: 'isHome',
      label: '首页显示',
    },
    {
      component: 'InputNumber',
      componentProps: {
        placeholder: '请输入排序',
        min: 0,
      },
      fieldName: 'sort',
      label: '排序',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入分类描述',
        rows: 3,
      },
      fieldName: 'description',
      label: '分类描述',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '启用', value: 1 },
          { label: '禁用', value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
  ],
  showDefaultActions: false,
});

type CategoryRecord = CategoryApi.CategoryRecord;

const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入分类名称',
      },
      fieldName: 'name',
      label: '分类名称',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择状态',
        options: [
          { label: '全部', value: '' },
          { label: '启用', value: 1 },
          { label: '禁用', value: 0 },
        ],
      },
      fieldName: 'status',
      label: '状态',
    },
  ],
  showCollapseButton: true,
  submitButtonOptions: {
    content: '查询',
  },
};

const gridOptions: any = {
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
  columns: [
    {
      type: 'checkbox',
      width: 50,
      fixed: 'left',
    },
    {
      title: '序号',
      type: 'seq',
      width: 60,
      fixed: 'left',
    },
    {
      field: 'name',
      title: '分类名称',
      fixed: 'left',
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'isHome',
      title: '是否首页显示',
      width: 140,
      slots: { default: 'isHome' },
    },
    {
      field: 'sort',
      title: '排序',
      width: 80,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'status',
      title: '状态',
      width: 100,
      slots: { default: 'status' },
    },
    {
      field: 'createTime',
      title: '创建时间',
      width: 180,
      sortable: true,
    },
    {
      field: 'updateTime',
      title: '更新时间',
      width: 180,
    },
    {
      field: 'action',
      fixed: 'right',
      title: '操作',
      width: 150,
      slots: { default: 'action' },
    },
  ],
  checkboxConfig: {
    highlight: true,
    reserve: true,
  },
  keepSource: true,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  printConfig: {
    sheetName: '信息分类列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '信息分类列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const resp = await getCategoryListApi({
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
  sortConfig: {
    multiple: true,
    trigger: 'cell',
  },
  toolbarConfig: {
    custom: true,
    export: true,
    print: false,
    refresh: true,
    zoom: true,
    search: true,
  },
};

const gridEvents: VxeGridListeners<CategoryRecord> = {
  editClosed: async ({ row, column }) => {
    const editableFields = ['name', 'sort'];
    if (!editableFields.includes(column.field)) {
      return;
    }

    const $grid = gridApi.grid;
    if (!$grid.isUpdateByRow(row)) {
      return;
    }

    try {
      const updateData: any = {};
      updateData[column.field] = row[column.field as keyof CategoryRecord];

      await updateCategoryInfoApi(row.id, updateData);
      message.success('更新成功');
      gridApi.reload();
    } catch (error) {
      message.error('更新失败');
      gridApi.reload();
    }
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents,
  gridOptions,
});

function handleAdd() {
  currentCategoryId.value = '';
  formApi.resetForm();
  drawerApi.open();
}

async function handleEdit(record: CategoryRecord) {
  currentCategoryId.value = record.id;

  formApi.setValues({
    name: record.name,
    isHome: record.isHome ?? 0,
    sort: record.sort,
    description: record.description,
    status: record.status,
  });

  drawerApi.open();
}

async function handleSubmit() {
  try {
    await formApi.validate();

    submitLoading.value = true;

    const formValues = await formApi.getValues();
    const payload = {
      name: formValues.name,
      isHome: formValues.isHome,
      sort: formValues.sort,
      description: formValues.description,
      status: formValues.status,
    };

    if (isEdit.value) {
      await updateCategoryInfoApi(currentCategoryId.value, payload);
      message.success('编辑分类成功');
    } else {
      await createCategoryApi(payload);
      message.success('新增分类成功');
    }

    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    message.error(isEdit.value ? '编辑分类失败' : '新增分类失败');
    console.error('操作失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

function handleDelete(record: CategoryRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除分类 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await deleteCategoryApi(record.id);
        message.success(`删除分类 ${record.name} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除分类失败');
      }
    },
  });
}

async function handleStatusChange(record: CategoryRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newStatus = isChecked ? 1 : 0;
  const statusText = newStatus === 1 ? '启用' : '禁用';
  const oldStatus = record.status;

  Modal.confirm({
    title: '确认操作',
    content: `确定要${statusText}分类 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await updateCategoryStatusApi(record.id, newStatus);
        record.status = newStatus;
        message.success(`${statusText}分类成功`);
        gridApi.reload();
      } catch (error) {
        message.error(`${statusText}分类失败`);
        record.status = oldStatus;
      }
    },
  });
}

async function handleIsHomeChange(record: CategoryRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newValue = isChecked ? 1 : 0;
  const oldValue = record.isHome ?? 0;
  const text = newValue === 1 ? '显示' : '不显示';

  Modal.confirm({
    title: '确认操作',
    content: `确定将分类 "${record.name}" 设置为首页${text}吗？`,
    onOk: async () => {
      try {
        await updateCategoryInfoApi(record.id, { isHome: newValue });
        record.isHome = newValue;
        message.success('设置成功');
        gridApi.reload();
      } catch (error) {
        message.error('设置失败');
        record.isHome = oldValue;
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的分类');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 个分类吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: CategoryRecord) => record.id);
        await batchDeleteCategoryApi(ids);
        message.success(`批量删除 ${selectRecords.length} 个分类成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除分类失败');
      }
    },
  });
}

</script>

<template>
  <div class="p-4">
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          新增分类
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <template #status="{ row }">
        <Switch :checked="row.status === 1" checked-children="启用" un-checked-children="禁用"
          @change="(checked) => handleStatusChange(row, checked)" />
      </template>

      <template #isHome="{ row }">
        <Switch :checked="(row.isHome ?? 0) === 1" checked-children="显示" un-checked-children="不显示"
          @change="(checked) => handleIsHomeChange(row, checked)" />
      </template>

      <template #action="{ row }">
        <div class="flex gap-2">
          <Button type="link" size="small" @click="handleEdit(row)">
            <template #icon><span class="i-ant-design:edit-outlined" /></template>
            编辑
          </Button>
          <Button type="link" size="small" danger @click="handleDelete(row)">
            <template #icon><span class="i-ant-design:delete-outlined" /></template>
            删除
          </Button>
        </div>
      </template>
    </Grid>

    <Drawer :title="drawerTitle" class="w-150">
      <CategoryForm />
      <template #footer>
        <div class="flex justify-end gap-4">
          <Button @click="drawerApi.close()">取消</Button>
          <Button type="primary" :loading="submitLoading" @click="handleSubmit">
            保存
          </Button>
        </div>
      </template>
    </Drawer>
  </div>
</template>

<style scoped>
.flex {
  display: flex;
}

.gap-2 {
  gap: 0.5rem;
}

.gap-4 {
  gap: 1rem;
}

.justify-end {
  justify-content: flex-end;
}

.w-150 {
  width: 600px;
}
</style>
