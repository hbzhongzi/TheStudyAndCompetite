-- 竞赛数据库完整更新脚本
-- 执行时间: 2024-01-01
-- 说明: 同步更新竞赛相关表结构，添加报名时间字段和优化现有结构

-- ========================================
-- 1. 更新竞赛表 (competitions)
-- ========================================

-- 检查并添加报名开始时间字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'registration_start') = 0,
    'ALTER TABLE competitions ADD COLUMN registration_start DATETIME NULL COMMENT "报名开始时间" AFTER organizer',
    'SELECT "registration_start 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加报名截止时间字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'registration_end') = 0,
    'ALTER TABLE competitions ADD COLUMN registration_end DATETIME NULL COMMENT "报名截止时间" AFTER registration_start',
    'SELECT "registration_end 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加院系限制字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'department_limit') = 0,
    'ALTER TABLE competitions ADD COLUMN department_limit VARCHAR(255) NULL COMMENT "院系限制（可选）" AFTER max_participants',
    'SELECT "department_limit 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加教师限制字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'teacher_limit') = 0,
    'ALTER TABLE competitions ADD COLUMN teacher_limit BOOLEAN DEFAULT FALSE COMMENT "是否需要绑定教师才能报名" AFTER department_limit',
    'SELECT "teacher_limit 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加获奖配置字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'award_config') = 0,
    'ALTER TABLE competitions ADD COLUMN award_config JSON NULL COMMENT "获奖配置" AFTER teacher_limit',
    'SELECT "award_config 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加竞赛地点字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'location') = 0,
    'ALTER TABLE competitions ADD COLUMN location VARCHAR(255) NULL COMMENT "竞赛地点" AFTER description',
    'SELECT "location 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加联系方式字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'contact') = 0,
    'ALTER TABLE competitions ADD COLUMN contact VARCHAR(255) NULL COMMENT "联系方式" AFTER location',
    'SELECT "contact 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加竞赛规则字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'rules') = 0,
    'ALTER TABLE competitions ADD COLUMN rules TEXT NULL COMMENT "竞赛规则" AFTER contact',
    'SELECT "rules 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加报名要求字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'requirements') = 0,
    'ALTER TABLE competitions ADD COLUMN requirements TEXT NULL COMMENT "报名要求" AFTER rules',
    'SELECT "requirements 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加作品提交截止时间字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'submission_deadline') = 0,
    'ALTER TABLE competitions ADD COLUMN submission_deadline DATETIME NULL COMMENT "作品提交截止时间" AFTER registration_end',
    'SELECT "submission_deadline 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加评审方式字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'judging_method') = 0,
    'ALTER TABLE competitions ADD COLUMN judging_method TEXT NULL COMMENT "评审方式" AFTER requirements',
    'SELECT "judging_method 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加重要注意事项字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'important_notes') = 0,
    'ALTER TABLE competitions ADD COLUMN important_notes TEXT NULL COMMENT "重要注意事项" AFTER judging_method',
    'SELECT "important_notes 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加相关链接字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'website') = 0,
    'ALTER TABLE competitions ADD COLUMN website VARCHAR(500) NULL COMMENT "相关链接" AFTER important_notes',
    'SELECT "website 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加QQ群字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'qq_group') = 0,
    'ALTER TABLE competitions ADD COLUMN qq_group VARCHAR(50) NULL COMMENT "QQ群" AFTER website',
    'SELECT "qq_group 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加文件格式要求字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'file_formats') = 0,
    'ALTER TABLE competitions ADD COLUMN file_formats VARCHAR(255) NULL COMMENT "支持的文件格式" AFTER qq_group',
    'SELECT "file_formats 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加文件大小限制字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competitions' 
     AND COLUMN_NAME = 'file_size_limit') = 0,
    'ALTER TABLE competitions ADD COLUMN file_size_limit VARCHAR(100) NULL COMMENT "文件大小限制" AFTER file_formats',
    'SELECT "file_size_limit 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================
-- 2. 更新竞赛报名表 (competition_registrations)
-- ========================================

-- 检查并添加团队名称字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND COLUMN_NAME = 'team_name') = 0,
    'ALTER TABLE competition_registrations ADD COLUMN team_name VARCHAR(100) NULL COMMENT "团队名称" AFTER teacher_id',
    'SELECT "team_name 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加团队负责人字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND COLUMN_NAME = 'team_leader') = 0,
    'ALTER TABLE competition_registrations ADD COLUMN team_leader BOOLEAN DEFAULT FALSE COMMENT "是否为团队负责人" AFTER team_name',
    'SELECT "team_leader 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加联系电话字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND COLUMN_NAME = 'contact_phone') = 0,
    'ALTER TABLE competition_registrations ADD COLUMN contact_phone VARCHAR(20) NULL COMMENT "联系电话" AFTER team_leader',
    'SELECT "contact_phone 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加联系邮箱字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND COLUMN_NAME = 'contact_email') = 0,
    'ALTER TABLE competition_registrations ADD COLUMN contact_email VARCHAR(100) NULL COMMENT "联系邮箱" AFTER contact_phone',
    'SELECT "contact_email 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加额外信息字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND COLUMN_NAME = 'additional_info') = 0,
    'ALTER TABLE competition_registrations ADD COLUMN additional_info JSON NULL COMMENT "额外信息" AFTER contact_email',
    'SELECT "additional_info 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================
