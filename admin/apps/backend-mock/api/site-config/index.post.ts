import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { userId, title, keywords, description, domain } = body;

  // 验证必填字段
  if (!userId) {
    return useResponseError('用户ID不能为空');
  }

  if (!title) {
    return useResponseError('网站标题不能为空');
  }

  if (!keywords) {
    return useResponseError('网站关键词不能为空');
  }

  if (!description) {
    return useResponseError('网站描述不能为空');
  }

  if (!domain) {
    return useResponseError('网站域名不能为空');
  }

  // 返回更新成功的数据
  return useResponseSuccess({
    userId,
    title,
    keywords,
    description,
    domain,
    message: '网站配置更新成功',
    updateTime: new Date().toISOString(),
  });
});
