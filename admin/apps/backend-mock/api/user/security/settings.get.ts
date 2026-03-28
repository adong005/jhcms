import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  // 返回模拟的安全设置数据
  return useResponseSuccess({
    accountPassword: true,
    securityPhone: true,
    securityPhoneNumber: '138****8293',
    securityQuestion: false,
    securityEmail: true,
    securityEmailAddress: 'ant***@sign.com',
    securityMfa: false,
    passwordStrength: '强',
  });
});
