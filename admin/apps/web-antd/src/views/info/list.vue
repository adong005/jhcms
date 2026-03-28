<script setup lang="ts">
import type { VxeGridListeners } from '#/adapter/vxe-table';
import { computed, nextTick, onBeforeUnmount, ref, shallowRef, watch } from 'vue';
import { Button, message, Modal, Spin, Switch } from 'ant-design-vue';
// @ts-ignore - 库的类型声明配置问题，与 publish 页一致
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor';
import '@wangeditor/editor/dist/css/style.css';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useVbenDrawer } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getCategoryListApi } from '#/api/info/category';
import {
  getInfoListApi,
  getInfoDetailApi,
  updateInfoApi,
  updateInfoStatusApi,
  deleteInfoApi,
  batchDeleteInfoApi,
} from '#/api/info/list';
import type { InfoApi } from '#/api/info/list';

type InfoRecord = InfoApi.InfoRecord;

const currentInfoId = ref('');
const isEditMode = ref(false);
const drawerDetailLoading = ref(false);
const drawerTitle = computed(() => {
  if (!currentInfoId.value) return '新建信息';
  return isEditMode.value ? '编辑信息' : '查看信息';
});

const contentValue = ref('');
const editorRef = shallowRef<IDomEditor>();

const editorConfig: Partial<IEditorConfig> = {
  placeholder: '请输入内容...',
  MENU_CONF: {},
};

const toolbarConfig: Partial<IToolbarConfig> = {
  toolbarKeys: [
    'headerSelect',
    'bold',
    'italic',
    'underline',
    'through',
    '|',
    'color',
    'bgColor',
    '|',
    'fontSize',
    'fontFamily',
    'lineHeight',
    '|',
    'bulletedList',
    'numberedList',
    'todo',
    '|',
    'justifyLeft',
    'justifyCenter',
    'justifyRight',
    'justifyJustify',
    '|',
    'insertLink',
    'insertImage',
    'insertTable',
    'codeBlock',
    'divider',
    '|',
    'undo',
    'redo',
    '|',
    'fullScreen',
  ],
};

function handleEditorCreated(editor: IDomEditor) {
  editorRef.value = editor;
  applyEditorReadonly();
}

function applyEditorReadonly() {
  const ed = editorRef.value;
  if (!ed) return;
  if (isEditMode.value) ed.enable();
  else ed.disable();
}

watch(isEditMode, () => {
  void nextTick(() => applyEditorReadonly());
});

const [InfoForm, formApi] = useVbenForm({
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
        placeholder: '请输入标题',
      },
      fieldName: 'title',
      label: '标题',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        placeholder: '加载分类中…',
        options: [],
      },
      fieldName: 'categoryId',
      label: '分类',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入作者',
      },
      fieldName: 'author',
      label: '作者',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入摘要',
        rows: 3,
      },
      fieldName: 'summary',
      label: '摘要',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '已发布', value: 1 },
          { label: '草稿', value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
  ],
  showDefaultActions: false,
});

const [Drawer, drawerApi] = useVbenDrawer({
  onOpenChange: (isOpen) => {
    if (!isOpen) {
      const ed = editorRef.value;
      if (ed) {
        try {
          ed.destroy();
        } catch {
          /* ignore */
        }
        editorRef.value = undefined as unknown as IDomEditor;
      }
      contentValue.value = '';
      formApi.resetForm();
      currentInfoId.value = '';
    }
  },
});

onBeforeUnmount(() => {
  const ed = editorRef.value;
  if (ed) {
    try {
      ed.destroy();
    } catch {
      /* ignore */
    }
  }
});

const drawerCategoryOptions = ref<{ label: string; value: string }[]>([]);

async function loadCategoryOptionsForDrawer() {
  const resp = await getCategoryListApi({ page: 1, pageSize: 500 });
  const opts = (resp.items || [])
    .filter((c) => Number(c.status) === 1)
    .sort((a, b) => (a.sort ?? 0) - (b.sort ?? 0))
    .map((c) => ({
      label: c.code ? `${c.name}（${c.code}）` : c.name,
      value: String(c.id),
    }));
  drawerCategoryOptions.value = opts;
  formApi.updateSchema([
    {
      fieldName: 'categoryId',
      componentProps: {
        placeholder: opts.length ? '请选择分类' : '暂无可用分类',
        options: opts,
        allowClear: true,
      },
    },
  ]);
}

