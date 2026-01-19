# 数据库结构对比文档

## 更新前后对比

### 1. 竞赛表 (competitions)

#### 更新前结构
```sql
CREATE TABLE competitions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL COMMENT '竞赛标题',
    type VARCHAR(50) COMMENT '竞赛类型',
    organizer VARCHAR(100) COMMENT '主办方',
    start_time DATETIME COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    description TEXT COMMENT '竞赛描述',
    attachment VARCHAR(255) COMMENT '附件URL',
    is_open BOOLEAN DEFAULT TRUE COMMENT '是否开放报名',
    max_participants INT COMMENT '最大参与人数',
    current_participants INT DEFAULT 0 COMMENT '当前参与人数',
    status ENUM('draft','registration','submission','review','completed') DEFAULT 'draft' COMMENT '竞赛状态',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 更新后结构
```sql
CREATE TABLE competitions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL COMMENT '竞赛标题',
    type VARCHAR(50) COMMENT '竞赛类型',
    organizer VARCHAR(100) COMMENT '主办方',
    -- 新增字段
    registration_start DATETIME NULL COMMENT '报名开始时间',
    registration_end DATETIME NULL COMMENT '报名截止时间',
    submission_deadline DATETIME NULL COMMENT '作品提交截止时间',
    start_time DATETIME COMMENT '比赛开始时间',
    end_time DATETIME COMMENT '比赛结束时间',
    description TEXT COMMENT '竞赛描述',
    -- 新增字段
    location VARCHAR(255) NULL COMMENT '竞赛地点',
    contact VARCHAR(255) NULL COMMENT '联系方式',
    rules TEXT NULL COMMENT '竞赛规则',
    requirements TEXT NULL COMMENT '报名要求',
    judging_method TEXT NULL COMMENT '评审方式',
    important_notes TEXT NULL COMMENT '重要注意事项',
    website VARCHAR(500) NULL COMMENT '相关链接',
    qq_group VARCHAR(50) NULL COMMENT 'QQ群',
    file_formats VARCHAR(255) NULL COMMENT '支持的文件格式',
    file_size_limit VARCHAR(100) NULL COMMENT '文件大小限制',
    attachment VARCHAR(255) COMMENT '附件URL',
    is_open BOOLEAN DEFAULT TRUE COMMENT '是否开放报名',
    max_participants INT COMMENT '最大参与人数',
    current_participants INT DEFAULT 0 COMMENT '当前参与人数',
    -- 新增字段
    department_limit VARCHAR(255) NULL COMMENT '院系限制（可选）',
    teacher_limit BOOLEAN DEFAULT FALSE COMMENT '是否需要绑定教师才能报名',
    status ENUM('draft','registration','submission','review','completed') DEFAULT 'draft' COMMENT '竞赛状态',
    award_config JSON NULL COMMENT '获奖配置',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 新增字段说明
| 字段名 | 类型 | 说明 | 用途 |
|--------|------|------|------|
| `registration_start` | DATETIME | 报名开始时间 | 控制报名开放时间 |
| `registration_end` | DATETIME | 报名截止时间 | 控制报名关闭时间 |
| `submission_deadline` | DATETIME | 作品提交截止时间 | 控制作品提交时间 |
| `location` | VARCHAR(255) | 竞赛地点 | 显示竞赛举办地点 |
| `contact` | VARCHAR(255) | 联系方式 | 提供联系信息 |
| `rules` | TEXT | 竞赛规则 | 详细规则说明 |
| `requirements` | TEXT | 报名要求 | 报名条件说明 |
| `judging_method` | TEXT | 评审方式 | 评审流程说明 |
| `important_notes` | TEXT | 重要注意事项 | 重要提醒信息 |
| `website` | VARCHAR(500) | 相关链接 | 官方链接 |
| `qq_group` | VARCHAR(50) | QQ群 | 交流群信息 |
| `file_formats` | VARCHAR(255) | 支持的文件格式 | 文件要求说明 |
| `file_size_limit` | VARCHAR(100) | 文件大小限制 | 文件大小要求 |
| `department_limit` | VARCHAR(255) | 院系限制 | 限制报名院系 |
| `teacher_limit` | BOOLEAN | 教师限制 | 是否需要教师指导 |
| `award_config` | JSON | 获奖配置 | 奖项设置信息 |

### 2. 竞赛报名表 (competition_registrations)

