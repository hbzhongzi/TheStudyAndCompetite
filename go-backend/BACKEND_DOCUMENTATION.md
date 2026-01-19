# 云梦高校科研竞赛管理系统 - 后端文档

## 概述
本文档整合了云梦高校科研竞赛管理系统后端的所有重要功能说明、API文档和开发指南。

## 目录
1. [系统架构](#系统架构)
2. [核心功能模块](#核心功能模块)
3. [API接口文档](#api接口文档)
4. [数据库设计](#数据库设计)
5. [部署指南](#部署指南)
6. [开发指南](#开发指南)
7. [常见问题](#常见问题)

## 系统架构

### 技术栈
- **编程语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM框架**: GORM
- **认证机制**: JWT
- **中间件**: CORS、认证、日志等

### 项目结构
```
go-backend/
├── config/                  # 配置管理
├── controllers/             # 控制器层
├── models/                  # 数据模型
├── services/                # 业务逻辑层
├── routes/                  # 路由配置
├── middlewares/             # 中间件
├── utils/                   # 工具函数
├── migrations/              # 数据库迁移
├── sql/                     # SQL脚本
└── docs/                    # 文档目录
```

## 核心功能模块

### 1. 用户认证与授权
- **JWT认证**: 基于Token的身份验证
- **角色权限**: 管理员、教师、学生三种角色
- **中间件保护**: 路由级别的权限控制
- **Token刷新**: 支持Token自动刷新

### 2. 项目管理模块
- **项目CRUD**: 完整的项目生命周期管理
- **状态管理**: 项目状态流转和审核
- **文件管理**: 项目文件上传和管理
- **进度跟踪**: 项目进度更新和监控
- **质量评估**: 多维度项目质量评估

### 3. 竞赛管理模块
- **竞赛信息**: 竞赛的创建、更新、删除
- **报名管理**: 学生报名和审核
- **评审系统**: 教师评审和打分
- **结果管理**: 竞赛结果和获奖信息

### 4. 用户管理模块
- **用户信息**: 用户基本信息的CRUD操作
- **角色管理**: 用户角色分配和权限管理
- **批量操作**: 支持批量用户导入和管理
- **状态管理**: 用户状态管理（激活/禁用）

### 5. 系统管理模块
- **系统监控**: 系统健康状态监控
- **日志管理**: 系统日志记录和查询
- **配置管理**: 系统配置参数管理
- **备份管理**: 数据备份和恢复

## API接口文档

### 认证接口
```
POST   /api/login                    # 用户登录
POST   /api/refresh-token            # 刷新Token
GET    /api/validate-token           # 验证Token
GET    /api/health                   # 健康检查
```

### 用户管理接口
```
GET    /api/users                    # 获取用户列表
GET    /api/users/:id                # 获取用户详情
POST   /api/users                    # 创建用户
PUT    /api/users/:id                # 更新用户
DELETE /api/users/:id                # 删除用户
PATCH  /api/users/:id/status         # 切换用户状态
POST   /api/users/:id/reset-password # 重置密码
POST   /api/users/batch-delete       # 批量删除用户
GET    /api/users/stats              # 获取用户统计
GET    /api/users/export             # 导出用户数据
```

### 项目管理接口
```
GET    /api/projects                 # 获取项目列表
GET    /api/projects/:id             # 获取项目详情
POST   /api/projects                 # 创建项目
PUT    /api/projects/:id             # 更新项目
POST   /api/projects/submit/:id      # 提交项目审核
PUT    /api/projects/:id/status      # 更新项目状态
GET    /api/projects/:id/status-history # 获取状态历史
POST   /api/projects/:id/milestones  # 创建里程碑
PUT    /api/projects/milestones/:id  # 更新里程碑
GET    /api/projects/:id/milestones  # 获取里程碑列表
POST   /api/projects/:id/extensions  # 申请延期
PUT    /api/projects/:id/progress    # 更新项目进度
POST   /api/projects/:id/files       # 上传项目文件
GET    /api/projects/:id/files       # 获取项目文件
```

### 教师专用接口
```
GET    /api/teachers/projects        # 获取指导的项目
GET    /api/teachers/students        # 获取指导的学生
PUT    /api/teacher-projects/:id/review # 审核项目
GET    /api/teacher-projects/:id/reviews # 获取审核记录
PUT    /api/teacher-projects/files/:id/review # 审核文件
POST   /api/teacher-projects/reviews/:id/delegate # 委托审核
GET    /api/teacher-projects/my-review-tasks # 获取审核任务
```

### 管理员专用接口
```
GET    /api/admin/dashboard          # 获取仪表板数据
GET    /api/admin/dashboard/stats    # 获取仪表板统计
GET    /api/admin/overview           # 获取用户概览
GET    /api/admin/logs               # 获取系统日志
GET    /api/admin/health             # 获取系统健康状态
GET    /api/admin/system/health      # 获取系统健康状态
GET    /api/admin/stats              # 获取系统统计
PUT    /api/admin/maintenance        # 更新维护模式
POST   /api/admin/logs/cleanup       # 清理过期日志
```

### 竞赛管理接口
```
GET    /api/competitions             # 获取竞赛列表
GET    /api/competitions/:id         # 获取竞赛详情
POST   /api/competitions             # 创建竞赛
PUT    /api/competitions/:id         # 更新竞赛
DELETE /api/competitions/:id         # 删除竞赛
GET    /api/competitions/stats       # 获取竞赛统计
```

### 学生竞赛接口
```
POST   /api/student-competitions/:id/register # 报名竞赛
GET    /api/student-competitions/my # 查看报名记录
POST   /api/student-competitions/:id/upload # 上传参赛作品
```

### 教师竞赛接口
```
GET    /api/teacher-competitions/:id/submissions # 查看提交作品
POST   /api/teacher-competitions/:id/feedback # 提交评语
POST   /api/teacher-competitions/submissions/:id/scores # 提交评分
GET    /api/teacher-competitions/submissions/:id/scores # 获取评分列表
```

### 管理员竞赛接口
```
GET    /api/admin/competitions/:id/registrations # 查看报名情况
POST   /api/admin/competitions/:id/result # 登记成绩
GET    /api/admin/competitions/:id/export # 导出竞赛数据
POST   /api/admin/competitions/:id/judges # 分配评审教师
GET    /api/admin/competitions/:id/judges # 获取评审教师
GET    /api/admin/competitions/:id/judging-progress # 获取评审进度
POST   /api/admin/competitions/:id/finalize # 最终确认成绩
```

### 通知系统接口
```
GET    /api/notifications            # 获取通知列表
GET    /api/notifications/unread-count # 获取未读数量
PUT    /api/notifications/:id/read   # 标记已读
PUT    /api/notifications/read-all   # 标记全部已读
DELETE /api/notifications/:id        # 删除通知
```

### 管理员通知接口
```
GET    /api/admin/notifications/templates # 获取通知模板
PUT    /api/admin/notifications/templates/:id # 更新通知模板
POST   /api/admin/notifications/send # 发送通知
```

## 数据库设计

### 核心表结构
- **users**: 用户信息表
- **projects**: 项目信息表
- **competitions**: 竞赛信息表
- **project_files**: 项目文件表
- **project_reviews**: 项目审核记录表
- **competition_registrations**: 竞赛报名表
- **competition_submissions**: 竞赛提交作品表
- **system_logs**: 系统日志表
- **notifications**: 通知消息表

### 关系设计
- 用户与项目：一对多关系
- 用户与竞赛：多对多关系（通过报名表）
- 项目与文件：一对多关系
- 项目与审核记录：一对多关系

## 部署指南

### 环境要求
- **操作系统**: Linux/Windows/macOS
- **Go版本**: 1.21+
- **MySQL版本**: 8.0+
- **内存**: 最少2GB
- **磁盘**: 最少10GB可用空间

### 安装步骤
```bash
# 1. 克隆项目
git clone <repository-url>
cd go-backend

# 2. 安装依赖
go mod tidy

# 3. 配置数据库
cp config/database.example.yml config/database.yml
# 编辑数据库配置

# 4. 运行数据库迁移
go run migrations/migrate.go

# 5. 启动服务
go run main.go
```

### 环境变量配置
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=cloud_dream_system

# 服务配置
PORT=8080
JWT_SECRET=your-secret-key
JWT_EXPIRE=24
```

### Docker部署
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
```

## 开发指南

### 代码结构
- **控制器层**: 处理HTTP请求和响应
- **服务层**: 业务逻辑处理
- **模型层**: 数据模型定义
- **中间件**: 请求预处理和后处理

### 开发规范
- **命名规范**: 使用驼峰命名法
- **错误处理**: 统一的错误处理机制
- **日志记录**: 关键操作记录日志
- **代码注释**: 重要功能添加注释

### 测试
```bash
# 运行单元测试
go test ./...

# 运行特定测试
go test ./controllers

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 常见问题

### 1. 数据库连接问题
- **连接失败**: 检查数据库配置和网络
- **权限不足**: 确保数据库用户有足够权限
- **字符集问题**: 使用UTF8MB4字符集

### 2. 认证问题
- **Token过期**: 检查JWT配置和过期时间
- **权限不足**: 检查用户角色和权限配置
- **CORS问题**: 检查跨域配置

### 3. 性能问题
- **响应慢**: 检查数据库查询和索引
- **内存占用**: 检查内存泄漏和资源管理
- **并发问题**: 检查并发控制和锁机制

## 更新日志

### v1.0.0 (2024-01-15)
- 初始版本发布
- 基础API实现
- 用户认证系统

### v1.1.0 (2024-01-20)
- 项目管理API完善
- 竞赛管理模块
- 文件上传功能

### v1.2.0 (2024-01-25)
- 教师端API增强
- 批量操作支持
- 质量评估系统
- 系统监控功能

## 技术支持

如有技术问题或建议，请联系开发团队：
- **邮箱**: dev@yunmeng.edu.cn
- **文档**: 查看项目README文件
- **问题反馈**: 提交GitHub Issue

---

*最后更新: 2024-01-25*
*版本: v1.2.0* 