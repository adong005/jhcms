import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { ids } = body;

  if (!ids || !Array.isArray(ids) || ids.length === 0) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '请选择要删除的信息',
    };
  }

  return useResponseSuccess({
    ids,
    count: ids.length,
    message: `成功删除 ${ids.length} 条信息`,
  });
});
