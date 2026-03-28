# Vben Admin 组件开发指南

## 组件概述

Vben Admin 提供了丰富的组件库，基于 Tailwind CSS 实现，可适用于不同 UI 组件库（Ant Design Vue、Element Plus、Naive UI、TDesign）。

### 组件分类

1. **布局组件** - 页面内容区域的顶层容器组件
2. **通用组件** - 常用业务组件（表单、表格、弹窗等）

## 一、布局组件

### Page 页面组件

页面内容区最常用的顶层布局容器，内置标题区、内容区和底部区三部分结构。

#### 基础用法

```vue
<template>
  <Page title="页面标题" description="页面描述">
    <!-- 页面内容 -->
  </Page>
</template>
```

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `title` | `string \| slot` | - | 页面标题 |
| `description` | `string \| slot` | - | 页面描述 |
| `extra` | `string \| slot` | - | 标题栏右侧额外内容 |
| `contentClass` | `string` | - | 内容区域类名 |
| `contentBackground` | `string` | - | 内容区域背景色 |
| `autoContentHeight` | `boolean` | `false` | 自动计算内容高度 |
| `contentPadding` | `number` | `0` | 内容区域内边距（px） |

#### Slots

- `default` - 默认内容区域
- `title` - 自定义标题
- `description` - 自定义描述
- `extra` - 自定义标题栏右侧内容
- `footer` - 底部区域

#### 注意事项

⚠️ 如果 `title`、`description`、`extra` 都没有提供内容，头部区域不会渲染

## 二、通用组件

### 1. Form 表单组件

基于适配器模式的表单组件，支持多种 UI 库。

#### 适配器配置

需要在 `#/adapter/form` 中配置适配器，注册表单组件。

#### 基础用法

```vue
<script setup lang="ts">
import { useVbenForm } from '#/adapter/form';

const [Form, formApi] = useVbenForm({
  // 所有表单项共用配置
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  // 提交函数
  handleSubmit: onSubmit,
  // 布局：horizontal（水平）| vertical（垂直）
  layout: 'horizontal',
  // 表单项配置
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
      component: 'InputPassword',
      componentProps: {
        placeholder: '请输入密码',
      },
      fieldName: 'password',
      label: '密码',
    },
  ],
  wrapperClass: 'grid-cols-1',
});

function onSubmit(values: Record<string, any>) {
  console.log('表单值:', values);
}
</script>

<template>
  <Form />
</template>
```

#### 支持的组件类型

- `Input` - 输入框
- `InputPassword` - 密码框
- `InputNumber` - 数字输入框
- `Select` - 下拉选择
- `RadioGroup` - 单选组
- `Radio` - 单选框
- `CheckboxGroup` - 多选组
- `Checkbox` - 复选框
- `DatePicker` - 日期选择器
- `RangePicker` - 范围选择器
- `TimePicker` - 时间选择器
- `TreeSelect` - 树选择
- `Switch` - 开关
- `Rate` - 评分
- `Mentions` - 提及

#### 表单校验

```typescript
schema: [
  {
    component: 'Input',
    fieldName: 'username',
    label: '用户名',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'email',
    label: '邮箱',
    rules: 'required|email',
  },
]
```

#### 表单联动

通过 `dependencies` 实现表单项之间的联动。

#### FormApi 方法

- `setValues(values)` - 设置表单值
- `getValues()` - 获取表单值
- `validate()` - 验证表单
- `reset()` - 重置表单
- `submit()` - 提交表单

### 2. Modal 模态框组件

可拖拽、自动计算高度的模态框组件。

#### 基础用法

```vue
<script setup lang="ts">
import { useVbenModal, VbenButton } from '@vben/common-ui';

const [Modal, modalApi] = useVbenModal();
</script>

<template>
  <VbenButton @click="modalApi.open()">打开弹窗</VbenButton>
  <Modal class="w-150" title="基础示例">
    弹窗内容
  </Modal>
</template>
```

#### 组件抽离

通过 `connectedComponent` 将弹窗内容抽离为独立组件：

