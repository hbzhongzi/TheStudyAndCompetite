# 迁移脚本故障排除指南

## 问题描述

在运行批处理文件时遇到以下错误：
```
'2.' is not recognized as an internal or external command
'.' is not recognized as an internal or external command
```

## 问题原因

这个错误通常由以下原因引起：
1. 批处理文件中的特殊字符编码问题
2. Windows系统对某些Unicode字符的支持问题
3. 批处理文件格式问题

## 解决方案

### 方案1：使用简化版脚本（推荐）

我们提供了简化版的迁移脚本，避免特殊字符：

**文件列表：**
- `sql/add_teacher_id_simple.sql` - 简化版SQL脚本
- `run_simple_migration.bat` - 简化版批处理文件

**执行步骤：**
```bash
# 直接运行简化版脚本
run_simple_migration.bat
```

### 方案2：手动执行SQL

如果批处理文件仍有问题，可以手动执行SQL：

1. **打开MySQL命令行：**
   ```bash
   mysql -hlocalhost -uroot -p123456
   ```

2. **选择数据库：**
   ```sql
   USE cloud_dream_system;
   ```

3. **执行迁移SQL：**
   ```sql
   -- 添加 teacher_id 列
   ALTER TABLE projects ADD COLUMN teacher_id BIGINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '指导老师ID' AFTER student_id;
   
   -- 添加外键约束
   ALTER TABLE projects ADD CONSTRAINT fk_projects_teacher_id FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE RESTRICT;
   
   -- 创建索引
   CREATE INDEX idx_projects_teacher_id ON projects(teacher_id);
   
   -- 更新现有数据
   UPDATE projects SET teacher_id = 2 WHERE teacher_id = 0 OR teacher_id IS NULL;
   ```

4. **验证结果：**
   ```sql
   DESCRIBE projects;
   SELECT COUNT(*) as total_projects FROM projects;
   ```

### 方案3：使用MySQL Workbench

1. 打开MySQL Workbench
2. 连接到数据库
3. 打开 `sql/add_teacher_id_simple.sql` 文件
4. 执行脚本

## 验证迁移是否成功

### 1. 检查表结构
```sql
DESCRIBE projects;
```
应该看到 `teacher_id` 列。

### 2. 检查外键约束
```sql
SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME 
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = 'cloud_dream_system' 
AND TABLE_NAME = 'projects' 
AND REFERENCED_TABLE_NAME IS NOT NULL;
```

### 3. 检查数据
```sql
SELECT COUNT(*) as total_projects, 
       COUNT(CASE WHEN teacher_id = 2 THEN 1 END) as projects_with_teacher 
FROM projects;
```

## 常见错误及解决方案

### 1. 连接错误
```
ERROR 2003 (HY000): Can't connect to MySQL server
```
**解决：** 检查MySQL服务是否启动

### 2. 权限错误
```
ERROR 1045 (28000): Access denied
```
**解决：** 检查用户名和密码

### 3. 数据库不存在
```
ERROR 1049 (42000): Unknown database
```
**解决：** 先创建数据库或修改数据库名称

### 4. 表不存在
```
ERROR 1146 (42S02): Table doesn't exist
```
**解决：** 先运行初始化脚本 `init_users.sql`

### 5. 列已存在
```
ERROR 1060 (42S21): Duplicate column name
```
**解决：** 这是正常情况，说明迁移已经完成

## 预防措施

### 1. 备份数据
执行迁移前备份数据库：
```bash
mysqldump -hlocalhost -uroot -p123456 cloud_dream_system > backup.sql
```

### 2. 测试环境
先在测试环境执行迁移脚本

### 3. 检查版本
确保MySQL版本兼容（5.7+ 或 8.0+）

## 联系支持

如果问题仍然存在，请提供：
1. MySQL版本信息
2. 操作系统版本
3. 完整的错误信息
4. 批处理文件内容

## 成功标志

迁移成功后，你应该能够：
1. 正常创建项目
2. 不再出现500错误
3. 可以正常选择指导老师
4. 项目与教师正确关联 