import { requestClient } from '#/api/request';

export namespace MenuApi {
  export interface MenuRecord {
    id: string;
    name: string;
    path: string;
    type: 'menu' | 'catalog' | 'button';
    icon?: string;
    component?: string;
    parentId?: string;
    order: number;
    isShow?: number;
    status: number;
    createTime: string;
    updateTime: string;
  }

  export interface MenuListParams {
    page: number;
    pageSize: number;
    name?: string;
    type?: string;
    status?: number;
  }

  export interface MenuListResult {
    items: MenuRecord[];
    total: number;
  }
}

/**
 * 获取菜单列表
 */
export async function getMenuListApi(params: MenuApi.MenuListParams) {
  return requestClient.post<MenuApi.MenuListResult>('/menu/list', params);
}

/**
 * 创建菜单
 */
export async function createMenuApi(data: Partial<MenuApi.MenuRecord>) {
  return requestClient.post('/menu/create', data);
}

/**
 * 更新菜单信息
 */
export async function updateMenuInfoApi(id: string, data: Partial<MenuApi.MenuRecord>) {
  return requestClient.post('/menu/update', { id: String(id), ...data });
}

/**
 * 更新菜单状态
 */
export async function updateMenuStatusApi(id: string, status: number) {
  return requestClient.post('/menu/status', { id: String(id), status });
}

/**
 * 更新菜单显示状态
 */
export async function updateMenuShowApi(id: string, isShow: number) {
  return requestClient.post('/menu/show', { id: String(id), isShow });
}

/**
 * 删除菜单
 */
export async function deleteMenuApi(id: string) {
  return requestClient.post('/menu/delete', { id: String(id) });
}

/**
 * 批量删除菜单
 */
export async function batchDeleteMenuApi(ids: string[]) {
  return requestClient.post('/menu/batch-delete', { ids: ids.map((id) => String(id)) });
}
