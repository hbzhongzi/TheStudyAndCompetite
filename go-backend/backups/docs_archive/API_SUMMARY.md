# 原始路径: go-backend/docs/API_SUMMARY.md
# 归档时间: 2026-01-19

````markdown
# 云梦高校科研竞赛管理系统 - 后端API接口总结（更新版）

## 系统概述

**基础URL**: `http://localhost:8080/api`  
**认证方式**: Bearer Token (JWT)  
**数据格式**: JSON  
**字符编码**: UTF-8

## API接口分类

### 1. 认证接口

#### 1.1 用户登录
- **接口**: `POST /api/login`
- **权限**: 无需认证
- **作用**: 用户登录获取JWT令牌
- **参数**: username, password, role
- **返回**: JWT token + 用户信息

---

### 2. 用户管理接口（管理员权限）

#### 2.1 用户列表管理
- **接口**: `GET /api/users`
- **权限**: 管理员
- **作用**: 获取用户列表，支持分页、搜索、筛选
- **参数**: page, size, search, role, status, department, sortBy, sortOrder

#### 2.2 用户详情
- **接口**: `GET /api/users/{id}`
- **权限**: 管理员
- **作用**: 获取指定用户的详细信息

#### 2.3 创建用户
- **接口**: `POST /api/users`
- **权限**: 管理员
- **作用**: 创建新用户账户
- **参数**: username, password, email, realName, roleKeys, phone, department, studentId等

#### 2.4 更新用户
- **接口**: `PUT /api/users/{id}`
- **权限**: 管理员
- **作用**: 更新用户信息

#### 2.5 删除用户
- **接口**: `DELETE /api/users/{id}`
- **权限**: 管理员
- **作用**: 删除指定用户

#### 2.6 切换用户状态
- **接口**: `PATCH /api/users/{id}/status`
- **权限**: 管理员
- **作用**: 激活/禁用用户账户

#### 2.7 重置用户密码
- **接口**: `POST /api/users/{id}/reset-password`
- **权限**: 管理员
- **作用**: 重置用户密码为默认密码

#### 2.8 批量删除用户
- **接口**: `POST /api/users/batch-delete`
- **权限**: 管理员
- **作用**: 批量删除多个用户

#### 2.9 用户统计信息
- **接口**: `GET /api/users/stats`
- **权限**: 管理员
- **作用**: 获取用户统计数据和图表信息

---

### 3. 项目模块接口

（内容已归档）

````
