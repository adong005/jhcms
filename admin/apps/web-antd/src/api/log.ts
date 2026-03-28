import { requestClient } from '#/api/request';

export namespace LogApi {
  export interface LogRecord {
    id: string;
    tenantId: string;
    username: string;
    action: string;
    module: string;
    description: string;
    ip: string;
    status: 'success' | 'fail';
    duration: number;
    createTime: string;
    errorMsg?: string;
    requestJson?: string;
  }

  export interface LogListParams {
    page: number;
    pageSize: number;
    tenantId?: string;
    username?: string;
    usernames?: string[];
    action?: string;
    status?: string;
    date?: string;
  }

  export interface LogListResult {
    items: LogRecord[];
    total: number;
  }
}

/**
 * 获取日志列表
 */
export async function getLogListApi(params: LogApi.LogListParams) {
  return requestClient.post<LogApi.LogListResult>('/system-logs/list', params);
}

/**
 * 删除日志
 */
export async function deleteLogApi(id: string) {
  return requestClient.delete(`/system-logs/${String(id)}`);
}

/**
 * 批量删除日志
 */
export async function batchDeleteLogApi(ids: string[]) {
  return requestClient.post('/system-logs/batch-delete', { ids: ids.map((id) => String(id)) });
}

/**
 * 清空日志
 */
export async function clearLogsApi() {
  return requestClient.post('/system-logs/clear');
}
