# 项目申报模块 API 文档

## 概述

项目申报模块为科研管理系统提供项目申报、审核、管理等功能。支持学生创建和提交项目，教师和管理员进行审核和管理。

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON

## 认证

所有API都需要在请求头中包含JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

## 学生项目API

### 1. 获取我的项目列表

**接口**: `GET /student/projects`

**权限**: 学生

**请求参数**:
```json
{
  "page": 1,           // 页码，默认1
  "size": 20,          // 每页数量，默认20
  "search": "",        // 搜索关键词（可选）
  "type": "科研",      // 项目类型筛选（可选）：科研、竞赛
  "status": "draft"    // 状态筛选（可选）：draft、pending、approved、rejected
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取我的项目成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "基于深度学习的图像识别系统",
        "description": "利用卷积神经网络实现高精度的图像分类和识别功能",
        "type": "科研",
        "status": "pending",
        "studentName": "张三",
        "studentId": "2021001",
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z",
        "memberCount": 3,
        "fileCount": 2,
        "reviewCount": 1
      }
    ],
    "total": 1,
    "page": 1,
    "size": 20
  }
}
```

### 2. 获取项目详情

**接口**: `GET /student/projects/{id}`

**权限**: 学生

**响应示例**:
```json
{
  "code": 200,
  "message": "获取项目详情成功",
  "data": {
    "id": 1,
    "title": "基于深度学习的图像识别系统",
    "description": "利用卷积神经网络实现高精度的图像分类和识别功能",
    "type": "科研",
    "status": "pending",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z",
    "student": {
      "id": 3,
      "username": "student001",
      "realName": "张三",
      "email": "zhangsan@yunmeng.edu.cn",
      "phone": "13800138002",
      "department": "计算机学院",
      "studentId": "2021001"
    },
    "members": [
      {
        "id": 1,
        "name": "张三",
        "studentNumber": "2021001",
        "role": "负责人"
      },
      {
        "id": 2,
        "name": "李四",
        "studentNumber": "2021002",
        "role": "成员"
      }
    ],
    "files": [
      {
        "id": 1,
        "fileName": "项目申请书.pdf",
        "fileUrl": "/uploads/projects/1/application.pdf",
        "uploadTime": "2024-01-15T10:30:00Z"
      }
    ],
    "reviews": [
      {
        "id": 1,
        "status": "pending",
        "comments": "项目方案较为完整，建议补充实验设计部分",
        "reviewTime": "2024-01-15T11:00:00Z",
        "reviewer": {
          "id": 2,
          "username": "teacher001",
          "realName": "李老师"
        }
      }
    ]
  }
}
```

### 3. 创建项目

**接口**: `POST /student/projects`

**权限**: 学生

