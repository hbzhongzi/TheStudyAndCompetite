-- 添加项目表缺失字段的迁移脚本（简化版）
-- 执行时间: 2025年8月
-- 说明: 为projects表添加缺失的字段，使其与Go模型保持一致

USE cloud_dream_system;

-- 添加缺失的字段到projects表
-- 注意：如果字段已存在，会报错，这是正常的

-- 添加 category_id 字段
ALTER TABLE projects ADD COLUMN category_id BIGINT UNSIGNED NULL COMMENT '项目分类ID';

-- 添加 expected_end_date 字段
ALTER TABLE projects ADD COLUMN expected_end_date DATETIME NULL COMMENT '预计完成时间';

-- 添加 actual_end_date 字段
ALTER TABLE projects ADD COLUMN actual_end_date DATETIME NULL COMMENT '实际完成时间';

-- 添加 progress 字段
ALTER TABLE projects ADD COLUMN progress INT DEFAULT 0 COMMENT '项目进度(0-100)';

-- 添加 is_extended 字段
ALTER TABLE projects ADD COLUMN is_extended BOOLEAN DEFAULT FALSE COMMENT '是否延期';

-- 添加 extension_count 字段
ALTER TABLE projects ADD COLUMN extension_count INT DEFAULT 0 COMMENT '延期次数';

-- 添加 force_status_reason 字段
ALTER TABLE projects ADD COLUMN force_status_reason TEXT NULL COMMENT '强制状态变更原因';

-- 添加 level 字段
ALTER TABLE projects ADD COLUMN level VARCHAR(50) NULL COMMENT '项目级别';

-- 添加 created_by 字段
ALTER TABLE projects ADD COLUMN created_by BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '创建者ID';

-- 添加 created_at 字段
ALTER TABLE projects ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- 添加 updated_at 字段
ALTER TABLE projects ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

-- 为现有项目设置默认值
UPDATE projects SET
    created_by = student_id,
    level = '普通',
    progress = 0,
    is_extended = FALSE,
    extension_count = 0
WHERE created_by IS NULL OR level IS NULL;

-- 添加外键约束（如果不存在的话）
-- 注意：如果约束已存在，会报错，这是正常的

-- 添加外键约束到 project_types 表
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_category 
FOREIGN KEY (category_id) REFERENCES project_types(id) ON DELETE SET NULL;

-- 添加外键约束到 users 表
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_creator 
FOREIGN KEY (created_by) REFERENCES users(id);

-- 添加索引（如果不存在的话）
-- 注意：如果索引已存在，会报错，这是正常的

-- 添加索引以提高查询性能
CREATE INDEX idx_projects_category_id ON projects(category_id);
CREATE INDEX idx_projects_created_by ON projects(created_by);
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_expected_end_date ON projects(expected_end_date);
CREATE INDEX idx_projects_level ON projects(level);

-- 显示更新后的表结构
DESCRIBE projects;

-- 显示添加的索引
SHOW INDEX FROM projects; 