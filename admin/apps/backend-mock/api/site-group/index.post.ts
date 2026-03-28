import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { keyword, subdomain, title, keywords, description } = body;

  if (!keyword || !subdomain || !title || !keywords || !description) {
    return useResponseError('请填写完整信息');
  }

  return useResponseSuccess({
    id: `sg_${Date.now()}`,
    keyword,
    subdomain,
    title,
    keywords,
    description,
    createTime: new Date().toISOString(),
    updateTime: new Date().toISOString(),
    message: '创建站群成功',
  });
});