#### 更新前结构
```sql
CREATE TABLE competition_registrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    teacher_id BIGINT UNSIGNED COMMENT '指导教师ID',
    register_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '报名时间',
    status ENUM('registered','withdrawn','approved','rejected') DEFAULT 'registered' COMMENT '报名状态',
    teacher_review_status ENUM('pending','approved','rejected') DEFAULT 'pending' COMMENT '教师审核状态',
    teacher_review_comment TEXT COMMENT '教师审核意见',
    teacher_review_time DATETIME COMMENT '教师审核时间'
);
```

#### 更新后结构
```sql
CREATE TABLE competition_registrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    teacher_id BIGINT UNSIGNED COMMENT '指导教师ID',
    register_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '报名时间',
    status ENUM('registered','withdrawn','approved','rejected') DEFAULT 'registered' COMMENT '报名状态',
    teacher_review_status ENUM('pending','approved','rejected') DEFAULT 'pending' COMMENT '教师审核状态',
    teacher_review_comment TEXT COMMENT '教师审核意见',
    teacher_review_time DATETIME COMMENT '教师审核时间',
    -- 新增字段
    team_name VARCHAR(100) NULL COMMENT '团队名称',
    team_leader BOOLEAN DEFAULT FALSE COMMENT '是否为团队负责人',
    contact_phone VARCHAR(20) NULL COMMENT '联系电话',
    contact_email VARCHAR(100) NULL COMMENT '联系邮箱',
    additional_info JSON NULL COMMENT '额外信息'
);
```

#### 新增字段说明
| 字段名 | 类型 | 说明 | 用途 |
|--------|------|------|------|
| `team_name` | VARCHAR(100) | 团队名称 | 团队报名时使用 |
| `team_leader` | BOOLEAN | 团队负责人 | 标识团队负责人 |
| `contact_phone` | VARCHAR(20) | 联系电话 | 联系信息 |
| `contact_email` | VARCHAR(100) | 联系邮箱 | 联系信息 |
| `additional_info` | JSON | 额外信息 | 存储其他报名信息 |

### 3. 竞赛提交表 (competition_submissions)

#### 更新前结构
```sql
CREATE TABLE competition_submissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    file_url VARCHAR(255) COMMENT '文件URL',
    file_name VARCHAR(100) COMMENT '文件名',
    file_size BIGINT COMMENT '文件大小',
    description TEXT COMMENT '成果描述',
    submit_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
    status ENUM('submitted','reviewing','approved','rejected') DEFAULT 'submitted' COMMENT '提交状态',
    review_comments TEXT COMMENT '评审意见'
);
```

#### 更新后结构
```sql
CREATE TABLE competition_submissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    student_id BIGINT UNSIGNED NOT NULL COMMENT '学生ID',
    file_url VARCHAR(255) COMMENT '文件URL',
    file_name VARCHAR(100) COMMENT '文件名',
    file_size BIGINT COMMENT '文件大小',
    description TEXT COMMENT '成果描述',
    -- 新增字段
    version INT DEFAULT 1 COMMENT '版本号',
    submit_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
    status ENUM('submitted','reviewing','approved','rejected') DEFAULT 'submitted' COMMENT '提交状态',
    review_comments TEXT COMMENT '评审意见',
    -- 新增字段
    locked BOOLEAN DEFAULT FALSE COMMENT '是否锁定',
    teacher_viewed BOOLEAN DEFAULT FALSE COMMENT '教师是否查看过作品',
    teacher_feedback TEXT NULL COMMENT '教师反馈',
    last_view_time DATETIME NULL COMMENT '最后查看时间'
);
```

#### 新增字段说明
| 字段名 | 类型 | 说明 | 用途 |
|--------|------|------|------|
| `version` | INT | 版本号 | 支持多版本提交 |
| `locked` | BOOLEAN | 是否锁定 | 防止重复修改 |
| `teacher_viewed` | BOOLEAN | 教师查看状态 | 跟踪查看状态 |
| `teacher_feedback` | TEXT | 教师反馈 | 存储教师意见 |
| `last_view_time` | DATETIME | 最后查看时间 | 记录查看时间 |

### 4. 新增表结构

