import { requestClient } from '#/api/request';

export namespace PermissionApi {
  export interface PermissionListParams {
    page: number;
    pageSize: number;
    name?: string;
    code?: string;
    module?: string;
  }

  export interface PermissionRecord {
    id: string;
    name: string;
    code: string;
    module: string;
    isDelegable: boolean;
    createTime: string;
    updateTime: string;
  }

  export interface PermissionListResult {
    items: PermissionRecord[];
    total: number;
  }
}

export async function getPermissionListApi(params: PermissionApi.PermissionListParams) {
  return requestClient.post<PermissionApi.PermissionListResult>('/permission/list', params);
}

export async function createPermissionApi(data: Partial<PermissionApi.PermissionRecord>) {
  return requestClient.post('/permission/create', data);
}

export async function updatePermissionApi(id: string, data: Partial<PermissionApi.PermissionRecord>) {
  return requestClient.post('/permission/update', { id: String(id), ...data });
}

export async function deletePermissionApi(id: string) {
  return requestClient.post('/permission/delete', { id: String(id) });
}

export async function batchDeletePermissionApi(ids: string[]) {
  return requestClient.post('/permission/batch-delete', { ids: ids.map((id) => String(id)) });
}
