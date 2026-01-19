# 云梦高校科研竞赛管理系统 - API文档

## 概述

本文档描述了云梦高校科研竞赛管理系统的后端API接口。系统采用RESTful API设计，使用JWT进行身份认证。

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: Bearer Token (JWT)
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

所有API接口都遵循统一的响应格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

### 状态码说明

- `200`: 操作成功
- `400`: 请求参数错误
- `401`: 未认证或认证失败
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

## 认证接口

### 用户登录

**接口地址**: `POST /api/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "123456",
  "role": "admin"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@yunmeng.edu.cn",
      "role": "admin"
    }
  }
}
```

**错误响应**:
```json
{
  "code": 401,
  "message": "账号或密码错误"
}
```

## 用户管理接口

### 获取用户列表

**接口地址**: `GET /api/users`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 20)
- `search`: 搜索关键词 (可选)
- `role`: 角色筛选 (可选: admin, teacher, student)
- `status`: 状态筛选 (可选: active, inactive)
- `department`: 部门筛选 (可选)
- `sortBy`: 排序字段 (可选)
- `sortOrder`: 排序方向 (可选: asc, desc)

**请求示例**:
```
GET /api/users?page=1&size=10&search=张三&role=student&status=active
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取用户列表成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@yunmeng.edu.cn",
        "realName": "系统管理员",
        "phone": "13800138000",
        "department": "信息技术部",
        "studentId": null,
        "status": "active",
        "roles": ["admin"],
        "createTime": "2024-01-01T00:00:00Z",
        "lastLogin": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10,
    "pages": 1
  }
}
```

### 获取用户详情

**接口地址**: `GET /api/users/{id}`

**请求头**:
```
Authorization: Bearer <token>
```

**路径参数**:
- `id`: 用户ID

**请求示例**:
```
GET /api/users/1
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取用户详情成功",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@yunmeng.edu.cn",
    "realName": "系统管理员",
    "phone": "13800138000",
    "department": "信息技术部",
    "studentId": null,
    "avatar": "https://example.com/avatar.jpg",
    "bio": "系统管理员",
    "interests": ["技术", "管理"],
    "status": "active",
    "roles": [
      {
        "id": 1,
        "roleKey": "admin",
        "roleName": "系统管理员",
        "description": "拥有系统所有权限"
      }
    ],
    "createTime": "2024-01-01T00:00:00Z",
    "updateTime": "2024-01-15T10:30:00Z",
    "lastLogin": "2024-01-15T10:30:00Z"
  }
}
```

### 创建用户

**接口地址**: `POST /api/users`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "username": "student001",
  "password": "123456",
  "email": "student001@yunmeng.edu.cn",
  "realName": "张三",
  "roleKeys": ["student"],
  "phone": "13800138001",
  "department": "计算机学院",
  "studentId": "2021001",
  "avatar": "https://example.com/avatar.jpg",
  "bio": "计算机专业学生",
  "interests": ["编程", "算法"]
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "用户创建成功",
  "data": {
    "id": 2,
    "username": "student001",
    "email": "student001@yunmeng.edu.cn"
  }
}
```

### 更新用户

**接口地址**: `PUT /api/users/{id}`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**路径参数**:
- `id`: 用户ID

**请求参数**:
```json
{
  "realName": "张三（已更新）",
  "email": "zhangsan_new@yunmeng.edu.cn",
  "phone": "13800138002",
  "department": "计算机科学与技术学院",
  "bio": "计算机专业学生，热爱编程"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户更新成功"
}
```

### 删除用户

**接口地址**: `DELETE /api/users/{id}`

**请求头**:
```
Authorization: Bearer <token>
```

**路径参数**:
- `id`: 用户ID

**请求示例**:
```
DELETE /api/users/2
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户删除成功"
}
```

### 切换用户状态

**接口地址**: `PATCH /api/users/{id}/status`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**路径参数**:
- `id`: 用户ID

**请求参数**:
```json
{
  "status": "inactive"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户状态更新成功"
}
```

### 重置用户密码

**接口地址**: `POST /api/users/{id}/reset-password`

**请求头**:
```
Authorization: Bearer <token>
```

**路径参数**:
- `id`: 用户ID

**请求示例**:
```
POST /api/users/2/reset-password
```

**响应示例**:
```json
{
  "code": 200,
  "message": "密码重置成功",
  "data": {
    "newPassword": "Abc123456"
  }
}
```

### 批量删除用户

**接口地址**: `POST /api/users/batch-delete`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "userIds": [2, 3, 4]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "批量删除成功",
  "data": {
    "deletedCount": 3
  }
}
```

### 获取用户统计信息

