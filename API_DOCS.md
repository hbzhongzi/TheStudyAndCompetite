# 云梦高校科研竞赛管理系统 - API接口文档

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: Bearer Token (JWT)
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": "2024-01-15T14:30:00.000Z"
}
```

### 分页响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "size": 20,
    "pages": 5
  },
  "timestamp": "2024-01-15T14:30:00.000Z"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": [
    {
      "field": "username",
      "message": "用户名不能为空"
    }
  ],
  "timestamp": "2024-01-15T14:30:00.000Z"
}
```

## 用户管理接口

### 1. 获取用户列表

**接口地址**: `GET /users`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | Integer | 否 | 页码，默认1 |
| size | Integer | 否 | 每页大小，默认20 |
| search | String | 否 | 搜索关键词（用户名、姓名、邮箱） |
| role | String | 否 | 角色筛选（admin/teacher/student） |
| status | String | 否 | 状态筛选（active/inactive） |
| sortBy | String | 否 | 排序字段（createTime/username/realName） |
| sortOrder | String | 否 | 排序方向（asc/desc） |

**请求示例**:
```
GET /api/users?page=1&size=20&search=张三&role=student&status=active
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
        "username": "student001",
        "realName": "张三",
        "email": "zhangsan@yunmeng.edu.cn",
        "role": "student",
        "status": "active",
        "createTime": "2024-01-01T00:00:00.000Z",
        "lastLogin": "2024-01-15T14:30:00.000Z",
        "phone": "13800138000",
        "department": "计算机学院",
        "studentId": "2021001"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20,
    "pages": 5
  }
}
```

### 2. 获取用户详情

**接口地址**: `GET /users/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Long | 是 | 用户ID |

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
    "username": "student001",
    "realName": "张三",
    "email": "zhangsan@yunmeng.edu.cn",
    "role": "student",
    "status": "active",
    "createTime": "2024-01-01T00:00:00.000Z",
    "lastLogin": "2024-01-15T14:30:00.000Z",
    "phone": "13800138000",
    "department": "计算机学院",
    "studentId": "2021001",
    "profile": {
      "avatar": "avatar.jpg",
      "bio": "计算机科学与技术专业学生",
      "interests": ["人工智能", "机器学习"]
    }
  }
}
```

### 3. 创建用户

**接口地址**: `POST /users`

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | String | 是 | 用户名，3-20个字符 |
| realName | String | 是 | 真实姓名 |
| email | String | 是 | 邮箱地址 |
| password | String | 是 | 密码，6-20个字符 |
| role | String | 是 | 角色（admin/teacher/student） |
| phone | String | 否 | 手机号码 |
| department | String | 否 | 所属部门/学院 |
| studentId | String | 否 | 学号（学生角色必填） |

**请求示例**:
```json
{
  "username": "student002",
  "realName": "李四",
  "email": "lisi@yunmeng.edu.cn",
  "password": "123456",
  "role": "student",
  "phone": "13800138001",
  "department": "计算机学院",
  "studentId": "2021002"
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "用户创建成功",
  "data": {
    "id": 2,
    "username": "student002",
    "realName": "李四",
    "email": "lisi@yunmeng.edu.cn",
    "role": "student",
    "status": "active",
    "createTime": "2024-01-15T14:30:00.000Z"
  }
}
```

### 4. 更新用户信息

**接口地址**: `PUT /users/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Long | 是 | 用户ID |

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| realName | String | 否 | 真实姓名 |
| email | String | 否 | 邮箱地址 |
| role | String | 否 | 角色（admin/teacher/student） |
| phone | String | 否 | 手机号码 |
| department | String | 否 | 所属部门/学院 |
| studentId | String | 否 | 学号 |

**请求示例**:
```json
{
  "realName": "李四（已更新）",
  "email": "lisi_new@yunmeng.edu.cn",
  "phone": "13800138002"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户信息更新成功",
  "data": {
    "id": 2,
    "username": "student002",
    "realName": "李四（已更新）",
    "email": "lisi_new@yunmeng.edu.cn",
    "role": "student",
    "status": "active",
    "updateTime": "2024-01-15T15:00:00.000Z"
  }
}
```

### 5. 删除用户

**接口地址**: `DELETE /users/{id}`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Long | 是 | 用户ID |

**请求示例**:
```
DELETE /api/users/2
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户删除成功",
  "data": null
}
```

### 6. 启用/禁用用户

**接口地址**: `PATCH /users/{id}/status`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Long | 是 | 用户ID |

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | String | 是 | 状态（active/inactive） |

**请求示例**:
```json
{
  "status": "inactive"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户状态更新成功",
  "data": {
    "id": 2,
    "status": "inactive",
    "updateTime": "2024-01-15T15:30:00.000Z"
  }
}
```

### 7. 重置用户密码

**接口地址**: `POST /users/{id}/reset-password`

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Long | 是 | 用户ID |

**请求示例**:
```
POST /api/users/2/reset-password
```

**响应示例**:
```json
{
  "code": 200,
  "message": "密码重置成功，新密码已发送到用户邮箱",
  "data": {
    "id": 2,
    "resetTime": "2024-01-15T16:00:00.000Z"
  }
}
```

### 8. 批量删除用户

**接口地址**: `POST /users/batch-delete`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userIds | Array<Long> | 是 | 用户ID数组 |

**请求示例**:
```json
{
  "userIds": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "批量删除成功",
  "data": {
    "deletedCount": 3,
    "deletedIds": [1, 2, 3]
  }
}
```

### 9. 导出用户数据

**接口地址**: `GET /users/export`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| format | String | 否 | 导出格式（excel/csv），默认excel |
| search | String | 否 | 搜索关键词 |
| role | String | 否 | 角色筛选 |
| status | String | 否 | 状态筛选 |

**请求示例**:
```
GET /api/users/export?format=excel&role=student
```

**响应**: 文件流（Excel或CSV文件）

### 10. 获取用户统计信息

**接口地址**: `GET /users/stats`

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
    "totalUsers": 1234,
    "activeUsers": 1189,
    "inactiveUsers": 45,
    "roleStats": {
      "admin": 5,
      "teacher": 89,
      "student": 1140
    },
    "departmentStats": {
      "计算机学院": 234,
      "数学学院": 156,
      "物理学院": 123
    },
    "monthlyGrowth": [
      {
        "month": "2024-01",
        "newUsers": 45,
        "activeUsers": 890
      }
    ]
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已存在） |
| 422 | 数据验证失败 |
| 500 | 服务器内部错误 |

## 数据模型

### User（用户）
```json
{
  "id": "Long",
  "username": "String",
  "realName": "String",
  "email": "String",
  "role": "String (admin/teacher/student)",
  "status": "String (active/inactive)",
  "createTime": "DateTime",
  "updateTime": "DateTime",
  "lastLogin": "DateTime",
  "phone": "String",
  "department": "String",
  "studentId": "String",
  "profile": {
    "avatar": "String",
    "bio": "String",
    "interests": "Array<String>"
  }
}
```

## 权限说明

- **admin**: 拥有所有用户管理权限
- **teacher**: 只能查看和管理学生用户
- **student**: 无用户管理权限

## 注意事项

1. 所有接口都需要JWT认证（除了登录接口）
2. 密码在传输和存储时都需要加密
3. 用户名和邮箱必须唯一
4. 删除用户前需要检查是否有关联数据
5. 批量操作建议限制数量，避免性能问题
6. 导出功能建议异步处理，避免超时 