/** 当前资讯的分类若已停用或不在列表里，补一条 option，否则 Select 无法显示默认值 */
function ensureDrawerCategoryOption(rawId: unknown, label: string) {
  if (rawId === undefined || rawId === null || String(rawId).trim() === '') return;
  const value = String(rawId);
  if (drawerCategoryOptions.value.some((o) => o.value === value)) return;
  const text = label.trim() || `分类 #${value}`;
  const next = [...drawerCategoryOptions.value, { label: text, value }];
  drawerCategoryOptions.value = next;
  formApi.updateSchema([
    {
      fieldName: 'categoryId',
      componentProps: {
        placeholder: next.length ? '请选择分类' : '暂无可用分类',
        options: next,
        allowClear: true,
      },
    },
  ]);
}

function setDrawerFormDisabled(disabled: boolean) {
  const opts = drawerCategoryOptions.value;
  formApi.updateSchema([
    { fieldName: 'title', componentProps: { disabled } },
    {
      fieldName: 'categoryId',
      componentProps: {
        placeholder: opts.length ? '请选择分类' : '暂无可用分类',
        options: opts,
        allowClear: true,
        disabled,
      },
    },
    { fieldName: 'author', componentProps: { disabled } },
    { fieldName: 'summary', componentProps: { disabled } },
    { fieldName: 'status', componentProps: { disabled } },
  ]);
}

function isEmptyHtml(html: string) {
  const t = html.replace(/\s/g, '');
  return !t || t === '<p><br></p>' || t === '<p></p>';
}

async function openInfoDrawer(record: InfoRecord, editMode: boolean) {
  contentValue.value = '';
  currentInfoId.value = record.id;
  isEditMode.value = editMode;
  drawerApi.open();
  drawerDetailLoading.value = true;
  try {
    await loadCategoryOptionsForDrawer();
    const detail = await getInfoDetailApi(record.id);
    const catRaw = detail.categoryId ?? record.categoryId;
    const catLabel = (detail.categoryName || record.categoryName || '').trim();
    ensureDrawerCategoryOption(catRaw, catLabel);
    setDrawerFormDisabled(!editMode);
    await nextTick();
    formApi.setValues({
      title: detail.title ?? '',
      categoryId:
        catRaw !== undefined && catRaw !== null && String(catRaw) !== '' ? String(catRaw) : undefined,
      author: detail.author ?? '',
      summary: detail.summary ?? '',
      status: Number(detail.status),
    });
    contentValue.value = detail.content ?? '';
    await nextTick();
    applyEditorReadonly();
  } catch {
    message.error('加载信息详情失败');
    drawerApi.close();
  } finally {
    drawerDetailLoading.value = false;
  }
}

