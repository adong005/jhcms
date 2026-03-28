import { requestClient } from '#/api/request';

export namespace UserApi {
  export interface UserListParams {
    page: number;
    pageSize: number;
    username?: string;
    nickName?: string;
    status?: number | string;
  }

  export interface UserRecord {
    id: string;
    username: string;
    nickName: string;
    email: string;
    phone: string;
    /** 后端可能返回 1/2/3 或 super_admin / admin / user */
    role: number | string;
    status: number;
    createTime: string;
    updateTime: string;
    /** 创建人用户 ID，可能为空 */
    createdBy?: string | number;
    lastLoginDate: string;
    expireDate: string;
  }

  export interface UserListResult {
    items: UserRecord[];
    total: number;
  }
}

/**
 * 获取用户列表
 */
export async function getUserListApi(params: UserApi.UserListParams) {
  return requestClient.post<UserApi.UserListResult>('/user/list', params);
}

/**
 * 更新用户状态
 */
export async function updateUserStatusApi(id: string, status: number) {
  return requestClient.post('/user/status', { id: String(id), status });
}

/**
 * 更新用户信息
 */
export async function updateUserInfoApi(id: string, data: Partial<UserApi.UserRecord>) {
  return requestClient.post('/user/update', { id: String(id), ...data });
}

/**
 * 创建用户
 */
export async function createUserApi(data: Partial<UserApi.UserRecord> & { password: string }) {
  return requestClient.post('/user/create', data);
}

/**
 * 删除用户
 */
export async function deleteUserApi(id: string) {
  return requestClient.post('/user/delete', { id: String(id) });
}

/** 管理员重置密码：走 /user/update + resetPassword；snapshot 用于旧后端误忽略 resetPassword 时尽量不清空其它字段 */
export async function resetUserPasswordApi(
  id: string,
  snapshot?: Partial<
    Pick<UserApi.UserRecord, 'nickName' | 'email' | 'phone' | 'status' | 'expireDate' | 'role'>
  >,
) {
  return requestClient.post<{ newPassword: string; email: string }>('/user/update', {
    id: String(id),
    resetPassword: true,
    ...snapshot,
  });
}

/**
 * 批量删除用户
 */
export async function batchDeleteUserApi(ids: string[]) {
  return requestClient.post('/user/batch-delete', { ids: ids.map((id) => String(id)) });
}

/**
 * 更新当前用户个人资料
 */
export async function updateProfileApi(data: {
  realName?: string;
  nickName?: string;
  username?: string;
}) {
  return requestClient.post('/user/profile/update', data);
}

/**
 * 获取安全设置
 */
export async function getSecuritySettingsApi() {
  return requestClient.get<{
    accountPassword: boolean;
    securityPhone: boolean;
    securityPhoneNumber?: string;
    securityQuestion: boolean;
    securityEmail: boolean;
    securityEmailAddress?: string;
    securityMfa: boolean;
    passwordStrength?: string;
  }>('/user/security/settings');
}

/**
 * 更新安全设置
 */
export async function updateSecuritySettingsApi(data: Record<string, boolean>) {
  return requestClient.post('/user/security/update', data);
}

/**
 * 获取密保手机设置
 */
export async function getPhoneSettingApi() {
  return requestClient.get<{ phone: string }>('/user/phone/settings');
}

/**
 * 更新密保手机
 */
export async function updatePhoneSettingApi(data: {
  newPhone: string;
  verifyCode: string;
}) {
  return requestClient.post('/user/phone/update', data);
}

/**
 * 获取密保问题设置
 */
export async function getQuestionSettingApi() {
  return requestClient.get<{
    question1?: string;
    answer1?: string;
    question2?: string;
    answer2?: string;
  }>('/user/question/settings');
}

/**
 * 更新密保问题
 */
export async function updateQuestionSettingApi(data: {
  question1: string;
  answer1: string;
  question2: string;
  answer2: string;
}) {
  return requestClient.post('/user/question/update', data);
}

/**
 * 获取邮箱设置
 */
export async function getEmailSettingApi() {
  return requestClient.get<{ email: string }>('/user/email/settings');
}

/**
 * 更新联系邮箱
 */
export async function updateEmailSettingApi(data: {
  newEmail: string;
  verifyCode: string;
}) {
  return requestClient.post('/user/email/update', data);
}

/**
 * 获取谷歌验证器设置
 */
export async function getGoogleAuthSettingApi() {
  return requestClient.get<{
    isBound: boolean;
    qrCodeUrl?: string;
    secretKey?: string;
  }>('/user/google-auth/settings');
}

/**
 * 绑定谷歌验证器
 */
export async function bindGoogleAuthApi(data: { verifyCode: string }) {
  return requestClient.post('/user/google-auth/bind', data);
}

/**
 * 解绑谷歌验证器
 */
export async function unbindGoogleAuthApi() {
  return requestClient.post('/user/google-auth/unbind');
}

/**
 * 修改当前用户密码
 */
export async function updatePasswordApi(data: { oldPassword: string; newPassword: string }) {
  return requestClient.post('/user/password/update', data);
}

/**
 * 获取通知设置
 */
export async function getNotificationSettingsApi() {
  return requestClient.get<{
    accountPassword: boolean;
    systemMessage: boolean;
    todoTask: boolean;
  }>('/user/notification/settings');
}

/**
 * 更新通知设置
 */
export async function updateNotificationSettingsApi(data: {
  accountPassword?: boolean;
  systemMessage?: boolean;
  todoTask?: boolean;
}) {
  return requestClient.post('/user/notification/update', data);
}
