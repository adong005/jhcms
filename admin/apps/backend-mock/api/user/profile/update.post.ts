import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { realName, username, introduction } = body;

  // 这里应该更新数据库中的用户信息
  // 为了演示，我们只返回成功响应
  
  return useResponseSuccess({
    realName,
    username,
    introduction,
  });
});
