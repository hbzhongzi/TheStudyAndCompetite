# 云梦高校科研竞赛管理系统 - 项目总结

## 项目概述

云梦高校科研竞赛管理系统是一个综合性的科研项目管理平台，采用前后端分离架构，支持项目申报、竞赛管理、成果展示等功能。系统已完成用户管理模块的完整实现，包括前端界面、后端API、数据库设计和权限控制。

## 技术架构

### 前端技术栈
- **Vue 3**: 现代化的前端框架，使用Composition API
- **Element Plus**: 企业级UI组件库，提供丰富的界面组件
- **ECharts**: 数据可视化图表库，用于统计报表展示
- **Vue Router**: 前端路由管理
- **Axios**: HTTP客户端，用于API调用
- **Vite**: 构建工具，提供快速的开发体验

### 后端技术栈
- **Go 1.20+**: 后端编程语言，高性能和并发处理
- **Gin**: 轻量级Web框架，提供RESTful API
- **GORM**: ORM框架，简化数据库操作
- **MySQL**: 关系型数据库，存储系统数据
- **JWT**: 身份认证和授权
- **bcrypt**: 密码加密算法

## 数据库设计

### 核心表结构

#### 1. 用户基础信息表 (users)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '加密密码',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '账户状态',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='用户账号信息表';
```

#### 2. 用户扩展信息表 (user_profiles)
```sql
CREATE TABLE user_profiles (
    user_id BIGINT PRIMARY KEY COMMENT '用户ID',
    real_name VARCHAR(50) COMMENT '真实姓名',
    phone VARCHAR(20) COMMENT '手机号',
    department VARCHAR(100) COMMENT '所属学院或部门',
    student_id VARCHAR(50) COMMENT '学号（仅限学生）',
    avatar VARCHAR(255) COMMENT '头像URL',
    bio TEXT COMMENT '简介',
    interests JSON COMMENT '兴趣（JSON数组）',
    last_login DATETIME COMMENT '最近登录时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT='用户详细信息表';
```

#### 3. 角色定义表 (roles)
```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_key VARCHAR(30) NOT NULL UNIQUE COMMENT '角色标识',
    role_name VARCHAR(50) NOT NULL COMMENT '角色名称',
    description TEXT COMMENT '描述'
) COMMENT='系统角色表';
```

#### 4. 用户角色关联表 (user_roles)
```sql
CREATE TABLE user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) COMMENT='用户与角色关联表';
```

#### 5. 登录记录表 (login_logs)
```sql
CREATE TABLE login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(50),
    user_agent TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT='用户登录日志表';
