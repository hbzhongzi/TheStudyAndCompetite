# 项目管理模块功能增强 - 数据库更新说明

## 📋 更新概述

本次数据库更新旨在完善项目管理模块的五个核心功能，提升系统的完整性和用户体验。

## 🎯 更新目标

1. **项目状态管理增强** - 支持更完整的项目生命周期
2. **项目生命周期管理** - 添加进度跟踪、里程碑、延期申请
3. **成果文件管理增强** - 支持文件类型、版本管理、审核状态
4. **项目分类管理增强** - 支持层级管理、路径查询、统计
5. **审核流程增强** - 支持多级审核、委托、提醒

## 🗄️ 数据库变更详情

### 1. 项目状态管理增强

#### 修改表: `projects`
- **扩展状态枚举**: 从5个状态扩展到9个状态
  - 新增: `in_progress`(进行中), `completed`(已完成), `suspended`(暂停), `need_revision`(待修改)
- **新增字段**:
  - `status_change_reason` (TEXT) - 状态变更原因
  - `status_changed_by` (BIGINT UNSIGNED) - 状态变更操作人ID
  - `status_changed_at` (DATETIME) - 状态变更时间

#### 影响分析
- ✅ **向后兼容**: 现有数据不受影响
- ✅ **数据完整性**: 新增外键约束确保数据一致性
- ⚠️ **性能影响**: 轻微，新增索引优化查询性能

### 2. 项目生命周期管理增强

#### 修改表: `projects`
- **新增字段**:
  - `start_date` (DATETIME) - 项目开始时间
  - `expected_end_date` (DATETIME) - 预计完成时间
  - `actual_end_date` (DATETIME) - 实际完成时间
  - `progress` (INT) - 项目进度(0-100)
  - `milestone_count` (INT) - 里程碑数量
  - `is_extended` (BOOLEAN) - 是否延期
  - `extension_reason` (TEXT) - 延期原因
  - `extension_count` (INT) - 延期次数

#### 新增表: `project_milestones`
- **用途**: 管理项目里程碑
- **主要字段**: 标题、描述、预计时间、完成时间、状态、进度

#### 新增表: `project_extensions`
- **用途**: 管理项目延期申请
- **主要字段**: 申请人、原因、原定时间、申请时间、审核状态

#### 影响分析
- ✅ **功能增强**: 支持完整的项目进度管理
- ✅ **数据扩展**: 不影响现有项目数据
- ⚠️ **存储增长**: 每个项目增加约200字节存储空间

### 3. 成果文件管理增强

#### 修改表: `project_files`
- **新增字段**:
  - `file_type` (ENUM) - 文件类型(开题报告、中期报告、结题报告、成果展示、其他)
  - `file_version` (VARCHAR(20)) - 文件版本
  - `review_status` (ENUM) - 审核状态(待审核、已通过、已驳回)
  - `review_comments` (TEXT) - 审核意见
  - `reviewed_by` (BIGINT UNSIGNED) - 审核人ID
  - `reviewed_at` (DATETIME) - 审核时间
  - `file_size` (BIGINT) - 文件大小(字节)
  - `download_count` (INT) - 下载次数
  - `is_public` (BOOLEAN) - 是否公开

#### 新增表: `file_type_configs`
- **用途**: 配置文件类型规则
- **主要字段**: 类型、显示名称、描述、是否必需、最大大小、允许扩展名

#### 影响分析
- ✅ **功能完善**: 支持文件分类和审核流程
- ✅ **向后兼容**: 现有文件自动归类为"其他"类型
- ⚠️ **存储增长**: 每个文件记录增加约300字节

### 4. 项目分类管理增强

#### 修改表: `project_types`
- **新增字段**:
  - `parent_id` (BIGINT UNSIGNED) - 父分类ID
  - `level` (INT) - 分类层级
  - `sort_order` (INT) - 排序
  - `is_active` (BOOLEAN) - 是否启用
  - `icon` (VARCHAR(100)) - 分类图标
  - `color` (VARCHAR(20)) - 分类颜色
  - `project_count` (INT) - 项目数量

#### 新增表: `project_type_paths`
- **用途**: 支持分类树快速查询
- **主要字段**: 祖先ID、后代ID、层级深度

