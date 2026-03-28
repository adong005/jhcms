import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  // 返回网站配置数据
  return useResponseSuccess({
    title: 'ADCMS 管理系统',
    keywords: '内容管理系统,CMS,后台管理,Vue3,Vben Admin',
    description: '基于 Vue3 和 Vben Admin 构建的现代化内容管理系统，提供完善的用户管理、权限管理、内容管理等功能。',
    domain: 'www.adcms.com',
  });
});