```vue
<script setup lang="ts">
import { useVbenModal } from '@vben/common-ui';
import ExtraModal from './modal.vue';

const [Modal, modalApi] = useVbenModal({
  connectedComponent: ExtraModal,
});
</script>

<template>
  <Modal />
</template>
```

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `title` | `string \| slot` | - | 弹窗标题 |
| `description` | `string \| slot` | - | 弹窗描述 |
| `draggable` | `boolean` | `false` | 是否可拖拽 |
| `loading` | `boolean` | `false` | 加载状态 |
| `closable` | `boolean` | `true` | 是否显示关闭按钮 |
| `closeOnClickModal` | `boolean` | `true` | 点击遮罩是否关闭 |
| `closeOnPressEscape` | `boolean` | `true` | 按ESC是否关闭 |
| `showConfirmButton` | `boolean` | `true` | 显示确认按钮 |
| `showCancelButton` | `boolean` | `true` | 显示取消按钮 |
| `confirmText` | `string` | `'确认'` | 确认按钮文本 |
| `cancelText` | `string` | `'取消'` | 取消按钮文本 |
| `fullscreen` | `boolean` | `false` | 全屏显示 |
| `appendToMain` | `boolean` | `false` | 挂载到内容区域 |
| `zIndex` | `number` | `1000` | 层级 |
| `animationType` | `'slide' \| 'scale'` | `'slide'` | 动画类型 |

#### Events

- `onCancel` - 取消回调，返回 `false` 可阻止关闭
- `onConfirm` - 确认回调
- `onBeforeClose` - 关闭前回调
- `onClosed` - 关闭后回调
- `onOpened` - 打开后回调
- `onOpenChange` - 打开状态变化

#### ModalApi 方法

- `open()` - 打开弹窗
- `close()` - 关闭弹窗
- `setState(state)` - 设置状态
- `setData(data)` - 设置数据
- `getData()` - 获取数据
- `lock(isLock)` - 锁定弹窗（防止重复提交）
- `unlock()` - 解锁弹窗

#### Slots

- `default` - 默认内容
- `title` - 自定义标题
- `description` - 自定义描述
- `extra` - 标题栏右侧内容
- `prepend-footer` - 底部前置内容
- `footer` - 自定义底部

### 3. Drawer 抽屉组件

从屏幕边缘滑出的抽屉组件，API 与 Modal 类似。

#### 基础用法

```vue
<script setup lang="ts">
import { useVbenDrawer, VbenButton } from '@vben/common-ui';

const [Drawer, drawerApi] = useVbenDrawer({
  title: '抽屉标题',
});
</script>

<template>
  <VbenButton @click="drawerApi.open()">打开抽屉</VbenButton>
  <Drawer>
    抽屉内容
  </Drawer>
</template>
```

#### Props

与 Modal 基本一致，额外支持：
- `placement` - 抽屉方向（`left` | `right` | `top` | `bottom`）

### 4. Vxe Table 表格组件

基于 vxe-table 的高性能表格组件。

#### 基础用法

```vue
<script setup lang="ts">
import { useVbenVxeGrid } from '#/adapter/vxe-table';

const [Grid] = useVbenVxeGrid({
  columns: [
    { field: 'name', title: '姓名' },
    { field: 'age', title: '年龄' },
    { field: 'address', title: '地址' },
  ],
  data: [
    { name: '张三', age: 18, address: '北京' },
    { name: '李四', age: 20, address: '上海' },
  ],
});
</script>

<template>
  <Grid />
</template>
```

#### 远程加载

```typescript
const [Grid] = useVbenVxeGrid({
  columns: [...],
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const res = await fetchData({
          page: page.currentPage,
          pageSize: page.pageSize,
        });
        return {
          records: res.data,
          total: res.total,
        };
      },
    },
  },
});
```

#### 特性

- ✅ 虚拟滚动（支持大数据量）
- ✅ 树形表格
- ✅ 固定列
- ✅ 单元格编辑
- ✅ 行编辑
- ✅ 自定义单元格
- ✅ 搜索表单集成

### 5. ApiComponent API组件包装器

用于包装需要异步加载数据的组件。

#### 基础用法

```vue
<template>
  <ApiComponent :api="fetchUserInfo">
    <template #default="{ data }">
      <div>用户名: {{ data.name }}</div>
    </template>
  </ApiComponent>
</template>

<script setup lang="ts">
import { ApiComponent } from '@vben/common-ui';

async function fetchUserInfo() {
  const res = await fetch('/api/user/info');
  return res.json();
}
</script>
```

#### Props

| 属性 | 类型 | 说明 |
|------|------|------|
| `api` | `() => Promise<T>` | 异步数据加载函数 |
| `immediate` | `boolean` | 是否立即加载 |
| `cache` | `boolean` | 是否缓存结果 |

#### Methods

- `reload()` - 重新加载数据

### 6. Alert 轻量提示框

轻量级的提示框组件。

#### 基础用法

```vue
<template>
  <VbenAlert type="success" title="成功提示" />
  <VbenAlert type="warning" title="警告提示" />
  <VbenAlert type="error" title="错误提示" />
  <VbenAlert type="info" title="信息提示" />
</template>
```

#### useAlertContext

在组件内部使用 `useAlertContext` 获取上下文：

