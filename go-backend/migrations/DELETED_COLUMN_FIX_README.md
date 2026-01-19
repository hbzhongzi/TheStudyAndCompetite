# 数据库缺失列修复说明

## 问题描述

前端调用 `GET /api/teachers/projects` 接口时，后端返回 500 错误：
```
Error 1054 (42S22): Unknown column 'deleted' in 'where clause'
```

## 问题原因

Go 模型 `Project` 和 `User` 中定义了 `Deleted` 字段，但数据库表中缺少对应的 `deleted` 列。这导致在查询时出现 "列不存在" 错误。

## 影响范围

以下查询受到影响：
- `GetTeacherProjects()` - 获取教师项目列表
- `GetTeacherListWithFilter()` - 获取教师列表（带过滤）
- `GetMyStudents()` - 获取我的学生列表

## 解决方案

### 1. 执行数据库迁移

运行以下命令添加缺失的列：

```bash
# Windows
migrations/execute_deleted_column_fix.bat

# Linux/Mac
mysql -hlocalhost -P3306 -uroot -p123456 cloud_dream_system < migrations/add_deleted_column.sql
```

### 2. 迁移内容

迁移脚本将添加以下列：

#### projects 表
- `deleted` - 软删除标记 (BOOLEAN, DEFAULT FALSE)
- 索引: `idx_projects_deleted`

#### users 表  
- `deleted` - 软删除标记 (BOOLEAN, DEFAULT FALSE)
- 索引: `idx_users_deleted`

### 3. 验证修复

执行迁移后，以下查询应该正常工作：
- `GET /api/teachers/projects` - 获取教师项目列表
- `GET /api/teachers` - 获取教师列表
- 其他涉及软删除功能的接口

## 技术细节

### 软删除机制

`deleted` 列用于实现软删除功能：
- `deleted = FALSE` - 记录正常
- `deleted = TRUE` - 记录已删除（软删除）

### 查询优化

添加索引 `idx_projects_deleted` 和 `idx_users_deleted` 以提高查询性能。

### 向后兼容

- 现有数据不受影响
- 新记录的 `deleted` 字段默认为 `FALSE`
- 现有查询逻辑保持不变

## 注意事项

1. **备份数据**: 执行迁移前请备份数据库
2. **权限要求**: 需要 ALTER TABLE 权限
3. **服务重启**: 迁移完成后无需重启后端服务
4. **测试验证**: 建议在测试环境先验证迁移脚本

## 相关文件

- `add_deleted_column.sql` - 迁移脚本
- `execute_deleted_column_fix.bat` - Windows 执行脚本
- `models/project.go` - Project 模型定义
- `models/user.go` - User 模型定义
- `services/project_service.go` - 项目服务层

## 联系支持

如果遇到问题，请检查：
1. 数据库连接参数是否正确
2. 用户权限是否足够
3. 数据库服务是否正常运行 