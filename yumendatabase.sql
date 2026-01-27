/*
 Navicat Premium Dump SQL

 Source Server         : text
 Source Server Type    : MySQL
 Source Server Version : 80043 (8.0.43)
 Source Host           : localhost:3306
 Source Schema         : cloud_dream_system

 Target Server Type    : MySQL
 Target Server Version : 80043 (8.0.43)
 File Encoding         : 65001

 Date: 27/01/2026 16:57:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for backup_records
-- ----------------------------
DROP TABLE IF EXISTS `backup_records`;
CREATE TABLE `backup_records`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '备份记录ID',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint NULL DEFAULT NULL COMMENT '文件大小（字节）',
  `backup_type` enum('full','incremental','manual') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'manual' COMMENT '备份类型',
  `status` enum('pending','in_progress','success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'pending' COMMENT '备份状态',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `created_by` bigint NULL DEFAULT NULL COMMENT '创建者ID',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `started_at` datetime NULL DEFAULT NULL COMMENT '开始时间',
  `completed_at` datetime NULL DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_backup_records_status`(`status` ASC) USING BTREE,
  INDEX `idx_backup_records_type`(`backup_type` ASC) USING BTREE,
  INDEX `idx_backup_records_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_backup_records_created_by`(`created_by` ASC) USING BTREE,
  INDEX `idx_backup_records_status_created_at`(`status` ASC, `created_at` ASC) USING BTREE,
  CONSTRAINT `backup_records_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '备份记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of backup_records
-- ----------------------------
INSERT INTO `backup_records` VALUES (1, 'initial_backup.sql', '/backups/initial_backup.sql', 1024000, 'manual', 'success', NULL, 1, '2026-01-21 10:07:25', NULL, '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for competition_audit_logs
-- ----------------------------
DROP TABLE IF EXISTS `competition_audit_logs`;
CREATE TABLE `competition_audit_logs`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '审计日志ID',
  `competition_id` bigint NOT NULL COMMENT '竞赛ID',
  `user_id` bigint NOT NULL COMMENT '操作用户ID',
  `action` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `details` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '操作详情',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_competition_audit_logs_competition_id`(`competition_id` ASC) USING BTREE,
  INDEX `idx_competition_audit_logs_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_competition_audit_logs_action`(`action` ASC) USING BTREE,
  INDEX `idx_competition_audit_logs_created_at`(`created_at` ASC) USING BTREE,
  CONSTRAINT `competition_audit_logs_ibfk_1` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_audit_logs_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛审计日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_audit_logs
-- ----------------------------

-- ----------------------------
-- Table structure for competition_feedback
-- ----------------------------
DROP TABLE IF EXISTS `competition_feedback`;
CREATE TABLE `competition_feedback`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评语ID',
  `submission_id` bigint NOT NULL COMMENT '提交ID',
  `teacher_id` bigint NOT NULL COMMENT '教师ID',
  `feedback` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评语内容',
  `rating` int NULL DEFAULT NULL COMMENT '评分(1-10)',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_submission_teacher`(`submission_id` ASC, `teacher_id` ASC) USING BTREE,
  INDEX `idx_competition_feedback_submission_id`(`submission_id` ASC) USING BTREE,
  INDEX `idx_competition_feedback_teacher_id`(`teacher_id` ASC) USING BTREE,
  INDEX `idx_competition_feedback_created_at`(`created_at` ASC) USING BTREE,
  CONSTRAINT `competition_feedback_ibfk_1` FOREIGN KEY (`submission_id`) REFERENCES `competition_submissions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_feedback_ibfk_2` FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛教师评语表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_feedback
-- ----------------------------

-- ----------------------------
-- Table structure for competition_judges
-- ----------------------------
DROP TABLE IF EXISTS `competition_judges`;
CREATE TABLE `competition_judges`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评审教师ID',
  `competition_id` bigint NOT NULL COMMENT '竞赛ID',
  `teacher_id` bigint NOT NULL COMMENT '教师ID',
  `assigned_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  `status` enum('active','inactive') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'active' COMMENT '是否参与',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_competition_teacher`(`competition_id` ASC, `teacher_id` ASC) USING BTREE,
  INDEX `idx_competition_judges_competition_id`(`competition_id` ASC) USING BTREE,
  INDEX `idx_competition_judges_teacher_id`(`teacher_id` ASC) USING BTREE,
  INDEX `idx_competition_judges_status`(`status` ASC) USING BTREE,
  CONSTRAINT `competition_judges_ibfk_1` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_judges_ibfk_2` FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛评审教师表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_judges
-- ----------------------------
INSERT INTO `competition_judges` VALUES (1, 1, 2, '2026-01-21 10:07:25', 'active', '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `competition_judges` VALUES (2, 1, 3, '2026-01-21 10:07:25', 'active', '2026-01-21 10:07:25', '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for competition_registrations
-- ----------------------------
DROP TABLE IF EXISTS `competition_registrations`;
CREATE TABLE `competition_registrations`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '报名ID',
  `competition_id` bigint NOT NULL COMMENT '竞赛ID',
  `student_id` bigint NOT NULL COMMENT '学生ID',
  `team_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '团队名称',
  `team_leader` bigint NULL DEFAULT NULL COMMENT '团队负责人ID',
  `registration_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '报名时间',
  `status` enum('pending','approved','rejected') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'pending' COMMENT '报名状态',
  `approved_by` bigint NULL DEFAULT NULL COMMENT '审批人ID',
  `approved_at` datetime NULL DEFAULT NULL COMMENT '审批时间',
  `rejection_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '拒绝原因',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_competition_student`(`competition_id` ASC, `student_id` ASC) USING BTREE,
  INDEX `approved_by`(`approved_by` ASC) USING BTREE,
  INDEX `idx_competition_registrations_competition_id`(`competition_id` ASC) USING BTREE,
  INDEX `idx_competition_registrations_student_id`(`student_id` ASC) USING BTREE,
  INDEX `idx_competition_registrations_team_leader`(`team_leader` ASC) USING BTREE,
  INDEX `idx_competition_registrations_status`(`status` ASC) USING BTREE,
  INDEX `idx_competition_registrations_registration_time`(`registration_time` ASC) USING BTREE,
  CONSTRAINT `competition_registrations_ibfk_1` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_registrations_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_registrations_ibfk_3` FOREIGN KEY (`team_leader`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT,
  CONSTRAINT `competition_registrations_ibfk_4` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛报名记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_registrations
-- ----------------------------
INSERT INTO `competition_registrations` VALUES (1, 1, 4, '编程小分队', 4, '2026-01-21 10:07:25', 'approved', 1, '2026-01-21 10:07:25', NULL);
INSERT INTO `competition_registrations` VALUES (2, 1, 5, '算法优化组', 5, '2026-01-21 10:07:25', 'approved', 1, '2026-01-21 10:07:25', NULL);
INSERT INTO `competition_registrations` VALUES (3, 1, 6, '代码工匠', 6, '2026-01-21 10:07:25', 'pending', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for competition_results
-- ----------------------------
DROP TABLE IF EXISTS `competition_results`;
CREATE TABLE `competition_results`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '获奖ID',
  `competition_id` bigint NOT NULL COMMENT '竞赛ID',
  `student_id` bigint NOT NULL COMMENT '学生ID',
  `award_level` enum('first','second','third','honorable') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '获奖等级',
  `award_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '获奖名称',
  `score` decimal(5, 2) NULL DEFAULT NULL COMMENT '最终得分',
  `ranking` int NULL DEFAULT NULL COMMENT '排名',
  `finalized_by` bigint NULL DEFAULT NULL COMMENT '最终确认成绩的管理员ID',
  `finalized_at` datetime NULL DEFAULT NULL COMMENT '确认时间',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_competition_student_result`(`competition_id` ASC, `student_id` ASC) USING BTREE,
  INDEX `idx_competition_results_competition_id`(`competition_id` ASC) USING BTREE,
  INDEX `idx_competition_results_student_id`(`student_id` ASC) USING BTREE,
  INDEX `idx_competition_results_award_level`(`award_level` ASC) USING BTREE,
  INDEX `idx_competition_results_ranking`(`ranking` ASC) USING BTREE,
  INDEX `idx_competition_results_finalized_by`(`finalized_by` ASC) USING BTREE,
  CONSTRAINT `competition_results_ibfk_1` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_results_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_results_ibfk_3` FOREIGN KEY (`finalized_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛获奖登记表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_results
-- ----------------------------

-- ----------------------------
-- Table structure for competition_scores
-- ----------------------------
DROP TABLE IF EXISTS `competition_scores`;
CREATE TABLE `competition_scores`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评分ID',
  `submission_id` bigint NOT NULL COMMENT '提交记录ID',
  `judge_id` bigint NOT NULL COMMENT '评审教师ID',
  `score` decimal(5, 2) NOT NULL COMMENT '评分',
  `comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '评语',
  `scored_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评分时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_submission_judge`(`submission_id` ASC, `judge_id` ASC) USING BTREE,
  INDEX `idx_competition_scores_submission_id`(`submission_id` ASC) USING BTREE,
  INDEX `idx_competition_scores_judge_id`(`judge_id` ASC) USING BTREE,
  INDEX `idx_competition_scores_scored_at`(`scored_at` ASC) USING BTREE,
  CONSTRAINT `competition_scores_ibfk_1` FOREIGN KEY (`submission_id`) REFERENCES `competition_submissions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_scores_ibfk_2` FOREIGN KEY (`judge_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛评分记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_scores
-- ----------------------------

-- ----------------------------
-- Table structure for competition_submissions
-- ----------------------------
DROP TABLE IF EXISTS `competition_submissions`;
CREATE TABLE `competition_submissions`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '提交ID',
  `competition_id` bigint NOT NULL COMMENT '竞赛ID',
  `student_id` bigint NOT NULL COMMENT '学生ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '作品标题',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '作品描述',
  `file_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '文件URL',
  `file_size` bigint NULL DEFAULT NULL COMMENT '文件大小',
  `version` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '1.0' COMMENT '版本号',
  `locked` tinyint(1) NULL DEFAULT 0 COMMENT '是否锁定',
  `submitted_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_competition_submissions_competition_id`(`competition_id` ASC) USING BTREE,
  INDEX `idx_competition_submissions_student_id`(`student_id` ASC) USING BTREE,
  INDEX `idx_competition_submissions_submitted_at`(`submitted_at` ASC) USING BTREE,
  INDEX `idx_competition_submissions_locked`(`locked` ASC) USING BTREE,
  CONSTRAINT `competition_submissions_ibfk_1` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `competition_submissions_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛成果提交表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competition_submissions
-- ----------------------------

-- ----------------------------
-- Table structure for competitions
-- ----------------------------
DROP TABLE IF EXISTS `competitions`;
CREATE TABLE `competitions`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '竞赛ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '竞赛标题',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '竞赛描述',
  `level` enum('school','provincial','national','international') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'school' COMMENT '竞赛级别',
  `category` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '竞赛类别',
  `registration_start` datetime NULL DEFAULT NULL COMMENT '报名开始时间',
  `registration_end` datetime NULL DEFAULT NULL COMMENT '报名结束时间',
  `submission_start` datetime NULL DEFAULT NULL COMMENT '提交开始时间',
  `submission_end` datetime NULL DEFAULT NULL COMMENT '提交结束时间',
  `max_participants` int NULL DEFAULT NULL COMMENT '最大参与人数',
  `current_participants` int NULL DEFAULT 0 COMMENT '当前参与人数',
  `status` enum('draft','registration','submission','review','completed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'draft' COMMENT '竞赛状态',
  `award_config` json NULL COMMENT '获奖配置',
  `created_by` bigint NOT NULL COMMENT '创建者ID',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_competitions_level`(`level` ASC) USING BTREE,
  INDEX `idx_competitions_category`(`category` ASC) USING BTREE,
  INDEX `idx_competitions_status`(`status` ASC) USING BTREE,
  INDEX `idx_competitions_registration_start`(`registration_start` ASC) USING BTREE,
  INDEX `idx_competitions_registration_end`(`registration_end` ASC) USING BTREE,
  INDEX `idx_competitions_submission_start`(`submission_start` ASC) USING BTREE,
  INDEX `idx_competitions_submission_end`(`submission_end` ASC) USING BTREE,
  INDEX `idx_competitions_created_by`(`created_by` ASC) USING BTREE,
  CONSTRAINT `competitions_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '竞赛信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of competitions
-- ----------------------------
INSERT INTO `competitions` VALUES (1, '2024年大学生程序设计竞赛（校级）', '校级程序设计竞赛，考察学生的编程能力和算法思维', 'school', '程序设计', '2024-01-01 00:00:00', '2024-01-31 23:59:59', '2024-02-01 00:00:00', '2024-02-28 23:59:59', 100, 2, 'registration', '{\"first_prize\": 3, \"third_prize\": 10, \"second_prize\": 6}', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `competitions` VALUES (2, '全国大学生数学建模竞赛（国家级）', '全国大学生数学建模竞赛，培养数学建模能力', 'national', '数学建模', '2024-03-01 00:00:00', '2024-03-31 23:59:59', '2024-04-01 00:00:00', '2024-04-30 23:59:59', 50, 0, 'draft', '{\"first_prize\": 1, \"third_prize\": 3, \"second_prize\": 2}', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for files
-- ----------------------------
DROP TABLE IF EXISTS `files`;
CREATE TABLE `files`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `original_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始文件名',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint NOT NULL COMMENT '文件大小',
  `file_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '文件类型',
  `mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'MIME类型',
  `uploaded_by` bigint NOT NULL COMMENT '上传者ID',
  `related_type` enum('project','competition','user') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '关联类型',
  `related_id` bigint NULL DEFAULT NULL COMMENT '关联ID',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_files_uploaded_by`(`uploaded_by` ASC) USING BTREE,
  INDEX `idx_files_related_type`(`related_type` ASC) USING BTREE,
  INDEX `idx_files_related_id`(`related_id` ASC) USING BTREE,
  INDEX `idx_files_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_files_file_type`(`file_type` ASC) USING BTREE,
  CONSTRAINT `files_ibfk_1` FOREIGN KEY (`uploaded_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of files
-- ----------------------------

-- ----------------------------
-- Table structure for login_logs
-- ----------------------------
DROP TABLE IF EXISTS `login_logs`;
CREATE TABLE `login_logs`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `login_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '用户代理',
  `status` enum('success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'success' COMMENT '登录状态',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_login_logs_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_login_logs_login_time`(`login_time` ASC) USING BTREE,
  INDEX `idx_login_logs_ip_address`(`ip_address` ASC) USING BTREE,
  INDEX `idx_login_logs_status`(`status` ASC) USING BTREE,
  CONSTRAINT `login_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '登录日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of login_logs
-- ----------------------------

-- ----------------------------
-- Table structure for project_types
-- ----------------------------
DROP TABLE IF EXISTS `project_types`;
CREATE TABLE `project_types`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `type_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '分类描述',
  `parent_id` bigint NULL DEFAULT NULL COMMENT '父分类ID',
  `sort_order` int NULL DEFAULT 0 COMMENT '排序',
  `is_active` tinyint(1) NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_project_types_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_project_types_sort_order`(`sort_order` ASC) USING BTREE,
  INDEX `idx_project_types_is_active`(`is_active` ASC) USING BTREE,
  CONSTRAINT `project_types_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `project_types` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '项目分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_types
-- ----------------------------
INSERT INTO `project_types` VALUES (1, '科研项目', '学术研究类项目', NULL, 1, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (2, '创新项目', '创新创业类项目', NULL, 2, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (3, '竞赛项目', '各类竞赛项目', NULL, 3, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (4, '实践项目', '社会实践类项目', NULL, 4, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (5, '自然科学', '自然科学类科研项目', 1, 1, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (6, '社会科学', '社会科学类科研项目', 1, 2, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (7, '工程技术', '工程技术类创新项目', 2, 1, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (8, '商业模式', '商业模式创新项目', 2, 2, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (9, '程序设计', '程序设计竞赛项目', 3, 1, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (10, '数学建模', '数学建模竞赛项目', 3, 2, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (11, '社会实践', '社会实践活动项目', 4, 1, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `project_types` VALUES (12, '志愿服务', '志愿服务活动项目', 4, 2, 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for projects
-- ----------------------------
DROP TABLE IF EXISTS `projects`;
CREATE TABLE `projects`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '项目ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '项目标题',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '项目描述',
  `type_id` bigint NULL DEFAULT NULL COMMENT '项目分类ID',
  `student_id` bigint NOT NULL COMMENT '学生ID',
  `teacher_id` bigint NULL DEFAULT NULL COMMENT '指导教师ID',
  `status` enum('draft','submitted','reviewing','approved','rejected','completed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'draft' COMMENT '项目状态',
  `submitted_at` datetime NULL DEFAULT NULL COMMENT '提交时间',
  `approved_at` datetime NULL DEFAULT NULL COMMENT '审批时间',
  `approved_by` bigint NULL DEFAULT NULL COMMENT '审批人ID',
  `rejection_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '拒绝原因',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint(1) NULL DEFAULT 0 COMMENT '是否删除',
  `is_approved` tinyint(1) NULL DEFAULT 0 COMMENT '是否通过审批\r\n',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '项目类型',
  `plan` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '预期成功',
  `progress` int NULL DEFAULT NULL COMMENT '进度率',
  `finish_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '预计完成时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `approved_by`(`approved_by` ASC) USING BTREE,
  INDEX `idx_projects_type_id`(`type_id` ASC) USING BTREE,
  INDEX `idx_projects_student_id`(`student_id` ASC) USING BTREE,
  INDEX `idx_projects_teacher_id`(`teacher_id` ASC) USING BTREE,
  INDEX `idx_projects_status`(`status` ASC) USING BTREE,
  INDEX `idx_projects_submitted_at`(`submitted_at` ASC) USING BTREE,
  INDEX `idx_projects_approved_at`(`approved_at` ASC) USING BTREE,
  CONSTRAINT `projects_ibfk_1` FOREIGN KEY (`type_id`) REFERENCES `project_types` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `projects_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `projects_ibfk_3` FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT,
  CONSTRAINT `projects_ibfk_4` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '项目表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of projects
-- ----------------------------
INSERT INTO `projects` VALUES (1, '基于深度学习的图像识别系统', '使用深度学习技术开发图像识别系统，提高识别准确率', 5, 4, 2, 'rejected', NULL, NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-26 15:41:53', 0, 0, '创新项目', '预计6个月完成，分为需求分析、设计、开发、测试四个阶段', 99, '2026-01-26 15:41:53');
INSERT INTO `projects` VALUES (2, '校园二手交易平台', '开发校园二手交易平台，促进资源循环利用', 8, 5, 3, 'reviewing', NULL, NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-26 14:35:56', 0, 0, '创新项目', '预计8个月完成，包括数据采集、预处理、分析、可视化等模块', 50, '2026-02-13 08:52:04');
INSERT INTO `projects` VALUES (3, '智能校园导航系统', '基于微信小程序的智能校园导航系统', 5, 6, 2, 'reviewing', NULL, NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-26 14:35:56', 0, 1, '创新项目', '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块', 80, '2026-01-29 08:52:10');

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色描述',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `role_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `role_name`(`role_name` ASC) USING BTREE,
  INDEX `idx_roles_name`(`role_name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '系统管理员', '2026-01-21 10:07:25', '2026-01-21 10:51:39', 'admin');
INSERT INTO `roles` VALUES (2, 'teacher', '教师', '2026-01-21 10:07:25', '2026-01-21 10:51:39', 'teacher');
INSERT INTO `roles` VALUES (3, 'student', '学生', '2026-01-21 10:07:25', '2026-01-21 10:51:39', 'student');

-- ----------------------------
-- Table structure for system_logs
-- ----------------------------
DROP TABLE IF EXISTS `system_logs`;
CREATE TABLE `system_logs`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint NULL DEFAULT NULL COMMENT '操作用户ID',
  `action` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '执行动作',
  `details` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '详细内容',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '用户代理（浏览器信息）',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_system_logs_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_system_logs_action`(`action` ASC) USING BTREE,
  INDEX `idx_system_logs_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_system_logs_ip_address`(`ip_address` ASC) USING BTREE,
  INDEX `idx_system_logs_action_created_at`(`action` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_system_logs_user_action`(`user_id` ASC, `action` ASC) USING BTREE,
  CONSTRAINT `system_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of system_logs
-- ----------------------------
INSERT INTO `system_logs` VALUES (1, 1, 'user_created', '创建新用户: admin (admin@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (2, 2, 'user_created', '创建新用户: teacher001 (teacher001@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (3, 3, 'user_created', '创建新用户: teacher002 (teacher002@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (4, 4, 'user_created', '创建新用户: student001 (student001@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (5, 5, 'user_created', '创建新用户: student002 (student002@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (6, 6, 'user_created', '创建新用户: student003 (student003@yunmeng.edu.cn)', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (7, 1, 'competition_created', '创建竞赛: 2024年大学生程序设计竞赛（校级）', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (8, 1, 'competition_created', '创建竞赛: 全国大学生数学建模竞赛（国家级）', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (9, 4, 'project_created', '创建项目: 基于深度学习的图像识别系统', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (10, 5, 'project_created', '创建项目: 校园二手交易平台', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (11, 6, 'project_created', '创建项目: 智能校园导航系统', NULL, NULL, '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (12, 1, 'system_startup', '系统启动', '127.0.0.1', 'System Setup', '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (13, 1, 'database_initialized', '数据库初始化完成', '127.0.0.1', 'System Setup', '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (14, 1, 'user_created', '创建测试用户', '127.0.0.1', 'System Setup', '2026-01-21 10:07:25');
INSERT INTO `system_logs` VALUES (15, 4, 'project_created', '创建项目: 12312', NULL, NULL, '2026-01-26 15:33:30');
INSERT INTO `system_logs` VALUES (16, 4, 'project_created', '创建项目: 12312', NULL, NULL, '2026-01-26 15:34:45');
INSERT INTO `system_logs` VALUES (17, 4, 'project_created', '创建项目: 12312', NULL, NULL, '2026-01-26 15:35:18');
INSERT INTO `system_logs` VALUES (18, 4, 'project_created', '创建项目: 12312', NULL, NULL, '2026-01-26 15:36:44');
INSERT INTO `system_logs` VALUES (19, 4, 'project_created', '创建项目: 12312', NULL, NULL, '2026-01-26 15:37:51');

-- ----------------------------
-- Table structure for system_settings
-- ----------------------------
DROP TABLE IF EXISTS `system_settings`;
CREATE TABLE `system_settings`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `setting_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置键',
  `setting_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '配置值',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置描述',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'general' COMMENT '配置分类',
  `is_public` tinyint(1) NULL DEFAULT 0 COMMENT '是否公开（前端可见）',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `setting_key`(`setting_key` ASC) USING BTREE,
  INDEX `idx_system_settings_key`(`setting_key` ASC) USING BTREE,
  INDEX `idx_system_settings_category`(`category` ASC) USING BTREE,
  INDEX `idx_system_settings_public`(`is_public` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of system_settings
-- ----------------------------
INSERT INTO `system_settings` VALUES (1, 'system_name', '云梦高校学生科研与竞赛项目管理系统', '系统名称', 'general', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (2, 'system_version', '1.0.0', '系统版本', 'general', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (3, 'max_file_size', '10485760', '最大文件上传大小（字节）', 'upload', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (4, 'allowed_file_types', 'jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx,zip,rar', '允许上传的文件类型', 'upload', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (5, 'session_timeout', '3600', '会话超时时间（秒）', 'security', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (6, 'max_login_attempts', '5', '最大登录尝试次数', 'security', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (7, 'backup_retention_days', '30', '备份保留天数', 'backup', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (8, 'auto_backup_enabled', 'true', '是否启用自动备份', 'backup', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (9, 'backup_schedule', '0 2 * * *', '备份计划（Cron表达式）', 'backup', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (10, 'email_notifications', 'true', '是否启用邮件通知', 'notification', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (11, 'smtp_host', '', 'SMTP服务器地址', 'email', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (12, 'smtp_port', '587', 'SMTP服务器端口', 'email', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (13, 'smtp_username', '', 'SMTP用户名', 'email', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (14, 'smtp_password', '', 'SMTP密码', 'email', 0, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (15, 'maintenance_mode', 'false', '维护模式', 'system', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `system_settings` VALUES (16, 'maintenance_message', '系统正在维护中，请稍后再试', '维护模式消息', 'system', 1, '2026-01-21 10:07:25', '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for user_profiles
-- ----------------------------
DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE `user_profiles`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '档案ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `real_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '真实姓名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '手机号',
  `department` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门/学院',
  `student_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '学号',
  `teacher_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '教师工号',
  `avatar_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像URL',
  `bio` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '个人简介',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_user_profiles_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_user_profiles_real_name`(`real_name` ASC) USING BTREE,
  INDEX `idx_user_profiles_department`(`department` ASC) USING BTREE,
  INDEX `idx_user_profiles_student_id`(`student_id` ASC) USING BTREE,
  INDEX `idx_user_profiles_teacher_id`(`teacher_id` ASC) USING BTREE,
  CONSTRAINT `user_profiles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户档案表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_profiles
-- ----------------------------
INSERT INTO `user_profiles` VALUES (1, 1, '系统管理员', '13800000000', '信息技术部', NULL, 'T001', NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `user_profiles` VALUES (2, 2, '张教授', '13800000001', '计算机学院', NULL, 'T002', NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `user_profiles` VALUES (3, 3, '李副教授', '13800000002', '数学学院', NULL, 'T003', NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `user_profiles` VALUES (4, 4, '王同学', '13800000003', '计算机学院', '2021001', NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `user_profiles` VALUES (5, 5, '刘同学', '13800000004', '数学学院', '2021002', NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');
INSERT INTO `user_profiles` VALUES (6, 6, '陈同学', '13800000005', '物理学院', '2021003', NULL, NULL, NULL, '2026-01-21 10:07:25', '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_user_role`(`user_id` ASC, `role_id` ASC) USING BTREE,
  INDEX `idx_user_roles_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_user_roles_role_id`(`role_id` ASC) USING BTREE,
  CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (1, 1, 1, '2026-01-21 10:07:25');
INSERT INTO `user_roles` VALUES (2, 2, 2, '2026-01-21 10:07:25');
INSERT INTO `user_roles` VALUES (3, 3, 2, '2026-01-21 10:07:25');
INSERT INTO `user_roles` VALUES (4, 4, 3, '2026-01-21 10:07:25');
INSERT INTO `user_roles` VALUES (5, 5, 3, '2026-01-21 10:07:25');
INSERT INTO `user_roles` VALUES (6, 6, 3, '2026-01-21 10:07:25');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `status` enum('active','inactive','suspended') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'active' COMMENT '用户状态',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `department` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '部门',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '职位',
  `grade` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '年级/等级',
  `major` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '专业',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间(GORM)',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(GORM)',
  `role_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `email`(`email` ASC) USING BTREE,
  INDEX `idx_users_username`(`username` ASC) USING BTREE,
  INDEX `idx_users_email`(`email` ASC) USING BTREE,
  INDEX `idx_users_status`(`status` ASC) USING BTREE,
  INDEX `idx_users_create_time`(`create_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'admin', 'admin@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 09:23:03', NULL, NULL, NULL, NULL, '2026-01-24 09:18:06', '2026-01-24 09:23:03', 'admin');
INSERT INTO `users` VALUES (2, 'teacher001', 'teacher001@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 09:22:58', NULL, NULL, NULL, NULL, '2026-01-24 09:18:06', '2026-01-24 09:22:58', 'teacher');
INSERT INTO `users` VALUES (3, 'teacher002', 'teacher002@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 09:22:57', NULL, NULL, NULL, NULL, '2026-01-24 09:18:06', '2026-01-24 09:22:57', 'teacher');
INSERT INTO `users` VALUES (4, 'student001', 'student001@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 10:36:15', NULL, NULL, '大一', '计算机科学', '2026-01-24 09:18:06', '2026-01-24 10:36:15', 'student');
INSERT INTO `users` VALUES (5, 'student002', 'student002@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 10:36:19', NULL, NULL, '大二', '计算机科学', '2026-01-24 09:18:06', '2026-01-24 10:36:19', 'student');
INSERT INTO `users` VALUES (6, 'student003', 'student003@yunmeng.edu.cn', '$2a$10$8U1x52zTNqbMgvbP.295Lu4iGpzF1OV7s8VFCxl2xDVXUUi5FtmUa', 'active', '2026-01-21 10:07:25', '2026-01-24 10:36:21', NULL, NULL, '大三', '计算机科学', '2026-01-24 09:18:06', '2026-01-24 10:36:21', 'student');

-- ----------------------------
-- View structure for backup_statistics
-- ----------------------------
DROP VIEW IF EXISTS `backup_statistics`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `backup_statistics` AS select `backup_records`.`backup_type` AS `backup_type`,`backup_records`.`status` AS `status`,count(0) AS `count`,avg(timestampdiff(SECOND,`backup_records`.`created_at`,`backup_records`.`completed_at`)) AS `avg_duration_seconds`,sum(`backup_records`.`file_size`) AS `total_size_bytes`,max(`backup_records`.`created_at`) AS `last_backup` from `backup_records` group by `backup_records`.`backup_type`,`backup_records`.`status`;

-- ----------------------------
-- View structure for competition_statistics
-- ----------------------------
DROP VIEW IF EXISTS `competition_statistics`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `competition_statistics` AS select `c`.`id` AS `id`,`c`.`title` AS `title`,`c`.`level` AS `level`,`c`.`status` AS `status`,count(distinct `cr`.`student_id`) AS `registered_count`,count(distinct `cs`.`student_id`) AS `submitted_count`,count(distinct `cj`.`teacher_id`) AS `judge_count`,count(distinct `cr2`.`student_id`) AS `awarded_count` from ((((`competitions` `c` left join `competition_registrations` `cr` on(((`c`.`id` = `cr`.`competition_id`) and (`cr`.`status` = 'approved')))) left join `competition_submissions` `cs` on((`c`.`id` = `cs`.`competition_id`))) left join `competition_judges` `cj` on(((`c`.`id` = `cj`.`competition_id`) and (`cj`.`status` = 'active')))) left join `competition_results` `cr2` on((`c`.`id` = `cr2`.`competition_id`))) group by `c`.`id`,`c`.`title`,`c`.`level`,`c`.`status`;

-- ----------------------------
-- View structure for system_logs_summary
-- ----------------------------
DROP VIEW IF EXISTS `system_logs_summary`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `system_logs_summary` AS select cast(`system_logs`.`created_at` as date) AS `log_date`,`system_logs`.`action` AS `action`,count(0) AS `action_count`,count(distinct `system_logs`.`user_id`) AS `unique_users`,count(distinct `system_logs`.`ip_address`) AS `unique_ips` from `system_logs` group by cast(`system_logs`.`created_at` as date),`system_logs`.`action` order by `log_date` desc,`action_count` desc;

-- ----------------------------
-- Procedure structure for CleanupOldBackupRecords
-- ----------------------------
DROP PROCEDURE IF EXISTS `CleanupOldBackupRecords`;
delimiter ;;
CREATE PROCEDURE `CleanupOldBackupRecords`(IN days_to_keep INT)
BEGIN
    DECLARE cutoff_date DATETIME;
    SET cutoff_date = DATE_SUB(NOW(), INTERVAL days_to_keep DAY);
    
    DELETE FROM backup_records WHERE created_at < cutoff_date AND status = 'success';
    
    SELECT ROW_COUNT() as deleted_count;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for CleanupOldLogs
-- ----------------------------
DROP PROCEDURE IF EXISTS `CleanupOldLogs`;
delimiter ;;
CREATE PROCEDURE `CleanupOldLogs`(IN days_to_keep INT)
BEGIN
    DECLARE cutoff_date DATETIME;
    SET cutoff_date = DATE_SUB(NOW(), INTERVAL days_to_keep DAY);
    
    DELETE FROM system_logs WHERE created_at < cutoff_date;
    
    SELECT ROW_COUNT() as deleted_count;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for CreateBackupRecord
-- ----------------------------
DROP PROCEDURE IF EXISTS `CreateBackupRecord`;
delimiter ;;
CREATE PROCEDURE `CreateBackupRecord`(IN p_file_name VARCHAR(255),
    IN p_file_path VARCHAR(500),
    IN p_backup_type ENUM('full','incremental','manual'),
    IN p_created_by BIGINT)
BEGIN
    INSERT INTO backup_records (file_name, file_path, backup_type, created_by, status)
    VALUES (p_file_name, p_file_path, p_backup_type, p_created_by, 'pending');
    
    SELECT LAST_INSERT_ID() as backup_id;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetSystemStats
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetSystemStats`;
delimiter ;;
CREATE PROCEDURE `GetSystemStats`()
BEGIN
    SELECT 
        (SELECT COUNT(*) FROM users) as total_users,
        (SELECT COUNT(*) FROM users WHERE status = 'active') as active_users,
        (SELECT COUNT(*) FROM projects) as total_projects,
        (SELECT COUNT(*) FROM competitions) as total_competitions,
        (SELECT COUNT(*) FROM system_logs WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)) as logs_24h,
        (SELECT COUNT(*) FROM backup_records WHERE status = 'success' AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)) as backups_7d;
END
;;
delimiter ;

-- ----------------------------
-- Event structure for cleanup_old_backup_records_weekly
-- ----------------------------
DROP EVENT IF EXISTS `cleanup_old_backup_records_weekly`;
delimiter ;;
CREATE EVENT `cleanup_old_backup_records_weekly`
ON SCHEDULE
EVERY '1' WEEK STARTS '2026-01-21 13:07:25'
DO BEGIN
    CALL CleanupOldBackupRecords(90);
END
;;
delimiter ;

-- ----------------------------
-- Event structure for cleanup_old_logs_daily
-- ----------------------------
DROP EVENT IF EXISTS `cleanup_old_logs_daily`;
delimiter ;;
CREATE EVENT `cleanup_old_logs_daily`
ON SCHEDULE
EVERY '1' DAY STARTS '2026-01-21 12:07:25'
DO BEGIN
    CALL CleanupOldLogs(30);
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table competition_registrations
-- ----------------------------
DROP TRIGGER IF EXISTS `update_competition_participants_insert`;
delimiter ;;
CREATE TRIGGER `update_competition_participants_insert` AFTER INSERT ON `competition_registrations` FOR EACH ROW BEGIN
    IF NEW.status = 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants + 1
        WHERE id = NEW.competition_id;
    END IF;
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table competition_registrations
-- ----------------------------
DROP TRIGGER IF EXISTS `update_competition_participants_update`;
delimiter ;;
CREATE TRIGGER `update_competition_participants_update` AFTER UPDATE ON `competition_registrations` FOR EACH ROW BEGIN
    IF OLD.status != 'approved' AND NEW.status = 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants + 1
        WHERE id = NEW.competition_id;
    ELSEIF OLD.status = 'approved' AND NEW.status != 'approved' THEN
        UPDATE competitions 
        SET current_participants = current_participants - 1
        WHERE id = NEW.competition_id;
    END IF;
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table competitions
-- ----------------------------
DROP TRIGGER IF EXISTS `log_competition_creation`;
delimiter ;;
CREATE TRIGGER `log_competition_creation` AFTER INSERT ON `competitions` FOR EACH ROW BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.created_by, 'competition_created', CONCAT('创建竞赛: ', NEW.title));
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table login_logs
-- ----------------------------
DROP TRIGGER IF EXISTS `log_user_login`;
delimiter ;;
CREATE TRIGGER `log_user_login` AFTER INSERT ON `login_logs` FOR EACH ROW BEGIN
    INSERT INTO system_logs (user_id, action, details, ip_address, user_agent)
    VALUES (NEW.user_id, 'user_login', CONCAT('用户登录成功，用户ID: ', NEW.user_id), NEW.ip_address, NEW.user_agent);
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table projects
-- ----------------------------
DROP TRIGGER IF EXISTS `log_project_creation`;
delimiter ;;
CREATE TRIGGER `log_project_creation` AFTER INSERT ON `projects` FOR EACH ROW BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.student_id, 'project_created', CONCAT('创建项目: ', NEW.title));
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table users
-- ----------------------------
DROP TRIGGER IF EXISTS `log_user_creation`;
delimiter ;;
CREATE TRIGGER `log_user_creation` AFTER INSERT ON `users` FOR EACH ROW BEGIN
    INSERT INTO system_logs (user_id, action, details)
    VALUES (NEW.id, 'user_created', CONCAT('创建新用户: ', NEW.username, ' (', NEW.email, ')'));
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table users
-- ----------------------------
DROP TRIGGER IF EXISTS `log_user_status_change`;
delimiter ;;
CREATE TRIGGER `log_user_status_change` AFTER UPDATE ON `users` FOR EACH ROW BEGIN
    IF OLD.status != NEW.status THEN
        INSERT INTO system_logs (user_id, action, details)
        VALUES (NEW.id, 'user_status_changed', CONCAT('用户状态从 ', OLD.status, ' 变更为 ', NEW.status));
    END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
