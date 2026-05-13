-- PR-2 回滚：恢复 users 表字段，删除 tenant_site_configs 表
ALTER TABLE users
  ADD COLUMN title         VARCHAR(255) DEFAULT '',
  ADD COLUMN keywords      VARCHAR(500) DEFAULT '',
  ADD COLUMN description   TEXT,
  ADD COLUMN domain        VARCHAR(255) DEFAULT '',
  ADD COLUMN logo          VARCHAR(255) DEFAULT '',
  ADD COLUMN icp_code      VARCHAR(120) DEFAULT '',
  ADD COLUMN contact_phone VARCHAR(64)  DEFAULT '',
  ADD COLUMN contact_address VARCHAR(255) DEFAULT '',
  ADD COLUMN contact_email VARCHAR(120) DEFAULT '';

-- 从 tenant_site_configs 回填
UPDATE users u
  JOIN tenant_site_configs c ON c.tenant_id = u.id
  SET u.title = c.title, u.keywords = c.keywords, u.description = c.description,
      u.domain = c.domain, u.logo = c.logo, u.icp_code = c.icp_code,
      u.contact_phone = c.contact_phone, u.contact_address = c.contact_address,
      u.contact_email = c.contact_email;

DROP TABLE IF EXISTS tenant_site_configs;