**接口地址**: `GET /api/users/stats`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/users/stats
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取统计信息成功",
  "data": {
    "totalUsers": 100,
    "activeUsers": 95,
    "inactiveUsers": 5,
    "adminUsers": 3,
    "teacherUsers": 20,
    "studentUsers": 77,
    "todayLogins": 45,
    "thisWeekLogins": 280,
    "departmentStats": [
      {
        "department": "计算机学院",
        "count": 30
      },
      {
        "department": "数学学院",
        "count": 25
      }
    ],
    "roleStats": [
      {
        "role": "admin",
        "count": 3
      },
      {
        "role": "teacher",
        "count": 20
      },
      {
        "role": "student",
        "count": 77
      }
    ]
  }
}
```

### 导出用户数据

**接口地址**: `GET /api/users/export`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `format`: 导出格式 (excel, csv, json)
- `role`: 角色筛选 (可选)
- `status`: 状态筛选 (可选)
- `department`: 部门筛选 (可选)

**请求示例**:
```
GET /api/users/export?format=excel&role=student
```

**响应**: 文件下载

## 管理员专用接口

### 获取仪表板数据

**接口地址**: `GET /api/admin/dashboard`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/admin/dashboard
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取仪表板数据成功",
  "data": {
    "userStats": {
      "totalUsers": 100,
      "activeUsers": 95,
      "inactiveUsers": 5,
      "adminUsers": 3,
      "teacherUsers": 20,
      "studentUsers": 77
    },
    "systemInfo": {
      "serverTime": "2024-01-15 10:30:00",
      "version": "1.0.0",
      "status": "running"
    },
    "quickActions": [
      {
        "name": "创建用户",
        "action": "create_user",
        "icon": "user-add"
      },
      {
        "name": "批量导入",
        "action": "batch_import",
        "icon": "upload"
      },
      {
        "name": "数据导出",
        "action": "export_data",
        "icon": "download"
      },
      {
        "name": "系统设置",
        "action": "system_settings",
        "icon": "setting"
      }
    ]
  }
}
```

### 获取用户概览

**接口地址**: `GET /api/admin/overview`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/admin/overview
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取用户概览成功",
  "data": {
    "recentUsers": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@yunmeng.edu.cn",
        "realName": "系统管理员",
        "status": "active",
        "createTime": "2024-01-01T00:00:00Z"
      }
    ],
    "activeUserCount": 95,
    "todayNewUsers": 5,
    "lastUpdated": "2024-01-15 10:30:00"
  }
}
```

### 获取系统日志

**接口地址**: `GET /api/admin/logs`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 10)
- `type`: 日志类型 (可选: info, warning, error, debug)

**请求示例**:
```
GET /api/admin/logs?page=1&size=10&type=error
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统日志成功",
  "data": {
    "list": [
      {
        "id": 1,
        "type": "info",
        "message": "系统启动成功",
        "timestamp": "2024-01-15 09:30:00",
        "user": "system"
      },
      {
        "id": 2,
        "type": "warning",
        "message": "用户登录失败",
        "timestamp": "2024-01-15 10:00:00",
        "user": "admin"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 10
  }
}
```

### 获取系统设置

**接口地址**: `GET /api/admin/settings`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/admin/settings
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统设置成功",
  "data": {
    "system": {
      "name": "云梦高校科研竞赛管理系统",
      "version": "1.0.0",
      "description": "高校学生科研与竞赛项目管理系统",
      "maintenance": false
    },
    "security": {
      "passwordMinLength": 6,
      "passwordComplexity": true,
      "sessionTimeout": 3600,
      "maxLoginAttempts": 5
    },
    "email": {
      "enabled": true,
      "smtpServer": "smtp.yunmeng.edu.cn",
      "smtpPort": 587,
      "fromAddress": "noreply@yunmeng.edu.cn"
    },
    "database": {
      "type": "MySQL",
      "version": "8.0",
      "backup": true,
      "backupTime": "02:00"
    }
  }
}
```

### 更新系统设置

**接口地址**: `PUT /api/admin/settings`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "system": {
    "maintenance": false
  },
  "security": {
    "passwordMinLength": 8,
    "sessionTimeout": 7200
  }
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "系统设置更新成功"
}
```

### 获取系统健康状态

**接口地址**: `GET /api/admin/health`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/admin/health
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统健康状态成功",
  "data": {
    "status": "healthy",
    "checks": [
      {
        "name": "数据库连接",
        "status": "healthy",
        "message": "连接正常",
        "latency": "5ms"
      },
      {
        "name": "内存使用",
        "status": "healthy",
        "message": "使用率正常",
        "usage": "45%"
      },
      {
        "name": "磁盘空间",
        "status": "healthy",
        "message": "空间充足",
        "usage": "30%"
      },
      {
        "name": "网络连接",
        "status": "healthy",
        "message": "连接正常",
        "latency": "10ms"
      }
    ],
    "lastChecked": "2024-01-15 10:30:00"
  }
}
```

