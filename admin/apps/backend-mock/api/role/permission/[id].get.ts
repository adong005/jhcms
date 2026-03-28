import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const roleId = getRouterParam(event, 'id');

  // Mock 数据：不同角色的权限配置
  const rolePermissions: Record<string, string[]> = {
    '1': ['1', '2', '5', '6'], // 超级管理员：工作台、用户管理、信息管理、个人中心
    '2': ['1', '2', '6'], // 管理员：工作台、用户管理、个人中心
    '3': ['1', '6'], // 普通用户：工作台、个人中心
  };

  const menuIds = rolePermissions[roleId || ''] || [];

  return useResponseSuccess({
    menuIds,
  });
});
