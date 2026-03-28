import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { name, code, tenantId, sort, description, status } = body;

  if (!name || !code) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '分类名称和编码不能为空',
    };
  }

  const newCategory = {
    id: `${Date.now()}`,
    name,
    code,
    tenantId: tenantId || 'tenant_001',
    tenantName: tenantId === 'tenant_002' ? '企业租户A' : tenantId === 'tenant_003' ? '企业租户B' : '默认租户',
    sort: sort || 0,
    description: description || '',
    status: status !== undefined ? status : 1,
    createTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
    updateTime: new Date().toISOString().replace('T', ' ').substring(0, 19),
  };

  return useResponseSuccess(newCategory);
});
