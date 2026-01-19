-- 添加缺失的deleted列到projects表
-- 执行时间: 2024年
-- 说明: 为projects表添加deleted列以支持软删除功能

-- 检查projects表是否存在deleted列
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'projects' 
     AND COLUMN_NAME = 'deleted') > 0,
    'SELECT "deleted列已存在" as message',
    'ALTER TABLE projects ADD COLUMN deleted BOOLEAN DEFAULT FALSE COMMENT "软删除标记"'
));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 为deleted列添加索引以提高查询性能
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'projects' 
     AND INDEX_NAME = 'idx_projects_deleted') > 0,
    'SELECT "deleted列索引已存在" as message',
    'CREATE INDEX idx_projects_deleted ON projects(deleted)'
));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查users表是否存在deleted列
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'users' 
     AND COLUMN_NAME = 'deleted') > 0,
    'SELECT "users表deleted列已存在" as message',
    'ALTER TABLE users ADD COLUMN deleted BOOLEAN DEFAULT FALSE COMMENT "软删除标记"'
));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 为users表的deleted列添加索引
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'users' 
     AND INDEX_NAME = 'idx_users_deleted') > 0,
    'SELECT "users表deleted列索引已存在" as message',
    'CREATE INDEX idx_users_deleted ON users(deleted)'
));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 显示执行结果
SELECT 'Migration completed successfully' as status; 