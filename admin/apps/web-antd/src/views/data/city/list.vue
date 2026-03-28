<script setup lang="ts">
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getCityListApi } from '#/api/site-group';

const formOptions: any = {
  collapsed: false,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: '城市',
      componentProps: {
        placeholder: '请输入城市名称或拼音',
      },
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
    keyField: 'cityCode',
  },
  columns: [
    { type: 'seq', title: '序号', width: 60 },
    { field: 'cityCode', title: '城市编码', width: 140 },
    { field: 'name', title: '城市名称', minWidth: 180 },
    { field: 'pinyin', title: '拼音', minWidth: 200 },
  ],
  pagerConfig: {
    pageSize: 20,
    pageSizes: [20, 50, 100, 200],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }: any, formValues: any) => {
        return await getCityListApi({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
      },
    },
  },
  toolbarConfig: {
    refresh: true,
    zoom: true,
    custom: true,
    search: true,
  },
};

const [Grid] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});
</script>

<template>
  <div class="p-4">
    <Grid />
  </div>
</template>
