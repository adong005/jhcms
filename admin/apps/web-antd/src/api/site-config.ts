import { requestClient } from '#/api/request';

export namespace SiteConfigApi {
  export interface SiteConfig {
    title: string;
    keywords: string;
    description: string;
    domain: string;
    logo?: string;
    icpCode: string;
    contactPhone: string;
    contactAddress: string;
    contactEmail: string;
  }

  export interface UpdateSiteConfigParams {
    userId: string;
    title: string;
    keywords: string;
    description: string;
    domain: string;
    logo?: string;
    icpCode: string;
    contactPhone: string;
    contactAddress: string;
    contactEmail: string;
  }
}

/**
 * 获取网站配置
 */
export async function getSiteConfigApi() {
  return requestClient.get<SiteConfigApi.SiteConfig>('/site-config');
}

/**
 * 更新网站配置
 */
export async function updateSiteConfigApi(data: SiteConfigApi.UpdateSiteConfigParams) {
  return requestClient.post('/site-config', data);
}

/**
 * 上传网站 Logo
 */
export async function uploadSiteLogoApi(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  return requestClient.post<{ logo: string }>('/site-config/logo/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}
