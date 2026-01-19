# 云梦高校科研竞赛管理系统 - 手动迁移指南

## 概述

如果自动迁移脚本遇到问题，您可以按照以下步骤手动执行数据库迁移。

## 前提条件

此文件内容已归档并移至 `go-backend/backups/docs_archive/MANUAL_MIGRATION_GUIDE.md`，可安全删除原件以释放空间。

## 手动迁移步骤

### 步骤1：连接到数据库

```bash
mysql -u root -p cloud_dream_system
```

### 步骤2：检查当前表结构

```sql
-- 查看projects表当前结构
DESCRIBE projects;

-- 查看是否存在student_teacher表
SHOW TABLES LIKE 'student_teacher';
```

### 步骤3：添加teacher_id字段

```sql
-- 添加teacher_id字段
ALTER TABLE projects ADD COLUMN teacher_id BIGINT NOT NULL DEFAULT 2 COMMENT '指导老师ID' AFTER student_id;
```

### 步骤4：添加submitted_at字段

```sql
-- 添加submitted_at字段
ALTER TABLE projects ADD COLUMN submitted_at DATETIME NULL COMMENT '提交时间' AFTER status;
```

### 步骤5：添加外键约束

```sql
-- 添加外键约束
ALTER TABLE projects ADD CONSTRAINT fk_projects_teacher FOREIGN KEY (teacher_id) REFERENCES users(id);
```

### 步骤6：创建student_teacher表

```sql
-- 创建student_teacher表
CREATE TABLE student_teacher (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  student_id BIGINT NOT NULL,
  teacher_id BIGINT NOT NULL,
  bind_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (student_id, teacher_id),
  FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生教师绑定关系表';
```

### 步骤7：创建索引

```sql
-- 为student_teacher表创建索引
CREATE INDEX idx_student_teacher_student_id ON student_teacher(student_id);
CREATE INDEX idx_student_teacher_teacher_id ON student_teacher(teacher_id);

-- 为projects表创建索引
CREATE INDEX idx_projects_teacher_id ON projects(teacher_id);
```

### 步骤8：插入示例数据

```sql
-- 插入示例学生教师绑定关系
INSERT INTO student_teacher (student_id, teacher_id) VALUES (3, 2);
```

### 步骤9：更新现有数据

```sql
-- 更新现有项目的teacher_id
UPDATE projects SET teacher_id = 2 WHERE teacher_id = 0 OR teacher_id IS NULL;

-- 更新现有项目的submitted_at时间
UPDATE projects SET submitted_at = created_at WHERE status IN ('pending', 'approved', 'rejected') AND submitted_at IS NULL;
```

### 步骤10：验证迁移结果

```sql
-- 查看projects表结构
DESCRIBE projects;

-- 查看student_teacher表
SELECT * FROM student_teacher;

-- 查看项目教师分配情况
SELECT p.id, p.title, p.teacher_id, u.real_name as teacher_name 
FROM projects p 
LEFT JOIN users u ON p.teacher_id = u.id;
```

## 错误处理

### 错误1：字段已存在
如果遇到 "Duplicate column name" 错误，说明字段已经存在，可以跳过该步骤。

### 错误2：表已存在
如果遇到 "Table already exists" 错误，说明表已经存在，可以跳过该步骤。

### 错误3：外键约束已存在
如果遇到 "Duplicate key name" 错误，说明外键约束已经存在，可以跳过该步骤。

### 错误4：索引已存在
如果遇到 "Duplicate key name" 错误，说明索引已经存在，可以跳过该步骤。

## 验证清单

迁移完成后，请检查以下项目：

- [ ] projects表包含teacher_id字段
- [ ] projects表包含submitted_at字段
- [ ] projects表有fk_projects_teacher外键约束
- [ ] student_teacher表存在
- [ ] 相关索引已创建
- [ ] 示例数据已插入
- [ ] 现有项目数据已更新

## 回滚方案

如果需要回滚迁移，可以执行以下SQL：

```sql
-- 删除外键约束
ALTER TABLE projects DROP FOREIGN KEY fk_projects_teacher;

-- 删除索引
DROP INDEX idx_projects_teacher_id ON projects;
DROP INDEX idx_student_teacher_student_id ON student_teacher;
DROP INDEX idx_student_teacher_teacher_id ON student_teacher;

-- 删除字段
ALTER TABLE projects DROP COLUMN teacher_id;
ALTER TABLE projects DROP COLUMN submitted_at;

-- 删除表
DROP TABLE student_teacher;
```

## 联系支持

如果在手动迁移过程中遇到问题，请：

1. 记录具体的错误信息
2. 检查MySQL版本和权限
3. 查看系统日志
4. 联系技术支持 