import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id, ...updateData } = body;

  console.log(`更新角色 ${id} 的信息:`, updateData);

  return useResponseSuccess({
    id,
    ...updateData,
    message: '角色信息更新成功',
  });
});
