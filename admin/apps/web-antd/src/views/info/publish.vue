<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, shallowRef } from 'vue';
import { Button, Card, message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';
// @ts-ignore - 库的类型声明配置问题，不影响功能
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor';
import '@wangeditor/editor/dist/css/style.css';
import { Page } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { getCategoryListApi } from '#/api/info/category';
import { createInfoApi, getInfoDetailApi, updateInfoApi } from '#/api/info/list';

const route = useRoute();
const router = useRouter();

const isEdit = ref(false);
const infoId = ref('');
const loading = ref(false);
const submitLoading = ref(false);
const contentValue = ref('');

// WangEditor 实例
const editorRef = shallowRef<IDomEditor>();

// 编辑器配置
const editorConfig: Partial<IEditorConfig> = {
  placeholder: '请输入内容...',
  MENU_CONF: {},
};

// 工具栏配置
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

// 编辑器创建回调
const handleCreated = (editor: IDomEditor) => {
  editorRef.value = editor;
};

// 组件销毁时，销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor == null) return;
  editor.destroy();
});

const [Form, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  layout: 'horizontal',
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
  ],
  showDefaultActions: false,
});

async function loadCategoryOptions() {
  try {
    const resp = await getCategoryListApi({ page: 1, pageSize: 500 });
    const opts = (resp.items || [])
      .filter((c) => Number(c.status) === 1)
      .sort((a, b) => (a.sort ?? 0) - (b.sort ?? 0))
      .map((c) => ({
        label: c.code ? `${c.name}（${c.code}）` : c.name,
        value: String(c.id),
      }));
    formApi.updateSchema([
      {
        fieldName: 'categoryId',
        componentProps: {
          placeholder: opts.length ? '请选择分类' : '暂无可用分类，请先在分类管理中添加',
          options: opts,
          allowClear: true,
        },
      },
    ]);
  } catch {
    message.error('加载分类失败');
  }
}

onMounted(async () => {
  await loadCategoryOptions();
  const id = route.query.id as string;
  if (id) {
    isEdit.value = true;
    infoId.value = id;
    await loadInfoDetail(id);
  }
});

async function loadInfoDetail(id: string) {
  try {
    loading.value = true;
    const detail = await getInfoDetailApi(id);
    formApi.setValues({
      title: detail.title,
      categoryId:
        detail.categoryId !== undefined && detail.categoryId !== null && detail.categoryId !== ''
          ? String(detail.categoryId)
          : undefined,
      author: detail.author,
      summary: detail.summary,
    });
    contentValue.value = detail.content;
  } catch (error) {
    message.error('加载信息详情失败');
    console.error('加载详情失败:', error);
  } finally {
    loading.value = false;
  }
}

async function handleSubmit(status: number) {
  try {
    await formApi.validate();

    if (!contentValue.value) {
      message.warning('请输入内容');
      return;
    }

    submitLoading.value = true;

    const formValues = await formApi.getValues();
    const rawCat = formValues.categoryId;
    let categoryId: number | undefined;
    if (rawCat !== '' && rawCat !== undefined && rawCat !== null) {
      const n = Number(rawCat);
      if (Number.isFinite(n)) categoryId = n;
    }
    const submitData = {
      ...formValues,
      categoryId,
      content: contentValue.value,
      status,
    };

    if (isEdit.value) {
      await updateInfoApi(infoId.value, {
        title: String(formValues.title ?? ''),
        content: contentValue.value,
        categoryId,
        status,
        author: String(formValues.author ?? ''),
        summary: String(formValues.summary ?? ''),
      });
      message.success(status === 1 ? '发布成功' : '保存草稿成功');
    } else {
      await createInfoApi(submitData);
      message.success(status === 1 ? '发布成功' : '保存草稿成功');
    }

    setTimeout(() => {
      router.push('/info/list');
    }, 1000);
  } catch (error) {
    if (error) {
      message.error(status === 1 ? '发布失败' : '保存失败');
      console.error('提交失败:', error);
    }
  } finally {
    submitLoading.value = false;
  }
}

function handlePublish() {
  handleSubmit(1);
}

function handleSaveDraft() {
  handleSubmit(0);
}

function handleCancel() {
  router.back();
}

const pageTitle = computed(() => isEdit.value ? '编辑信息' : '发布信息');

</script>

<template>
  <Page :loading="loading" :title="pageTitle" content-class="p-0">
    <Card :bordered="false">
      <Form />

      <div class="mt-6">
        <div class="ant-row ant-form-item">
          <div class="ant-col ant-col-4 ant-form-item-label">
            <label>
              <span class="text-red-500">*</span> 内容
            </label>
          </div>
          <div class="ant-col ant-col-18 ant-form-item-control">
            <div class="editor-container">
              <Toolbar :editor="editorRef" :defaultConfig="toolbarConfig" mode="default" class="editor-toolbar" />
              <Editor v-model="contentValue" :defaultConfig="editorConfig" mode="default" class="editor-content"
                @onCreated="handleCreated" />
            </div>
          </div>
        </div>
      </div>

      <div class="mt-8 pt-6 border-t flex justify-end gap-3">
        <Button size="large" @click="handleCancel">取消</Button>
        <Button size="large" @click="handleSaveDraft" :loading="submitLoading">
          保存草稿
        </Button>
        <Button type="primary" size="large" @click="handlePublish" :loading="submitLoading">
          发布
        </Button>
      </div>
    </Card>
  </Page>
</template>

<style scoped>
.editor-container {
  border: 1px solid var(--vben-color-border, #d9d9d9);
  border-radius: var(--vben-border-radius, 4px);
  overflow: hidden;
  background-color: transparent;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
}

.editor-toolbar {
  border-bottom: 1px solid var(--vben-color-border, #d9d9d9);
}

.editor-content {
  min-height: 500px;
  overflow-y: auto;
}

/* WangEditor 样式覆盖 - 跟随系统主题 */
:deep(.w-e-text-container) {
  min-height: 500px !important;
  background-color: transparent !important;
  color: var(--vben-color-text-1) !important;
}

:deep(.w-e-text-container [data-slate-editor]) {
  min-height: 500px !important;
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

/* 下拉菜单样式 */
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

/* 模态框样式 */
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

/* 表格样式 */
:deep(.w-e-text-container table) {
  border-color: var(--vben-color-border) !important;
}

:deep(.w-e-text-container table td),
:deep(.w-e-text-container table th) {
  border-color: var(--vben-color-border) !important;
}

/* 代码块样式 */
:deep(.w-e-text-container pre) {
  background-color: var(--vben-color-bg-2) !important;
  border: 1px solid var(--vben-color-border) !important;
}

:deep(.w-e-text-container code) {
  background-color: var(--vben-color-bg-2) !important;
  color: var(--vben-color-primary) !important;
}

/* 引用块样式 */
:deep(.w-e-text-container blockquote) {
  border-left-color: var(--vben-color-primary) !important;
  background-color: var(--vben-color-primary-1) !important;
  color: var(--vben-color-text-1) !important;
}

/* 链接样式 */
:deep(.w-e-text-container a) {
  color: var(--vben-color-primary) !important;
}

:deep(.w-e-text-container a:hover) {
  color: var(--vben-color-primary-hover) !important;
}

/* 分割线样式 */
:deep(.w-e-text-container hr) {
  border-color: var(--vben-color-border) !important;
}

/* 滚动条样式 */
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
