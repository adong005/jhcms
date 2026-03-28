import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { id } = body;

  if (!id) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '信息ID不能为空',
    };
  }

  return useResponseSuccess({
    id,
    message: '删除成功',
  });
});
