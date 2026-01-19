# 云梦高校科研竞赛管理系统 - 更新版API文档（含新增功能）
此文件内容已归档并移至 `go-backend/backups/docs_archive/UPDATED_API_WITH_NEW_FEATURES.md`，可安全删除原件以释放空间。
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 20)
- `search`: 搜索关键词 (可选)
- `type`: 项目类型筛选 (可选: 科研, 竞赛)
- `status`: 项目状态筛选 (可选: draft, pending, approved, rejected)
- `studentId`: 学生ID筛选 (可选)
- `sortBy`: 排序字段 (可选)
- `sortOrder`: 排序方向 (可选: asc, desc)

**响应示例**:
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

### 2. 教师管理功能

#### 2.1 获取教师列表（支持院系筛选）
```http
GET /api/teachers/filter
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 20)
- `search`: 搜索关键词 (可选)
- `department`: 院系筛选 (可选)
- `sortBy`: 排序字段 (可选)
- `sortOrder`: 排序方向 (可选: asc, desc)

**响应示例**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 2,
        "username": "teacher001",
        "realName": "李老师",
        "email": "li.teacher@yunmeng.edu.cn",
        "phone": "13800138001",
        "department": "计算机学院",
        "bio": "计算机科学教授"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 20
  }
}
```

#### 2.2 学生绑定教师（学生端接口）
```http
POST /api/students/bind-teacher
Authorization: Bearer <token>
Content-Type: application/json

{
  "teacherId": 2
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "绑定成功",
  "data": {
    "id": 1,
    "studentId": 3,
    "teacherId": 2,
    "bindTime": "2024-01-15T10:30:00Z",
    "teacher": {
      "id": 2,
      "username": "teacher001",
      "realName": "李老师",
      "email": "li.teacher@yunmeng.edu.cn",
      "phone": "13800138001",
      "department": "计算机学院",
      "bio": "计算机科学教授"
    }
  }
}
```

### 3. 项目更新验证

#### 3.1 更新项目（带验证）
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

**验证规则**:
- 只有项目创建者可以修改
- 只有草稿状态的项目可以修改
- 项目提交后禁止修改指导老师

**错误响应示例**:
```json
{
  "code": 400,
  "message": "项目已提交，无法修改指导老师"
}
```

---

## 权限控制

### 教师权限
- 查看自己的指导项目列表
- 获取教师列表（支持筛选）
- 管理师生绑定关系
- 审核项目

### 学生权限
- 绑定指导教师
- 创建和修改自己的项目（仅草稿状态）
- 提交项目审核
- 查看项目详情

### 管理员权限
- 所有教师权限
- 用户管理
- 系统管理

---

## 数据库约束

### 项目状态约束
- `draft`: 草稿状态，可以修改所有字段
- `pending`: 待审核状态，不能修改指导老师
- `approved`: 已批准状态，不能修改
- `rejected`: 已拒绝状态，不能修改

### 师生绑定约束
- 一个学生可以绑定多个教师
- 同一学生和教师只能绑定一次（唯一约束）
- 删除用户时级联删除绑定关系

---

## 使用示例

### 1. 教师查看指导项目
```bash
curl -X GET "http://localhost:8080/api/teachers/projects?status=pending" \
  -H "Authorization: Bearer <token>"
```

### 2. 学生按院系筛选教师
```bash
curl -X GET "http://localhost:8080/api/teachers/filter?department=计算机学院" \
  -H "Authorization: Bearer <token>"
```

### 3. 学生绑定教师
```bash
curl -X POST "http://localhost:8080/api/students/bind-teacher" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"teacherId": 2}'
```

### 4. 更新项目（带验证）
```bash
curl -X PUT "http://localhost:8080/api/projects/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "新标题", "teacherId": 3}'
```

---

## 错误处理

### 常见错误码
- `400`: 请求参数错误或业务逻辑错误
- `401`: 未认证或认证失败
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

### 业务逻辑错误
- "项目已提交，无法修改指导老师"
- "该学生已经绑定此教师"
- "教师不存在或不是教师角色"
- "只有草稿状态的项目可以修改"

---

## 技术实现

### 新增模型
- `TeacherProjectQueryParams`: 教师项目查询参数
- `TeacherProjectResponse`: 教师项目响应
- `TeacherQueryParams`: 教师查询参数
- `StudentBindTeacherRequest`: 学生绑定教师请求
- `StudentBindTeacherResponse`: 学生绑定教师响应

### 新增服务方法
- `GetTeacherProjects()`: 获取教师指导项目
- `GetTeacherListWithFilter()`: 获取教师列表（支持筛选）
- `BindStudentToTeacher()`: 学生绑定教师
- `ValidateProjectUpdate()`: 验证项目更新权限

### 新增控制器方法
- `GetTeacherProjects()`: 教师指导项目接口
- `GetTeacherListWithFilter()`: 教师列表筛选接口
- `BindStudentToTeacher()`: 学生绑定教师接口
- `UpdateProjectWithValidation()`: 带验证的项目更新接口

### 路由更新
- 新增 `/api/teachers/projects` 路由
- 新增 `/api/teachers/filter` 路由
- 新增 `/api/students/bind-teacher` 路由
- 更新 `/api/projects/{id}` 路由使用验证版本 