const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入标题',
      },
      fieldName: 'title',
      label: '标题',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择分类',
        options: [
          { label: '全部', value: '' },
          { label: '新闻资讯', value: '1' },
          { label: '技术文章', value: '2' },
          { label: '产品介绍', value: '3' },
          { label: '公司动态', value: '4' },
          { label: '行业资讯', value: '5' },
        ],
      },
      fieldName: 'categoryId',
      label: '分类',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        placeholder: '请选择状态',
        options: [
          { label: '全部', value: '' },
          { label: '已发布', value: 1 },
          { label: '草稿', value: 0 },
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
      field: 'title',
      title: '标题',
      width: 200,
      fixed: 'left',
    },
    {
      field: 'categoryName',
      title: '分类',
      width: 120,
    },
    {
      field: 'author',
      title: '作者',
      width: 100,
    },
    {
      field: 'ownerNickName',
      title: '归属',
      width: 160,
      showOverflow: true,
    },
    {
      field: 'viewCount',
      title: '浏览量',
      width: 100,
    },
    {
      field: 'summary',
      title: '摘要',
      minWidth: 200,
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
      field: 'action',
      fixed: 'right',
      title: '操作',
      width: 180,
      slots: { default: 'action' },
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
    sheetName: '信息列表',
    mode: 'current',
  },
  exportConfig: {
    filename: '信息列表',
    type: 'xlsx',
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        const resp = await getInfoListApi({
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

const gridEvents: VxeGridListeners<InfoRecord> = {};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents,
  gridOptions,
});

function handleView(record: InfoRecord) {
  void openInfoDrawer(record, false);
}

function handleEdit(record: InfoRecord) {
  void openInfoDrawer(record, true);
}

function handleDelete(record: InfoRecord) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除信息 "${record.title}" 吗？`,
    onOk: async () => {
      try {
        await deleteInfoApi(record.id);
        message.success(`删除信息 ${record.title} 成功`);
        gridApi.reload();
      } catch (error) {
        message.error('删除信息失败');
      }
    },
  });
}

async function handleStatusChange(record: InfoRecord, checked: any) {
  const isChecked = checked === true || checked === 1 || checked === '1' || checked === 'true';
  const newStatus = isChecked ? 1 : 0;
  const statusText = newStatus === 1 ? '发布' : '撤回';
  const oldStatus = record.status;

  Modal.confirm({
    title: '确认操作',
    content: `确定要${statusText}信息 "${record.title}" 吗？`,
    onOk: async () => {
      try {
        await updateInfoStatusApi(record.id, newStatus);
        record.status = newStatus;
        message.success(`${statusText}信息成功`);
        gridApi.reload();
      } catch (error) {
        message.error(`${statusText}信息失败`);
        record.status = oldStatus;
      }
    },
  });
}

function handleBatchDelete() {
  const selectRecords = gridApi.grid.getCheckboxRecords();
  if (selectRecords.length === 0) {
    message.warning('请先选择要删除的信息');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectRecords.length} 条信息吗？`,
    onOk: async () => {
      try {
        const ids = selectRecords.map((record: InfoRecord) => record.id);
        await batchDeleteInfoApi(ids);
        message.success(`批量删除 ${selectRecords.length} 条信息成功`);
        gridApi.reload();
      } catch (error) {
        message.error('批量删除信息失败');
      }
    },
  });
}

async function handleSave() {
  try {
    await formApi.validate();
    if (isEmptyHtml(contentValue.value)) {
      message.warning('请输入内容');
      return;
    }
    const v = await formApi.getValues();
    const rawCat = v.categoryId;
    let categoryId: number | undefined;
    if (rawCat !== '' && rawCat !== undefined && rawCat !== null) {
      const n = Number(rawCat);
      if (Number.isFinite(n)) categoryId = n;
    }
    await updateInfoApi(currentInfoId.value, {
      title: String(v.title ?? ''),
      content: contentValue.value,
      categoryId,
      status: Number(v.status),
      author: String(v.author ?? ''),
      summary: String(v.summary ?? ''),
    });
    message.success('保存成功');
    drawerApi.close();
    gridApi.reload();
  } catch {
    message.error('保存失败');
  }
}

function handlePublish() {
  window.open('/info/publish', '_blank');
}

</script>

<template>
  <div class="p-4">
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="handlePublish">
          <template #icon><span class="i-ant-design:plus-outlined" /></template>
          发布信息
        </Button>
        <Button danger @click="handleBatchDelete">
          <template #icon><span class="i-ant-design:delete-outlined" /></template>
          批量删除
        </Button>
      </template>

      <template #status="{ row }">
        <Switch :checked="row.status === 1" checked-children="已发布" un-checked-children="草稿"
          @change="(checked) => handleStatusChange(row, checked)" />
      </template>

      <template #action="{ row }">
        <div class="flex gap-2">
          <Button type="link" size="small" @click="handleView(row)">
            <template #icon><span class="i-ant-design:eye-outlined" /></template>
            查看
          </Button>
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

    <Drawer :title="drawerTitle" class="info-list-drawer">
      <Spin :spinning="drawerDetailLoading">
        <template v-if="currentInfoId">
          <InfoForm />
          <div class="mt-4">
            <div class="mb-2 text-sm font-medium">
              <span v-if="isEditMode" class="text-red-500">*</span>
              内容
            </div>
            <div class="editor-container">
              <Toolbar
                v-if="isEditMode"
                :editor="editorRef"
                :default-config="toolbarConfig"
                mode="default"
                class="editor-toolbar"
              />
              <Editor
                :key="currentInfoId"
                v-model="contentValue"
                :default-config="editorConfig"
                mode="default"
                class="editor-content"
                @on-created="handleEditorCreated"
              />
            </div>
          </div>
        </template>
      </Spin>
      <template #footer>
        <div class="flex justify-end gap-4">
          <Button @click="drawerApi.close()">{{ isEditMode ? '取消' : '关闭' }}</Button>
          <Button v-if="isEditMode" type="primary" @click="handleSave">保存</Button>
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

