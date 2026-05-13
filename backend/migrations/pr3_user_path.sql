-- PR-3: users 表加 path 字段（物化路径）
-- 回滚：pr3_user_path_down.sql

-- 1. 加列
ALTER TABLE users
  ADD COLUMN path VARCHAR(512) NOT NULL DEFAULT '' COMMENT '物化路径，如 /rootAdminId/agentId/userId/',
  ADD INDEX idx_users_path (path);

-- 2. 回填 super_admin：path = '/'
UPDATE users SET path = '/' WHERE is_admin = 1;

-- 3. 回填 admin（顶层租户管理员）：path = '/{自己id}/'
UPDATE users SET path = CONCAT('/', id, '/') WHERE role = 'admin' AND (parent_id IS NULL OR parent_id = '');

-- 4. 回填第一层子用户（parent 是 admin）：path = '/{parent_id}/{自己id}/'
UPDATE users u
  JOIN users p ON p.id = u.parent_id
  SET u.path = CONCAT(p.path, u.id, '/')
  WHERE u.path = '' AND p.path <> '';

-- 5. 若层级更深，重复执行步骤4直到没有空 path（当前数据最多两层，执行一次即可）
UPDATE users u
  JOIN users p ON p.id = u.parent_id
  SET u.path = CONCAT(p.path, u.id, '/')
  WHERE u.path = '' AND p.path <> '';

-- 6. 兜底：仍为空的用户（孤立数据）设置为 '/{自己id}/'
UPDATE users SET path = CONCAT('/', id, '/') WHERE path = '';
