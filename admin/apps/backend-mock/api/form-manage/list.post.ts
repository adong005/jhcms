import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { page = 1, pageSize = 10, contact, phone, company } = body;

  // 生成表单数据
  const generateForms = () => {
    const forms = [];
    for (let i = 1; i <= 100; i++) {
      forms.push({
        id: `form_${i}`,
        contact: `联系人${i}`,
        phone: `138${String(i).padStart(8, '0')}`,
        company: `公司${i}`,
        ip: `192.168.1.${(i % 255) + 1}`,
        handleStatus: (i % 3) as 0 | 1 | 2, // 0=未处理 1=已分配 2=已完成
        createTime: `2024-03-${String((i % 28) + 1).padStart(2, '0')} ${String((i % 24)).padStart(2, '0')}:${String((i % 60)).padStart(2, '0')}:00`,
        remark: i % 3 === 0 ? `这是备注信息${i}` : '',
      });
    }
    return forms;
  };

  const allData = generateForms();

  // 过滤数据
  let filteredData = allData;

  if (contact) {
    filteredData = filteredData.filter(item =>
      item.contact.includes(contact as string)
    );
  }

  if (phone) {
    filteredData = filteredData.filter(item =>
      item.phone.includes(phone as string)
    );
  }

  if (company) {
    filteredData = filteredData.filter(item =>
      item.company.includes(company as string)
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
