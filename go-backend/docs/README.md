# 云梦高校科研竞赛管理系统 - 后端服务

## 项目概述

这是云梦高校科研竞赛管理系统的后端服务，使用Go语言和Gin框架开发，提供RESTful API接口。

## 技术栈

- **Go 1.20+**: 后端编程语言
- **Gin**: Web框架
- **GORM**: ORM框架
- **MySQL**: 数据库
- **JWT**: 身份认证
- **bcrypt**: 密码加密

## 快速开始

## 文档索引

本目录的主索引已集中到 `DOCS_INDEX.md`，请优先查看该文件以获得合并后的文档导航和归档说明。


### 环境要求

- Go 1.20+
- MySQL 8.0+
- Git

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-backend
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

确保MySQL服务已启动，并创建数据库：

```sql
CREATE DATABASE cloud_dream_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 初始化数据库

#### Windows:
```cmd
init_database.bat
```

#### Linux/Mac:
```bash
chmod +x init_database.sh
./init_database.sh
```

#### 手动执行:
```bash
mysql -u root -p < sql/migrate_existing.sql
```

### 5. 启动服务

#### Windows:
```cmd
start.bat
```

#### Linux/Mac:
```bash
chmod +x start.sh
./start.sh
```

#### 手动启动:
```bash
go run main.go
```

### 6. 验证服务

服务启动后，访问以下地址验证：

- 后端API: http://localhost:8080/api
- 登录接口: http://localhost:8080/api/login

## 故障排除

### 常见问题

#### 1. 数据库连接失败
```
错误: 数据库连接失败
```
**解决方案:**
- 检查MySQL服务是否启动
- 验证数据库连接参数
- 确认数据库用户权限

#### 2. 数据库迁移错误
```
错误: Incorrect table definition; there can be only one auto column and it must be defined as a key
```
**解决方案:**
- 执行数据库初始化脚本：`init_database.bat` 或 `./init_database.sh`
- 或者手动执行：`mysql -u root -p < sql/migrate_existing.sql`

#### 3. 端口被占用
```
错误: 服务器启动失败: listen tcp :8080: bind: address already in use
```
**解决方案:**
- 修改PORT环境变量：`set PORT=8081` (Windows) 或 `export PORT=8081` (Linux/Mac)
- 或停止占用端口的进程

#### 4. 依赖安装失败
```
错误: go mod tidy 失败
```
**解决方案:**
- 检查网络连接
- 配置Go代理：`go env -w GOPROXY=https://goproxy.cn,direct`

### 测试数据库连接

运行简化测试程序：

```bash
go run test_simple.go
```

### 测试用户服务

运行用户服务测试：

```bash
go run test_user_service.go
```

## 项目结构

```
go-backend/
├── config/              # 配置文件
│   └── database.go      # 数据库配置
├── controllers/         # 控制器层
│   ├── login.go         # 登录控制器
│   └── user_controller.go # 用户管理控制器
├── middlewares/         # 中间件
│   └── auth.go          # 认证中间件
├── models/              # 数据模型
│   └── user.go          # 用户相关模型
├── routes/              # 路由配置
│   └── routes.go        # 路由定义
├── services/            # 服务层
│   └── user_service.go  # 用户管理服务
├── sql/                 # 数据库脚本
│   ├── init_users.sql   # 用户表初始化脚本
│   └── migrate_existing.sql # 数据库迁移脚本
├── utils/               # 工具函数
│   ├── jwt.go           # JWT工具
│   └── password.go      # 密码工具
├── main.go              # 主程序入口
├── test_simple.go       # 简化测试程序
├── test_user_service.go # 用户服务测试
├── init_database.bat    # Windows数据库初始化脚本
├── init_database.sh     # Linux/Mac数据库初始化脚本
├── start.bat            # Windows启动脚本
├── start.sh             # Linux/Mac启动脚本
├── go.mod               # Go模块文件
└── README.md            # 项目说明文档
```

## API接口

### 认证接口

#### 用户登录
```
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456",
  "role": "admin"
}
```

### 用户管理接口（需要管理员权限）

#### 获取用户列表
```
GET /api/users?page=1&size=20&search=张三&role=student&status=active
Authorization: Bearer <token>
```

#### 获取用户详情
```
GET /api/users/{id}
Authorization: Bearer <token>
```

#### 创建用户
```
POST /api/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "student002",
  "password": "123456",
  "email": "student002@yunmeng.edu.cn",
  "realName": "李四",
  "roleKeys": ["student"],
  "phone": "13800138001",
  "department": "计算机学院",
  "studentId": "2021002"
}
```

#### 更新用户
```
PUT /api/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "realName": "李四（已更新）",
  "email": "lisi_new@yunmeng.edu.cn",
  "phone": "13800138002"
}
```

#### 删除用户
```
DELETE /api/users/{id}
Authorization: Bearer <token>
```

#### 切换用户状态
```
PATCH /api/users/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "inactive"
}
```

#### 重置用户密码
```
POST /api/users/{id}/reset-password
Authorization: Bearer <token>
```

#### 批量删除用户
```
POST /api/users/batch-delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "userIds": [1, 2, 3]
}
```

#### 获取用户统计信息
```
GET /api/users/stats
Authorization: Bearer <token>
```

#### 导出用户数据
```
GET /api/users/export?format=excel&role=student
Authorization: Bearer <token>
```

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

## 权限控制

系统采用基于角色的权限控制（RBAC）：

- **admin**: 系统管理员，拥有所有权限
- **teacher**: 教师，可以管理项目和指导学生
- **student**: 学生，可以参与项目和竞赛

## 开发指南

### 添加新的API接口

1. 在 `models/` 目录下定义数据模型
2. 在 `services/` 目录下实现业务逻辑
3. 在 `controllers/` 目录下实现控制器
4. 在 `routes/routes.go` 中注册路由

### 数据库迁移

系统使用GORM的AutoMigrate功能自动创建和更新数据库表结构。

### 日志记录

系统使用Go标准库的log包记录日志，包括：
- 数据库连接状态
- API请求日志
- 错误信息

## 部署

### 开发环境

```bash
go run main.go
```

### 生产环境

```bash
# 编译
go build -o yunmeng-backend main.go

# 运行
./yunmeng-backend
```

### Docker部署

```dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o yunmeng-backend main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/yunmeng-backend .
EXPOSE 8080
CMD ["./yunmeng-backend"]
```

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [Issues]
- 邮箱: admin@yunmeng.edu.cn

## 许可证

MIT License 