import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { ids } = body;

  if (!ids || !Array.isArray(ids) || ids.length === 0) {
    return useResponseError('请选择要删除的站群');
  }

  return useResponseSuccess({
    message: `批量删除 ${ids.length} 个站群成功`,
  });
});
