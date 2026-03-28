import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { username, nickName, email, phone, role, status, expireDate, password } = body;

  // 模拟创建用户
  // 实际项目中这里应该插入数据库
  const newUser = {
    id: String(Date.now()),
    username,
    nickName,
    email,
    phone,
    role,
    status,
    createTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
    lastLoginDate: '',
    expireDate,
  };

  console.log('创建用户:', newUser);

  return useResponseSuccess({
    ...newUser,
    message: '用户创建成功',
  });
});
