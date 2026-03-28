import { verifyAccessToken } from '~/utils/jwt-utils';
import { unAuthorizedResponse } from '~/utils/response';

export default eventHandler(async (event) => {
  const userinfo = verifyAccessToken(event);
  if (!userinfo) {
    return unAuthorizedResponse(event);
  }

  const body = await readBody(event);
  const { question1, answer1, question2, answer2 } = body;

  return useResponseSuccess({
    question1,
    answer1,
    question2,
    answer2,
    message: '密保问题设置成功',
  });
});
