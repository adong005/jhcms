import { requestClient } from '#/api/request';

export namespace RoleApi {
  export interface RoleListParams {
    page: number;
    pageSize: number;
    name?: string;
    code?: string;
    status?: number | string;
  }

  export interface RoleRecord {
    id: string;
    name: string;
    code: string;
    description: string;
    status: number;
    createTime: string;
    updateTime: string;
  }

  export interface RoleListResult {
    items: RoleRecord[];
    total: number;
  }
}

/**
 * 获取角色列表
 */
export async function getRoleListApi(params: RoleApi.RoleListParams) {
  return requestClient.post<RoleApi.RoleListResult>('/role/list', params);
}

/**
 * 更新角色状态
 */
export async function updateRoleStatusApi(id: string, status: number) {
  return requestClient.post('/role/status', { id: String(id), status });
}

/**
 * 更新角色信息
 */
export async function updateRoleInfoApi(id: string, data: Partial<RoleApi.RoleRecord>) {
  return requestClient.post('/role/update', { id: String(id), ...data });
}

/**
 * 创建角色
 */
export async function createRoleApi(data: Partial<RoleApi.RoleRecord>) {
  return requestClient.post('/role/create', data);
}

/**
 * 删除角色
 */
export async function deleteRoleApi(id: string) {
  return requestClient.post('/role/delete', { id: String(id) });
}

/**
 * 批量删除角色
 */
export async function batchDeleteRoleApi(ids: string[]) {
  return requestClient.post('/role/batch-delete', { ids: ids.map((id) => String(id)) });
}

/**
 * 获取角色权限
 */
export async function getRolePermissionApi(roleId: string) {
  return requestClient.get<{ menuIds?: string[]; permissionIds?: string[] }>(
    `/role/permission/${String(roleId)}`,
  );
}

/**
 * 更新角色权限
 */
export async function updateRolePermissionApi(roleId: string, permissionIds: string[]) {
  return requestClient.post('/role/permission', {
    roleId: String(roleId),
    permissionIds: permissionIds.map((id) => String(id)),
  });
}
