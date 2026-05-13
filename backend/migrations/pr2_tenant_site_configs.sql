-- PR-2: 新建 tenant_site_configs 表，迁移 users 表站点配置字段
-- 执行顺序：先建表+迁数据，再 drop 列
-- 回滚：pr2_tenant_site_configs_down.sql

-- 1. 建新表
CREATE TABLE IF NOT EXISTS tenant_site_configs (
  id            CHAR(36)     NOT NULL,
  tenant_id     CHAR(36)     NOT NULL COMMENT '顶层 admin 的 user_id',
  title         VARCHAR(255) DEFAULT '' COMMENT '站点标题',
  keywords      VARCHAR(500) DEFAULT '' COMMENT 'SEO 关键词',
  description   TEXT         COMMENT 'SEO 描述',
  domain        VARCHAR(255) DEFAULT '' COMMENT '绑定域名',
  logo          VARCHAR(255) DEFAULT '' COMMENT 'Logo 路径',
  icp_code      VARCHAR(120) DEFAULT '' COMMENT 'ICP 备案号',
  contact_phone VARCHAR(64)  DEFAULT '' COMMENT '联系电话',
  contact_address VARCHAR(255) DEFAULT '' COMMENT '联系地址',
  contact_email VARCHAR(120) DEFAULT '' COMMENT '联系邮箱',
  created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_tenant_id (tenant_id),
  INDEX idx_domain (domain)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户站点配置';

-- 2. 从 users 迁移数据（只迁 admin 级别）
INSERT IGNORE INTO tenant_site_configs
  (id, tenant_id, title, keywords, description, domain, logo, icp_code, contact_phone, contact_address, contact_email)
SELECT
  UUID(), id,
  COALESCE(title, ''), COALESCE(keywords, ''), description,
  COALESCE(domain, ''), COALESCE(logo, ''),
  COALESCE(icp_code, ''), COALESCE(contact_phone, ''),
  COALESCE(contact_address, ''), COALESCE(contact_email, '')
FROM users
WHERE role IN ('super_admin', 'admin');

-- 3. 删除 users 表中的站点配置列
ALTER TABLE users
  DROP COLUMN title,
  DROP COLUMN keywords,
  DROP COLUMN description,
  DROP COLUMN domain,
  DROP COLUMN logo,
  DROP COLUMN icp_code,
  DROP COLUMN contact_phone,
  DROP COLUMN contact_address,
  DROP COLUMN contact_email;
