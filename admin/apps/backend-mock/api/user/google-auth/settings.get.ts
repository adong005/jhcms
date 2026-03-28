import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  return useResponseSuccess({
    isBound: false,
    qrCodeUrl: 'https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=otpauth://totp/ADCMS:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=ADCMS',
    secretKey: 'JBSWY3DPEHPK3PXP',
  });
});
