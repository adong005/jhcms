import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  // 获取请求参数
  const body = await readBody(event);
  const { page = 1, pageSize = 10, name, code, status } = body;

  // 生成角色数据
  const generateRoles = () => {
    const roles = [
      {
        id: '1',
        name: '超级管理员',
        code: 'super_admin',
        description: '拥有系统所有权限',
        status: 1,
        createTime: '2024-01-01 10:00:00',
        updateTime: '2024-03-20 14:00:00',
      },
      {
        id: '2',
        name: '管理员',
        code: 'admin',
        description: '拥有系统大部分权限',
        status: 1,
        createTime: '2024-01-02 10:00:00',
        updateTime: '2024-03-21 14:00:00',
      },
      {
        id: '3',
        name: '用户',
        code: 'user',
        description: '拥有基本权限',
        status: 1,
        createTime: '2024-01-03 10:00:00',
        updateTime: '2024-03-22 14:00:00',
      },
    ];
    
    return roles;
  };

  const allRoles = generateRoles();

  // 过滤数据
  let filteredData = allRoles;

  if (name) {
    filteredData = filteredData.filter(item => 
      item.name.includes(name)
    );
  }

  if (code) {
    filteredData = filteredData.filter(item => 
      item.code.includes(code)
    );
  }

  if (status !== undefined && status !== '') {
    filteredData = filteredData.filter(item => 
      item.status === Number(status)
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
