import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { name, path, type, icon, component, order, status } = body;

  // 模拟创建菜单
  const newMenu = {
    id: String(Date.now()),
    name,
    path,
    type,
    icon,
    component,
    order,
    status,
    createTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
    updateTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
  };

  console.log('创建菜单:', newMenu);

  return useResponseSuccess({
    ...newMenu,
    message: '菜单创建成功',
  });
});
