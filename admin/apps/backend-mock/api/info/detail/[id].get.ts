import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const id = getRouterParam(event, 'id');

  const mockInfoDetails: Record<string, any> = {
    '1': {
      id: '1',
      title: 'Vue 3.5 新特性详解',
      categoryId: '2',
      categoryName: '技术文章',
      author: '张三',
      tenantId: 'tenant_001',
      tenantName: '默认租户',
      summary: '深入解析 Vue 3.5 版本带来的新特性和改进',
      content: '<h1>Vue 3.5 新特性</h1><p>Vue 3.5 版本带来了许多令人兴奋的新特性...</p><h2>响应式系统优化</h2><p>新版本对响应式系统进行了全面优化...</p>',
      viewCount: 1520,
      status: 1,
      createTime: '2024-03-19 15:00:00',
      updateTime: '2024-03-20 10:00:00',
    },
  };

  const detail = mockInfoDetails[id || ''] || {
    id: id || '',
    title: '示例信息',
    categoryId: '1',
    categoryName: '新闻资讯',
    author: '系统',
    tenantId: 'tenant_001',
    tenantName: '默认租户',
    summary: '这是一条示例信息',
    content: '<h1>示例标题</h1><p>这是示例内容</p>',
    viewCount: 0,
    status: 0,
    createTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
    updateTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
  };

  return useResponseSuccess(detail);
});
