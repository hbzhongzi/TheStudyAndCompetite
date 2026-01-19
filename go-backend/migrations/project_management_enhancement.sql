-- 项目管理模块功能增强 - 数据库增量更新脚本
-- 执行时间: 2024年
-- 说明: 完善项目状态管理、生命周期管理、成果文件管理、分类管理、审核流程

-- =============================================
-- 1. 项目状态管理增强
-- =============================================

-- 修改projects表的status字段，添加新的状态值
ALTER TABLE projects 
MODIFY COLUMN status ENUM('draft','submitted','approved','rejected','archived','in_progress','completed','suspended','need_revision') DEFAULT 'draft' COMMENT '项目状态';

-- 添加状态变更原因字段
ALTER TABLE projects 
ADD COLUMN status_change_reason TEXT NULL COMMENT '状态变更原因',
ADD COLUMN status_changed_by BIGINT UNSIGNED NULL COMMENT '状态变更操作人ID',
ADD COLUMN status_changed_at DATETIME NULL COMMENT '状态变更时间';

-- 添加外键约束
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_status_changer 
FOREIGN KEY (status_changed_by) REFERENCES users(id);

-- =============================================
-- 2. 项目生命周期管理增强
-- =============================================

-- 添加项目进度相关字段
ALTER TABLE projects 
ADD COLUMN start_date DATETIME NULL COMMENT '项目开始时间',
ADD COLUMN expected_end_date DATETIME NULL COMMENT '预计完成时间',
ADD COLUMN actual_end_date DATETIME NULL COMMENT '实际完成时间',
ADD COLUMN progress INT DEFAULT 0 COMMENT '项目进度(0-100)',
ADD COLUMN milestone_count INT DEFAULT 0 COMMENT '里程碑数量',
ADD COLUMN is_extended BOOLEAN DEFAULT FALSE COMMENT '是否延期',
ADD COLUMN extension_reason TEXT NULL COMMENT '延期原因',
ADD COLUMN extension_count INT DEFAULT 0 COMMENT '延期次数';

-- 创建项目里程碑表
CREATE TABLE project_milestones (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    title VARCHAR(200) NOT NULL COMMENT '里程碑标题',
    description TEXT NULL COMMENT '里程碑描述',
    due_date DATETIME NOT NULL COMMENT '预计完成时间',
    completed_date DATETIME NULL COMMENT '实际完成时间',
    status ENUM('pending','in_progress','completed','overdue') DEFAULT 'pending' COMMENT '里程碑状态',
    progress INT DEFAULT 0 COMMENT '完成进度(0-100)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- 创建项目延期申请表
CREATE TABLE project_extensions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    applicant_id BIGINT UNSIGNED NOT NULL COMMENT '申请人ID',
    reason TEXT NOT NULL COMMENT '延期原因',
    original_end_date DATETIME NOT NULL COMMENT '原定结束时间',
    requested_end_date DATETIME NOT NULL COMMENT '申请延期到的时间',
    status ENUM('pending','approved','rejected') DEFAULT 'pending' COMMENT '申请状态',
    reviewer_id BIGINT UNSIGNED NULL COMMENT '审核人ID',
    review_comments TEXT NULL COMMENT '审核意见',
    reviewed_at DATETIME NULL COMMENT '审核时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (applicant_id) REFERENCES users(id),
    FOREIGN KEY (reviewer_id) REFERENCES users(id)
);

-- =============================================
-- 3. 成果文件管理增强
-- =============================================

-- 修改project_files表，添加文件类型和审核状态
ALTER TABLE project_files 
ADD COLUMN file_type ENUM('proposal','midterm','final','achievement','other') DEFAULT 'other' COMMENT '文件类型',
ADD COLUMN file_version VARCHAR(20) DEFAULT '1.0' COMMENT '文件版本',
ADD COLUMN review_status ENUM('pending','approved','rejected') DEFAULT 'pending' COMMENT '审核状态',
ADD COLUMN review_comments TEXT NULL COMMENT '审核意见',
ADD COLUMN reviewed_by BIGINT UNSIGNED NULL COMMENT '审核人ID',
ADD COLUMN reviewed_at DATETIME NULL COMMENT '审核时间',
ADD COLUMN file_size BIGINT DEFAULT 0 COMMENT '文件大小(字节)',
ADD COLUMN download_count INT DEFAULT 0 COMMENT '下载次数',
ADD COLUMN is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开';

-- 添加外键约束
ALTER TABLE project_files 
ADD CONSTRAINT fk_project_files_reviewer 
FOREIGN KEY (reviewed_by) REFERENCES users(id);

-- 创建文件类型配置表
CREATE TABLE file_type_configs (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    display_name VARCHAR(100) NOT NULL COMMENT '显示名称',
    description TEXT NULL COMMENT '描述',
    is_required BOOLEAN DEFAULT FALSE COMMENT '是否必需',
    max_file_size BIGINT DEFAULT 10485760 COMMENT '最大文件大小(字节)',
    allowed_extensions TEXT NULL COMMENT '允许的文件扩展名(逗号分隔)',
    sort_order INT DEFAULT 0 COMMENT '排序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_file_type (file_type)
);