### 获取数据报表

**接口地址**: `GET /api/admin/reports`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `type`: 报表类型 (user_growth, user_activity, department_distribution)
- `startDate`: 开始日期 (可选)
- `endDate`: 结束日期 (可选)

**请求示例**:
```
GET /api/admin/reports?type=user_growth&startDate=2024-01-01&endDate=2024-01-15
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取数据报表成功",
  "data": {
    "title": "用户增长趋势",
    "data": [
      {
        "date": "2024-01-01",
        "count": 100
      },
      {
        "date": "2024-01-02",
        "count": 120
      },
      {
        "date": "2024-01-03",
        "count": 150
      }
    ]
  }
}
```

### 导出数据

**接口地址**: `GET /api/admin/export`

**请求头**:
```
Authorization: Bearer <token>
```

**查询参数**:
- `type`: 数据类型 (users, logs, reports)
- `format`: 导出格式 (excel, csv, json)

**请求示例**:
```
GET /api/admin/export?type=users&format=excel
```

**响应**: 文件下载

### 获取备份状态

**接口地址**: `GET /api/admin/backup/status`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
GET /api/admin/backup/status
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取备份状态成功",
  "data": {
    "lastBackup": "2024-01-14 02:00:00",
    "nextBackup": "2024-01-15 02:00:00",
    "backupSize": "256MB",
    "backupStatus": "success",
    "autoBackup": true,
    "backupRetention": 30
  }
}
```

### 创建备份

**接口地址**: `POST /api/admin/backup`

**请求头**:
```
Authorization: Bearer <token>
```

**请求示例**:
```
POST /api/admin/backup
```

**响应示例**:
```json
{
  "code": 200,
  "message": "数据备份创建成功",
  "data": {
    "backupId": "backup_20240115103000",
    "backupTime": "2024-01-15 10:30:00",
    "status": "completed"
  }
}
```

## 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求参数格式和必填字段 |
| 401 | 未认证 | 检查Authorization头是否正确 |
| 403 | 权限不足 | 确认用户角色是否有相应权限 |
| 404 | 资源不存在 | 检查资源ID是否正确 |
| 500 | 服务器错误 | 联系系统管理员 |

### 错误响应示例

```json
{
  "code": 400,
  "message": "参数错误: 用户名不能为空"
}
```

## 认证说明

### JWT Token格式

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token有效期

- 默认有效期: 24小时
- 过期后需要重新登录

### 权限说明

- **admin**: 拥有所有权限，可以管理所有用户和系统
- **teacher**: 可以查看学生信息，管理自己的项目
- **student**: 只能查看自己的信息

## 测试工具

### 使用curl测试

```bash
# 登录获取token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456","role":"admin"}'

# 使用token获取用户列表
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 获取管理员仪表板数据
curl -X GET http://localhost:8080/api/admin/dashboard \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 使用测试脚本

```bash
# 运行用户API测试脚本
go run test_api.go

# 运行完整系统测试脚本
go run test_full_system.go

# 运行管理员API测试脚本
go run test_admin_api.go
```

## 数据模型

### 用户基础信息 (users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 用户ID，主键 |
| username | VARCHAR(50) | 用户名，唯一 |
| password | VARCHAR(100) | 加密密码 |
| email | VARCHAR(100) | 邮箱，唯一 |
| status | ENUM | 账户状态 (active/inactive) |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |

### 用户详细信息 (user_profiles)

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | BIGINT | 用户ID，外键 |
| real_name | VARCHAR(50) | 真实姓名 |
| phone | VARCHAR(20) | 手机号 |
| department | VARCHAR(100) | 所属学院或部门 |
| student_id | VARCHAR(50) | 学号（仅限学生） |
| avatar | VARCHAR(255) | 头像URL |
| bio | TEXT | 简介 |
| interests | JSON | 兴趣（JSON数组） |
| last_login | DATETIME | 最近登录时间 |

### 角色信息 (roles)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 角色ID，主键 |
| role_key | VARCHAR(30) | 角色标识，唯一 |
| role_name | VARCHAR(50) | 角色名称 |
| description | TEXT | 描述 |

## 部署说明

### 环境要求

- Go 1.20+
- MySQL 8.0+
- 内存: 512MB+
- 磁盘: 1GB+

### 快速启动

```bash
# Windows
quick_start.bat

# Linux/Mac
chmod +x quick_start.sh
./quick_start.sh
```

### 配置说明

环境变量配置：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=cloud_dream_system

# 服务配置
PORT=8080
JWT_SECRET=your_jwt_secret
```

## 联系信息

- **技术支持**: admin@yunmeng.edu.cn
- **项目地址**: [GitHub Repository]
- **文档更新**: 2024年1月 