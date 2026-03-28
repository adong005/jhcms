<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, onMounted, ref } from 'vue';
import { Button, message, Modal, Tag } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getLogListApi, deleteLogApi, batchDeleteLogApi, clearLogsApi } from '#/api/log';
import type { LogApi } from '#/api/log';
import { getUserListApi } from '#/api/user';
import { useUserStore } from '@vben/stores';

// 使用 API 中的类型定义
type LogRecord = LogApi.LogRecord;
type UserNode = {
  key: string;
  title: string;
  value: string;
  children?: UserNode[];
};

const userStore = useUserStore();
const userTreeOptions = ref<UserNode[]>([]);
const selectedUsernames = ref<string[]>([]);

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

function buildUserTree(
  users: Array<{ id: string; username: string; nickName?: string; createdBy?: string | number }>,
): UserNode[] {
  const nodeMap = new Map<string, UserNode>();
  const roots: UserNode[] = [];

  users.forEach((u) => {
    const label = u.nickName ? `${u.nickName} (${u.username})` : u.username;
    nodeMap.set(u.id, { key: u.id, value: u.username, title: label, children: [] });
  });

  users.forEach((u) => {
    const node = nodeMap.get(u.id);
    const parentId = u.createdBy ? String(u.createdBy) : '';
    if (!node) return;
    if (parentId && nodeMap.has(parentId)) {
      nodeMap.get(parentId)!.children!.push(node);
    } else {
      roots.push(node);
    }
  });

  const normalize = (nodes: UserNode[]): UserNode[] => nodes.map((n) => ({
    ...n,
    children: n.children && n.children.length > 0 ? normalize(n.children) : undefined,
  }));
  return normalize(roots);
}

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: isSuperAdmin.value ? 'TreeSelect' : 'Input',
      componentProps: isSuperAdmin.value
        ? {
            allowClear: true,
            maxTagCount: 2,
            multiple: true,
            showSearch: true,
            treeCheckable: true,
            treeData: userTreeOptions,
            placeholder: '请选择操作用户（可多选）',
            onChange: (val: string[]) => {
              selectedUsernames.value = Array.isArray(val) ? val : [];
            },
          }
        : {
            placeholder: '请输入操作用户',
          },
      fieldName: isSuperAdmin.value ? 'usernames' : 'username',
      label: '操作用户',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择操作类型',
        options: [
          { label: '全部', value: '' },
          { label: '登录', value: 'login' },
          { label: '登出', value: 'logout' },
          { label: '新增', value: 'create' },
          { label: '编辑', value: 'update' },
          { label: '删除', value: 'delete' },
          { label: '查询', value: 'query' },
          { label: '导出', value: 'export' },
        ],
      },
      fieldName: 'action',
      label: '操作类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择状态',
        options: [
          { label: '全部', value: '' },
          { label: '成功', value: 'success' },
          { label: '失败', value: 'fail' },
        ],
      },
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'DatePicker',
      componentProps: {
        placeholder: '请选择日期',
        style: { width: '100%' },
      },
      fieldName: 'date',
      label: '操作日期',
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
      field: 'username',
      title: '操作用户',
      width: 120,
      fixed: 'left',
    },
    {
      field: 'action',
      title: '操作类型',
      width: 100,
      slots: { default: 'action' },
    },
    {
      field: 'module',
      title: '操作模块',
      width: 120,
    },
    {
      field: 'description',
      title: '操作描述',
      minWidth: 200,
    },
    {
      field: 'requestJson',
      title: '请求JSON',
      minWidth: 280,
      showOverflow: 'tooltip',
    },
    {
      field: 'ip',
      title: 'IP地址',
      width: 140,
    },
    {
      field: 'status',
      title: '状态',
      width: 80,
      slots: { default: 'status' },
    },
    {
      field: 'duration',
      title: '耗时(ms)',
      width: 100,
    },
    {
      field: 'createTime',
      title: '操作时间',
      width: 180,
      sortable: true,
    },
    {
      field: 'action_btn',
      fixed: 'right',
      title: '操作',
      width: 120,
      slots: { default: 'action_btn' },
    },
  ],
  checkboxConfig: {
    highlight: true,
    reserve: true,
  },
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  printConfig: {
    sheetName: '日志列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '日志列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const payload = { ...formValues };
        if (isSuperAdmin.value) {
          payload.usernames = selectedUsernames.value;
          delete payload.username;
        }
        const resp = await getLogListApi({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...payload,
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
const gridEvents: VxeGridListeners<LogRecord> = {};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents,
  gridOptions,
});

onMounted(async () => {
  if (!isSuperAdmin.value) return;
  try {
    const resp = await getUserListApi({ page: 1, pageSize: 1000 });
    userTreeOptions.value = buildUserTree(resp.items || []);
  } catch (error) {
    console.error('加载用户树失败:', error);
  }
});

// 操作方法
function handleView(record: LogRecord) {
  const statusText = record.status === 'success' ? '成功' : '失败';
  const errorInfo = record.errorMsg ? `\n错误信息：${record.errorMsg}` : '';
  const requestInfo = record.requestJson ? `\n请求JSON：${record.requestJson}` : '';

  Modal.info({
    title: '日志详情',
    width: 600,
    content: `操作用户：${record.username}
操作类型：${getActionText(record.action)}
操作模块：${record.module}
操作描述：${record.description}
IP地址：${record.ip}
状态：${statusText}
耗时：${record.duration}ms
操作时间：${record.createTime}${errorInfo}${requestInfo}`,
  });
}

function handleDelete(record: LogRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除该日志记录吗？`,
    onOk: async () => {
      try {
        await deleteLogApi(record.id);
        message.success('删除日志成功');
        gridApi.reload();
      } catch (error) {
        message.error('删除日志失败');
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的日志');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 条日志吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: LogRecord) => record.id);
        await batchDeleteLogApi(ids);
        message.success(`批量删除 ${selectRecords.length} 条日志成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除日志失败');
      }
    },
  });
}

