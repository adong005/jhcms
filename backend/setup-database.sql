-- ADCMS 数据库初始化脚本
-- 在 MySQL 服务器上执行此脚本

-- 1. 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS adcms DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 2. 创建用户并授权（允许从任何主机连接）
CREATE USER IF NOT EXISTS 'jhcms'@'%' IDENTIFIED BY '4XFiGi8simGJAzK4';
GRANT ALL PRIVILEGES ON adcms.* TO 'jhcms'@'%';

-- 3. 如果只允许特定主机连接，使用以下命令（替换 YOUR_HOST_IP）
-- CREATE USER IF NOT EXISTS 'jhcms'@'YOUR_HOST_IP' IDENTIFIED BY '4XFiGi8simGJAzK4';
-- GRANT ALL PRIVILEGES ON adcms.* TO 'jhcms'@'YOUR_HOST_IP';

-- 4. 刷新权限
FLUSH PRIVILEGES;

-- 5. 验证用户和权限
SELECT user, host FROM mysql.user WHERE user = 'jhcms';
SHOW GRANTS FOR 'jhcms'@'%';
