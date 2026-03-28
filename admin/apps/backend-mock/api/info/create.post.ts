import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { title, categoryId, author, tenantId, summary, content, status } = body;

  if (!title || !categoryId || !content) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '标题、分类和内容不能为空',
    };
  }

  const categoryNames: Record<string, string> = {
    '1': '新闻资讯',
    '2': '技术文章',
    '3': '产品介绍',
    '4': '公司动态',
    '5': '行业资讯',
  };

  const tenantNames: Record<string, string> = {
    'tenant_001': '默认租户',
    'tenant_002': '企业租户A',
    'tenant_003': '企业租户B',
  };

  const newInfo = {
    id: `${Date.now()}`,
    title,
    categoryId,
    categoryName: categoryNames[categoryId] || '未知分类',
    author: author || '匿名',
    tenantId: tenantId || 'tenant_001',
    tenantName: tenantNames[tenantId] || '默认租户',
    summary: summary || '',
    content,
    viewCount: 0,
    status: status !== undefined ? status : 0,
    createTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
    updateTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
  };

  return useResponseSuccess(newInfo);
});
