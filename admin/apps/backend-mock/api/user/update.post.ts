import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id, ...updateData } = body;

  // 模拟更新用户信息
  // 实际项目中这里应该更新数据库
  console.log(`更新用户 ${id} 的信息:`, updateData);

  return useResponseSuccess({
    id,
    ...updateData,
    message: '用户信息更新成功',
  });
});