```typescript
import { useAlertContext } from '@vben/common-ui';

const { close } = useAlertContext();
```

### 7. CountToAnimator 数字动画

数字滚动动画组件。

#### 基础用法

```vue
<template>
  <CountToAnimator :end-value="1000" :duration="2000" />
</template>
```

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `startValue` | `number` | `0` | 起始值 |
| `endValue` | `number` | - | 结束值 |
| `duration` | `number` | `2000` | 动画时长（ms） |
| `prefix` | `string` | - | 前缀 |
| `suffix` | `string` | - | 后缀 |
| `separator` | `string` | `,` | 千分位分隔符 |
| `decimals` | `number` | `0` | 小数位数 |

#### Events

- `onStarted` - 动画开始
- `onFinished` - 动画结束

### 8. EllipsisText 省略文本

支持省略和展开的文本组件。

#### 基础用法

```vue
<template>
  <EllipsisText :line="2" tooltip>
    这是一段很长的文本内容...
  </EllipsisText>
</template>
```

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `line` | `number` | `1` | 最大行数 |
| `tooltip` | `boolean` | `false` | 是否显示 tooltip |
| `expand` | `boolean` | `false` | 是否可展开 |
| `expandText` | `string` | `'展开'` | 展开按钮文本 |
| `collapseText` | `string` | `'收起'` | 收起按钮文本 |

## 三、开发规范

### 1. 组件引入

```typescript
// 从适配器引入（推荐）
import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';

// 从通用组件库引入
import { useVbenModal, VbenButton } from '@vben/common-ui';
```

### 2. 适配器模式

框架使用适配器模式适配不同 UI 库，需要在 `#/adapter` 目录下配置：

- `#/adapter/form` - 表单适配器
- `#/adapter/vxe-table` - 表格适配器
- `#/adapter/component` - 组件适配器

### 3. 样式使用

组件基于 Tailwind CSS 实现，推荐使用 Tailwind 类名：

```vue
<template>
  <Modal class="w-150 h-100">
    <div class="p-4 space-y-4">
      <!-- 内容 -->
    </div>
  </Modal>
</template>
```

### 4. 类型安全

所有组件都提供完整的 TypeScript 类型定义：

```typescript
import type { VbenFormSchema, FormApi } from '@vben/common-ui';

const schema: VbenFormSchema[] = [
  // ...
];
```

## 四、最佳实践

### 1. 表单开发

```vue
<script setup lang="ts">
import { useVbenForm } from '#/adapter/form';

const [Form, formApi] = useVbenForm({
  schema: [
    {
      component: 'Input',
      fieldName: 'username',
      label: '用户名',
      rules: 'required',
    },
  ],
  handleSubmit: async (values) => {
    try {
      await submitData(values);
      message.success('提交成功');
    } catch (error) {
      message.error('提交失败');
    }
  },
});
</script>
```

### 2. 弹窗表单组合

```vue
<script setup lang="ts">
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    await formApi.validate();
    const values = formApi.getValues();
    await submitData(values);
    modalApi.close();
  },
});

const [Form, formApi] = useVbenForm({
  schema: [/* ... */],
});
</script>

<template>
  <Modal title="编辑">
    <Form />
  </Modal>
</template>
```

### 3. 表格 CRUD

```vue
<script setup lang="ts">
import { useVbenVxeGrid } from '#/adapter/vxe-table';

const [Grid, gridApi] = useVbenVxeGrid({
  columns: [
    { field: 'name', title: '姓名' },
    {
      field: 'actions',
      title: '操作',
      slots: { default: 'actions' },
    },
  ],
  proxyConfig: {
    ajax: {
      query: fetchList,
    },
  },
});

async function handleEdit(row) {
  // 打开编辑弹窗
}

async function handleDelete(row) {
  await deleteData(row.id);
  gridApi.reload();
}
</script>
```

## 五、注意事项

1. **组件注册**: 使用的组件必须在适配器中注册
2. **缓存清理**: 更改配置后需清空浏览器缓存
3. **类型导入**: 优先使用 `#/` 别名导入适配器
4. **响应式**: 组件内部已处理响应式，无需额外包装
5. **性能优化**: 大数据量表格使用虚拟滚动

## 六、相关链接

- 组件文档：https://doc.vben.pro/components/introduction.html
- Tailwind CSS：https://tailwindcss.com
- Vxe Table：https://vxetable.cn

## 七、自定义组件

如果现有组件不满足需求，可以：

1. 直接使用原生 UI 组件库组件
2. 基于现有组件二次封装
3. 完全自定义新组件

框架提供的组件并非强制使用，完全取决于业务需求。
