-- PR-1 回滚：从备份恢复（无法自动重建，需从备份文件 restore）
-- mysqldump 备份位置：backend/backup/YYYYMMDD_HHMMSS_before_phase1.sql
-- 恢复命令示例：
-- mysql -u jhcms -p jhcms < backend/backup/YYYYMMDD_HHMMSS_before_phase1.sql
SELECT 'Please restore from backup file at backend/backup/' AS notice;