-- 插入默认文件类型配置
INSERT INTO file_type_configs (file_type, display_name, description, is_required, max_file_size, allowed_extensions, sort_order) VALUES
('proposal', '开题报告', '项目开题报告，包含项目背景、目标、计划等', TRUE, 20971520, 'pdf,doc,docx', 1),
('midterm', '中期报告', '项目中期进展报告', FALSE, 20971520, 'pdf,doc,docx', 2),
('final', '结题报告', '项目结题总结报告', TRUE, 20971520, 'pdf,doc,docx', 3),
('achievement', '成果展示', '项目成果展示材料', FALSE, 52428800, 'pdf,doc,docx,ppt,pptx,zip,rar', 4),
('other', '其他材料', '其他相关材料', FALSE, 10485760, 'pdf,doc,docx,jpg,jpeg,png,gif', 5);

-- =============================================
-- 4. 项目分类管理增强
-- =============================================

-- 修改project_types表，添加层级管理
ALTER TABLE project_types 
ADD COLUMN parent_id BIGINT UNSIGNED NULL COMMENT '父分类ID',
ADD COLUMN level INT DEFAULT 1 COMMENT '分类层级',
ADD COLUMN sort_order INT DEFAULT 0 COMMENT '排序',
ADD COLUMN is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
ADD COLUMN icon VARCHAR(100) NULL COMMENT '分类图标',
ADD COLUMN color VARCHAR(20) NULL COMMENT '分类颜色',
ADD COLUMN project_count INT DEFAULT 0 COMMENT '项目数量';

-- 添加自引用外键
ALTER TABLE project_types 
ADD CONSTRAINT fk_project_types_parent 
FOREIGN KEY (parent_id) REFERENCES project_types(id);

-- 创建分类路径表（用于快速查询分类树）
CREATE TABLE project_type_paths (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    ancestor_id BIGINT UNSIGNED NOT NULL COMMENT '祖先分类ID',
    descendant_id BIGINT UNSIGNED NOT NULL COMMENT '后代分类ID',
    depth INT NOT NULL COMMENT '层级深度',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ancestor_id) REFERENCES project_types(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES project_types(id) ON DELETE CASCADE,
    UNIQUE KEY unique_path (ancestor_id, descendant_id)
);

-- 创建分类统计表
CREATE TABLE project_type_stats (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    type_id BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    total_projects INT DEFAULT 0 COMMENT '总项目数',
    draft_projects INT DEFAULT 0 COMMENT '草稿项目数',
    submitted_projects INT DEFAULT 0 COMMENT '已提交项目数',
    approved_projects INT DEFAULT 0 COMMENT '已通过项目数',
    in_progress_projects INT DEFAULT 0 COMMENT '进行中项目数',
    completed_projects INT DEFAULT 0 COMMENT '已完成项目数',
    rejected_projects INT DEFAULT 0 COMMENT '已驳回项目数',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (type_id) REFERENCES project_types(id) ON DELETE CASCADE,
    UNIQUE KEY unique_type (type_id)
);

-- =============================================
-- 5. 审核流程增强
-- =============================================

-- 修改project_reviews表，添加多级审核支持
ALTER TABLE project_reviews 
ADD COLUMN review_level INT DEFAULT 1 COMMENT '审核级别(1:指导教师,2:学院,3:学校)',
ADD COLUMN review_order INT DEFAULT 1 COMMENT '审核顺序',
ADD COLUMN is_required BOOLEAN DEFAULT TRUE COMMENT '是否必需审核',
ADD COLUMN deadline DATETIME NULL COMMENT '审核截止时间',
ADD COLUMN auto_approve BOOLEAN DEFAULT FALSE COMMENT '超时是否自动通过',
ADD COLUMN review_duration INT DEFAULT 0 COMMENT '审核耗时(分钟)',
ADD COLUMN is_urgent BOOLEAN DEFAULT FALSE COMMENT '是否紧急';

-- 创建审核流程配置表
CREATE TABLE project_review_flows (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    project_type_id BIGINT UNSIGNED NULL COMMENT '项目类型ID',
    review_level INT NOT NULL COMMENT '审核级别',
    reviewer_role VARCHAR(50) NOT NULL COMMENT '审核角色',
    reviewer_department VARCHAR(100) NULL COMMENT '审核部门',
    review_order INT NOT NULL COMMENT '审核顺序',
    is_required BOOLEAN DEFAULT TRUE COMMENT '是否必需',
    deadline_hours INT DEFAULT 72 COMMENT '审核时限(小时)',
    auto_approve BOOLEAN DEFAULT FALSE COMMENT '超时自动通过',
    can_delegate BOOLEAN DEFAULT FALSE COMMENT '是否可以委托他人审核',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_type_id) REFERENCES project_types(id)
);

