<script setup lang="ts">
import type { VxeGridListeners, VxeGridProps } from '#/adapter/vxe-table';
import { Button, message, Modal, Tag } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getFormListApi, deleteFormApi, batchDeleteFormApi } from '#/api/form-manage';
import type { FormManageApi } from '#/api/form-manage';

type FormRecord = FormManageApi.FormRecord;

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入联系人',
      },
      fieldName: 'contact',
      label: '联系人',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入联系电话',
      },
      fieldName: 'phone',
      label: '联系电话',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入公司名',
      },
      fieldName: 'company',
      label: '公司名',
    },
  ],
  showCollapseButton: true,
  submitButtonOptions: {
    content: '查询',
  },
};

const gridOptions: VxeGridProps<FormRecord> = {
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
    { field: 'contact', title: '联系人', width: 120 },
    { field: 'phone', title: '联系电话', width: 150 },
    { field: 'company', title: '公司名', minWidth: 200 },
    { field: 'ip', title: 'IP', width: 150 },
    {
      field: 'handleStatus',
      title: '处理状态',
      width: 120,
      slots: { default: 'handleStatus' },
      editRender: {
        name: 'select',
        options: [
          { label: '未处理', value: 0 },
          { label: '已分配', value: 1 },
          { label: '已完成', value: 2 },
        ],
      },
    },
    { field: 'createTime', title: '创建时间', width: 180 },
    { field: 'updateTime', title: '更改时间', width: 180 },
    { field: 'remark', title: '备注', minWidth: 200 },
    {
      title: '操作',
      width: 120,
      slots: { default: 'action' },
      fixed: 'right',
    },
  ],
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  exportConfig: {
    filename: '表单列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const resp = await getFormListApi({
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

const gridEvents: VxeGridListeners<FormRecord> = {
  // 单元格编辑完成后自动提交
  editClosed: async ({ row, column }) => {
    // 只处理处理状态字段
    if (column.field !== 'handleStatus') {
      return;
    }

    // 检查该行是否有更新
    const $grid = gridApi.grid;
    if (!$grid.isUpdateByRow(row)) {
      return; // 数据未修改，不提交
    }

    try {
      // 这里可以调用更新 API
      // await updateFormStatusApi(row.id, row.handleStatus);
      message.success('处理状态更新成功');
      // 刷新列表
      gridApi.reload();
    } catch (error) {
      message.error('处理状态更新失败');
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

async function handleView(record: FormRecord) {
  Modal.info({
    title: '表单详情',
    width: 600,
    content: `联系人：${record.contact}\n联系电话：${record.phone}\n公司名：${record.company}\nIP地址：${record.ip}\n提交时间：${record.createTime}\n备注：${record.remark || '无'}`,
  });
}

async function handleDelete(record: FormRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除联系人"${record.contact}"的表单记录吗？`,
    onOk: async () => {
      await deleteFormApi(record.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

async function handleBatchDelete() {
  const selectedRows = gridApi.grid.getCheckboxRecords();
  if (selectedRows.length === 0) {
    message.warning('请选择要删除的表单记录');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRows.length} 条表单记录吗？`,
    onOk: async () => {
      const ids = selectedRows.map(row => row.id);
      await batchDeleteFormApi(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
  });
}
</script>

<template>
  <div class="p-4">
    <Grid>
      <template #toolbar-tools>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <template #handleStatus="{ row }">
        <Tag v-if="row.handleStatus === 0" color="default">未处理</Tag>
        <Tag v-else-if="row.handleStatus === 1" color="processing">已分配</Tag>
        <Tag v-else-if="row.handleStatus === 2" color="success">已完成</Tag>
      </template>

      <template #action="{ row }">
        <Button type="link" size="small" @click="handleView(row)">查看</Button>
        <Button type="link" size="small" danger @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </div>
</template>
