export interface SiteMeta {
  title: string;
  keywords: string;
  description: string;
}

function parseHost() {
  const host = window.location.hostname.toLowerCase();
  const parts = host.split('.');
  const subdomain = parts.length > 2 ? parts[0] : '';
  const rootDomain = parts.length > 2 ? parts.slice(1).join('.') : host;
  return { host, subdomain, rootDomain };
}

export async function loadSiteMetaByHost(): Promise<SiteMeta | null> {
  const { host, subdomain, rootDomain } = parseHost();
  const apiBase = (import.meta.env.VITE_PUBLIC_API_BASE || '').trim();
  if (!apiBase) return null;

  try {
    const url = new URL('/api/public/site-meta', apiBase);
    url.searchParams.set('host', host);
    url.searchParams.set('subdomain', subdomain);
    url.searchParams.set('rootDomain', rootDomain);
    const resp = await fetch(url.toString());
    if (!resp.ok) return null;
    const data = await resp.json();
    return data?.data || null;
  } catch {
    return null;
  }
}