```

### 关系设计
- **一个用户 ↔ 一个详细信息**: 1对1关系
- **一个用户 ↔ 多个角色**: 多对多关系（通过user_roles表）
- **一个用户 ↔ 多条登录记录**: 1对多关系

## 功能模块

### 1. 用户管理模块

#### 前端功能
- **用户列表展示**: 支持分页、搜索、筛选
- **用户详情查看**: 完整的用户信息展示
- **用户创建**: 表单验证和角色分配
- **用户编辑**: 信息更新和角色管理
- **用户删除**: 单个和批量删除
- **状态管理**: 启用/禁用用户账户
- **密码重置**: 管理员重置用户密码
- **数据导出**: 支持Excel格式导出
- **统计报表**: 用户数据可视化展示

#### 后端API
- **GET /api/users**: 获取用户列表（支持分页、搜索、筛选）
- **GET /api/users/{id}**: 获取用户详情
- **POST /api/users**: 创建用户
- **PUT /api/users/{id}**: 更新用户信息
- **DELETE /api/users/{id}**: 删除用户
- **PATCH /api/users/{id}/status**: 切换用户状态
- **POST /api/users/{id}/reset-password**: 重置密码
- **POST /api/users/batch-delete**: 批量删除用户
- **GET /api/users/stats**: 获取用户统计信息
- **GET /api/users/export**: 导出用户数据

### 2. 认证授权模块

#### JWT认证
- 基于JWT的身份认证
- Token自动刷新机制
- 安全的密码加密存储

#### 权限控制
- 基于角色的权限控制（RBAC）
- 中间件级别的权限验证
- 细粒度的API权限管理

#### 角色定义
- **admin**: 系统管理员，拥有所有权限
- **teacher**: 教师，可以管理项目和指导学生
- **student**: 学生，可以参与项目和竞赛

### 3. 系统管理模块

#### 管理员界面
- **仪表板**: 系统概览和快速操作
- **用户管理**: 完整的用户CRUD操作
- **项目管理**: 科研项目管理（待实现）
- **竞赛管理**: 竞赛活动管理（待实现）
- **系统设置**: 系统配置管理（待实现）
- **数据统计**: 系统数据分析和报表（待实现）
- **系统日志**: 操作日志查看和管理（待实现）

## 项目结构

```
云梦高校科研竞赛管理系统/
├── yunmeng-frontend/          # 前端项目
│   ├── src/
│   │   ├── views/            # 页面组件
│   │   │   ├── user/         # 用户相关页面
│   │   │   │   ├── AdminView.vue           # 管理员主界面
│   │   │   │   ├── UserManagement.vue      # 用户管理
│   │   │   │   ├── ProjectManagement.vue   # 项目管理
│   │   │   │   ├── CompetitionManagement.vue # 竞赛管理
│   │   │   │   ├── SystemSettings.vue      # 系统设置
│   │   │   │   ├── DataReports.vue         # 数据统计
│   │   │   │   └── SystemLogs.vue          # 系统日志
│   │   │   ├── login/        # 登录页面
│   │   │   ├── dashboard/    # 仪表板
│   │   │   ├── student/      # 学生页面
│   │   │   └── teacher/      # 教师页面
│   │   ├── services/         # API服务
│   │   │   └── userService.js # 用户管理API
│   │   ├── router/           # 路由配置
│   │   ├── store/            # 状态管理
│   │   └── components/       # 公共组件
│   ├── package.json
│   └── README.md
├── go-backend/               # 后端项目
│   ├── config/              # 配置文件
│   │   └── database.go      # 数据库配置
│   ├── controllers/         # 控制器层
│   │   ├── login.go         # 登录控制器
│   │   └── user_controller.go # 用户管理控制器
│   ├── middlewares/         # 中间件
│   │   └── auth.go          # 认证中间件
│   ├── models/              # 数据模型
│   │   └── user.go          # 用户相关模型
│   ├── routes/              # 路由配置
│   │   └── routes.go        # 路由定义
│   ├── services/            # 服务层
│   │   └── user_service.go  # 用户管理服务
│   ├── sql/                 # 数据库脚本
│   │   └── init_users.sql   # 用户表初始化脚本
│   ├── utils/               # 工具函数
│   │   ├── jwt.go           # JWT工具
│   │   └── password.go      # 密码工具
│   ├── main.go              # 主程序入口
│   ├── go.mod               # Go模块文件
│   ├── start.sh             # Linux/Mac启动脚本
│   ├── start.bat            # Windows启动脚本
│   └── README.md            # 项目说明文档
├── sql/                     # 数据库脚本
├── mock-server.js           # 模拟API服务器
├── package.json             # 模拟服务器配置
├── API_DOCS.md             # API接口文档
├── README.md               # 项目说明文档
└── PROJECT_SUMMARY.md      # 项目总结文档
```

## 核心特性

### 1. 现代化架构
- **前后端分离**: 清晰的职责分离，便于开发和维护
- **RESTful API**: 标准的API设计，易于集成和扩展
- **组件化开发**: 可复用的前端组件，提高开发效率

### 2. 安全性
- **JWT认证**: 无状态的用户认证机制
- **密码加密**: 使用bcrypt算法加密存储密码
- **权限控制**: 基于角色的细粒度权限管理
- **输入验证**: 前后端双重数据验证

### 3. 用户体验
- **响应式设计**: 适配不同屏幕尺寸
- **直观界面**: 基于Element Plus的美观界面
- **实时反馈**: 操作结果即时反馈
- **数据可视化**: 图表展示统计数据

### 4. 可扩展性
- **模块化设计**: 清晰的代码结构，便于扩展
- **配置化管理**: 环境变量配置，便于部署
- **数据库抽象**: ORM层抽象，支持多种数据库

## 部署说明

### 环境要求
- Node.js 16+
- Go 1.20+
- MySQL 8.0+

### 快速部署

#### 1. 启动后端服务
```bash
cd go-backend
go mod tidy
go run main.go
```

#### 2. 启动前端服务
```bash
cd yunmeng-frontend
npm install
npm run dev
```

#### 3. 启动模拟服务器（可选）
```bash
npm install
npm start
```

### 访问地址
- **前端**: http://localhost:5173
- **后端API**: http://localhost:8080/api
- **模拟服务器**: http://localhost:8080/api

## 开发进度

### 已完成功能 ✅
- [x] 用户管理模块（完整CRUD操作）
- [x] 认证授权系统
- [x] 角色权限控制
- [x] 数据库设计和迁移
- [x] API接口文档
- [x] 前端用户管理界面
- [x] 数据统计和可视化
- [x] 系统日志记录

### 待开发功能 🔄
- [ ] 项目管理模块
- [ ] 竞赛管理模块
- [ ] 系统设置模块
- [ ] 文件上传功能
- [ ] 邮件通知系统
- [ ] 移动端适配
- [ ] 性能优化
- [ ] 单元测试

## 技术亮点

### 1. 数据库设计
- **规范化设计**: 符合第三范式，避免数据冗余
- **外键约束**: 保证数据完整性
- **索引优化**: 提高查询性能
- **JSON字段**: 灵活存储复杂数据

### 2. 后端架构
- **分层架构**: 控制器、服务、模型清晰分离
- **中间件机制**: 统一的认证和权限控制
- **事务管理**: 保证数据一致性
- **错误处理**: 统一的错误响应格式

### 3. 前端架构
- **组件化**: 可复用的Vue组件
- **状态管理**: 响应式数据管理
- **路由管理**: 单页应用路由
- **API封装**: 统一的接口调用

### 4. 开发体验
- **热重载**: 开发时实时更新
- **类型安全**: TypeScript支持
- **代码规范**: ESLint和Prettier
- **文档完善**: 详细的API文档

## 总结

云梦高校科研竞赛管理系统已完成用户管理模块的完整实现，包括：

1. **完整的用户管理功能**: 从数据库设计到前端界面的完整实现
2. **安全的认证系统**: JWT认证和权限控制
3. **现代化的技术栈**: Vue 3 + Go + MySQL
4. **良好的代码结构**: 清晰的分层架构和模块化设计
5. **完善的文档**: API文档和项目说明

系统具有良好的可扩展性和维护性，为后续功能模块的开发奠定了坚实的基础。通过前后端分离的架构，可以独立开发和部署，提高了开发效率和系统稳定性。 