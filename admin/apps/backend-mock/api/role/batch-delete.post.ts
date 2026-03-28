import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { ids } = body;

  console.log(`批量删除角色 IDs: ${ids.join(', ')}`);

  return useResponseSuccess({
    count: ids.length,
    message: `成功删除 ${ids.length} 个角色`,
  });
});
