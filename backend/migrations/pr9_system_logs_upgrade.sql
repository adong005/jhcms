-- PR-9: system_logs 表升级
-- 新增字段：user_id, path, method, url, target_id, user_agent, request_id, log_type, status_code
-- 新增索引：path, (tenant_id, created_at), (action, status), user_id

ALTER TABLE system_logs
  ADD COLUMN request_id   CHAR(36)      NOT NULL DEFAULT '' AFTER id,
  ADD COLUMN user_id      CHAR(36)      NOT NULL DEFAULT '' AFTER tenant_id,
  ADD COLUMN path         VARCHAR(512)  NOT NULL DEFAULT '' AFTER user_id,
  ADD COLUMN method       VARCHAR(10)   NOT NULL DEFAULT '' AFTER ip,
  ADD COLUMN url          VARCHAR(500)  NOT NULL DEFAULT '' AFTER method,
  ADD COLUMN user_agent   VARCHAR(255)  NOT NULL DEFAULT '' AFTER url,
  ADD COLUMN target_id    VARCHAR(64)   NOT NULL DEFAULT '' AFTER description,
  ADD COLUMN log_type     VARCHAR(20)   NOT NULL DEFAULT 'api' AFTER status,
  ADD COLUMN status_code  SMALLINT      NOT NULL DEFAULT 0 AFTER log_type;

-- 补充索引
ALTER TABLE system_logs
  ADD INDEX idx_logs_path        (path(191)),
  ADD INDEX idx_logs_tenant_time (tenant_id, created_at),
  ADD INDEX idx_logs_action_st   (action, status),
  ADD INDEX idx_logs_user_id     (user_id);

-- 回填 user_id / path（根据 username 关联 users 表）
UPDATE system_logs sl
  LEFT JOIN users u ON u.username = sl.username
  SET sl.user_id = COALESCE(u.id, ''),
      sl.path    = COALESCE(u.path, '');
