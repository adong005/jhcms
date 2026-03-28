import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { page = 1, pageSize = 10, keyword, subdomain } = body;

  // 生成站群数据
  const generateSiteGroups = () => {
    const siteGroups = [];
    for (let i = 1; i <= 50; i++) {
      siteGroups.push({
        id: `sg_${i}`,
        keyword: `关键词${i}`,
        subdomain: `site${i}`,
        title: `站群网站${i}`,
        keywords: `关键词${i},SEO,优化`,
        description: `这是站群网站${i}的描述信息，用于展示网站的主要内容和特点。`,
        createTime: `2024-03-${String((i % 28) + 1).padStart(2, '0')} ${String((i % 24)).padStart(2, '0')}:${String((i % 60)).padStart(2, '0')}:00`,
        updateTime: `2024-03-25 ${String((i % 24)).padStart(2, '0')}:${String((i % 60)).padStart(2, '0')}:00`,
      });
    }
    return siteGroups;
  };

  const allData = generateSiteGroups();

  // 过滤数据
  let filteredData = allData;

  if (keyword) {
    filteredData = filteredData.filter(item =>
      item.keyword.includes(keyword)
    );
  }

  if (subdomain) {
    filteredData = filteredData.filter(item =>
      item.subdomain.includes(subdomain)
    );
  }

  // 分页
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = filteredData.slice(start, end);

  return useResponseSuccess({
    items,
    total: filteredData.length,
  });
});
