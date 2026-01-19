# 云梦高校科研竞赛管理系统 - 迁移问题总结

## 问题描述

在执行数据库迁移脚本时遇到了以下问题：

1. **MySQL语法错误**：`ADD COLUMN IF NOT EXISTS` 语法在某些MySQL版本中不被支持
2. **表已存在错误**：`student_teacher` 表可能已经存在
3. **批处理文件格式问题**：Windows批处理文件中的多行SQL语句格式不正确

## 解决方案

### 方案1：使用错误处理迁移脚本（推荐）

我们已经创建了 `sql/error_handled_migration.sql` 脚本，它会：
此文件内容已归档并移至 `go-backend/backups/docs_archive/MIGRATION_ISSUE_SUMMARY.md`，可安全删除原件以释放空间。
3. 确认数据库连接参数
4. 查看 `MANUAL_MIGRATION_GUIDE.md` 进行手动迁移 