import { requestClient } from '#/api/request';

export namespace LogApi {
  export interface LogRecord {
    id: string;
    requestId?: string;
    tenantId: string;
    userId?: string;
    username: string;
    action: string;
    module: string;
    description: string;
    targetId?: string;
    ip: string;
    method?: string;
    url?: string;
    userAgent?: string;
    status: 'success' | 'fail';
    logType?: string;
    statusCode?: number;
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
    module?: string;
    ip?: string;
    logType?: string;
    startTime?: string;
    endTime?: string;
  }

  export interface LogListResult {
    items: LogRecord[];
    total: number;
  }
}

export async function getLogListApi(params: LogApi.LogListParams) {
  return requestClient.post<LogApi.LogListResult>('/system-logs/list', params);
}

export async function deleteLogApi(id: string) {
  return requestClient.delete(`/system-logs/${String(id)}`);
}

export async function batchDeleteLogApi(ids: string[]) {
  return requestClient.post('/system-logs/batch-delete', { ids: ids.map((id) => String(id)) });
}

export async function clearLogsApi(force = false) {
  return requestClient.post('/system-logs/clear', { force });
}

export async function purgeOldLogsApi(days: number) {
  return requestClient.post('/system-logs/purge', { days });
}

export async function exportLogsApi(params: Omit<LogApi.LogListParams, 'page' | 'pageSize'>) {
  return requestClient.post('/system-logs/export', params, { responseType: 'blob' });
}
