import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id } = body;

  // 模拟删除用户
  // 实际项目中这里应该从数据库删除
  console.log(`删除用户 ID: ${id}`);

  return useResponseSuccess({
    id,
    message: '用户删除成功',
  });
});
