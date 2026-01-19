-- 添加项目表缺失字段的迁移脚本
-- 执行时间: 2025年8月
-- 说明: 为projects表添加缺失的字段，使其与Go模型保持一致

USE cloud_dream_system;

-- 检查projects表是否存在
SELECT COUNT(*) FROM information_schema.tables 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects';

-- 添加缺失的字段到projects表（逐个添加以避免重复字段错误）
SET @sql = '';

-- 检查并添加 category_id 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'category_id';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN category_id BIGINT UNSIGNED NULL COMMENT ''项目分类ID'';');
END IF;

-- 检查并添加 expected_end_date 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'expected_end_date';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN expected_end_date DATETIME NULL COMMENT ''预计完成时间'';');
END IF;

-- 检查并添加 actual_end_date 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'actual_end_date';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN actual_end_date DATETIME NULL COMMENT ''实际完成时间'';');
END IF;

-- 检查并添加 progress 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'progress';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN progress INT DEFAULT 0 COMMENT ''项目进度(0-100)'';');
END IF;

-- 检查并添加 is_extended 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'is_extended';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN is_extended BOOLEAN DEFAULT FALSE COMMENT ''是否延期'';');
END IF;

-- 检查并添加 extension_count 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'extension_count';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN extension_count INT DEFAULT 0 COMMENT ''延期次数'';');
END IF;

-- 检查并添加 force_status_reason 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'force_status_reason';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN force_status_reason TEXT NULL COMMENT ''强制状态变更原因'';');
END IF;

-- 检查并添加 level 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'level';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN level VARCHAR(50) NULL COMMENT ''项目级别'';');
END IF;

-- 检查并添加 created_by 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'created_by';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN created_by BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT ''创建者ID'';');
END IF;

-- 检查并添加 created_at 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'created_at';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'';');
END IF;

-- 检查并添加 updated_at 字段
SELECT COUNT(*) INTO @exists FROM information_schema.columns 
WHERE table_schema = 'cloud_dream_system' AND table_name = 'projects' AND column_name = 'updated_at';
IF @exists = 0 THEN
    SET @sql = CONCAT(@sql, 'ALTER TABLE projects ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''更新时间'';');
END IF;

-- 执行SQL语句
IF @sql != '' THEN
    SET @sql = CONCAT('SET @sql = '''', @sql, '''';');
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END IF;

-- 为现有项目设置默认值
UPDATE projects SET
    created_by = student_id,
    level = '普通',
    progress = 0,
    is_extended = FALSE,
    extension_count = 0
WHERE created_by IS NULL OR level IS NULL;

-- 添加外键约束
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_category 
FOREIGN KEY (category_id) REFERENCES project_types(id) ON DELETE SET NULL;

ALTER TABLE projects 
ADD CONSTRAINT fk_projects_creator 
FOREIGN KEY (created_by) REFERENCES users(id);

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_projects_category_id ON projects(category_id);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
CREATE INDEX IF NOT EXISTS idx_projects_expected_end_date ON projects(expected_end_date);
CREATE INDEX IF NOT EXISTS idx_projects_level ON projects(level);

-- 显示更新后的表结构
DESCRIBE projects;

-- 显示添加的索引
SHOW INDEX FROM projects; 