function handleClearLogs() {
  Modal.confirm({
    title: '确认清空',
    content: '确定要清空所有日志吗？此操作不可恢复！',
    okText: '确定清空',
    okType: 'danger',
    onOk: async () => {
      try {
        await clearLogsApi();
        message.success('清空日志成功');
        gridApi.reload();
      } catch (error) {
        message.error('清空日志失败');
      }
    },
  });
}

// 获取操作类型标签颜色
function getActionColor(action: string) {
  const colorMap: Record<string, string> = {
    login: 'green',
    logout: 'blue',
    create: 'cyan',
    update: 'orange',
    delete: 'red',
    query: 'default',
    export: 'purple',
  };
  return colorMap[action] || 'default';
}

// 获取操作类型文本
function getActionText(action: string) {
  const textMap: Record<string, string> = {
    login: '登录',
    logout: '登出',
    create: '新增',
    update: '编辑',
    delete: '删除',
    query: '查询',
    export: '导出',
  };
  return textMap[action] || action;
}
</script>

<template>
  <div class="p-4">
    <Grid>
      <!-- 工具栏按钮 -->
      <template #toolbar-tools>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
        <Button danger @click="handleClearLogs">
          <template #icon><span class="i-ant-design:clear-outlined" /></template>
          清空日志
        </Button>
      </template>

      <!-- 操作类型列 -->
      <template #action="{ row }">
        <Tag :color="getActionColor(row.action)">
          {{ getActionText(row.action) }}
        </Tag>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <Tag :color="row.status === 'success' ? 'success' : 'error'">
          {{ row.status === 'success' ? '成功' : '失败' }}
        </Tag>
      </template>

      <!-- 操作列 -->
      <template #action_btn="{ row }">
        <div class="flex gap-2">
          <Button type="link" size="small" @click="handleView(row)">
            <template #icon><span class="i-ant-design:eye-outlined" /></template>
            查看
          </Button>
          <Button type="link" size="small" danger @click="handleDelete(row)">
            <template #icon><span class="i-ant-design:delete-outlined" /></template>
            删除
          </Button>
        </div>
      </template>
    </Grid>
  </div>
</template>

<style scoped>
.flex {
  display: flex;
}

.gap-2 {
  gap: 0.5rem;
}

.space-y-2>*+* {
  margin-top: 0.5rem;
}

.p-4 {
  padding: 1rem;
}

.mt-2 {
  margin-top: 0.5rem;
}

.p-2 {
  padding: 0.5rem;
}

.bg-gray-100 {
  background-color: #f3f4f6;
}

.rounded {
  border-radius: 0.25rem;
}
</style>
