import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id } = body;

  console.log(`删除角色 ID: ${id}`);

  return useResponseSuccess({
    id,
    message: '角色删除成功',
  });
});
