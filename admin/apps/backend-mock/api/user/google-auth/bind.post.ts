import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { verifyCode } = body;

  // 验证验证码（6位数字）
  if (!verifyCode || !/^\d{6}$/.test(verifyCode)) {
    return useResponseError('验证码格式不正确，请输入6位数字');
  }

  return useResponseSuccess({
    isBound: true,
    message: '谷歌验证器绑定成功',
  });
});
