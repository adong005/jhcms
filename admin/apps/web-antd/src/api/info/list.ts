import { requestClient } from '#/api/request';

export namespace InfoApi {
  export interface InfoListParams {
    page: number;
    pageSize: number;
    title?: string;
    categoryId?: string;
    status?: number | string;
  }

  export interface InfoRecord {
    id: string;
    title: string;
    categoryId: string;
    categoryName: string;
    author: string;
    tenantId?: string;
    /** 归属：后端按租户解析的会员昵称等 */
    ownerNickName?: string;
    summary: string;
    content: string;
    viewCount: number;
    status: number;
    createTime: string;
    updateTime: string;
  }

  export interface InfoListResult {
    items: InfoRecord[];
    total: number;
  }
}

export async function getInfoListApi(params: InfoApi.InfoListParams) {
  return requestClient.post<InfoApi.InfoListResult>('/info/list', params);
}

export async function updateInfoStatusApi(id: string, status: number) {
  return requestClient.post('/info/status', { id, status });
}

export async function updateInfoApi(
  id: string,
  data: {
    title: string;
    content: string;
    categoryId?: string;
    status: number;
    author?: string;
    summary?: string;
  },
) {
  return requestClient.post('/info/update', {
    id,
    title: data.title,
    content: data.content,
    categoryId: data.categoryId,
    status: data.status,
    author: data.author ?? '',
    summary: data.summary ?? '',
  });
}

export async function deleteInfoApi(id: string) {
  return requestClient.post('/info/delete', { id });
}

export async function batchDeleteInfoApi(ids: string[]) {
  return requestClient.post('/info/batch-delete', { ids });
}

export async function getInfoDetailApi(id: string) {
  return requestClient.get<InfoApi.InfoRecord>(`/info/detail/${String(id)}`);
}

export async function createInfoApi(data: Partial<InfoApi.InfoRecord>) {
  return requestClient.post('/info/create', data);
}
