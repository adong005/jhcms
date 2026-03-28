<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenForm } from '#/adapter/form';
import {
  batchDeletePermissionApi,
  createPermissionApi,
  deletePermissionApi,
  getPermissionListApi,
  updatePermissionApi,
} from '#/api/permission';
import type { PermissionApi } from '#/api/permission';

import { useVbenDrawer } from '@vben/common-ui';
import { Button, message, Modal, Switch } from 'ant-design-vue';

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentPermissionId.value = '';
    }
  },
});

const currentPermissionId = ref('');
const isEdit = computed(() => !!currentPermissionId.value);
const drawerTitle = computed(() => (isEdit.value ? '编辑权限' : '新增权限'));
const submitLoading = ref(false);

const [PermissionForm, formApi] = useVbenForm({
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
        placeholder: '请输入权限名称',
      },
      fieldName: 'name',
      label: '权限名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入权限编码（如 system:user:list）',
      },
      fieldName: 'code',
      label: '权限编码',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入模块（如 system）',
      },
      fieldName: 'module',
      label: '模块',
    },
    {
      component: 'Switch',
      componentProps: {
        checkedChildren: '可委派',
        unCheckedChildren: '不可委派',
      },
      defaultValue: true,
      fieldName: 'isDelegable',
      label: '可委派',
    },
  ],
  showDefaultActions: false,
});

type PermissionRecord = PermissionApi.PermissionRecord;
type PermissionTreeRow = PermissionRecord & {
  id: string;
  isGroup?: boolean;
  parentId?: null | string;
};

const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入权限名称',
      },
      fieldName: 'name',
      label: '权限名称',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入权限编码',
      },
      fieldName: 'code',
      label: '权限编码',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入模块',
      },
      fieldName: 'module',
      label: '模块',
    },
  ],
  showCollapseButton: true,
  submitButtonOptions: {
    content: '查询',
  },
};

function buildPermissionTreeRows(items: PermissionRecord[]): PermissionTreeRow[] {
  const moduleMap = new Map<string, Map<string, PermissionRecord[]>>();
  for (const item of items) {
    const [module = 'other', resource = 'misc'] = String(item.code || '').split(':');
    if (!moduleMap.has(module)) moduleMap.set(module, new Map());
    const resourceMap = moduleMap.get(module)!;
    if (!resourceMap.has(resource)) resourceMap.set(resource, []);
    resourceMap.get(resource)!.push(item);
  }

  const rows: PermissionTreeRow[] = [];
  for (const [module, resourceMap] of moduleMap.entries()) {
    const moduleId = `module:${module}`;
    rows.push({
      id: moduleId,
      name: module,
      code: '',
      module,
      isDelegable: false,
      createTime: '',
      updateTime: '',
      isGroup: true,
      parentId: null,
    });

    for (const [resource, perms] of resourceMap.entries()) {
      const resourceId = `resource:${module}:${resource}`;
      rows.push({
        id: resourceId,
        name: resource,
        code: '',
        module,
        isDelegable: false,
        createTime: '',
        updateTime: '',
        isGroup: true,
        parentId: moduleId,
      });

      for (const perm of perms) {
        rows.push({
          ...perm,
          id: String(perm.id),
          parentId: resourceId,
        });
      }
    }
  }
  return rows;
}

