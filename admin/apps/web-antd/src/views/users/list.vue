<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, h, ref } from 'vue';
import { Button, message, Modal, Switch, Tag } from 'ant-design-vue';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm, z } from '#/adapter/form';
import dayjs from 'dayjs';
import {
  getUserListApi,
  updateUserStatusApi,
  updateUserInfoApi,
  createUserApi,
  deleteUserApi,
  batchDeleteUserApi,
  resetUserPasswordApi,
} from '#/api/user';
import type { UserApi } from '#/api/user';
import { useUserStore } from '@vben/stores';
import { useAuthStore } from '#/store';
import { getRoleListApi } from '#/api/role';

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      formApi.resetForm();
      currentUserId.value = '';
    }
  },
});

const userStore = useUserStore();
const authStore = useAuthStore();

const currentUserId = ref('');
const isEdit = computed(() => !!currentUserId.value);

/** 可「进入他人后台」：与后端 canImpersonate 一致（平台标志或 roles 含 super_admin） */
const canEnterUserBackend = computed(() => {
  const u = userStore.userInfo as {
    isPlatformSuperAdmin?: boolean;
    roles?: string[];
  } | null;
  if (!u) return false;
  if (u.isPlatformSuperAdmin) return true;
  return (u.roles ?? []).includes('super_admin');
});
const drawerTitle = computed(() => isEdit.value ? '编辑用户' : '新增用户');
const submitLoading = ref(false);
const roleOptions = ref<Array<{ label: string; value: string }>>([
  { label: '用户', value: 'user' },
]);

function defaultExpireDate() {
  return dayjs().add(1, 'month');
}

const [UserForm, formApi] = useVbenForm({
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
        placeholder: '请输入用户名',
      },
      fieldName: 'username',
      label: '用户名',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入昵称',
      },
      fieldName: 'nickName',
      label: '昵称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入邮箱',
        type: 'email',
      },
      fieldName: 'email',
      label: '邮箱',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入手机号',
      },
      fieldName: 'phone',
      label: '手机号',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: roleOptions,
        placeholder: '请选择角色',
      },
      fieldName: 'role',
      label: '角色',
      rules: 'required',
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
    {
      component: 'DatePicker',
      componentProps: {
        format: 'YYYY-MM-DD',
        placeholder: '请选择过期日期',
        style: { width: '100%' },
      },
      fieldName: 'expireDate',
      label: '过期日期',
      rules: 'required',
    },
    {
      component: 'InputPassword',
      componentProps: {
        placeholder: '请输入密码',
      },
      fieldName: 'password',
      label: '密码',
      rules: 'required',
    },
  ],
  showDefaultActions: false,
});

// 使用 API 中的类型定义
type UserRecord = UserApi.UserRecord;

// 角色：表格展示（后端列表多为字符串 super_admin / admin / user）
const roleColorMap: Record<number, string> = {
  1: 'red',
  2: 'blue',
  3: 'default',
};
const roleStringColorMap: Record<string, string> = {
  super_admin: 'red',
  admin: 'blue',
  user: 'default',
};
const roleNameMap: Record<number, string> = {
  1: '超级管理员',
  2: '管理员',
  3: '用户',
};
const roleStringNameMap: Record<string, string> = {
  super_admin: '超级管理员',
  admin: '管理员',
  user: '用户',
};

function roleTagColor(role: string | number | undefined | null): string {
  if (role === undefined || role === null || role === '') return 'default';
  if (typeof role === 'number') return roleColorMap[role] ?? 'default';
  const s = String(role);
  if (/^\d+$/.test(s)) return roleColorMap[Number(s)] ?? 'default';
  return roleStringColorMap[s] ?? 'default';
}

function roleDisplayName(role: string | number | undefined | null): string {
  if (role === undefined || role === null || role === '') return '-';
  if (typeof role === 'number') return roleNameMap[role] ?? String(role);
  const s = String(role);
  if (/^\d+$/.test(s)) return roleNameMap[Number(s)] ?? s;
  return roleStringNameMap[s] ?? s;
}

