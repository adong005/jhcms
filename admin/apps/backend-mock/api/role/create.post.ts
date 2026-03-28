import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { name, code, description, status } = body;

  // 模拟创建角色
  const newRole = {
    id: String(Date.now()),
    name,
    code,
    description,
    status,
    createTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
    updateTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
  };

  console.log('创建角色:', newRole);

  return useResponseSuccess({
    ...newRole,
    message: '角色创建成功',
  });
});
