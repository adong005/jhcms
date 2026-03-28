import { requestClient } from '#/api/request';

export namespace CategoryApi {
  export interface CategoryListParams {
    page: number;
    pageSize: number;
    name?: string;
    status?: number | string;
  }

  export interface CategoryRecord {
    id: string;
    name: string;
    code?: string;
    isHome?: number;
    tenantId?: string | number;
    /** 归属展示：后端按租户解析的会员昵称等 */
    ownerNickName?: string;
    sort: number;
    description: string;
    status: number;
    createTime: string;
    updateTime: string;
  }

  export interface CategoryListResult {
    items: CategoryRecord[];
    total: number;
  }
}

export async function getCategoryListApi(params: CategoryApi.CategoryListParams) {
  return requestClient.post<CategoryApi.CategoryListResult>('/info/category/list', params);
}

export async function updateCategoryStatusApi(id: string, status: number) {
  return requestClient.post('/info/category/status', { id: String(id), status });
}

export async function updateCategoryInfoApi(id: string, data: Partial<CategoryApi.CategoryRecord>) {
  return requestClient.post('/info/category/update', { id: String(id), ...data });
}

export async function createCategoryApi(data: Partial<CategoryApi.CategoryRecord>) {
  return requestClient.post('/info/category/create', data);
}

export async function deleteCategoryApi(id: string) {
  return requestClient.post('/info/category/delete', { id: String(id) });
}

export async function batchDeleteCategoryApi(ids: string[]) {
  return requestClient.post('/info/category/batch-delete', { ids: ids.map((id) => String(id)) });
}
