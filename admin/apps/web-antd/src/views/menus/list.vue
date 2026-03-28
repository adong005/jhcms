<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, ref } from 'vue';
import { Button, message, Modal, Switch, Tag } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getMenuListApi, updateMenuStatusApi, updateMenuShowApi, updateMenuInfoApi, createMenuApi, deleteMenuApi, batchDeleteMenuApi } from '#/api/menu';
import type { MenuApi } from '#/api/menu';

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentMenuId.value = '';
    }
  },
});

const currentMenuId = ref('');
const isEdit = computed(() => !!currentMenuId.value);
const drawerTitle = computed(() => isEdit.value ? '编辑菜单' : '新增菜单');
const submitLoading = ref(false);

const [MenuForm, formApi] = useVbenForm({
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
        placeholder: '请输入菜单名称',
      },
      fieldName: 'name',
      label: '菜单名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入菜单路径',
      },
      fieldName: 'path',
      label: '菜单路径',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        placeholder: '请选择菜单类型',
        options: [
          { label: '目录', value: 'catalog' },
          { label: '菜单', value: 'menu' },
          { label: '按钮', value: 'button' },
        ],
      },
      fieldName: 'type',
      label: '菜单类型',
      rules: 'required',
    },
    {
      component: 'IconPicker',
      componentProps: {
        placeholder: '请选择图标',
      },
      fieldName: 'icon',
      label: '图标',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入组件路径',
      },
      fieldName: 'component',
      label: '组件路径',
    },
    {
      component: 'InputNumber',
      componentProps: {
        placeholder: '请输入排序',
        min: 0,
      },
      defaultValue: 0,
      fieldName: 'order',
      label: '排序',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '显示', value: 1 },
          { label: '隐藏', value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'isShow',
      label: '是否显示',
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

// 使用 API 中的类型定义
type MenuRecord = MenuApi.MenuRecord;

// 菜单类型映射
const menuTypeMap: Record<string, string> = {
  catalog: '目录',
  menu: '菜单',
  button: '按钮',
};

const menuTypeColorMap: Record<string, string> = {
  catalog: 'blue',
  menu: 'green',
  button: 'orange',
};

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入菜单名称',
      },
      fieldName: 'name',
      label: '菜单名称',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择菜单类型',
        options: [
          { label: '全部', value: '' },
          { label: '目录', value: 'catalog' },
          { label: '菜单', value: 'menu' },
          { label: '按钮', value: 'button' },
        ],
      },
      fieldName: 'type',
      label: '菜单类型',
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

