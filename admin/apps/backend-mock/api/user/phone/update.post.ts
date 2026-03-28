import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { newPhone, verifyCode } = body;

  // 验证验证码（这里简单模拟）
  if (!verifyCode || verifyCode.length !== 6) {
    return useResponseError('验证码格式不正确');
  }

  return useResponseSuccess({
    phone: newPhone,
    message: '手机号更新成功',
  });
});