/** 编辑表单里 Select 的 value 统一为角色 code 字符串 */
function roleToFormValue(role: string | number | undefined | null): string | undefined {
  if (role === undefined || role === null || role === '') return undefined;
  if (typeof role === 'number' && role >= 1 && role <= 3) {
    return ({ 1: 'super_admin', 2: 'admin', 3: 'user' } as Record<number, string>)[role];
  }
  const s = String(role);
  if (/^\d+$/.test(s)) {
    const n = Number(s);
    return n >= 1 && n <= 3 ? ({ 1: 'super_admin', 2: 'admin', 3: 'user' } as Record<number, string>)[n] : undefined;
  }
  return s;
}

const currentRoleCode = computed(() => {
  const u = userStore.userInfo as { role?: string; roles?: string[] } | null;
  return (u?.role || u?.roles?.[0] || '').trim();
});

async function loadRoleOptions() {
  try {
    const operatorRole = currentRoleCode.value;
    const options: Array<{ label: string; value: string }> = [];
    if (operatorRole === 'super_admin') {
      options.push(
        { label: '超级管理员', value: 'super_admin' },
        { label: '管理员', value: 'admin' },
        { label: '用户', value: 'user' },
      );
    } else {
      // 管理员最小可选角色：用户
      options.push({ label: '用户', value: 'user' });
    }
    const resp = await getRoleListApi({ page: 1, pageSize: 1000, status: 1 });
    for (const item of resp.items || []) {
      const code = String(item.code || '').trim();
      if (!code || ['super_admin', 'admin', 'user'].includes(code)) {
        continue;
      }
      options.push({ label: item.name || code, value: code });
    }
    const dedup = new Map<string, { label: string; value: string }>();
    options.forEach((o) => dedup.set(o.value, o));
    roleOptions.value = [...dedup.values()];
    formApi.updateSchema([
      {
        fieldName: 'role',
        componentProps: {
          allowClear: true,
          options: roleOptions,
          placeholder: '请选择角色',
        },
      },
    ]);
  } catch (error) {
    roleOptions.value = [{ label: '用户', value: 'user' }];
  }
}

function formatSubmitError(error: unknown, fallback: string): string {
  if (error && typeof error === 'object') {
    const err = error as {
      message?: string;
      response?: {
        status?: number;
        data?: unknown;
      };
      errors?: Record<string, unknown>;
    };
    const status = err.response?.status;
    if (status === 404) {
      return `${fallback}（请求 404：请确认后端已重新编译部署）`;
    }
    const data = err.response?.data;
    if (typeof data === 'string' && data.trim()) return data.trim();
    if (data && typeof data === 'object') {
      const d = data as { message?: string; error?: string };
      const fromResp = d.message ?? d.error;
      if (typeof fromResp === 'string' && fromResp.trim()) return fromResp.trim();
    }
    if (err.errors && typeof err.errors === 'object') {
      const first = Object.values(err.errors).flat()[0];
      if (typeof first === 'string' && first.trim()) return first.trim();
    }
    if (typeof err.message === 'string' && err.message && err.message !== '[object Object]') {
      return err.message;
    }
  }
  if (error instanceof Error && error.message) return error.message;
  return fallback;
}

