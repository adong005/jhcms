-- PR-4 回滚
ALTER TABLE menus DROP INDEX idx_menus_path_chain;
ALTER TABLE menus DROP COLUMN path_chain;
ALTER TABLE menus DROP INDEX idx_menus_parent_id;
ALTER TABLE menus CHANGE parent_id parent_id_menu CHAR(36) NULL DEFAULT NULL;
