-- 云梦高校学生科研与竞赛项目管理系统 - 完整数据库设置脚本
-- 包含所有表结构、索引、触发器、存储过程、视图和测试数据

-- 创建数据库
CREATE DATABASE IF NOT EXISTS cloud_dream_system
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE cloud_dream_system;

-- ==================== 1. 用户管理模块 ====================

-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    status ENUM('active','inactive','suspended') DEFAULT 'active' COMMENT '用户状态',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_status (status),
    INDEX idx_users_create_time (create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 用户角色表
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    role_name VARCHAR(50) UNIQUE NOT NULL COMMENT '角色名称',
    description VARCHAR(255) COMMENT '角色描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_roles_name (role_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色表';

-- 用户角色关联表
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_user_role (user_id, role_id),
    INDEX idx_user_roles_user_id (user_id),
    INDEX idx_user_roles_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 用户档案表
CREATE TABLE user_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '档案ID',
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID',
    real_name VARCHAR(50) COMMENT '真实姓名',
    phone VARCHAR(20) COMMENT '手机号',
    department VARCHAR(100) COMMENT '所属部门/学院',
    student_id VARCHAR(20) COMMENT '学号',
    teacher_id VARCHAR(20) COMMENT '教师工号',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    bio TEXT COMMENT '个人简介',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_user_profiles_user_id (user_id),
    INDEX idx_user_profiles_real_name (real_name),
    INDEX idx_user_profiles_department (department),
    INDEX idx_user_profiles_student_id (student_id),
    INDEX idx_user_profiles_teacher_id (teacher_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户档案表';

-- 登录日志表
CREATE TABLE login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    status ENUM('success','failed') DEFAULT 'success' COMMENT '登录状态',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_login_logs_user_id (user_id),
    INDEX idx_login_logs_login_time (login_time),
    INDEX idx_login_logs_ip_address (ip_address),
    INDEX idx_login_logs_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ==================== 2. 项目管理模块 ====================

-- 项目分类表
CREATE TABLE project_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
    type_name VARCHAR(100) NOT NULL COMMENT '分类名称',
    description TEXT COMMENT '分类描述',
    parent_id BIGINT DEFAULT NULL COMMENT '父分类ID',
    sort_order INT DEFAULT 0 COMMENT '排序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (parent_id) REFERENCES project_types(id) ON DELETE SET NULL,
    
    INDEX idx_project_types_parent_id (parent_id),
    INDEX idx_project_types_sort_order (sort_order),
    INDEX idx_project_types_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目分类表';

-- 项目表
CREATE TABLE projects (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '项目ID',
    title VARCHAR(255) NOT NULL COMMENT '项目标题',
    description TEXT COMMENT '项目描述',
    type_id BIGINT NOT NULL COMMENT '项目分类ID',
    student_id BIGINT NOT NULL COMMENT '学生ID',
    teacher_id BIGINT COMMENT '指导教师ID',
    status ENUM('draft','submitted','reviewing','approved','rejected','completed') DEFAULT 'draft' COMMENT '项目状态',
    submitted_at DATETIME COMMENT '提交时间',
    approved_at DATETIME COMMENT '审批时间',
    approved_by BIGINT COMMENT '审批人ID',
    rejection_reason TEXT COMMENT '拒绝原因',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (type_id) REFERENCES project_types(id) ON DELETE RESTRICT,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    
    INDEX idx_projects_type_id (type_id),
    INDEX idx_projects_student_id (student_id),
    INDEX idx_projects_teacher_id (teacher_id),
    INDEX idx_projects_status (status),
    INDEX idx_projects_submitted_at (submitted_at),
    INDEX idx_projects_approved_at (approved_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- ==================== 3. 竞赛管理模块 ====================

-- 竞赛信息表
CREATE TABLE competitions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '竞赛ID',
    title VARCHAR(255) NOT NULL COMMENT '竞赛标题',
    description TEXT COMMENT '竞赛描述',
    level ENUM('school','provincial','national','international') DEFAULT 'school' COMMENT '竞赛级别',
    category VARCHAR(100) COMMENT '竞赛类别',
    registration_start DATETIME COMMENT '报名开始时间',
    registration_end DATETIME COMMENT '报名结束时间',
    submission_start DATETIME COMMENT '提交开始时间',
    submission_end DATETIME COMMENT '提交结束时间',
    max_participants INT COMMENT '最大参与人数',
    current_participants INT DEFAULT 0 COMMENT '当前参与人数',
    status ENUM('draft','registration','submission','review','completed') DEFAULT 'draft' COMMENT '竞赛状态',
    award_config JSON COMMENT '获奖配置',
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_competitions_level (level),
    INDEX idx_competitions_category (category),
    INDEX idx_competitions_status (status),
    INDEX idx_competitions_registration_start (registration_start),
    INDEX idx_competitions_registration_end (registration_end),
    INDEX idx_competitions_submission_start (submission_start),
    INDEX idx_competitions_submission_end (submission_end),
    INDEX idx_competitions_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛信息表';

-- 竞赛报名记录表
CREATE TABLE competition_registrations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '报名ID',
    competition_id BIGINT NOT NULL COMMENT '竞赛ID',
    student_id BIGINT NOT NULL COMMENT '学生ID',
    team_name VARCHAR(100) COMMENT '团队名称',
    team_leader BIGINT COMMENT '团队负责人ID',
    registration_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '报名时间',
    status ENUM('pending','approved','rejected') DEFAULT 'pending' COMMENT '报名状态',
    approved_by BIGINT COMMENT '审批人ID',
    approved_at DATETIME COMMENT '审批时间',
    rejection_reason TEXT COMMENT '拒绝原因',
    
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (team_leader) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_competition_student (competition_id, student_id),
    INDEX idx_competition_registrations_competition_id (competition_id),
    INDEX idx_competition_registrations_student_id (student_id),
    INDEX idx_competition_registrations_team_leader (team_leader),
    INDEX idx_competition_registrations_status (status),
    INDEX idx_competition_registrations_registration_time (registration_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛报名记录表';

-- 竞赛成果提交表
CREATE TABLE competition_submissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '提交ID',
    competition_id BIGINT NOT NULL COMMENT '竞赛ID',
    student_id BIGINT NOT NULL COMMENT '学生ID',
    title VARCHAR(255) NOT NULL COMMENT '作品标题',
    description TEXT COMMENT '作品描述',
    file_url VARCHAR(500) COMMENT '文件URL',
    file_size BIGINT COMMENT '文件大小',
    version VARCHAR(20) DEFAULT '1.0' COMMENT '版本号',
    locked BOOLEAN DEFAULT FALSE COMMENT '是否锁定',
    submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_competition_submissions_competition_id (competition_id),
    INDEX idx_competition_submissions_student_id (student_id),
    INDEX idx_competition_submissions_submitted_at (submitted_at),
    INDEX idx_competition_submissions_locked (locked)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛成果提交表';

-- 竞赛教师评语表
CREATE TABLE competition_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评语ID',
    submission_id BIGINT NOT NULL COMMENT '提交ID',
    teacher_id BIGINT NOT NULL COMMENT '教师ID',
    feedback TEXT NOT NULL COMMENT '评语内容',
    rating INT COMMENT '评分(1-10)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (submission_id) REFERENCES competition_submissions(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_submission_teacher (submission_id, teacher_id),
    INDEX idx_competition_feedback_submission_id (submission_id),
    INDEX idx_competition_feedback_teacher_id (teacher_id),
    INDEX idx_competition_feedback_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛教师评语表';

-- 竞赛评审教师表
CREATE TABLE competition_judges (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评审教师ID',
    competition_id BIGINT NOT NULL COMMENT '竞赛ID',
    teacher_id BIGINT NOT NULL COMMENT '教师ID',
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    status ENUM('active','inactive') DEFAULT 'active' COMMENT '是否参与',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_competition_teacher (competition_id, teacher_id),
    INDEX idx_competition_judges_competition_id (competition_id),
    INDEX idx_competition_judges_teacher_id (teacher_id),
    INDEX idx_competition_judges_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛评审教师表';

-- 竞赛评分记录表
CREATE TABLE competition_scores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评分ID',
    submission_id BIGINT NOT NULL COMMENT '提交记录ID',
    judge_id BIGINT NOT NULL COMMENT '评审教师ID',
    score DECIMAL(5,2) NOT NULL COMMENT '评分',
    comment TEXT COMMENT '评语',
    scored_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '评分时间',
    
    FOREIGN KEY (submission_id) REFERENCES competition_submissions(id) ON DELETE CASCADE,
    FOREIGN KEY (judge_id) REFERENCES users(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_submission_judge (submission_id, judge_id),
    INDEX idx_competition_scores_submission_id (submission_id),
    INDEX idx_competition_scores_judge_id (judge_id),
    INDEX idx_competition_scores_scored_at (scored_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛评分记录表';

-- 竞赛获奖登记表
CREATE TABLE competition_results (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '获奖ID',
    competition_id BIGINT NOT NULL COMMENT '竞赛ID',
    student_id BIGINT NOT NULL COMMENT '学生ID',
    award_level ENUM('first','second','third','honorable') COMMENT '获奖等级',
    award_name VARCHAR(100) COMMENT '获奖名称',
    score DECIMAL(5,2) COMMENT '最终得分',
    ranking INT COMMENT '排名',
    finalized_by BIGINT COMMENT '最终确认成绩的管理员ID',
    finalized_at DATETIME COMMENT '确认时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (finalized_by) REFERENCES users(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_competition_student_result (competition_id, student_id),
    INDEX idx_competition_results_competition_id (competition_id),
    INDEX idx_competition_results_student_id (student_id),
    INDEX idx_competition_results_award_level (award_level),
    INDEX idx_competition_results_ranking (ranking),
    INDEX idx_competition_results_finalized_by (finalized_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛获奖登记表';

-- 竞赛审计日志表
CREATE TABLE competition_audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '审计日志ID',
    competition_id BIGINT NOT NULL COMMENT '竞赛ID',
    user_id BIGINT NOT NULL COMMENT '操作用户ID',
    action VARCHAR(100) NOT NULL COMMENT '操作类型',
    details TEXT COMMENT '操作详情',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_competition_audit_logs_competition_id (competition_id),
    INDEX idx_competition_audit_logs_user_id (user_id),
    INDEX idx_competition_audit_logs_action (action),
    INDEX idx_competition_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛审计日志表';

-- ==================== 4. 系统管理模块 ====================

-- 系统日志表
CREATE TABLE system_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    user_id BIGINT COMMENT '操作用户ID',
    action VARCHAR(255) NOT NULL COMMENT '执行动作',
    details TEXT COMMENT '详细内容',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理（浏览器信息）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    
    INDEX idx_system_logs_user_id (user_id),
    INDEX idx_system_logs_action (action),
    INDEX idx_system_logs_created_at (created_at),
    INDEX idx_system_logs_ip_address (ip_address),
    INDEX idx_system_logs_action_created_at (action, created_at),
    INDEX idx_system_logs_user_action (user_id, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统日志表';

-- 系统配置表
CREATE TABLE system_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    setting_key VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    setting_value TEXT COMMENT '配置值',
    description VARCHAR(255) COMMENT '配置描述',
    category VARCHAR(50) DEFAULT 'general' COMMENT '配置分类',
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开（前端可见）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_system_settings_key (setting_key),
    INDEX idx_system_settings_category (category),
    INDEX idx_system_settings_public (is_public)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 备份记录表
CREATE TABLE backup_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '备份记录ID',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件路径',
    file_size BIGINT COMMENT '文件大小（字节）',
    backup_type ENUM('full','incremental','manual') DEFAULT 'manual' COMMENT '备份类型',
    status ENUM('pending','in_progress','success','failed') DEFAULT 'pending' COMMENT '备份状态',
    error_message TEXT COMMENT '错误信息',
    created_by BIGINT COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    started_at DATETIME COMMENT '开始时间',
    completed_at DATETIME COMMENT '完成时间',
    
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    
    INDEX idx_backup_records_status (status),
    INDEX idx_backup_records_type (backup_type),
    INDEX idx_backup_records_created_at (created_at),
    INDEX idx_backup_records_created_by (created_by),
    INDEX idx_backup_records_status_created_at (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='备份记录表';

-- ==================== 5. 文件管理模块 ====================

-- 文件表
CREATE TABLE files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '文件ID',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    original_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件路径',
    file_size BIGINT NOT NULL COMMENT '文件大小',
    file_type VARCHAR(100) COMMENT '文件类型',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    uploaded_by BIGINT NOT NULL COMMENT '上传者ID',
    related_type ENUM('project','competition','user') COMMENT '关联类型',
    related_id BIGINT COMMENT '关联ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_files_uploaded_by (uploaded_by),
    INDEX idx_files_related_type (related_type),
    INDEX idx_files_related_id (related_id),
    INDEX idx_files_created_at (created_at),
    INDEX idx_files_file_type (file_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- ==================== 6. 视图创建 ====================

-- 系统日志统计视图
CREATE VIEW system_logs_summary AS
SELECT 
    DATE(created_at) as log_date,
    action,
    COUNT(*) as action_count,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT ip_address) as unique_ips
FROM system_logs
GROUP BY DATE(created_at), action
ORDER BY log_date DESC, action_count DESC;

-- 备份统计视图
CREATE VIEW backup_statistics AS
SELECT 
    backup_type,
    status,
    COUNT(*) as count,
    AVG(TIMESTAMPDIFF(SECOND, created_at, completed_at)) as avg_duration_seconds,
    SUM(file_size) as total_size_bytes,
    MAX(created_at) as last_backup
FROM backup_records
GROUP BY backup_type, status;

-- 竞赛统计视图
CREATE VIEW competition_statistics AS
SELECT 
    c.id,
    c.title,
    c.level,
    c.status,
    COUNT(DISTINCT cr.student_id) as registered_count,
    COUNT(DISTINCT cs.student_id) as submitted_count,
    COUNT(DISTINCT cj.teacher_id) as judge_count,
    COUNT(DISTINCT cr2.student_id) as awarded_count
FROM competitions c
LEFT JOIN competition_registrations cr ON c.id = cr.competition_id AND cr.status = 'approved'
LEFT JOIN competition_submissions cs ON c.id = cs.competition_id
LEFT JOIN competition_judges cj ON c.id = cj.competition_id AND cj.status = 'active'
LEFT JOIN competition_results cr2 ON c.id = cr2.competition_id
GROUP BY c.id, c.title, c.level, c.status;

-- ==================== 7. 触发器创建 ====================

DELIMITER //

-- 竞赛报名触发器 - 自动更新参与人数
CREATE TRIGGER update_competition_participants_insert
AFTER INSERT ON competition_registrations
FOR EACH ROW
BEGIN
    IF NEW.status = 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants + 1
        WHERE id = NEW.competition_id;
    END IF;
END//

CREATE TRIGGER update_competition_participants_update
AFTER UPDATE ON competition_registrations
FOR EACH ROW
BEGIN
    IF OLD.status != 'approved' AND NEW.status = 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants + 1
        WHERE id = NEW.competition_id;
    ELSEIF OLD.status = 'approved' AND NEW.status != 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants - 1
        WHERE id = NEW.competition_id;
    END IF;
END//

-- 系统日志触发器
CREATE TRIGGER log_user_login
AFTER INSERT ON login_logs
FOR EACH ROW
BEGIN
    INSERT INTO system_logs (user_id, action, details, ip_address, user_agent)
    VALUES (NEW.user_id, 'user_login', CONCAT('用户登录成功，用户ID: ', NEW.user_id), NEW.ip_address, NEW.user_agent);
END//

CREATE TRIGGER log_user_creation
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.id, 'user_created', CONCAT('创建新用户: ', NEW.username, ' (', NEW.email, ')'));
END//

CREATE TRIGGER log_user_status_change
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    IF OLD.status != NEW.status THEN
        INSERT INTO system_logs (user_id, action, details)
        VALUES (NEW.id, 'user_status_changed', CONCAT('用户状态从 ', OLD.status, ' 变更为 ', NEW.status));
    END IF;
END//

CREATE TRIGGER log_competition_creation
AFTER INSERT ON competitions
FOR EACH ROW
BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.created_by, 'competition_created', CONCAT('创建竞赛: ', NEW.title));
END//

CREATE TRIGGER log_project_creation
AFTER INSERT ON projects
FOR EACH ROW
BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.student_id, 'project_created', CONCAT('创建项目: ', NEW.title));
END//

DELIMITER ;

-- ==================== 8. 存储过程创建 ====================

DELIMITER //

-- 清理过期日志
CREATE PROCEDURE CleanupOldLogs(IN days_to_keep INT)
BEGIN
    DECLARE cutoff_date DATETIME;
    SET cutoff_date = DATE_SUB(NOW(), INTERVAL days_to_keep DAY);
    
    DELETE FROM system_logs WHERE created_at < cutoff_date;
    
    SELECT ROW_COUNT() as deleted_count;
END//

-- 清理过期备份记录
CREATE PROCEDURE CleanupOldBackupRecords(IN days_to_keep INT)
BEGIN
    DECLARE cutoff_date DATETIME;
    SET cutoff_date = DATE_SUB(NOW(), INTERVAL days_to_keep DAY);
    
    DELETE FROM backup_records WHERE created_at < cutoff_date AND status = 'success';
    
    SELECT ROW_COUNT() as deleted_count;
END//

-- 获取系统统计信息
CREATE PROCEDURE GetSystemStats()
BEGIN
    SELECT 
        (SELECT COUNT(*) FROM users) as total_users,
        (SELECT COUNT(*) FROM users WHERE status = 'active') as active_users,
        (SELECT COUNT(*) FROM projects) as total_projects,
        (SELECT COUNT(*) FROM competitions) as total_competitions,
        (SELECT COUNT(*) FROM system_logs WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)) as logs_24h,
        (SELECT COUNT(*) FROM backup_records WHERE status = 'success' AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)) as backups_7d;
END//

-- 创建系统备份记录
CREATE PROCEDURE CreateBackupRecord(
    IN p_file_name VARCHAR(255),
    IN p_file_path VARCHAR(500),
    IN p_backup_type ENUM('full','incremental','manual'),
    IN p_created_by BIGINT
)
BEGIN
    INSERT INTO backup_records (file_name, file_path, backup_type, created_by, status)
    VALUES (p_file_name, p_file_path, p_backup_type, p_created_by, 'pending');
    
    SELECT LAST_INSERT_ID() as backup_id;
END//

DELIMITER ;

-- ==================== 9. 事件创建 ====================

DELIMITER //

-- 每天凌晨2点清理30天前的日志
CREATE EVENT IF NOT EXISTS cleanup_old_logs_daily
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_TIMESTAMP + INTERVAL 2 HOUR
DO
BEGIN
    CALL CleanupOldLogs(30);
END//

-- 每周日凌晨3点清理90天前的备份记录
CREATE EVENT IF NOT EXISTS cleanup_old_backup_records_weekly
ON SCHEDULE EVERY 1 WEEK
STARTS CURRENT_TIMESTAMP + INTERVAL 3 HOUR
DO
BEGIN
    CALL CleanupOldBackupRecords(90);
END//

DELIMITER ;

-- ==================== 10. 测试数据插入 ====================

-- 插入角色数据
INSERT INTO roles (role_name, description) VALUES
('admin', '系统管理员'),
('teacher', '教师'),
('student', '学生');

-- 插入用户数据（密码都是123456的bcrypt哈希）
INSERT INTO users (username, email, password_hash, status) VALUES
('admin', 'admin@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('teacher001', 'teacher001@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('teacher002', 'teacher002@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('student001', 'student001@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('student002', 'student002@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('student003', 'student003@yunmeng.edu.cn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active');

-- 插入用户角色关联
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- admin -> admin
(2, 2), -- teacher001 -> teacher
(3, 2), -- teacher002 -> teacher
(4, 3), -- student001 -> student
(5, 3), -- student002 -> student
(6, 3); -- student003 -> student

-- 插入用户档案
INSERT INTO user_profiles (user_id, real_name, phone, department, student_id, teacher_id) VALUES
(1, '系统管理员', '13800000000', '信息技术部', NULL, 'T001'),
(2, '张教授', '13800000001', '计算机学院', NULL, 'T002'),
(3, '李副教授', '13800000002', '数学学院', NULL, 'T003'),
(4, '王同学', '13800000003', '计算机学院', '2021001', NULL),
(5, '刘同学', '13800000004', '数学学院', '2021002', NULL),
(6, '陈同学', '13800000005', '物理学院', '2021003', NULL);

-- 插入项目分类
INSERT INTO project_types (type_name, description, sort_order) VALUES
('科研项目', '学术研究类项目', 1),
('创新项目', '创新创业类项目', 2),
('竞赛项目', '各类竞赛项目', 3),
('实践项目', '社会实践类项目', 4);

-- 插入子分类
INSERT INTO project_types (type_name, description, parent_id, sort_order) VALUES
('自然科学', '自然科学类科研项目', 1, 1),
('社会科学', '社会科学类科研项目', 1, 2),
('工程技术', '工程技术类创新项目', 2, 1),
('商业模式', '商业模式创新项目', 2, 2),
('程序设计', '程序设计竞赛项目', 3, 1),
('数学建模', '数学建模竞赛项目', 3, 2),
('社会实践', '社会实践活动项目', 4, 1),
('志愿服务', '志愿服务活动项目', 4, 2);

-- 插入竞赛数据
INSERT INTO competitions (title, description, level, category, registration_start, registration_end, submission_start, submission_end, max_participants, status, award_config, created_by) VALUES
('2024年大学生程序设计竞赛（校级）', '校级程序设计竞赛，考察学生的编程能力和算法思维', 'school', '程序设计', '2024-01-01 00:00:00', '2024-01-31 23:59:59', '2024-02-01 00:00:00', '2024-02-28 23:59:59', 100, 'registration', '{"first_prize": 3, "second_prize": 6, "third_prize": 10}', 1),
('全国大学生数学建模竞赛（国家级）', '全国大学生数学建模竞赛，培养数学建模能力', 'national', '数学建模', '2024-03-01 00:00:00', '2024-03-31 23:59:59', '2024-04-01 00:00:00', '2024-04-30 23:59:59', 50, 'draft', '{"first_prize": 1, "second_prize": 2, "third_prize": 3}', 1);

-- 插入系统配置
INSERT INTO system_settings (setting_key, setting_value, description, category, is_public) VALUES
('system_name', '云梦高校学生科研与竞赛项目管理系统', '系统名称', 'general', TRUE),
('system_version', '1.0.0', '系统版本', 'general', TRUE),
('max_file_size', '10485760', '最大文件上传大小（字节）', 'upload', TRUE),
('allowed_file_types', 'jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx,zip,rar', '允许上传的文件类型', 'upload', TRUE),
('session_timeout', '3600', '会话超时时间（秒）', 'security', FALSE),
('max_login_attempts', '5', '最大登录尝试次数', 'security', FALSE),
('backup_retention_days', '30', '备份保留天数', 'backup', FALSE),
('auto_backup_enabled', 'true', '是否启用自动备份', 'backup', FALSE),
('backup_schedule', '0 2 * * *', '备份计划（Cron表达式）', 'backup', FALSE),
('email_notifications', 'true', '是否启用邮件通知', 'notification', FALSE),
('smtp_host', '', 'SMTP服务器地址', 'email', FALSE),
('smtp_port', '587', 'SMTP服务器端口', 'email', FALSE),
('smtp_username', '', 'SMTP用户名', 'email', FALSE),
('smtp_password', '', 'SMTP密码', 'email', FALSE),
('maintenance_mode', 'false', '维护模式', 'system', TRUE),
('maintenance_message', '系统正在维护中，请稍后再试', '维护模式消息', 'system', TRUE);

-- 插入示例项目
INSERT INTO projects (title, description, type_id, student_id, teacher_id, status) VALUES
('基于深度学习的图像识别系统', '使用深度学习技术开发图像识别系统，提高识别准确率', 5, 4, 2, 'submitted'),
('校园二手交易平台', '开发校园二手交易平台，促进资源循环利用', 8, 5, 3, 'draft'),
('智能校园导航系统', '基于微信小程序的智能校园导航系统', 5, 6, 2, 'reviewing');

-- 插入竞赛报名记录
INSERT INTO competition_registrations (competition_id, student_id, team_name, team_leader, status, approved_by, approved_at) VALUES
(1, 4, '编程小分队', 4, 'approved', 1, NOW()),
(1, 5, '算法优化组', 5, 'approved', 1, NOW()),
(1, 6, '代码工匠', 6, 'pending', NULL, NULL);

-- 插入竞赛评审教师
INSERT INTO competition_judges (competition_id, teacher_id, status) VALUES
(1, 2, 'active'),
(1, 3, 'active');

-- 插入示例系统日志
INSERT INTO system_logs (user_id, action, details, ip_address, user_agent) VALUES
(1, 'system_startup', '系统启动', '127.0.0.1', 'System Setup'),
(1, 'database_initialized', '数据库初始化完成', '127.0.0.1', 'System Setup'),
(1, 'user_created', '创建测试用户', '127.0.0.1', 'System Setup');

-- 插入示例备份记录
INSERT INTO backup_records (file_name, file_path, file_size, backup_type, status, created_by, completed_at) VALUES
('initial_backup.sql', '/backups/initial_backup.sql', 1024000, 'manual', 'success', 1, NOW());

-- ==================== 11. 验证脚本 ====================

-- 验证表创建
SELECT 'Database setup completed successfully!' as status;

-- 显示表统计
SELECT 
    'Tables created' as info,
    COUNT(*) as count
FROM information_schema.tables 
WHERE table_schema = 'cloud_dream_system';

-- 显示测试数据统计
SELECT 'Test data summary' as info;
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Roles' as table_name, COUNT(*) as count FROM roles
UNION ALL
SELECT 'User Roles' as table_name, COUNT(*) as count FROM user_roles
UNION ALL
SELECT 'Project Types' as table_name, COUNT(*) as count FROM project_types
UNION ALL
SELECT 'Projects' as table_name, COUNT(*) as count FROM projects
UNION ALL
SELECT 'Competitions' as table_name, COUNT(*) as count FROM competitions
UNION ALL
SELECT 'Competition Registrations' as table_name, COUNT(*) as count FROM competition_registrations
UNION ALL
SELECT 'Competition Judges' as table_name, COUNT(*) as count FROM competition_judges
UNION ALL
SELECT 'System Settings' as table_name, COUNT(*) as count FROM system_settings
UNION ALL
SELECT 'System Logs' as table_name, COUNT(*) as count FROM system_logs
UNION ALL
SELECT 'Backup Records' as table_name, COUNT(*) as count FROM backup_records; 