import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id, status } = body;

  console.log(`更新菜单 ${id} 的状态为 ${status}`);

  return useResponseSuccess({
    id,
    status,
    message: '状态更新成功',
  });
});