**请求体**:
```json
{
  "title": "基于深度学习的图像识别系统",
  "description": "利用卷积神经网络实现高精度的图像分类和识别功能",
  "type": "科研",
  "members": [
    {
      "name": "张三",
      "studentNumber": "2021001",
      "role": "负责人"
    },
    {
      "name": "李四",
      "studentNumber": "2021002",
      "role": "成员"
    }
  ],
  "files": [
    {
      "fileName": "项目申请书.pdf",
      "fileUrl": "/uploads/projects/1/application.pdf"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "项目创建成功",
  "data": {
    "id": 1,
    "title": "基于深度学习的图像识别系统",
    "description": "利用卷积神经网络实现高精度的图像分类和识别功能",
    "type": "科研",
    "status": "draft",
    "studentId": 3,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

### 4. 更新项目

**接口**: `PUT /student/projects/{id}`

**权限**: 学生

**请求体**:
```json
{
  "title": "基于深度学习的图像识别系统（更新版）",
  "description": "利用卷积神经网络实现高精度的图像分类和识别功能，包含实验设计部分",
  "type": "科研",
  "members": [
    {
      "name": "张三",
      "studentNumber": "2021001",
      "role": "负责人"
    },
    {
      "name": "李四",
      "studentNumber": "2021002",
      "role": "成员"
    },
    {
      "name": "王五",
      "studentNumber": "2021003",
      "role": "成员"
    }
  ],
  "files": [
    {
      "fileName": "项目申请书.pdf",
      "fileUrl": "/uploads/projects/1/application.pdf"
    },
    {
      "fileName": "技术方案.docx",
      "fileUrl": "/uploads/projects/1/technical_plan.docx"
    }
  ]
}
```

### 5. 删除项目

**接口**: `DELETE /student/projects/{id}`

**权限**: 学生

**响应示例**:
```json
{
  "code": 200,
  "message": "项目删除成功"
}
```

### 6. 提交项目审核

**接口**: `POST /student/projects/{id}/submit`

**权限**: 学生

**响应示例**:
```json
{
  "code": 200,
  "message": "项目提交成功，等待审核"
}
```

## 教师/管理员项目API

### 1. 获取项目列表

**接口**: `GET /projects`

**权限**: 教师、管理员

**请求参数**:
```json
{
  "page": 1,           // 页码，默认1
  "size": 20,          // 每页数量，默认20
  "search": "",        // 搜索关键词（可选）
  "type": "科研",      // 项目类型筛选（可选）：科研、竞赛
  "status": "pending", // 状态筛选（可选）：draft、pending、approved、rejected
  "studentId": 3       // 学生ID筛选（可选）
}
```

### 2. 审核项目

**接口**: `POST /projects/{id}/review`

**权限**: 教师、管理员

**请求体**:
```json
{
  "status": "approved",  // 审核结果：approved、rejected
  "comments": "项目目标明确，技术路线合理，同意立项"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "项目审核完成"
}
```

### 3. 获取项目统计

**接口**: `GET /projects/stats`

**权限**: 教师、管理员

**响应示例**:
```json
{
  "code": 200,
  "message": "获取项目统计成功",
  "data": {
    "totalProjects": 10,
    "draftProjects": 3,
    "pendingProjects": 4,
    "approvedProjects": 2,
    "rejectedProjects": 1,
    "researchProjects": 7,
    "competitionProjects": 3,
    "typeStats": {
      "科研": 7,
      "竞赛": 3
    },
    "statusStats": {
      "draft": 3,
      "pending": 4,
      "approved": 2,
      "rejected": 1
    },
    "monthlyStats": {
      "2024-01": 5,
      "2023-12": 3,
      "2023-11": 2
    }
  }
}
```

### 4. 导出项目数据

**接口**: `POST /projects/export`

**权限**: 教师、管理员

**请求体**:
```json
{
  "format": "excel",  // 导出格式：excel、csv
  "filters": {
    "type": "科研",
    "status": "approved"
  }
}
```

## 数据模型

### Project（项目）

```json
{
  "id": 1,
  "title": "项目标题",
  "description": "项目描述",
  "type": "科研|竞赛",
  "status": "draft|pending|approved|rejected",
  "studentId": 3,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

### ProjectMember（项目成员）

```json
{
  "id": 1,
  "projectId": 1,
  "name": "成员姓名",
  "studentNumber": "学号",
  "role": "角色"
}
```

### ProjectFile（项目文件）

```json
{
  "id": 1,
  "projectId": 1,
  "fileName": "文件名",
  "fileUrl": "文件URL",
  "uploadTime": "2024-01-15T10:30:00Z"
}
```

### ProjectReview（项目审核）

```json
{
  "id": 1,
  "projectId": 1,
  "reviewerId": 2,
  "status": "approved|rejected",
  "comments": "审核意见",
  "reviewTime": "2024-01-15T11:00:00Z"
}
```

## 状态说明

### 项目状态

- **draft**: 草稿状态，学生可以编辑
- **pending**: 待审核状态，等待教师或管理员审核
- **approved**: 已通过审核
- **rejected**: 审核被拒绝

### 项目类型

- **科研**: 科研项目
- **竞赛**: 竞赛项目

## 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 注意事项

1. 学生只能操作自己创建的项目
2. 项目一旦提交审核，学生就不能再修改
3. 只有教师和管理员可以审核项目
4. 项目删除会同时删除相关的成员、文件和审核记录
5. 文件上传需要单独的文件上传接口（本文档未包含） 