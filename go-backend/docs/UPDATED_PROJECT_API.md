# 云梦高校科研竞赛管理系统 - 项目模块API文档（更新版）

## 概述

本文档描述了项目模块的API接口，包含最新的数据库结构更新和新增功能。

## 数据库更新内容

### 1. projects表更新
- 新增 `teacher_id` 字段：指导老师ID（必填）
- 新增 `submitted_at` 字段：项目提交时间
- 新增外键约束：`fk_projects_teacher`

### 2. 新增 student_teacher 表
- 用于管理学生和教师的绑定关系
- 支持一对多的师生关系

### 3. project_reviews 表（已存在）
- 保存项目审核历史记录

## API接口

### 1. 项目基础操作

#### 1.1 创建项目
```http
POST /api/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "项目标题",
  "description": "项目描述",
  "type": "科研",
  "status": "draft",
  "teacherId": 2,
  "members": [
    {
      "name": "成员姓名",
      "studentNumber": "学号",
      "role": "角色"
    }
  ],
  "attachments": [
    {
      "fileName": "文件名",
      "fileUrl": "文件URL"
    }
  ]
}
```

**响应：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "projectId": 1
  }
}
```

#### 1.2 更新项目
```http
PUT /api/projects/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "type": "竞赛",
  "teacherId": 3
}
```

#### 1.3 提交项目审核
```http
POST /api/projects/submit/{id}
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 200,
  "message": "项目提交成功，等待审核"
}
```

#### 1.4 获取项目详情
```http
GET /api/projects/{id}
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 200,
  "message": "获取项目详情成功",
  "data": {
    "id": 1,
    "title": "项目标题",
    "description": "项目描述",
    "type": "科研",
    "status": "pending",
    "submittedAt": "2024-01-15T10:30:00Z",
    "createdAt": "2024-01-10T09:00:00Z",
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
    "teacher": {
      "id": 2,
      "username": "teacher001",
      "realName": "李老师",
      "email": "li.teacher@yunmeng.edu.cn",
      "phone": "13800138001",
      "department": "计算机学院"
    },
    "members": [...],
    "files": [...],
    "reviews": [...]
  }
}
```

### 2. 教师管理功能

#### 2.1 获取教师列表
```http
GET /api/teachers
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 200,
  "message": "获取教师列表成功",
  "data": [
    {
      "id": 2,
      "username": "teacher001",
      "realName": "李老师",
      "email": "li.teacher@yunmeng.edu.cn",
      "phone": "13800138001",
      "department": "计算机学院",
      "bio": "计算机学院教师，专注于人工智能研究"
    }
  ]
}
```

#### 2.2 绑定学生和教师
```http
POST /api/teachers/bind
Authorization: Bearer <token>
Content-Type: application/json

{
  "studentId": 3,
  "teacherId": 2
}
```

**响应：**
```json
{
  "code": 200,
  "message": "绑定成功",
  "data": {
    "id": 1,
    "studentId": 3,
    "teacherId": 2,
    "bindTime": "2024-01-15T11:00:00Z",
    "student": {
      "id": 3,
      "username": "student001",
      "realName": "张三"
    },
    "teacher": {
      "id": 2,
      "username": "teacher001",
      "realName": "李老师"
    }
  }
}
```

#### 2.3 获取学生的指导教师
```http
GET /api/teachers/students/{studentId}
Authorization: Bearer <token>
```

#### 2.4 解绑学生和教师
```http
DELETE /api/teachers/students/{studentId}/teachers/{teacherId}
Authorization: Bearer <token>
```

### 3. 项目审核功能

#### 3.1 审核项目
```http
PUT /api/projects/{id}/review
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "approved",
  "comments": "项目方案完整，技术路线清晰，同意立项"
}
```

#### 3.2 获取审核记录
```http
GET /api/projects/{id}/reviews
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "reviewer": "李老师",
      "status": "approved",
      "comments": "项目方案完整，技术路线清晰，同意立项",
      "reviewTime": "2024-01-15T14:30:00Z"
    }
  ]
}
```

### 4. 项目查询功能

#### 4.1 获取项目列表（教师/管理员）
```http
GET /api/projects?page=1&size=20&type=科研&status=pending&search=关键词
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "项目标题",
        "description": "项目描述",
        "type": "科研",
        "status": "pending",
        "studentName": "张三",
        "studentId": "2021001",
        "teacherName": "李老师",
        "teacherId": 2,
        "submittedAt": "2024-01-15T10:30:00Z",
        "createdAt": "2024-01-10T09:00:00Z",
        "updatedAt": "2024-01-15T10:30:00Z",
        "memberCount": 3,
        "fileCount": 2,
        "reviewCount": 1
      }
    ],
    "total": 1
  }
}
```

#### 4.2 获取我的项目（学生）
```http
GET /api/projects/my?status=pending
Authorization: Bearer <token>
```

## 数据模型

### Project 模型
此文件内容已归档并移至 `go-backend/backups/docs_archive/UPDATED_PROJECT_API.md`，可安全删除原件以释放空间。
    TeacherID uint      `json:"teacherId"`
    BindTime  time.Time `json:"bindTime"`
    
    Student *User `json:"student,omitempty"`
    Teacher *User `json:"teacher,omitempty"`
}
```

## 状态说明

### 项目状态 (status)
- `draft`: 草稿状态，可以修改
- `pending`: 待审核状态，已提交等待审核
- `approved`: 已通过审核
- `rejected`: 已拒绝

### 项目类型 (type)
- `科研`: 科研项目
- `竞赛`: 竞赛项目

### 审核状态 (review status)
- `approved`: 通过
- `rejected`: 拒绝

## 权限说明

1. **学生权限**：
   - 创建、修改自己的草稿项目
   - 提交项目审核
   - 查看自己的项目

2. **教师权限**：
   - 查看所有项目
   - 审核项目
   - 管理学生教师绑定关系

3. **管理员权限**：
   - 所有教师权限
   - 系统管理功能

## 迁移说明

执行数据库迁移脚本：
```bash
# Windows
migrate_projects.bat

# Linux/Mac
chmod +x migrate_projects.sh
./migrate_projects.sh
```

迁移内容包括：
1. 为projects表添加teacher_id和submitted_at字段
2. 创建student_teacher表
3. 创建相关索引
4. 插入示例数据 