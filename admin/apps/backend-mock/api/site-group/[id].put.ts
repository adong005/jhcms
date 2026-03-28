import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const id = getRouterParam(event, 'id');
  const body = await readBody(event);
  const { keyword, subdomain, title, keywords, description } = body;

  if (!keyword || !subdomain || !title || !keywords || !description) {
    return useResponseError('请填写完整信息');
  }

  return useResponseSuccess({
    id,
    keyword,
    subdomain,
    title,
    keywords,
    description,
    updateTime: new Date().toISOString(),
    message: '更新站群成功',
  });
});