-- 创建审核委托表
CREATE TABLE review_delegations (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    original_reviewer_id BIGINT UNSIGNED NOT NULL COMMENT '原审核人ID',
    delegated_reviewer_id BIGINT UNSIGNED NOT NULL COMMENT '被委托人ID',
    project_id BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    reason TEXT NOT NULL COMMENT '委托原因',
    start_date DATETIME NOT NULL COMMENT '委托开始时间',
    end_date DATETIME NOT NULL COMMENT '委托结束时间',
    status ENUM('active','expired','cancelled') DEFAULT 'active' COMMENT '委托状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (original_reviewer_id) REFERENCES users(id),
    FOREIGN KEY (delegated_reviewer_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- 创建审核提醒表
CREATE TABLE review_reminders (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    review_id BIGINT UNSIGNED NOT NULL COMMENT '审核记录ID',
    reviewer_id BIGINT UNSIGNED NOT NULL COMMENT '审核人ID',
    reminder_type ENUM('deadline_approaching','overdue','urgent') NOT NULL COMMENT '提醒类型',
    message TEXT NOT NULL COMMENT '提醒消息',
    is_sent BOOLEAN DEFAULT FALSE COMMENT '是否已发送',
    sent_at DATETIME NULL COMMENT '发送时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (review_id) REFERENCES project_reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES users(id)
);

-- =============================================
-- 6. 通知系统基础结构
-- =============================================

-- 创建项目通知表
CREATE TABLE project_notifications (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '接收用户ID',
    type VARCHAR(50) NOT NULL COMMENT '通知类型',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content TEXT NOT NULL COMMENT '通知内容',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    read_at DATETIME NULL COMMENT '阅读时间',
    priority ENUM('low','normal','high','urgent') DEFAULT 'normal' COMMENT '优先级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建通知模板表
CREATE TABLE notification_templates (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    template_key VARCHAR(100) NOT NULL COMMENT '模板键',
    title_template VARCHAR(200) NOT NULL COMMENT '标题模板',
    content_template TEXT NOT NULL COMMENT '内容模板',
    variables TEXT NULL COMMENT '可用变量说明',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_template_key (template_key)
);

-- 插入默认通知模板
INSERT INTO notification_templates (template_key, title_template, content_template, variables) VALUES
('project_submitted', '项目提交通知', '您的项目《{project_title}》已成功提交，等待审核。', 'project_title,student_name,submitted_at'),
('project_approved', '项目通过通知', '恭喜！您的项目《{project_title}》已通过审核。', 'project_title,student_name,approved_at,reviewer_name'),
('project_rejected', '项目驳回通知', '您的项目《{project_title}》未通过审核，请查看驳回原因并修改后重新提交。', 'project_title,student_name,rejected_at,reviewer_name,reject_reason'),
('review_assigned', '审核任务分配', '您有一个新的审核任务：项目《{project_title}》，请在{deadline}前完成审核。', 'project_title,student_name,deadline'),
('deadline_reminder', '截止时间提醒', '项目《{project_title}》的{task_type}即将截止，请及时处理。', 'project_title,task_type,deadline');

-- =============================================
-- 7. 索引优化
-- =============================================

-- 为新增字段添加索引
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_progress ON projects(progress);
CREATE INDEX idx_projects_start_date ON projects(start_date);
CREATE INDEX idx_projects_expected_end_date ON projects(expected_end_date);
CREATE INDEX idx_projects_status_changed_at ON projects(status_changed_at);

CREATE INDEX idx_project_files_type ON project_files(file_type);
CREATE INDEX idx_project_files_review_status ON project_files(review_status);
CREATE INDEX idx_project_files_upload_time ON project_files(upload_time);

CREATE INDEX idx_project_types_parent ON project_types(parent_id);
CREATE INDEX idx_project_types_level ON project_types(level);
CREATE INDEX idx_project_types_active ON project_types(is_active);

CREATE INDEX idx_project_reviews_level ON project_reviews(review_level);
CREATE INDEX idx_project_reviews_deadline ON project_reviews(deadline);
CREATE INDEX idx_project_reviews_urgent ON project_reviews(is_urgent);

CREATE INDEX idx_project_notifications_user ON project_notifications(user_id);
CREATE INDEX idx_project_notifications_type ON project_notifications(type);
CREATE INDEX idx_project_notifications_read ON project_notifications(is_read);

-- =============================================
-- 8. 数据迁移和初始化
-- =============================================

-- 更新现有项目的状态变更信息
UPDATE projects 
SET status_changed_by = created_by, 
    status_changed_at = created_at 
WHERE status_changed_by IS NULL;

-- 初始化分类统计
INSERT INTO project_type_stats (type_id, total_projects)
SELECT category_id, COUNT(*) 
FROM projects 
WHERE category_id IS NOT NULL AND deleted = FALSE
GROUP BY category_id
ON DUPLICATE KEY UPDATE 
    total_projects = VALUES(total_projects);

-- 创建分类路径数据（为现有分类创建自引用路径）
INSERT INTO project_type_paths (ancestor_id, descendant_id, depth)
SELECT id, id, 0 FROM project_types;

-- =============================================
-- 执行完成提示
-- =============================================
SELECT '项目管理模块功能增强数据库更新完成！' AS message; 