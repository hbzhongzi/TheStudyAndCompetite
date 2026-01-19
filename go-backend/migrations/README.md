# 竞赛数据库同步更新说明

## 概述

本目录包含了竞赛管理系统的数据库同步更新脚本，用于将现有数据库结构更新到与新的Go后端模型完全一致的状态。

## 文件说明

### 1. 主要更新脚本

- **`update_competition_database_complete.sql`** - 完整的数据库更新脚本
  - 添加报名时间相关字段
  - 添加竞赛详细信息字段
  - 创建竞赛关联表
  - 添加性能优化索引
  - 设置外键约束
  - 更新现有数据

### 2. 执行脚本

- **`execute_database_update.bat`** - Windows批处理脚本
- **`execute_database_update.sh`** - Linux/Mac Shell脚本

### 3. 文档

- **`DATABASE_STRUCTURE_COMPARISON.md`** - 数据库结构对比文档
- **`README.md`** - 本文档

## 更新内容

### 新增字段

#### 竞赛表 (competitions)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| `registration_start` | DATETIME | 报名开始时间 |
| `registration_end` | DATETIME | 报名截止时间 |
| `submission_deadline` | DATETIME | 作品提交截止时间 |
| `location` | VARCHAR(255) | 竞赛地点 |
| `contact` | VARCHAR(255) | 联系方式 |
| `rules` | TEXT | 竞赛规则 |
| `requirements` | TEXT | 报名要求 |
| `judging_method` | TEXT | 评审方式 |
| `important_notes` | TEXT | 重要注意事项 |
| `website` | VARCHAR(500) | 相关链接 |
| `qq_group` | VARCHAR(50) | QQ群 |
| `file_formats` | VARCHAR(255) | 支持的文件格式 |
| `file_size_limit` | VARCHAR(100) | 文件大小限制 |
| `department_limit` | VARCHAR(255) | 院系限制 |
| `teacher_limit` | BOOLEAN | 教师限制 |
| `award_config` | JSON | 获奖配置 |

#### 竞赛报名表 (competition_registrations)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| `team_name` | VARCHAR(100) | 团队名称 |
| `team_leader` | BOOLEAN | 团队负责人 |
| `contact_phone` | VARCHAR(20) | 联系电话 |
| `contact_email` | VARCHAR(100) | 联系邮箱 |
| `additional_info` | JSON | 额外信息 |

#### 竞赛提交表 (competition_submissions)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| `version` | INT | 版本号 |
| `locked` | BOOLEAN | 是否锁定 |
| `teacher_viewed` | BOOLEAN | 教师查看状态 |
| `teacher_feedback` | TEXT | 教师反馈 |
| `last_view_time` | DATETIME | 最后查看时间 |

### 新增表

- **`competition_feedback`** - 竞赛反馈表
- **`competition_results`** - 竞赛结果表
- **`competition_judges`** - 竞赛评审教师表
- **`competition_scores`** - 竞赛评分表

### 性能优化

- 添加了复合索引优化查询性能
- 设置了外键约束确保数据完整性
- 优化了表结构布局

## 执行步骤

### 1. 准备工作

1. **备份数据库**
   ```bash
   # 使用mysqldump备份
   mysqldump -u用户名 -p密码 数据库名 > backup.sql
   ```

2. **检查MySQL客户端**
   - 确保MySQL客户端已安装
   - 确保有足够的数据库权限

3. **选择执行脚本**
   - Windows用户：使用 `execute_database_update.bat`
   - Linux/Mac用户：使用 `execute_database_update.sh`

### 2. 执行更新

#### Windows用户
```cmd
cd go-backend/migrations
execute_database_update.bat
```

#### Linux/Mac用户
```bash
cd go-backend/migrations
chmod +x execute_database_update.sh
./execute_database_update.sh
```

### 3. 手动执行（可选）

如果自动脚本无法使用，可以手动执行：

```bash
mysql -u用户名 -p密码 数据库名 < update_competition_database_complete.sql
```

## 验证更新结果

### 1. 检查表结构
```sql
-- 检查竞赛表结构
DESCRIBE competitions;

-- 检查竞赛报名表结构
DESCRIBE competition_registrations;

-- 检查竞赛提交表结构
DESCRIBE competition_submissions;
```

### 2. 检查新字段
```sql
-- 检查新字段是否有数据
SELECT 
    COUNT(registration_start) as has_reg_start,
    COUNT(registration_end) as has_reg_end,
    COUNT(location) as has_location,
    COUNT(contact) as has_contact
FROM competitions;
```

### 3. 检查索引
```sql
-- 检查索引是否创建成功
SHOW INDEX FROM competitions;
SHOW INDEX FROM competition_registrations;
SHOW INDEX FROM competition_submissions;
```

### 4. 检查外键约束
```sql
-- 检查外键约束
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
```

## 注意事项

### 1. 安全注意事项
- **务必先备份数据库**
- 在测试环境验证脚本
- 确保有足够的数据库权限
- 在生产环境执行前充分测试

### 2. 兼容性说明
- 新字段设置为 `NULL`，兼容现有数据
- 现有数据会自动设置合理的默认值
- 不会影响现有的业务逻辑

### 3. 性能影响
- 添加字段和索引会短暂影响性能
- 建议在业务低峰期执行
- 大量数据时可能需要较长时间

### 4. 回滚方案
如果更新失败，可以使用备份文件恢复：
```bash
mysql -u用户名 -p密码 数据库名 < backup.sql
```

## 常见问题

### 1. 权限不足
```
ERROR 1142 (42000): ALTER command denied to user
```
**解决方案**：确保用户有 `ALTER` 权限

### 2. 字段已存在
```
Duplicate column name 'registration_start'
```
**解决方案**：脚本会自动检查，不会重复添加

### 3. 外键约束失败
```
Cannot add foreign key constraint
```
**解决方案**：检查引用的表是否存在，字段类型是否匹配

### 4. 索引创建失败
```
Duplicate key name 'idx_competitions_registration_time'
```
**解决方案**：脚本使用 `IF NOT EXISTS`，不会重复创建

## 技术支持

如果在执行过程中遇到问题：

1. 检查错误日志
2. 确认数据库版本兼容性
3. 验证数据库连接信息
4. 检查磁盘空间是否充足

## 更新日志

- **2024-01-01** - 初始版本，添加报名时间字段和竞赛详细信息
- 支持完整的竞赛管理功能
- 优化查询性能
- 增强数据完整性

---

**重要提醒**：执行数据库更新前，请务必备份数据库，并在测试环境验证脚本的正确性。 