.info-list-drawer {
  width: min(92vw, 880px) !important;
  max-width: 880px;
}

.editor-container {
  border: 1px solid var(--vben-color-border, #d9d9d9);
  border-radius: var(--vben-border-radius, 4px);
  overflow: hidden;
  background-color: transparent;
  box-shadow:
    0 1px 2px 0 rgba(0, 0, 0, 0.03),
    0 1px 6px -1px rgba(0, 0, 0, 0.02),
    0 2px 4px 0 rgba(0, 0, 0, 0.02);
}

.editor-toolbar {
  border-bottom: 1px solid var(--vben-color-border, #d9d9d9);
}

.editor-content {
  min-height: 360px;
  overflow-y: auto;
}

:deep(.w-e-text-container) {
  min-height: 360px !important;
  background-color: transparent !important;
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-text-container [data-slate-editor]) {
  min-height: 360px !important;
  padding: 16px !important;
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-toolbar) {
  background-color: transparent !important;
  border-bottom: 1px solid var(--vben-color-border) !important;
  flex-wrap: wrap !important;
}

:deep(.w-e-text-placeholder) {
  color: var(--vben-color-text-3) !important;
  top: 16px !important;
  left: 16px !important;
}

:deep(.w-e-bar-item button) {
  color: var(--vben-color-text-1) !important;
  transition: all 0.2s;
}

:deep(.w-e-bar-item button:hover) {
  background-color: var(--vben-color-bg-3) !important;
}

:deep(.w-e-bar-item-active button) {
  background-color: var(--vben-color-primary-1) !important;
  color: var(--vben-color-primary) !important;
}

:deep(.w-e-bar-divider) {
  border-left: 1px solid var(--vben-color-border) !important;
}

:deep(.w-e-select-list) {
  background-color: var(--vben-color-bg-1) !important;
  border: 1px solid var(--vben-color-border) !important;
  box-shadow: var(--vben-shadow-base) !important;
}

:deep(.w-e-select-list .w-e-select-list-item) {
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-select-list .w-e-select-list-item:hover) {
  background-color: var(--vben-color-bg-3) !important;
}

:deep(.w-e-modal) {
  background-color: var(--vben-color-bg-1) !important;
  border: 1px solid var(--vben-color-border) !important;
}

:deep(.w-e-modal-title) {
  color: var(--vben-color-text-1) !important;
  border-bottom: 1px solid var(--vben-color-border) !important;
}

:deep(.w-e-modal-content) {
  color: var(--vben-color-text-2) !important;
}

:deep(.w-e-modal input),
:deep(.w-e-modal textarea) {
  background-color: var(--vben-color-bg-1) !important;
  border-color: var(--vben-color-border) !important;
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-text-container table) {
  border-color: var(--vben-color-border) !important;
}

:deep(.w-e-text-container table td),
:deep(.w-e-text-container table th) {
  border-color: var(--vben-color-border) !important;
}

:deep(.w-e-text-container pre) {
  background-color: var(--vben-color-bg-2) !important;
  border: 1px solid var(--vben-color-border) !important;
}

:deep(.w-e-text-container code) {
  background-color: var(--vben-color-bg-2) !important;
  color: var(--vben-color-primary) !important;
}

:deep(.w-e-text-container blockquote) {
  border-left-color: var(--vben-color-primary) !important;
  background-color: var(--vben-color-primary-1) !important;
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-text-container a) {
  color: var(--vben-color-primary) !important;
}

:deep(.w-e-text-container a:hover) {
  color: var(--vben-color-primary-hover) !important;
}

:deep(.w-e-text-container hr) {
  border-color: var(--vben-color-border) !important;
}

:deep(.w-e-text-container::-webkit-scrollbar) {
  width: 8px;
  height: 8px;
}

:deep(.w-e-text-container::-webkit-scrollbar-thumb) {
  background-color: var(--vben-color-border);
  border-radius: 4px;
}

:deep(.w-e-text-container::-webkit-scrollbar-thumb:hover) {
  background-color: var(--vben-color-text-3);
}

:deep(.w-e-text-container::-webkit-scrollbar-track) {
  background-color: var(--vben-color-bg-2);
}
</style>
