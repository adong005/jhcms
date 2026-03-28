-- 填充菜单数据到 menus 表
-- 参考 Mock 数据结构

-- 清空现有数据
TRUNCATE TABLE menus;

-- 插入一级菜单
INSERT INTO menus (id, name, path, component, icon, parent_id_menu, sort_order, status, created_at, updated_at) VALUES
(1, '工作台', '/analytics', '/dashboard/analytics/index', 'mdi:view-dashboard', NULL, 1, 1, '2024-01-01 10:00:00', '2024-03-20 14:00:00'),
(2, '用户管理', '/users/list', '/users/list', 'lucide:users', NULL, 2, 1, '2024-01-02 10:00:00', '2024-03-21 14:00:00'),
(3, '角色管理', '/roles/list', '/roles/list', 'lucide:shield', NULL, 3, 1, '2024-01-03 10:00:00', '2024-03-22 14:00:00'),
(4, '菜单管理', '/menus/list', '/menus/list', 'lucide:menu', NULL, 4, 1, '2024-01-04 10:00:00', '2024-03-23 14:00:00'),
(25, '网站配置', '/site-config', '/site-config/index', 'lucide:settings', NULL, 5, 1, '2024-03-25 12:00:00', '2024-03-25 12:00:00'),
(5, '信息管理', '/info', '', 'mdi:information', NULL, 6, 1, '2024-01-05 10:00:00', '2024-03-24 14:00:00'),
(7, '站群管理', '/site-group', '', 'lucide:network', NULL, 7, 1, '2024-03-25 12:00:00', '2024-03-25 12:00:00'),
(8, '表单管理', '/form-manage', '', 'lucide:file-text', NULL, 8, 1, '2024-03-25 12:20:00', '2024-03-25 12:20:00');

-- 插入二级菜单（信息管理的子菜单）
INSERT INTO menus (id, name, path, component, icon, parent_id_menu, sort_order, status, created_at, updated_at) VALUES
(51, '信息分类', '/info/category/list', '/info/category/list', 'mdi:folder-multiple', 5, 1, 1, '2024-01-05 10:10:00', '2024-03-24 14:10:00'),
(52, '信息列表', '/info/list', '/info/list', 'mdi:file-document-multiple', 5, 2, 1, '2024-01-05 11:00:00', '2024-03-24 15:00:00'),
(53, '发布信息', '/info/publish', '/info/publish', 'mdi:plus-circle', 5, 3, 1, '2024-01-05 12:00:00', '2024-03-24 16:00:00');

-- 插入二级菜单（站群管理的子菜单）
INSERT INTO menus (id, name, path, component, icon, parent_id_menu, sort_order, status, created_at, updated_at) VALUES
(71, '站群列表', '/site-group/list', '/site-group/list', 'lucide:list', 7, 1, 1, '2024-03-25 12:10:00', '2024-03-25 12:10:00');

-- 插入二级菜单（表单管理的子菜单）
INSERT INTO menus (id, name, path, component, icon, parent_id_menu, sort_order, status, created_at, updated_at) VALUES
(81, '表单列表', '/form-manage/list', '/form-manage/list', 'lucide:list', 8, 1, 1, '2024-03-25 12:30:00', '2024-03-25 12:30:00');

-- 验证插入结果
SELECT COUNT(*) as total_menus FROM menus;
SELECT id, name, path, parent_id_menu, sort_order, status FROM menus ORDER BY sort_order, id;
