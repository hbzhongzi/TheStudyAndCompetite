# SQL文件说明

## 文件整理完成

根据您的要求，已经整理了SQL文件夹，只保留了数据库初始化文件和插入模拟数据文件，删除了所有迁移脚本。

## 保留的文件

### `database_setup.sql` - 完整数据库设置脚本 (MySQL 8.0+)
**用途**: 完整的数据库初始化，包含所有表结构和测试数据
**特点**:
- 包含所有表结构创建
- 包含完整的索引创建
- 包含触发器创建
- 包含存储过程和函数
- 包含视图创建
- 包含事件创建
- 包含完整的测试数据
- 包含系统管理相关表结构
- 包含竞赛评审相关表结构
- 兼容MySQL 8.0+版本

**使用方法**:
```bash
mysql -u root -p < database_setup.sql
```

### `database_setup_compatible.sql` - 兼容版数据库设置脚本 (MySQL 5.7)
**用途**: 兼容MySQL 5.7及以下版本的数据库初始化
**特点**:
- 移除所有COMMENT字段定义（MySQL 5.7不支持）
- 将JSON字段改为TEXT字段
- 保持所有功能完整性
- 兼容MySQL 5.7及以下版本

**使用方法**:
```bash
mysql -u root -p < database_setup_compatible.sql
```

## 已删除的文件

以下迁移脚本文件已被删除：

- `fix_competition_award_config.sql` - 竞赛奖项配置修复脚本
- `database_enhancement_migration.sql` - 数据库增强迁移脚本
- `system_management_enhancement.sql` - 系统管理增强脚本
- `competition_enhancement_migration.sql` - 竞赛增强迁移脚本

这些脚本的功能已经整合到主要的数据库初始化文件中。

## 版本兼容性

### MySQL 8.0+ 版本
- 使用 `database_setup.sql`
- 支持所有新特性（COMMENT、JSON字段等）

### MySQL 5.7 及以下版本
- 使用 `database_setup_compatible.sql`
- 移除不兼容的特性
- 保持功能完整性

## 数据库结构

### 用户管理模块 (5个表)
1. **users** - 用户表
2. **roles** - 角色表
3. **user_roles** - 用户角色关联表
4. **user_profiles** - 用户档案表
5. **login_logs** - 登录日志表

### 项目管理模块 (2个表)
1. **project_types** - 项目分类表
2. **projects** - 项目表

### 竞赛管理模块 (8个表)
1. **competitions** - 竞赛信息表
2. **competition_registrations** - 竞赛报名记录表
3. **competition_submissions** - 竞赛成果提交表
4. **competition_feedback** - 竞赛教师评语表
5. **competition_judges** - 竞赛评审教师表
6. **competition_scores** - 竞赛评分记录表
7. **competition_results** - 竞赛获奖登记表
8. **competition_audit_logs** - 竞赛审计日志表

### 系统管理模块 (6个表)
1. **system_logs** - 系统日志表
2. **system_settings** - 系统配置表
3. **system_health_logs** - 系统健康监控表
4. **system_performance_logs** - 系统性能监控表
5. **system_alerts** - 系统告警表
6. **system_diagnostics** - 系统诊断表
7. **backup_records** - 备份记录表

### 文件管理模块 (1个表)
1. **files** - 文件表

## 测试数据

### 默认用户
- **admin** (密码: 123456) - 系统管理员
- **teacher001** (密码: 123456) - 张教授
- **teacher002** (密码: 123456) - 李副教授
- **student001** (密码: 123456) - 王同学
- **student002** (密码: 123456) - 刘同学

### 默认角色
- **admin** - 系统管理员
- **teacher** - 教师
- **student** - 学生

### 项目类型
- 学术研究
- 创新项目
- 艺术项目
- 技术项目
- 社会实践

### 系统配置
- 系统名称和版本
- 文件上传配置
- 安全配置
- 备份配置
- 邮件配置
- 监控配置

## 快速开始

1. **选择适合的脚本**:
   - MySQL 8.0+: 使用 `database_setup.sql`
   - MySQL 5.7: 使用 `database_setup_compatible.sql`

2. **执行数据库初始化**:
   ```bash
   mysql -u root -p < database_setup.sql
   ```

3. **验证安装**:
   ```sql
   USE cloud_dream_system;
   SHOW TABLES;
   SELECT COUNT(*) FROM users;
   ```

## 注意事项

- 执行脚本前请确保MySQL服务正在运行
- 脚本会自动创建 `cloud_dream_system` 数据库
- 如果数据库已存在，脚本会安全地添加缺失的表和数据
- 所有测试数据都使用 `INSERT IGNORE` 语句，不会重复插入 