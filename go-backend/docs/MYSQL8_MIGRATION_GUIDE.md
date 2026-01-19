此文件内容已归档并移至 `go-backend/backups/docs_archive/MYSQL8_MIGRATION_GUIDE.md`，可安全删除原件以释放空间。
- `run_teacher_id_migration_mysql8.sh` - Linux/Mac 执行脚本

## 执行步骤

### Windows 用户
```bash
# 直接运行批处理文件
run_teacher_id_migration_mysql8.bat
```

### Linux/Mac 用户
```bash
# 设置执行权限
chmod +x run_teacher_id_migration_mysql8.sh

# 执行脚本
./run_teacher_id_migration_mysql8.sh
```

## 迁移内容

### 1. 表结构修改
```sql
-- 添加 teacher_id 列
ALTER TABLE projects 
ADD COLUMN teacher_id BIGINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '指导老师ID' AFTER student_id;
```

### 2. 外键约束
```sql
-- 添加外键约束
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_teacher_id 
FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE RESTRICT;
```

### 3. 索引创建
```sql
-- 创建索引
CREATE INDEX idx_projects_teacher_id ON projects(teacher_id);
```

### 4. 数据更新
```sql
-- 更新现有项目
UPDATE projects SET teacher_id = 2 WHERE teacher_id = 0 OR teacher_id IS NULL;
```

### 5. 默认教师创建
```sql
-- 创建默认教师（如果不存在）
INSERT IGNORE INTO users (id, username, password, email, status) VALUES
(2, 'teacher001', '$2a$12$JDDRq/VWcwIoTV2mNdxUeOGB5ZMzbs4Ye540Zj/vKrqGLd7vyPTV2', 'li.teacher@yunmeng.edu.cn', 'active');
```

## 验证输出

迁移脚本会输出以下验证信息：

### 1. 表结构验证
```
+------------------------+-------------+-------------+------+-----+-------------------+----------------+
| TABLE_NAME             | COLUMN_NAME | DATA_TYPE   | NULL | KEY | DEFAULT           | EXTRA          |
+------------------------+-------------+-------------+------+-----+-------------------+----------------+
| projects               | id          | bigint      | NO   | PRI | NULL              | auto_increment |
| projects               | student_id  | bigint      | NO   | MUL | NULL              |                |
| projects               | teacher_id  | bigint      | NO   | MUL | 2                 |                |
+------------------------+-------------+-------------+------+-----+-------------------+----------------+
```

### 2. 外键约束验证
```
+------------------------+-------------+-------------+------------------------+------------------------+
| CONSTRAINT_NAME        | COLUMN_NAME | REFERENCED_TABLE_NAME | REFERENCED_COLUMN_NAME |
+------------------------+-------------+-------------+------------------------+------------------------+
| fk_projects_student_id | student_id  | users                 | id                     |
| fk_projects_teacher_id | teacher_id  | users                 | id                     |
+------------------------+-------------+-------------+------------------------+------------------------+
```

### 3. 数据统计验证
```
+------------------+------------------------+------------------------+
| total_projects   | projects_with_default_teacher | draft_projects |
+------------------+------------------------+------------------------+
| 3                | 3                      | 1                |
+------------------+------------------------+------------------------+
```

## 错误处理

### 常见错误及解决方案

1. **连接错误**
   ```
   ERROR 2003 (HY000): Can't connect to MySQL server
   ```
   - 检查MySQL服务是否启动
   - 验证连接参数

2. **权限错误**
   ```
   ERROR 1045 (28000): Access denied
   ```
   - 检查用户名和密码
   - 验证用户权限

3. **数据库不存在**
   ```
   ERROR 1049 (42000): Unknown database
   ```
   - 先创建数据库
   - 或修改数据库名称

4. **外键约束错误**
   ```
   ERROR 1452 (23000): Cannot add or update a child row
   ```
   - 确保引用的用户ID存在
   - 检查外键约束设置

## 性能优化

### 1. 索引优化
- 为 `teacher_id` 创建索引
- 为常用查询字段创建复合索引

### 2. 数据类型优化
- 使用 `BIGINT UNSIGNED` 提高性能
- 使用 `ENUM` 类型节省存储空间

### 3. 约束优化
- 使用 `ON DELETE RESTRICT` 保护数据完整性
- 使用 `ON UPDATE CASCADE` 自动更新

## 备份建议

执行迁移前，建议备份数据库：

```bash
# 备份整个数据库
mysqldump -hlocalhost -uroot -p123456 cloud_dream_system > backup_before_migration.sql

# 备份特定表
mysqldump -hlocalhost -uroot -p123456 cloud_dream_system projects > projects_backup.sql
```

## 回滚方案

如果需要回滚迁移：

```sql
-- 删除外键约束
ALTER TABLE projects DROP FOREIGN KEY fk_projects_teacher_id;

-- 删除索引
DROP INDEX idx_projects_teacher_id ON projects;

-- 删除列
ALTER TABLE projects DROP COLUMN teacher_id;
```

## 注意事项

1. **版本要求**: 仅支持 MySQL 8.0+
2. **权限要求**: 需要 ALTER, CREATE, INSERT, SELECT 权限
3. **数据备份**: 执行前请备份重要数据
4. **服务重启**: 迁移后需要重启后端服务
5. **测试验证**: 迁移后请测试项目创建功能

## 联系支持

如遇到问题，请检查：
1. MySQL版本是否为8.0+
2. 数据库连接参数是否正确
3. 用户权限是否足够
4. 数据库是否存在

迁移完成后，项目创建功能应该正常工作，不再出现500错误。 