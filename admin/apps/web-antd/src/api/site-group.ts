import { requestClient } from '#/api/request';

export namespace SiteGroupApi {
  export interface SiteGroup {
    id: string;
    adminId?: string;
    adminName?: string;
    keyword: string;
    subdomain: string;
    title: string;
    keywords: string;
    description: string;
    createTime: string;
    updateTime: string;
  }

  export interface SiteGroupListParams {
    page: number;
    pageSize: number;
    keyword?: string;
    subdomain?: string;
    adminId?: string;
  }

  export interface AdminOption {
    userId: string;
    username: string;
    nickName: string;
  }

  export interface SiteGroupListResult {
    items: SiteGroup[];
    total: number;
  }

  export interface CityRecord {
    cityCode: number;
    name: string;
    pinyin: string;
  }

  export interface CityListParams {
    page: number;
    pageSize: number;
    name?: string;
  }

  export interface CityListResult {
    items: CityRecord[];
    total: number;
  }

  export interface CreateSiteGroupParams {
    keyword: string;
    subdomain: string;
    title: string;
    keywords: string;
    description: string;
  }

  export interface UpdateSiteGroupParams {
    id: string;
    keyword: string;
    subdomain: string;
    title: string;
    keywords: string;
    description: string;
  }
}

/**
 * 获取站群列表
 */
export async function getSiteGroupListApi(params: SiteGroupApi.SiteGroupListParams) {
  return requestClient.post<SiteGroupApi.SiteGroupListResult>('/site-group/list', params);
}

/**
 * 创建站群
 */
export async function createSiteGroupApi(data: SiteGroupApi.CreateSiteGroupParams) {
  return requestClient.post('/site-group', data);
}

/**
 * 更新站群
 */
export async function updateSiteGroupApi(data: SiteGroupApi.UpdateSiteGroupParams) {
  return requestClient.put(`/site-group/${String(data.id)}`, { ...data, id: String(data.id) });
}

/**
 * 删除站群
 */
export async function deleteSiteGroupApi(id: string) {
  return requestClient.delete(`/site-group/${String(id)}`);
}

/**
 * 批量删除站群
 */
export async function batchDeleteSiteGroupApi(ids: string[]) {
  return requestClient.post('/site-group/batch-delete', { ids: ids.map((id) => String(id)) });
}

/**
 * 获取管理员选项
 */
export async function getSiteGroupAdminsApi() {
  return requestClient.get<SiteGroupApi.AdminOption[]>('/site-group/admins');
}

/**
 * 获取城市列表（超管）
 */
export async function getCityListApi(params: SiteGroupApi.CityListParams) {
  return requestClient.post<SiteGroupApi.CityListResult>('/site-group/cities', params);
}
