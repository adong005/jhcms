import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id, status } = body;

  // 模拟更新用户状态
  // 实际项目中这里应该更新数据库
  console.log(`更新用户 ${id} 的状态为 ${status}`);

  return useResponseSuccess({
    id,
    status,
    message: '状态更新成功',
  });
});
