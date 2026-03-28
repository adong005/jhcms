<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, ref } from 'vue';
import { Button, message, Modal, Switch, Tree } from 'ant-design-vue';
import type { TreeProps } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getRoleListApi, updateRoleStatusApi, updateRoleInfoApi, createRoleApi, deleteRoleApi, batchDeleteRoleApi, getRolePermissionApi, updateRolePermissionApi } from '#/api/role';
import type { RoleApi } from '#/api/role';
import { getPermissionListApi } from '#/api/permission';

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentRoleId.value = '';
    }
  },
});

// 权限抽屉
const [PermissionDrawer, permissionDrawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      currentRoleId.value = '';
      checkedKeys.value = [];
    }
  },
});

const currentRoleId = ref('');
const isEdit = computed(() => !!currentRoleId.value);
const drawerTitle = computed(() => isEdit.value ? '编辑角色' : '新增角色');
const submitLoading = ref(false);

// 权限相关状态
const permissionTreeData = ref<TreeProps['treeData']>([]);
const checkedKeys = ref<string[]>([]);
const permissionLoading = ref(false);
const currentRoleName = ref('');

const [RoleForm, formApi] = useVbenForm({
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
        placeholder: '请输入角色名称',
      },
      fieldName: 'name',
      label: '角色名称',
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入角色描述',
        rows: 3,
      },
      fieldName: 'description',
      label: '角色描述',
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
type RoleRecord = RoleApi.RoleRecord;

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入角色名称',
      },
      fieldName: 'name',
      label: '角色名称',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入角色编码',
      },
      fieldName: 'code',
      label: '角色编码',
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
      title: '角色名称',
      width: 150,
      fixed: 'left',
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'code',
      title: '角色编码',
      width: 150,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'description',
      title: '角色描述',
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
    enabled: false,
  },
  printConfig: {
    sheetName: '角色列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '角色列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async (_params: any, formValues: any) => {
        const resp = await getRoleListApi({
          page: 1,
          pageSize: 1000,
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

// 表格事件
const gridEvents: VxeGridListeners<RoleRecord> = {
  // 单元格编辑完成后自动提交
  editClosed: async ({ row, column }) => {
    // 只处理可编辑的字段
    const editableFields = ['name', 'description'];
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
      updateData[column.field] = row[column.field as keyof RoleRecord];

      await updateRoleInfoApi(row.id, updateData);
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

// 操作方法
function handleAdd() {
  currentRoleId.value = '';
  formApi.resetForm();
  drawerApi.open();
}

async function handleEdit(record: RoleRecord) {
  currentRoleId.value = record.id;

  formApi.setValues({
    name: record.name,
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
    const submitData = {
      name: formValues.name,
      description: formValues.description,
      status: formValues.status,
    };

    if (isEdit.value) {
      // 编辑角色
      await updateRoleInfoApi(currentRoleId.value, submitData);
      message.success('编辑角色成功');
    } else {
      // 新增角色
      await createRoleApi(submitData);
      message.success('新增角色成功');
    }

    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    message.error(isEdit.value ? '编辑角色失败' : '新增角色失败');
    console.error('操作失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

function handleDelete(record: RoleRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除角色 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await deleteRoleApi(record.id);
        message.success(`删除角色 ${record.name} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除角色失败');
      }
    },
  });
}

async function handleStatusChange(record: RoleRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newStatus = isChecked ? 1 : 0;
  const statusText = newStatus === 1 ? '启用' : '禁用';
  const oldStatus = record.status;

  Modal.confirm({
    title: '确认操作',
    content: `确定要${statusText}角色 "${record.name}" 吗？`,
    onOk: async () => {
      try {
        await updateRoleStatusApi(record.id, newStatus);
        record.status = newStatus;
        message.success(`${statusText}角色成功`);
        gridApi.reload();
      } catch (error) {
        message.error(`${statusText}角色失败`);
        record.status = oldStatus;
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的角色');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 个角色吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: RoleRecord) => record.id);
        await batchDeleteRoleApi(ids);
        message.success(`批量删除 ${selectRecords.length} 个角色成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除角色失败');
      }
    },
  });
}

function buildPermissionTree(items: Array<{ code: string; id: string; name: string }>) {
  const moduleMap = new Map<string, Map<string, Array<{ code: string; id: string; name: string }>>>();
  for (const item of items) {
    const [module = 'other', resource = 'misc'] = String(item.code || '').split(':');
    if (!moduleMap.has(module)) {
      moduleMap.set(module, new Map());
    }
    const resourceMap = moduleMap.get(module)!;
    if (!resourceMap.has(resource)) {
      resourceMap.set(resource, []);
    }
    resourceMap.get(resource)!.push(item);
  }

  const tree: any[] = [];
  for (const [module, resourceMap] of moduleMap.entries()) {
    const moduleNode: any = {
      key: `module:${module}`,
      title: module,
      disableCheckbox: true,
      children: [],
    };
    for (const [resource, perms] of resourceMap.entries()) {
      const resourceNode: any = {
        key: `resource:${module}:${resource}`,
        title: resource,
        disableCheckbox: true,
        children: perms.map((p) => ({
          key: String(p.id),
          title: `${p.name} (${p.code})`,
        })),
      };
      moduleNode.children.push(resourceNode);
    }
    tree.push(moduleNode);
  }
  return tree;
}

// 获取权限点数据并组装树
async function loadPermissionData() {
  try {
    const response = await getPermissionListApi({ page: 1, pageSize: 1000 });
    const permissions = (response.items || []).map((item: any) => ({
      code: String(item.code || ''),
      id: String(item.id),
      name: String(item.name || ''),
    }));
    permissionTreeData.value = buildPermissionTree(permissions);
  } catch (error) {
    message.error('加载权限树失败');
    console.error('加载权限树失败:', error);
  }
}

// 打开权限配置抽屉
async function handlePermission(record: RoleRecord) {
  currentRoleId.value = String((record as any).id ?? '');
  currentRoleName.value = (record as any).name || (record as any).roleName || '';

  try {
    // 加载权限树
    await loadPermissionData();

    // 获取该角色已有的权限
    const permissionResult = await getRolePermissionApi(String((record as any).id ?? ''));
    const ids = permissionResult.permissionIds || permissionResult.menuIds || [];
    checkedKeys.value = ids.map((id: string | number) => String(id));

    permissionDrawerApi.open();
  } catch (error) {
    message.error('加载角色权限失败');
    console.error('加载角色权限失败:', error);
  }
}

// 保存权限配置
async function handleSavePermission() {
  if (checkedKeys.value.length === 0) {
    message.warning('请至少选择一个权限点');
    return;
  }

  try {
    permissionLoading.value = true;

    await updateRolePermissionApi(
      currentRoleId.value,
      checkedKeys.value.map((id) => String(id)),
    );

    message.success(
      `已为角色"${currentRoleName.value}"配置 ${checkedKeys.value.length} 个权限点`,
    );
    permissionDrawerApi.close();
  } catch (error) {
    message.error('权限配置保存失败');
    console.error('保存权限失败:', error);
  } finally {
    permissionLoading.value = false;
  }
}

// 树节点选中变化
function onTreeCheck(checked: any) {
  const keys = checked.checked || checked;
  checkedKeys.value = (keys || []).map((id: string | number) => String(id));
}

</script>

<template>
  <div class="p-4">
    <Grid>
      <!-- 工具栏按钮 -->
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          新增角色
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <Switch :checked="row.status === 1" checked-children="启用" un-checked-children="禁用"
          @change="(checked) => handleStatusChange(row, checked)" />
      </template>

      <!-- 操作列 -->
      <template #action="{ row }">
        <div class="flex gap-2">
          <Button type="link" size="small" @click="handleEdit(row)">
            <template #icon><span class="i-ant-design:edit-outlined" /></template>
            编辑
          </Button>
          <Button type="link" size="small" @click="handlePermission(row)">
            <template #icon><span class="i-ant-design:safety-outlined" /></template>
            权限
          </Button>
          <Button type="link" size="small" danger @click="handleDelete(row)">
            <template #icon><span class="i-ant-design:delete-outlined" /></template>
            删除
          </Button>
        </div>
      </template>
    </Grid>

    <!-- 角色表单 Drawer -->
    <Drawer :title="drawerTitle" class="w-150">
      <RoleForm />
      <template #footer>
        <div class="flex justify-end gap-4">
          <Button @click="drawerApi.close()">取消</Button>
          <Button type="primary" :loading="submitLoading" @click="handleSubmit">
            保存
          </Button>
        </div>
      </template>
    </Drawer>

    <!-- 权限配置 Drawer -->
    <PermissionDrawer :title="`配置权限 - ${currentRoleName}`" class="w-150">
      <div class="p-4">
        <div class="mb-4 text-gray-600">
          请选择该角色可以访问的权限点
        </div>
        <Tree v-model:checkedKeys="checkedKeys" checkable :tree-data="permissionTreeData" :default-expand-all="true"
          @check="onTreeCheck" />
      </div>
      <template #footer>
        <div class="flex justify-end gap-4">
          <Button @click="permissionDrawerApi.close()">取消</Button>
          <Button type="primary" :loading="permissionLoading" @click="handleSavePermission">
            保存
          </Button>
        </div>
      </template>
    </PermissionDrawer>
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

.mb-4 {
  margin-bottom: 1rem;
}

.text-gray-600 {
  color: #6b7280;
}

.p-4 {
  padding: 1rem;
}
</style>
