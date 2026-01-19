# 云梦高校科研竞赛管理系统 - 最终迁移指南

## 迁移状态

根据之前的错误信息，我们发现：

此文件内容已归档并移至 `go-backend/backups/docs_archive/FINAL_MIGRATION_GUIDE.md`，可安全删除原件以释放空间。
### 方案1：使用简单分步迁移（推荐）

我们已经创建了 `sql/simple_step_migration.sql` 脚本，它会：
- 逐步执行每个迁移步骤
- 如果某个步骤失败（如字段已存在），会显示错误但继续执行
- 最终完成所有必要的迁移

**执行命令：**
```bash
cd go-backend
migrate_projects.bat
```

### 方案2：手动执行缺失的步骤

如果自动迁移仍有问题，可以手动执行以下SQL：

```sql
-- 连接到数据库
mysql -u root -p cloud_dream_system

-- 检查当前状态
DESCRIBE projects;
SHOW TABLES LIKE 'student_teacher';

-- 如果submitted_at字段不存在，添加它
ALTER TABLE projects ADD COLUMN submitted_at DATETIME NULL COMMENT '提交时间' AFTER status;

-- 如果外键约束不存在，添加它
ALTER TABLE projects ADD CONSTRAINT fk_projects_teacher FOREIGN KEY (teacher_id) REFERENCES users(id);

-- 如果student_teacher表不存在，创建它
CREATE TABLE student_teacher (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  student_id BIGINT NOT NULL,
  teacher_id BIGINT NOT NULL,
  bind_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (student_id, teacher_id),
  FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生教师绑定关系表';

-- 创建索引
CREATE INDEX idx_student_teacher_student_id ON student_teacher(student_id);
CREATE INDEX idx_student_teacher_teacher_id ON student_teacher(teacher_id);
CREATE INDEX idx_projects_teacher_id ON projects(teacher_id);

-- 插入示例数据
INSERT IGNORE INTO student_teacher (student_id, teacher_id) VALUES (3, 2);

-- 更新数据
UPDATE projects SET teacher_id = 2 WHERE teacher_id = 0 OR teacher_id IS NULL;
UPDATE projects SET submitted_at = created_at WHERE status IN ('pending', 'approved', 'rejected') AND submitted_at IS NULL;
```

## 验证迁移结果

### 方法1：使用检查脚本
```bash
simple_check.bat
```

### 方法2：手动验证
```sql
-- 检查projects表结构
DESCRIBE projects;

-- 检查student_teacher表
SELECT * FROM student_teacher;

-- 检查项目教师分配
SELECT p.id, p.title, p.teacher_id, u.real_name as teacher_name 
FROM projects p 
LEFT JOIN users u ON p.teacher_id = u.id;
```

## 预期结果

迁移完成后，您应该看到：

### projects表结构
```
+-------------+--------------+------+-----+-------------------+----------------+
| Field       | Type         | Null | Key | Default           | Extra          |
+-------------+--------------+------+-----+-------------------+----------------+
| id          | bigint       | NO   | PRI | NULL              | auto_increment |
| title       | varchar(100) | NO   |     | NULL              |                |
| description | text         | YES  |     | NULL              |                |
| type        | enum(...)    | YES  |     | 科研              |                |
| status      | enum(...)    | YES  |     | draft             |                |
| teacher_id  | bigint       | NO   | MUL | 2                 |                |
| submitted_at| datetime     | YES  |     | NULL              |                |
| student_id  | bigint       | NO   | MUL | NULL              |                |
| created_at  | datetime     | YES  |     | CURRENT_TIMESTAMP |                |
| updated_at  | datetime     | YES  |     | CURRENT_TIMESTAMP | on update ...  |
+-------------+--------------+------+-----+-------------------+----------------+
```

### student_teacher表
```
+----+------------+-----------+---------------------+
| id | student_id | teacher_id| bind_time           |
+----+------------+-----------+---------------------+
|  1 |          3 |         2 | 2024-01-15 10:00:00 |
+----+------------+-----------+---------------------+
```

## 启动服务

迁移完成后，启动后端服务：

```bash
cd go-backend
go run main.go
```

## 测试系统

使用测试脚本验证系统功能：

```bash
test_refactored_system.bat
```

## 常见问题

### Q1: 字段已存在错误
**A**: 这是正常的，说明该字段已经迁移过了，可以忽略这个错误。

### Q2: 表已存在错误
**A**: 这是正常的，说明该表已经创建过了，可以忽略这个错误。

### Q3: 外键约束已存在错误
**A**: 这是正常的，说明该约束已经添加过了，可以忽略这个错误。

### Q4: 索引已存在错误
**A**: 这是正常的，说明该索引已经创建过了，可以忽略这个错误。

## 成功标志

迁移成功的标志：
1. 所有SQL语句执行完成（即使有错误信息）
2. 最终显示 "Migration completed!" 消息
3. 能够看到完整的projects表结构
4. student_teacher表存在且有数据

## 下一步

1. 验证迁移结果
2. 启动后端服务
3. 测试API接口
4. 开始使用新的师生关系管理功能 