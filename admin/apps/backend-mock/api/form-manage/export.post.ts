import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);

  // 模拟导出功能
  return useResponseSuccess({
    message: '导出成功',
    data: {
      filename: `表单数据_${new Date().toISOString().split('T')[0]}.xlsx`,
      url: '/download/forms.xlsx',
    },
  });
});