const gridOptions: any = {
  border: true,
  stripe: true,
  rowConfig: {
    keyField: 'id',
  },
  treeConfig: {
    rowField: 'id',
    parentField: 'parentId',
    transform: true,
    expandAll: true,
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
      title: '权限名称',
      minWidth: 220,
      fixed: 'left',
      treeNode: true,
      slots: { default: 'name' },
    },
    {
      field: 'code',
      title: '权限编码',
      minWidth: 220,
    },
    {
      field: 'module',
      title: '模块',
      width: 140,
    },
    {
      field: 'isDelegable',
      title: '可委派',
      width: 120,
      slots: { default: 'isDelegable' },
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
      width: 140,
      slots: { default: 'action' },
    },
  ],
  checkboxConfig: {
    highlight: true,
    reserve: true,
    checkMethod: ({ row }: { row: PermissionTreeRow }) => !row.isGroup,
  },
  pagerConfig: {
    enabled: false,
  },
  proxyConfig: {
    ajax: {
      query: async (_params: any, formValues: any) => {
        const resp = await getPermissionListApi({
          page: 1,
          pageSize: 1000,
          ...formValues,
        });
        const items = (resp.items || []).map((item: any) => ({
          ...item,
          id: String(item.id),
        })) as PermissionRecord[];
        const treeRows = buildPermissionTreeRows(items);
        return {
          ...resp,
          items: treeRows,
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

function handleAdd() {
  currentPermissionId.value = '';
  formApi.resetForm();
  drawerApi.open();
}

async function handleEdit(record: PermissionRecord) {
  if ((record as PermissionTreeRow).isGroup) {
    return;
  }
  currentPermissionId.value = record.id;
  formApi.setValues({
    name: record.name,
    code: record.code,
    module: record.module,
    isDelegable: record.isDelegable,
  });
  drawerApi.open();
}

async function handleSubmit() {
  try {
    await formApi.validate();
    submitLoading.value = true;
    const values = await formApi.getValues();

    if (isEdit.value) {
      await updatePermissionApi(currentPermissionId.value, values);
      message.success('编辑权限成功');
    } else {
      await createPermissionApi(values);
      message.success('新增权限成功');
    }
    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    message.error(isEdit.value ? '编辑权限失败' : '新增权限失败');
  } finally {
    submitLoading.value = false;
  }
}

function handleDelete(record: PermissionRecord) {
  if ((record as PermissionTreeRow).isGroup) {
    return;
  }
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除权限 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await deletePermissionApi(record.id);
        message.success(`删除权限 ${record.name} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除权限失败');
      }
    },
  });
}

function handleBatchDelete() {
  const rows = (gridApi.grid.getCheckboxRecords() || []).filter(
    (row: PermissionTreeRow) => !row.isGroup,
  );
  if (rows.length === 0) {
    message.warning('请先选择要删除的权限');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${rows.length} 个权限吗？`,
    onOk: async () => {
      try {
        const ids = rows.map((row: PermissionRecord) => row.id);
        await batchDeletePermissionApi(ids);
        message.success(`批量删除 ${rows.length} 个权限成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除权限失败');
      }
    },
  });
}

async function handleDelegableChange(record: PermissionRecord, checked: boolean) {
  if ((record as PermissionTreeRow).isGroup) {
    return;
  }
  const oldValue = record.isDelegable;
  record.isDelegable = checked;
  try {
    await updatePermissionApi(record.id, { isDelegable: checked });
    message.success('更新成功');
    gridApi.reload();
  } catch (error) {
    record.isDelegable = oldValue;
    message.error('更新失败');
  }
}
</script>

<template>
  <div class="p-4">
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          新增权限
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <template #name="{ row }">
        <span :style="{ fontWeight: row.isGroup ? 600 : 400 }">
          {{ row.name }}
        </span>
      </template>

      <template #isDelegable="{ row }">
        <Switch
          :checked="!!row.isDelegable"
          :disabled="!!row.isGroup"
          checked-children="是"
          un-checked-children="否"
          @change="(checked) => handleDelegableChange(row, !!checked)"
        />
      </template>

      <template #action="{ row }">
        <div v-if="!row.isGroup" class="flex gap-2">
          <Button type="link" size="small" @click="handleEdit(row)">
            <template #icon><span class="i-ant-design:edit-outlined" /></template>
            编辑
          </Button>
          <Button type="link" size="small" danger @click="handleDelete(row)">
            <template #icon><span class="i-ant-design:delete-outlined" /></template>
            删除
          </Button>
        </div>
        <span v-else class="text-gray-400">-</span>
      </template>
    </Grid>

    <Drawer :title="drawerTitle" class="w-150">
      <PermissionForm />
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
