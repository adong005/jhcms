import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { newEmail, verifyCode } = body;

  // 验证验证码
  if (!verifyCode || verifyCode.length !== 6) {
    return useResponseError('验证码格式不正确');
  }

  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(newEmail)) {
    return useResponseError('邮箱格式不正确');
  }

  return useResponseSuccess({
    email: newEmail,
    message: '邮箱更新成功',
  });
});