#### 4.1 竞赛反馈表 (competition_feedback)
```sql
CREATE TABLE competition_feedback (
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
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 4.2 竞赛结果表 (competition_results)
```sql
CREATE TABLE competition_results (
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
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 4.3 竞赛评审教师表 (competition_judges)
```sql
CREATE TABLE competition_judges (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    competition_id BIGINT UNSIGNED NOT NULL COMMENT '竞赛ID',
    teacher_id BIGINT UNSIGNED NOT NULL COMMENT '教师ID',
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 4.4 竞赛评分表 (competition_scores)
```sql
CREATE TABLE competition_scores (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    submission_id BIGINT UNSIGNED NOT NULL COMMENT '提交作品ID',
    judge_id BIGINT UNSIGNED NOT NULL COMMENT '评审教师ID',
    score DECIMAL(5,2) NOT NULL COMMENT '评分',
    comment TEXT NULL COMMENT '评语',
    scored_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '评分时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 索引优化

### 新增索引
| 表名 | 索引名 | 字段 | 说明 |
|------|--------|------|------|
| competitions | idx_competitions_registration_time | registration_start, registration_end | 报名时间查询优化 |
| competitions | idx_competitions_competition_time | start_time, end_time | 比赛时间查询优化 |
| competitions | idx_competitions_status | status | 状态查询优化 |
| competitions | idx_competitions_is_open | is_open | 开放状态查询优化 |
| competitions | idx_competitions_type | type | 类型查询优化 |
| competitions | idx_competitions_created_by | created_by | 创建者查询优化 |
| competitions | idx_competitions_created_at | created_at | 创建时间查询优化 |
| competition_registrations | idx_registrations_competition | competition_id | 竞赛关联查询优化 |
| competition_registrations | idx_registrations_student | student_id | 学生关联查询优化 |
| competition_registrations | idx_registrations_teacher | teacher_id | 教师关联查询优化 |
| competition_registrations | idx_registrations_status | status | 状态查询优化 |
| competition_registrations | idx_registrations_register_time | register_time | 报名时间查询优化 |
| competition_submissions | idx_submissions_competition | competition_id | 竞赛关联查询优化 |
| competition_submissions | idx_submissions_student | student_id | 学生关联查询优化 |
| competition_submissions | idx_submissions_status | status | 状态查询优化 |
| competition_submissions | idx_submissions_submit_time | submit_time | 提交时间查询优化 |

## 外键约束

### 新增外键
| 表名 | 约束名 | 引用表 | 引用字段 | 说明 |
|------|--------|--------|----------|------|
| competition_registrations | fk_registrations_competition | competitions | id | 竞赛关联约束 |
| competition_submissions | fk_submissions_competition | competitions | id | 竞赛关联约束 |

## 数据迁移策略

### 1. 字段添加策略
- 使用 `IF NOT EXISTS` 检查，避免重复添加
- 新字段设置为 `NULL`，允许现有数据兼容
- 使用 `AFTER` 指定字段位置，保持表结构清晰

### 2. 默认值设置策略
- 报名时间：基于现有数据自动计算
- 院系限制：设置为"不限"
- 文件要求：设置常用默认值

### 3. 数据完整性保证
- 外键约束确保引用完整性
- 索引优化确保查询性能
- 字段注释确保可维护性

## 执行顺序

1. **字段添加**：按表顺序添加新字段
2. **表创建**：创建缺失的关联表
3. **约束添加**：添加外键约束
4. **索引创建**：创建性能优化索引
5. **数据更新**：设置默认值和迁移数据
6. **结果验证**：验证更新结果和数据结构

## 注意事项

1. **备份数据**：执行前务必备份数据库
2. **测试环境**：先在测试环境验证脚本
3. **执行权限**：确保有足够的数据库权限
4. **监控执行**：关注执行过程中的错误信息
5. **回滚准备**：准备回滚脚本以防意外

## 验证方法

### 1. 结构验证
```sql
-- 检查表结构
DESCRIBE competitions;
DESCRIBE competition_registrations;
DESCRIBE competition_submissions;

-- 检查新表是否存在
SHOW TABLES LIKE 'competition_%';
```

### 2. 数据验证
```sql
-- 检查数据完整性
SELECT COUNT(*) FROM competitions;
SELECT COUNT(*) FROM competition_registrations;
SELECT COUNT(*) FROM competition_submissions;

-- 检查新字段是否有数据
SELECT 
    COUNT(registration_start) as has_reg_start,
    COUNT(registration_end) as has_reg_end
FROM competitions;
```

### 3. 性能验证
```sql
-- 检查索引是否创建成功
SHOW INDEX FROM competitions;
SHOW INDEX FROM competition_registrations;
SHOW INDEX FROM competition_submissions;

-- 检查外键约束
SELECT * FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME IN ('competition_registrations', 'competition_submissions');
```

这个数据库结构更新将显著提升竞赛管理系统的功能性和性能，为后续的功能扩展奠定坚实的基础。 