-- 3. 更新竞赛提交表 (competition_submissions)
-- ========================================

-- 检查并添加版本号字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND COLUMN_NAME = 'version') = 0,
    'ALTER TABLE competition_submissions ADD COLUMN version INT DEFAULT 1 COMMENT "版本号" AFTER description',
    'SELECT "version 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加锁定状态字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND COLUMN_NAME = 'locked') = 0,
    'ALTER TABLE competition_submissions ADD COLUMN locked BOOLEAN DEFAULT FALSE COMMENT "是否锁定" AFTER version',
    'SELECT "locked 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加教师查看状态字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND COLUMN_NAME = 'teacher_viewed') = 0,
    'ALTER TABLE competition_submissions ADD COLUMN teacher_viewed BOOLEAN DEFAULT FALSE COMMENT "教师是否查看过作品" AFTER locked',
    'SELECT "teacher_viewed 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加教师反馈字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND COLUMN_NAME = 'teacher_feedback') = 0,
    'ALTER TABLE competition_submissions ADD COLUMN teacher_feedback TEXT NULL COMMENT "教师反馈" AFTER teacher_viewed',
    'SELECT "teacher_feedback 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并添加最后查看时间字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND COLUMN_NAME = 'last_view_time') = 0,
    'ALTER TABLE competition_submissions ADD COLUMN last_view_time DATETIME NULL COMMENT "最后查看时间" AFTER teacher_feedback',
    'SELECT "last_view_time 字段已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================
-- 4. 创建缺失的表
-- ========================================

