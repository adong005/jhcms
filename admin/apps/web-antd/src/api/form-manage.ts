import { requestClient } from '#/api/request';

export namespace FormManageApi {
  export interface FormRecord {
    id: string;
    contact: string;
    phone: string;
    company: string;
    ip: string;
    handleStatus: 0 | 1 | 2; // 0=未处理 1=已分配 2=已完成
    createTime: string;
    updateTime: string;
    createdBy?: string | number;
    remark: string;
  }

  export interface FormListParams {
    page: number;
    pageSize: number;
    contact?: string;
    phone?: string;
    company?: string;
  }

  export interface FormListResult {
    items: FormRecord[];
    total: number;
  }
}

/**
 * 获取表单列表
 */
export async function getFormListApi(params: FormManageApi.FormListParams) {
  return requestClient.post<FormManageApi.FormListResult>('/form-manage/list', params);
}

/**
 * 删除表单记录
 */
export async function deleteFormApi(id: string) {
  return requestClient.delete(`/form-manage/${String(id)}`);
}

/**
 * 批量删除表单记录
 */
export async function batchDeleteFormApi(ids: string[]) {
  return requestClient.post('/form-manage/batch-delete', { ids: ids.map((id) => String(id)) });
}

/**
 * 导出表单数据
 */
export async function exportFormApi(params: FormManageApi.FormListParams) {
  return requestClient.post('/form-manage/export', params, {
    responseType: 'blob',
  });
}
