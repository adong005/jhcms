import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id, status } = body;

  if (!id) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '信息ID不能为空',
    };
  }

  return useResponseSuccess({
    id,
    status,
    updateTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
  });
});
