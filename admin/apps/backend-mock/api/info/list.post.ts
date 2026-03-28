import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { page = 1, pageSize = 10, title, categoryId, tenantId, status } = body;

  const generateInfoList = () => {
    const infoList = [
      {
        id: '1',
        title: 'Vue 3.5 新特性详解',
        categoryId: '2',
        categoryName: '技术文章',
        author: '张三',
        tenantId: 'tenant_001',
        tenantName: '默认租户',
        summary: '深入解析 Vue 3.5 版本带来的新特性和改进',
        content: '<h1>Vue 3.5 新特性</h1><p>本文详细介绍...</p>',
        viewCount: 1520,
        status: 1,
        createTime: '2024-03-19 15:00:00',
        updateTime: '2024-03-20 10:00:00',
      },
      {
        id: '2',
        title: '公司年度总结大会成功举办',
        categoryId: '4',
        categoryName: '公司动态',
        author: '李四',
        tenantId: 'tenant_002',
        tenantName: '企业租户A',
        summary: '2024年度总结大会圆满落幕，展望新一年发展规划',
        content: '<h1>年度总结</h1><p>会议内容...</p>',
        viewCount: 856,
        status: 1,
        createTime: '2024-03-18 10:00:00',
        updateTime: '2024-03-18 14:00:00',
      },
      {
        id: '3',
        title: '新产品发布预告',
        categoryId: '3',
        categoryName: '产品介绍',
        author: '王五',
        tenantId: 'tenant_002',
        tenantName: '企业租户A',
        summary: '即将发布的新产品功能预览和特性介绍',
        content: '<h1>新产品</h1><p>产品特性...</p>',
        viewCount: 2340,
        status: 0,
        createTime: '2024-03-22 09:00:00',
        updateTime: '2024-03-24 16:00:00',
      },
      {
        id: '4',
        title: '行业发展趋势分析报告',
        categoryId: '5',
        categoryName: '行业资讯',
        author: '赵六',
        tenantId: 'tenant_003',
        tenantName: '企业租户B',
        summary: '2024年行业发展趋势深度分析',
        content: '<h1>趋势分析</h1><p>详细内容...</p>',
        viewCount: 3210,
        status: 1,
        createTime: '2024-03-14 14:00:00',
        updateTime: '2024-03-15 11:00:00',
      },
      {
        id: '5',
        title: '最新政策解读',
        categoryId: '1',
        categoryName: '新闻资讯',
        author: '孙七',
        tenantId: 'tenant_001',
        tenantName: '默认租户',
        summary: '解读最新发布的行业政策及其影响',
        content: '<h1>政策解读</h1><p>政策内容...</p>',
        viewCount: 1890,
        status: 1,
        createTime: '2024-03-21 10:00:00',
        updateTime: '2024-03-21 16:00:00',
      },
      {
        id: '6',
        title: 'TypeScript 5.0 实战指南',
        categoryId: '2',
        categoryName: '技术文章',
        author: '周八',
        tenantId: 'tenant_001',
        tenantName: '默认租户',
        summary: 'TypeScript 5.0 新特性和最佳实践',
        content: '<h1>TypeScript 5.0</h1><p>实战案例...</p>',
        viewCount: 2670,
        status: 1,
        createTime: '2024-03-16 15:00:00',
        updateTime: '2024-03-17 13:00:00',
      },
    ];

    return infoList;
  };

  let infoList = generateInfoList();

  if (title) {
    infoList = infoList.filter(item => item.title.includes(title));
  }

  if (categoryId) {
    infoList = infoList.filter(item => item.categoryId === categoryId);
  }

  if (tenantId) {
    infoList = infoList.filter(item => item.tenantId === tenantId);
  }

  if (status !== undefined && status !== '') {
    infoList = infoList.filter(item => item.status === Number(status));
  }

  const total = infoList.length;
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = infoList.slice(start, end);

  return useResponseSuccess({
    items,
    total,
  });
});
