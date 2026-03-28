import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse, useResponseSuccess } from '~/utils/response';

export default defineEventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { roleId, menuIds } = body;

  if (!roleId) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '角色ID不能为空',
    };
  }

  if (!menuIds || !Array.isArray(menuIds)) {
    setResponseStatus(event, 400);
    return {
      code: 400,
      message: '菜单权限列表格式不正确',
    };
  }

  // Mock: 这里应该将权限保存到数据库
  // 实际项目中需要：
  // 1. 验证 roleId 是否存在
  // 2. 验证 menuIds 中的菜单ID是否都有效
  // 3. 删除该角色的旧权限
  // 4. 插入新的权限关系

  console.log(`更新角色 ${roleId} 的权限:`, menuIds);

  return useResponseSuccess({
    roleId,
    menuIds,
    message: '权限配置更新成功',
  });
});
