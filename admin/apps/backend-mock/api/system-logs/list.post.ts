import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { page = 1, pageSize = 10, username, action, status, date } = body;

  // 生成日志数据
  const generateLogs = () => {
    const logs = [];
    const actions = ['登录', '登出', '创建', '更新', '删除', '查询', '导出'];
    const modules = ['用户管理', '角色管理', '菜单管理', '信息管理', '站群管理', '表单管理'];
    const statuses: ('success' | 'fail')[] = ['success', 'fail'];
    
    for (let i = 1; i <= 200; i++) {
      const actionType = actions[i % actions.length];
      const moduleType = modules[i % modules.length];
      const statusType = statuses[i % 10 === 0 ? 1 : 0];
      
      logs.push({
        id: `log_${i}`,
        username: `user${(i % 20) + 1}`,
        action: actionType,
        module: moduleType,
        description: `${actionType}${moduleType}`,
        ip: `192.168.1.${(i % 255) + 1}`,
        status: statusType,
        duration: Math.floor(Math.random() * 1000) + 100,
        createTime: `2024-03-${String((i % 28) + 1).padStart(2, '0')} ${String((i % 24)).padStart(2, '0')}:${String((i % 60)).padStart(2, '0')}:00`,
        errorMsg: statusType === 'fail' ? '操作失败：权限不足' : undefined,
      });
    }
    return logs;
  };

  const allData = generateLogs();

  // 过滤数据
  let filteredData = allData;

  if (username) {
    filteredData = filteredData.filter(item =>
      item.username.includes(username)
    );
  }

  if (action) {
    filteredData = filteredData.filter(item =>
      item.action.includes(action)
    );
  }

  if (status) {
    filteredData = filteredData.filter(item =>
      item.status === status
    );
  }

  if (date) {
    filteredData = filteredData.filter(item =>
      item.createTime.startsWith(date)
    );
  }

  // 分页
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = filteredData.slice(start, end);

  return useResponseSuccess({
    items,
    total: filteredData.length,
  });
});
