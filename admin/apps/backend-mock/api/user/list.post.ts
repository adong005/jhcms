import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  // 获取请求参数
  const body = await readBody(event);
  const { page = 1, pageSize = 10, username, nickName, status } = body;

  // 生成200条用户数据
  const generateUsers = () => {
    const users = [];
    const roles = [1, 2, 3]; // 1=超级管理员, 2=管理员, 3=用户
    const statuses = [0, 1]; // 0=禁用, 1=启用
    
    for (let i = 1; i <= 200; i++) {
      const roleIndex = i <= 5 ? 0 : i <= 30 ? 1 : 2; // 前5个超管，6-30管理员，其余用户
      const role = roles[roleIndex];
      const statusValue = i % 5 === 0 ? 0 : 1; // 每5个有1个禁用
      
      const user = {
        id: String(i),
        username: `user${i.toString().padStart(3, '0')}`,
        nickName: role === 1 ? `超级管理员${i}` : role === 2 ? `管理员${i}` : `用户${i}`,
        email: `user${i}@example.com`,
        phone: `138${String(i).padStart(8, '0')}`,
        role,
        status: statusValue,
        createTime: new Date(2024, 0, (i % 28) + 1, 10, 0, 0).toISOString().slice(0, 19).replace('T', ' '),
        lastLoginDate: new Date(2024, 2, (i % 28) + 1, (i % 24), (i % 60), 0).toISOString().slice(0, 19).replace('T', ' '),
        expireDate: statusValue === 1 ? '2024-12-31' : '2024-03-31',
      };
      
      users.push(user);
    }
    
    return users;
  };

  const allUsers = generateUsers();

  // 过滤数据
  let filteredData = allUsers;

  if (username) {
    filteredData = filteredData.filter(item => 
      item.username.includes(username)
    );
  }

  if (nickName) {
    filteredData = filteredData.filter(item => 
      item.nickName.includes(nickName)
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
