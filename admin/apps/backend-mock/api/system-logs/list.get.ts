import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const query = getQuery(event);
  const page = Number(query.page) || 1;
  const pageSize = Number(query.pageSize) || 10;

  // 模拟日志数据
  const allLogs = [
    {
      id: '1',
      username: 'admin',
      action: 'login',
      module: '系统登录',
      description: '用户登录系统',
      ip: '192.168.1.100',
      status: 'success',
      duration: 125,
      createTime: '2024-03-25 10:30:15',
    },
    {
      id: '2',
      username: 'admin',
      action: 'create',
      module: '用户管理',
      description: '创建新用户: test_user',
      ip: '192.168.1.100',
      status: 'success',
      duration: 89,
      createTime: '2024-03-25 10:35:20',
    },
    {
      id: '3',
      username: 'admin',
      action: 'update',
      module: '角色管理',
      description: '更新角色: 管理员',
      ip: '192.168.1.100',
      status: 'success',
      duration: 156,
      createTime: '2024-03-25 10:40:30',
    },
    {
      id: '4',
      username: 'user1',
      action: 'delete',
      module: '菜单管理',
      description: '删除菜单: 测试菜单',
      ip: '192.168.1.101',
      status: 'fail',
      duration: 45,
      createTime: '2024-03-25 10:45:10',
      errorMsg: '权限不足，无法删除菜单',
    },
    {
      id: '5',
      username: 'admin',
      action: 'query',
      module: '日志管理',
      description: '查询系统日志',
      ip: '192.168.1.100',
      status: 'success',
      duration: 234,
      createTime: '2024-03-25 10:50:25',
    },
    {
      id: '6',
      username: 'admin',
      action: 'export',
      module: '用户管理',
      description: '导出用户列表',
      ip: '192.168.1.100',
      status: 'success',
      duration: 1523,
      createTime: '2024-03-25 11:00:15',
    },
    {
      id: '7',
      username: 'user2',
      action: 'login',
      module: '系统登录',
      description: '用户登录系统',
      ip: '192.168.1.102',
      status: 'fail',
      duration: 89,
      createTime: '2024-03-25 11:05:30',
      errorMsg: '用户名或密码错误',
    },
    {
      id: '8',
      username: 'admin',
      action: 'update',
      module: '用户管理',
      description: '更新用户信息: user1',
      ip: '192.168.1.100',
      status: 'success',
      duration: 178,
      createTime: '2024-03-25 11:10:45',
    },
  ];

  // 过滤
  let filteredLogs = allLogs;
  if (query.username) {
    filteredLogs = filteredLogs.filter(log => 
      log.username.includes(query.username as string)
    );
  }
  if (query.action) {
    filteredLogs = filteredLogs.filter(log => log.action === query.action);
  }
  if (query.status) {
    filteredLogs = filteredLogs.filter(log => log.status === query.status);
  }

  const total = filteredLogs.length;
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = filteredLogs.slice(start, end);

  return useResponseSuccess({
    items,
    total,
  });
});