// 搜索表单配置
const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入用户名',
      },
      fieldName: 'username',
      label: '用户名',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入昵称',
      },
      fieldName: 'nickName',
      label: '昵称',
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
      field: 'username',
      title: '用户名',
      width: 120,
      fixed: 'left',
    },
    {
      field: 'nickName',
      title: '昵称',
      width: 120,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'email',
      title: '邮箱',
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'phone',
      title: '手机号',
      width: 130,
      editRender: {
        name: 'input',
      },
    },
    {
      field: 'role',
      title: '角色',
      width: 120,
      slots: { default: 'role' },
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
      title: '更改时间',
      width: 180,
      sortable: true,
    },
    {
      field: 'lastLoginDate',
      title: '最后登录日期',
      width: 150,
    },
    {
      field: 'expireDate',
      title: '过期日期',
      width: 150,
    },
    {
      field: 'action',
      fixed: 'right',
      title: '操作',
      width: 280,
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
    sheetName: '用户列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '用户列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const resp = await getUserListApi({
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

// 表格事件
const gridEvents: VxeGridListeners<UserRecord> = {
  // 单元格编辑完成后自动提交
  editClosed: async ({ row, column }) => {
    // 只处理可编辑的字段
    const editableFields = ['nickName', 'email', 'phone'];
    if (!editableFields.includes(column.field)) {
      return;
    }

    // 仅在有真实变更时才提交，避免点击进入编辑即触发“保存成功”
    const $grid = gridApi.grid;
    if (!$grid.isUpdateByRow(row)) {
      return;
    }

    try {
      const updateData: any = {};
      updateData[column.field] = row[column.field as keyof UserRecord];

      await updateUserInfoApi(row.id, updateData);
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
  void loadRoleOptions();
  currentUserId.value = '';
  formApi.resetForm();
  formApi.setValues({
    role: 'user',
    expireDate: defaultExpireDate(),
  });
  formApi.updateSchema([
    {
      fieldName: 'password',
      rules: 'required',
      componentProps: { placeholder: '请输入密码' },
    },
  ]);
  drawerApi.open();
}

async function handleEdit(record: UserRecord) {
  await loadRoleOptions();
  currentUserId.value = record.id;

  formApi.updateSchema([
    {
      fieldName: 'password',
      rules: z.string().optional(),
      componentProps: { placeholder: '留空表示不修改密码' },
    },
  ]);

  // 设置表单值，将日期字符串转换为 dayjs 对象；角色后端为字符串时需映射为 Select 的 1/2/3
  formApi.setValues({
    username: record.username,
    nickName: record.nickName,
    email: record.email,
    phone: record.phone,
    role: roleToFormValue(record.role),
    status: record.status,
    expireDate: record.expireDate ? dayjs(record.expireDate) : defaultExpireDate(),
    password: '',
  });

  drawerApi.open();
}

async function handleSubmit() {
  try {
    await formApi.validate();

    submitLoading.value = true;

    const formValues = await formApi.getValues();

    // 处理日期格式；编辑时密码留空则不提交
    const { password, ...rest } = formValues as Record<string, any>;
    const submitData: any = {
      ...rest,
      expireDate: formValues.expireDate ? formValues.expireDate.format('YYYY-MM-DD') : '',
    };
    if (!isEdit.value || (password !== undefined && String(password).trim() !== '')) {
      submitData.password = password;
    }

    if (isEdit.value) {
      // 编辑用户
      await updateUserInfoApi(currentUserId.value, submitData);
      message.success('编辑用户成功');
    } else {
      // 新增用户
      await createUserApi(submitData);
      message.success('新增用户成功');
    }

    drawerApi.close();
    gridApi.reload();
  } catch (error) {
    const fallback = isEdit.value ? '编辑用户失败' : '新增用户失败';
    message.error(formatSubmitError(error, fallback));
    console.error('操作失败', formatSubmitError(error, fallback), error);
  } finally {
    submitLoading.value = false;
  }
}

function handleDelete(record: UserRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除用户 "${record.username}" 吗？`,
    onOk: async () => {
      try {
        await deleteUserApi(record.id);
        message.success(`删除用户 ${record.username} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除用户失败');
      }
    },
  });
}

function handleResetPassword(record: UserRecord) {
  Modal.confirm({
    title: '确认重置密码',
    content: `确定要重置用户 "${record.username}" 的密码吗？重置后系统将通知该用户邮箱，并在此展示新密码。`,
    async onOk() {
      try {
        const data = await resetUserPasswordApi(record.id, {
          nickName: record.nickName,
          email: record.email,
          phone: record.phone,
          status: record.status,
          expireDate: record.expireDate,
          role: record.role,
        });
        if (!data?.newPassword) {
          message.warning('未获取到新密码，请确认后端已更新并包含重置密码逻辑');
          return;
        }
        const emailTip = data.email?.trim()
          ? `系统已向该用户邮箱「${data.email}」发送通知邮件，请提醒客户查收。`
          : '该用户未填写邮箱，请通过其他方式将新密码告知客户。';
        Modal.success({
          title: '重置密码成功',
          width: 520,
          okText: '知道了',
          content: h('div', { style: { textAlign: 'left' } }, [
            h('p', { style: { marginBottom: '12px', lineHeight: '1.6' } }, emailTip),
            h('p', { style: { marginBottom: '8px', color: 'rgba(0,0,0,0.65)' } }, '新登录密码：'),
            h(
              'pre',
              {
                style: {
                  margin: '0 0 12px',
                  padding: '10px 12px',
                  background: 'rgba(0,0,0,0.04)',
                  borderRadius: '6px',
                  userSelect: 'all',
                  wordBreak: 'break-all',
                  fontSize: '14px',
                  lineHeight: '1.5',
                },
              },
              data.newPassword,
            ),
            h(
              Button,
              {
                size: 'small',
                type: 'primary',
                onClick: () => {
                  void navigator.clipboard
                    .writeText(data.newPassword)
                    .then(() => message.success('已复制到剪贴板'))
                    .catch(() => message.warning('复制失败，请手动选中上方密码复制'));
                },
              },
              { default: () => '复制密码' },
            ),
          ]),
        });
      } catch (error) {
        message.error(formatSubmitError(error, '重置密码失败'));
        return Promise.reject(error);
      }
    },
  });
}

async function handleStatusChange(record: UserRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newStatus = isChecked ? 1 : 0;
  const statusText = newStatus === 1 ? '启用' : '禁用';
  const oldStatus = record.status;

  Modal.confirm({
    title: '确认操作',
    content: `确定要${statusText}用户 "${record.username}" 吗？`,
    onOk: async () => {
      try {
        // 调用 API 更新状态
        await updateUserStatusApi(record.id, newStatus);
        record.status = newStatus;
        message.success(`${statusText}用户成功`);
        // 刷新列表
        gridApi.reload();
      } catch (error) {
        message.error(`${statusText}用户失败`);
        // 恢复原状态
        record.status = oldStatus;
      }
    },
  });
}

async function handleEnterBackend(record: UserRecord) {
  Modal.confirm({
    title: '进入该用户后台',
    content: `将以「${record.username}」的身份重新登录后台，当前会话将被替换。是否继续？`,
    okText: '确定进入',
    async onOk() {
      try {
        await authStore.enterUserBackend(String(record.id));
        message.success('已切换到该用户后台');
      } catch (error) {
        message.error(formatSubmitError(error, '进入失败，请确认已为超级管理员且后端已更新部署'));
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的用户');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 个用户吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: UserRecord) => record.id);
        await batchDeleteUserApi(ids);
        message.success(`批量删除 ${selectRecords.length} 个用户成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除用户失败');
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
          新增用户
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <!-- 角色列 -->
      <template #role="{ row }">
        <Tag :color="roleTagColor(row.role)">
          {{ roleDisplayName(row.role) }}
        </Tag>
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
          <Button
            v-if="canEnterUserBackend"
            type="link"
            size="small"
            @click="handleEnterBackend(row)"
          >
            <template #icon><span class="i-ant-design:login-outlined" /></template>
            进入后台
          </Button>
          <Button type="link" size="small" @click="handleResetPassword(row)">
            <template #icon><span class="i-ant-design:key-outlined" /></template>
            重置密码
          </Button>
          <Button type="link" size="small" danger @click="handleDelete(row)">
            <template #icon><span class="i-ant-design:delete-outlined" /></template>
            删除
          </Button>
        </div>
      </template>
    </Grid>

    <!-- 用户表单 Drawer -->
    <Drawer :title="drawerTitle" class="w-150">
      <UserForm />
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
