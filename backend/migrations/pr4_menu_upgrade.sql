-- PR-4: 菜单表升级——重命名 parent_id_menu 为 parent_id，加 path_chain 物化路径
-- 回滚：pr4_menu_upgrade_down.sql

-- 1. 重命名字段
ALTER TABLE menus CHANGE parent_id_menu parent_id CHAR(36) NULL DEFAULT NULL;
ALTER TABLE menus ADD INDEX idx_menus_parent_id (parent_id);

-- 2. 加 path_chain 列
ALTER TABLE menus
  ADD COLUMN path_chain VARCHAR(512) NOT NULL DEFAULT '' COMMENT '菜单物化路径，如 /rootId/subId/',
  ADD INDEX idx_menus_path_chain (path_chain);

-- 3. 回填根菜单（parent_id IS NULL）
UPDATE menus SET path_chain = CONCAT('/', id, '/') WHERE parent_id IS NULL;

-- 4. 逐层回填（执行多次直到没有 affected rows；当前数据最多2层，执行3次以覆盖）
UPDATE menus c JOIN menus p ON p.id = c.parent_id
  SET c.path_chain = CONCAT(p.path_chain, c.id, '/')
  WHERE c.path_chain = '' AND p.path_chain <> '';

UPDATE menus c JOIN menus p ON p.id = c.parent_id
  SET c.path_chain = CONCAT(p.path_chain, c.id, '/')
  WHERE c.path_chain = '' AND p.path_chain <> '';

UPDATE menus c JOIN menus p ON p.id = c.parent_id
  SET c.path_chain = CONCAT(p.path_chain, c.id, '/')
  WHERE c.path_chain = '' AND p.path_chain <> '';

-- 5. 兜底：仍为空（孤立菜单）
UPDATE menus SET path_chain = CONCAT('/', id, '/') WHERE path_chain = '';