// 表格配置
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
      title: '菜单名称',
      minWidth: 200,
      fixed: 'left',
      treeNode: true,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'path',
      title: '菜单路径',
      width: 200,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'type',
      title: '菜单类型',
      width: 100,
      slots: { default: 'type' },
    },
    {
      field: 'icon',
      title: '图标',
      width: 100,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'component',
      title: '组件路径',
      width: 200,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'order',
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
      field: 'isShow',
      title: '是否显示',
      width: 110,
      slots: { default: 'isShow' },
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
    enabled: false,
  },
  treeConfig: {
    transform: true,
    parentField: 'parentId',
    rowField: 'id',
    expandAll: true,
  },
  printConfig: {
    sheetName: '菜单列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '菜单列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async (_params: any, formValues: any) => {
        const resp = await getMenuListApi({
          page: 1,
          pageSize: 1000,
          ...formValues,
        });
        return {
          ...resp,
          items: (resp.items || []).map((item: any) => ({
            ...item,
            id: String(item.id),
            parentId: item.parentId === null || item.parentId === undefined ? null : String(item.parentId),
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

// 表格事件
const gridEvents: VxeGridListeners<MenuRecord> = {
  // 单元格编辑完成后自动提交
  editClosed: async ({ row, column }) => {
    // 只处理可编辑的字段
    const editableFields = ['name', 'path', 'icon', 'component', 'order'];
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
      updateData[column.field] = row[column.field as keyof MenuRecord];

      await updateMenuInfoApi(row.id, updateData);
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

// 展开所有节点
const expandAll = () => {
  gridApi.grid?.setAllTreeExpand(true);
};

// 折叠所有节点
const collapseAll = () => {
  gridApi.grid?.setAllTreeExpand(false);
};

// 操作方法
function handleAdd() {
  currentMenuId.value = '';
  formApi.resetForm();
  drawerApi.open();
}

async function handleEdit(record: MenuRecord) {
  currentMenuId.value = record.id;

  formApi.setValues({
    name: record.name,
    path: record.path,
    type: record.type,
    icon: record.icon,
    component: record.component,
    order: record.order,
    status: record.status,
    isShow: record.isShow ?? 1,
  });

  drawerApi.open();
}

async function handleSubmit() {
  try {
    await formApi.validate();

    submitLoading.value = true;

    const formValues = await formApi.getValues();

    if (isEdit.value) {
      // 编辑菜单
      await updateMenuInfoApi(currentMenuId.value, formValues);
      message.success('编辑菜单成功');
    } else {
      // 新增菜单
      await createMenuApi(formValues);
      message.success('新增菜单成功');
    }

    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    message.error(isEdit.value ? '编辑菜单失败' : '新增菜单失败');
    console.error('操作失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

function handleDelete(record: MenuRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除菜单 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await deleteMenuApi(record.id);
        message.success(`删除菜单 ${record.name} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除菜单失败');
      }
    },
  });
}

async function handleStatusChange(record: MenuRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newStatus = isChecked ? 1 : 0;
  const statusText = newStatus === 1 ? '启用' : '禁用';
  const oldStatus = record.status;

  Modal.confirm({
    title: '确认操作',
    content: `确定要${statusText}菜单 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await updateMenuStatusApi(record.id, newStatus);
        record.status = newStatus;
        message.success(`${statusText}菜单成功`);
        gridApi.reload();
      } catch (error) {
        message.error(`${statusText}菜单失败`);
        record.status = oldStatus;
      }
    },
  });
}

async function handleShowChange(record: MenuRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newValue = isChecked ? 1 : 0;
  const text = newValue === 1 ? '显示' : '隐藏';
  const oldValue = record.isShow ?? 1;
  Modal.confirm({
    title: '确认操作',
    content: `确定要${text}菜单 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await updateMenuShowApi(record.id, newValue);
        record.isShow = newValue;
        message.success(`${text}菜单成功`);
        gridApi.reload();
      } catch (error) {
        message.error(`${text}菜单失败`);
        record.isShow = oldValue;
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的菜单');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 个菜单吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: MenuRecord) => record.id);
        await batchDeleteMenuApi(ids);
        message.success(`批量删除 ${selectRecords.length} 个菜单成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除菜单失败');
      }
    },
  });
}

</script>

<template>
  <div class="p-4">
    <Grid>
      <!-- 工具栏按钮 -->
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          新增菜单
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
        <Button @click="expandAll">
          <template #icon><span class="i-ant-design:plus-square-outlined" /></template>
          展开全部
        </Button>
        <Button @click="collapseAll">
          <template #icon><span class="i-ant-design:minus-square-outlined" /></template>
          折叠全部
        </Button>
      </template>

      <!-- 菜单类型列 -->
      <template #type="{ row }">
        <Tag :color="menuTypeColorMap[row.type] || 'default'">
          {{ menuTypeMap[row.type] || row.type }}
        </Tag>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <Switch :checked="row.status === 1" checked-children="启用" un-checked-children="禁用"
          @change="(checked) => handleStatusChange(row, checked)" />
      </template>

      <template #isShow="{ row }">
        <Switch :checked="(row.isShow ?? 1) === 1" checked-children="显示" un-checked-children="隐藏"
          @change="(checked) => handleShowChange(row, checked)" />
      </template>

      <!-- 操作列 -->
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

    <!-- 菜单表单 Drawer -->
    <Drawer :title="drawerTitle" class="w-150">
      <MenuForm />
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
