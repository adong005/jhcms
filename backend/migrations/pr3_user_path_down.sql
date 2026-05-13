-- PR-3 回滚：删除 path 列
ALTER TABLE users DROP INDEX idx_users_path;
ALTER TABLE users DROP COLUMN path;
