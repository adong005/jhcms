-- PR-9 回滚：移除 system_logs 新增字段与索引
ALTER TABLE system_logs
  DROP INDEX idx_logs_path,
  DROP INDEX idx_logs_tenant_time,
  DROP INDEX idx_logs_action_st,
  DROP INDEX idx_logs_user_id,
  DROP COLUMN request_id,
  DROP COLUMN user_id,
  DROP COLUMN path,
  DROP COLUMN method,
  DROP COLUMN url,
  DROP COLUMN user_agent,
  DROP COLUMN target_id,
  DROP COLUMN log_type,
  DROP COLUMN status_code;