#### 新增表: `project_type_stats`
- **用途**: 分类统计信息
- **主要字段**: 各状态项目数量、最后更新时间

#### 影响分析
- ✅ **层级支持**: 支持无限层级的分类结构
- ✅ **性能优化**: 路径表支持快速查询分类树
- ⚠️ **复杂度增加**: 需要维护分类路径数据

### 5. 审核流程增强

#### 修改表: `project_reviews`
- **新增字段**:
  - `review_level` (INT) - 审核级别(1:指导教师,2:学院,3:学校)
  - `review_order` (INT) - 审核顺序
  - `is_required` (BOOLEAN) - 是否必需审核
  - `deadline` (DATETIME) - 审核截止时间
  - `auto_approve` (BOOLEAN) - 超时是否自动通过
  - `review_duration` (INT) - 审核耗时(分钟)
  - `is_urgent` (BOOLEAN) - 是否紧急

#### 新增表: `project_review_flows`
- **用途**: 配置审核流程
- **主要字段**: 项目类型、审核级别、审核角色、时限、自动通过

#### 新增表: `review_delegations`
- **用途**: 管理审核委托
- **主要字段**: 原审核人、被委托人、项目、原因、时间范围

#### 新增表: `review_reminders`
- **用途**: 审核提醒
- **主要字段**: 提醒类型、消息、发送状态

#### 影响分析
- ✅ **流程完善**: 支持多级审核和委托机制
- ✅ **自动化**: 支持超时自动通过和提醒
- ⚠️ **配置复杂**: 需要配置审核流程规则

## 📊 数据统计

### 新增表数量: 10个
1. `project_milestones` - 项目里程碑
2. `project_extensions` - 项目延期申请
3. `file_type_configs` - 文件类型配置
4. `project_type_paths` - 分类路径
5. `project_type_stats` - 分类统计
6. `project_review_flows` - 审核流程配置
7. `review_delegations` - 审核委托
8. `review_reminders` - 审核提醒
9. `project_notifications` - 项目通知
10. `notification_templates` - 通知模板

### 修改表数量: 3个
1. `projects` - 项目主表
2. `project_files` - 项目文件表
3. `project_types` - 项目分类表
4. `project_reviews` - 项目审核表

### 新增字段总数: 约25个
- 项目表: 8个字段
- 文件表: 9个字段
- 分类表: 7个字段
- 审核表: 7个字段

## 🔧 执行步骤

### 1. 准备工作
- 确保MySQL服务正在运行
- 确保有足够的数据库权限(ALTER, CREATE, INSERT, INDEX)
- 建议在非生产环境先测试

### 2. 执行更新
```bash
# Windows
execute_project_enhancement.bat

# Linux/Mac
mysql -h[host] -P[port] -u[user] -p[password] [database] < project_management_enhancement.sql
```

### 3. 验证结果
- 检查新表是否创建成功
- 检查字段是否添加成功
- 检查索引是否创建成功
- 测试基本功能是否正常

## ⚠️ 注意事项

### 1. 备份要求
- **必须**: 执行前创建完整数据库备份
- **建议**: 在测试环境先验证脚本正确性

### 2. 权限要求
- 需要ALTER权限修改表结构
- 需要CREATE权限创建新表
- 需要INDEX权限创建索引

### 3. 兼容性
- 现有数据完全保留
- 现有API接口不受影响
- 需要更新Go模型文件以支持新字段

### 4. 性能影响
- 新增字段会增加存储空间
- 新增索引会提升查询性能
- 建议在低峰期执行更新

## 🚀 后续工作

### 1. 代码更新
- 更新Go模型文件(`models/project.go`)
- 更新服务层代码(`services/project_service.go`)
- 更新控制器代码
- 更新路由配置

### 2. 功能测试
- 测试新状态转换逻辑
- 测试里程碑管理功能
- 测试文件分类和审核
- 测试多级审核流程

### 3. 用户培训
- 更新用户手册
- 培训管理员使用新功能
- 通知用户新功能可用

## 📞 技术支持

如果在执行过程中遇到问题，请：

1. 检查错误日志
2. 确认数据库权限
3. 验证SQL语法
4. 联系技术支持团队

---

**更新时间**: 2024年  
**版本**: 1.0  
**作者**: 云梦高校项目管理系统开发团队 