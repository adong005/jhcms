import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { page = 1, pageSize = 10, name, tenantId, status } = body;

  const generateCategories = () => {
    const categories = [
      {
        id: '1',
        name: '新闻资讯',
        code: 'news',
        tenantId: 'tenant_001',
        tenantName: '默认租户',
        sort: 1,
        description: '新闻资讯分类',
        status: 1,
        createTime: '2024-01-01 10:00:00',
        updateTime: '2024-03-25 10:00:00',
      },
      {
        id: '2',
        name: '技术文章',
        code: 'tech',
        tenantId: 'tenant_001',
        tenantName: '默认租户',
        sort: 2,
        description: '技术相关文章',
        status: 1,
        createTime: '2024-01-02 10:00:00',
        updateTime: '2024-03-25 10:00:00',
      },
      {
        id: '3',
        name: '产品介绍',
        code: 'product',
        tenantId: 'tenant_002',
        tenantName: '企业租户A',
        sort: 3,
        description: '产品介绍相关',
        status: 1,
        createTime: '2024-01-03 10:00:00',
        updateTime: '2024-03-25 10:00:00',
      },
      {
        id: '4',
        name: '公司动态',
        code: 'company',
        tenantId: 'tenant_002',
        tenantName: '企业租户A',
        sort: 4,
        description: '公司动态信息',
        status: 0,
        createTime: '2024-01-04 10:00:00',
        updateTime: '2024-03-25 10:00:00',
      },
      {
        id: '5',
        name: '行业资讯',
        code: 'industry',
        tenantId: 'tenant_003',
        tenantName: '企业租户B',
        sort: 5,
        description: '行业相关资讯',
        status: 1,
        createTime: '2024-01-05 10:00:00',
        updateTime: '2024-03-25 10:00:00',
      },
    ];

    return categories;
  };

  let categories = generateCategories();

  if (name) {
    categories = categories.filter(item => item.name.includes(name));
  }

  if (tenantId) {
    categories = categories.filter(item => item.tenantId === tenantId);
  }

  if (status !== undefined && status !== '') {
    categories = categories.filter(item => item.status === Number(status));
  }

  const total = categories.length;
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = categories.slice(start, end);

  return useResponseSuccess({
    items,
    total,
  });
});
