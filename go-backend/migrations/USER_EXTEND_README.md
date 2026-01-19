# 用户表扩展迁移说明

## 概述
本次迁移扩展了 `users` 表结构，添加了新的字段以支持更丰富的用户信息管理。

## 新增字段

### 基础信息字段
- `department` VARCHAR(100) - 部门/院系
- `title` VARCHAR(50) - 职称/职位  
- `grade` VARCHAR(20) - 年级
- `major` VARCHAR(100) - 专业

### 时间字段
- `created_at` TIMESTAMP - 创建时间
- `updated_at` TIMESTAMP - 更新时间

## 迁移文件
- `extend_user_table.sql` - 数据库迁移脚本
- `execute_user_extend.bat` - Windows批处理执行脚本
- `execute_user_extend.ps1` - PowerShell执行脚本

## 执行方法

### 方法1: 使用批处理脚本（推荐）
```bash
# 在 go-backend 目录下执行
scripts/execute_user_extend.bat
```

### 方法2: 使用PowerShell脚本
```powershell
# 在 go-backend 目录下执行
scripts/execute_user_extend.ps1
```

### 方法3: 手动执行SQL
```bash
mysql -u root -p < migrations/extend_user_table.sql
```

## 注意事项

1. **备份数据**: 执行迁移前请备份数据库
2. **MySQL服务**: 确保MySQL服务正在运行
3. **权限要求**: 需要数据库管理员权限
4. **默认值**: 新字段会为现有用户设置默认值"未设置"

## 索引优化
迁移脚本会自动创建以下索引以提高查询性能：
- `idx_users_department` - 部门索引
- `idx_users_title` - 职称索引  
- `idx_users_grade` - 年级索引
- `idx_users_major` - 专业索引
- `idx_users_created_at` - 创建时间索引

## 回滚方案
如果需要回滚，可以执行以下SQL：
```sql
ALTER TABLE users 
DROP COLUMN department,
DROP COLUMN title,
DROP COLUMN grade,
DROP COLUMN major,
DROP COLUMN created_at,
DROP COLUMN updated_at;
```

## 验证迁移
迁移完成后，可以使用以下命令验证表结构：
```sql
DESCRIBE users;
```

## 相关代码更新
- `models/user.go` - 已更新User结构体
- `services/project_service.go` - 已修复字段访问问题
- 其他相关服务文件 - 需要相应更新以使用新字段 