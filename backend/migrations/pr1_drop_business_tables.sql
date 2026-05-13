-- PR-1: 删除已废弃的业务表（info/info_category/site_group/form/city）
-- 执行前请确保已 mysqldump 备份
-- 回滚脚本：pr1_drop_business_tables_down.sql

DROP TABLE IF EXISTS forms;
DROP TABLE IF EXISTS infos;
DROP TABLE IF EXISTS info_categories;
DROP TABLE IF EXISTS site_groups;
DROP TABLE IF EXISTS city_list;