-- 检查并创建竞赛反馈表
CREATE TABLE IF NOT EXISTS competition_feedback (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    teacher_id BIGINT UNSIGNED NOT NULL COMMENT '教师ID',
    reviewer_id BIGINT UNSIGNED NULL COMMENT '评审教师ID',
    submission_id BIGINT UNSIGNED NULL COMMENT '提交作品ID',
    comment TEXT NOT NULL COMMENT '评语',
    score DECIMAL(5,2) NULL COMMENT '评分',
    is_final BOOLEAN DEFAULT FALSE COMMENT '是否为最终评分',
    feedback_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '评语时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_competition_feedback_competition (competition_id),
    INDEX idx_competition_feedback_student (student_id),
    INDEX idx_competition_feedback_teacher (teacher_id),
    INDEX idx_competition_feedback_submission (submission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛反馈表';

-- 检查并创建竞赛结果表
CREATE TABLE IF NOT EXISTS competition_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    submission_id BIGINT UNSIGNED NULL COMMENT '提交作品ID',
    award_level VARCHAR(50) NOT NULL COMMENT '获奖等级',
    final_score INT NULL COMMENT '最终得分',
    certificate_url VARCHAR(500) NULL COMMENT '证书URL',
    publish_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    finalized_by BIGINT UNSIGNED NULL COMMENT '最终确认者ID',
    finalized_at DATETIME NULL COMMENT '最终确认时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_competition_results_competition (competition_id),
    INDEX idx_competition_results_student (student_id),
    INDEX idx_competition_results_submission (submission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛结果表';

-- 检查并创建竞赛评审教师表
CREATE TABLE IF NOT EXISTS competition_judges (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    teacher_id BIGINT UNSIGNED NOT NULL COMMENT '教师ID',
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_competition_judge (competition_id, teacher_id),
    INDEX idx_competition_judges_competition (competition_id),
    INDEX idx_competition_judges_teacher (teacher_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛评审教师表';

-- 检查并创建竞赛评分表
CREATE TABLE IF NOT EXISTS competition_scores (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    submission_id BIGINT UNSIGNED NOT NULL COMMENT '提交作品ID',
    judge_id BIGINT UNSIGNED NOT NULL COMMENT '评审教师ID',
    score DECIMAL(5,2) NOT NULL COMMENT '评分',
    comment TEXT NULL COMMENT '评语',
    scored_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '评分时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_submission_judge (submission_id, judge_id),
    INDEX idx_competition_scores_submission (submission_id),
    INDEX idx_competition_scores_judge (judge_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞赛评分表';

-- ========================================
-- 5. 添加外键约束
-- ========================================

-- 添加竞赛报名表的外键约束
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_registrations' 
     AND CONSTRAINT_NAME = 'fk_registrations_competition') = 0,
    'ALTER TABLE competition_registrations ADD CONSTRAINT fk_registrations_competition FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE',
    'SELECT "competition_registrations 外键约束已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加竞赛提交表的外键约束
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
     WHERE TABLE_SCHEMA = DATABASE() 
     AND TABLE_NAME = 'competition_submissions' 
     AND CONSTRAINT_NAME = 'fk_submissions_competition') = 0,
    'ALTER TABLE competition_submissions ADD CONSTRAINT fk_submissions_competition FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE',
    'SELECT "competition_submissions 外键约束已存在" AS message'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================
-- 6. 创建性能优化索引
-- ========================================

-- 竞赛表索引
CREATE INDEX IF NOT EXISTS idx_competitions_registration_time ON competitions(registration_start, registration_end);
CREATE INDEX IF NOT EXISTS idx_competitions_competition_time ON competitions(start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_competitions_status ON competitions(status);
CREATE INDEX IF NOT EXISTS idx_competitions_is_open ON competitions(is_open);
CREATE INDEX IF NOT EXISTS idx_competitions_type ON competitions(type);
CREATE INDEX IF NOT EXISTS idx_competitions_created_by ON competitions(created_by);
CREATE INDEX IF NOT EXISTS idx_competitions_created_at ON competitions(created_at);

-- 竞赛报名表索引
CREATE INDEX IF NOT EXISTS idx_registrations_competition ON competition_registrations(competition_id);
CREATE INDEX IF NOT EXISTS idx_registrations_student ON competition_registrations(student_id);
CREATE INDEX IF NOT EXISTS idx_registrations_teacher ON competition_registrations(teacher_id);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON competition_registrations(status);
CREATE INDEX IF NOT EXISTS idx_registrations_register_time ON competition_registrations(register_time);

-- 竞赛提交表索引
CREATE INDEX IF NOT EXISTS idx_submissions_competition ON competition_submissions(competition_id);
CREATE INDEX IF NOT EXISTS idx_submissions_student ON competition_submissions(student_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON competition_submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_submit_time ON competition_submissions(submit_time);

-- ========================================
-- 7. 更新现有数据
-- ========================================

-- 更新现有竞赛数据，设置默认的报名时间
UPDATE competitions 
SET registration_start = created_at,
    registration_end = start_time
WHERE registration_start IS NULL 
  AND start_time IS NOT NULL;

-- 对于没有设置报名时间的竞赛，将报名开始时间设置为创建时间，报名截止时间设置为当前时间+30天
UPDATE competitions 
SET registration_start = created_at,
    registration_end = DATE_ADD(NOW(), INTERVAL 30 DAY)
WHERE registration_start IS NULL 
  AND start_time IS NULL;

-- 设置默认的院系限制（如果为空）
UPDATE competitions 
SET department_limit = '不限'
WHERE department_limit IS NULL OR department_limit = '';

-- 设置默认的文件格式要求
UPDATE competitions 
SET file_formats = 'PDF、PPT、DOC、ZIP'
WHERE file_formats IS NULL OR file_formats = '';

-- 设置默认的文件大小限制
UPDATE competitions 
SET file_size_limit = '单个文件不超过50MB'
WHERE file_size_limit IS NULL OR file_size_limit = '';

-- ========================================
-- 8. 验证更新结果
-- ========================================

-- 显示竞赛表结构
DESCRIBE competitions;

-- 显示竞赛报名表结构
DESCRIBE competition_registrations;

-- 显示竞赛提交表结构
DESCRIBE competition_submissions;

-- 验证数据完整性
SELECT 
    'competitions' AS table_name,
    COUNT(*) AS total_records,
    COUNT(registration_start) AS has_registration_start,
    COUNT(registration_end) AS has_registration_end,
    COUNT(location) AS has_location,
    COUNT(contact) AS has_contact
FROM competitions
UNION ALL
SELECT 
    'competition_registrations' AS table_name,
    COUNT(*) AS total_records,
    COUNT(team_name) AS has_team_name,
    COUNT(contact_phone) AS has_contact_phone,
    COUNT(contact_email) AS has_contact_email,
    COUNT(additional_info) AS has_additional_info
FROM competition_registrations
UNION ALL
SELECT 
    'competition_submissions' AS table_name,
    COUNT(*) AS total_records,
    COUNT(version) AS has_version,
    COUNT(locked) AS has_locked,
    COUNT(teacher_viewed) AS has_teacher_viewed,
    COUNT(teacher_feedback) AS has_teacher_feedback
FROM competition_submissions;

-- 显示索引信息
SHOW INDEX FROM competitions;
SHOW INDEX FROM competition_registrations;
SHOW INDEX FROM competition_submissions;

-- 显示外键约束
SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    CONSTRAINT_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME IN ('competition_registrations', 'competition_submissions')
  AND REFERENCED_TABLE_NAME IS NOT NULL;

-- 完成提示
SELECT '数据库更新完成！' AS message; 