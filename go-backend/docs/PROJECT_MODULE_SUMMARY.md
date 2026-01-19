# 云梦高校科研竞赛管理系统 - 项目模块重构总结

## 重构概述

根据数据库增量更新需求，对项目模块进行了全面重构，主要包括数据库结构更新、模型调整、服务层扩展和API接口优化。

## 数据库更新内容

### 1. projects表结构更新
- **新增字段**：
  - `teacher_id BIGINT NOT NULL`：指导老师ID，必填字段
  - `submitted_at DATETIME NULL`：项目提交时间
- **新增约束**：
  - `fk_projects_teacher`：外键约束，关联users表
- **新增索引**：
  - `idx_projects_teacher_id`：优化教师查询性能

### 2. 新增student_teacher表
- **用途**：管理学生和教师的绑定关系
- **字段**：
  - `id BIGINT PRIMARY KEY AUTO_INCREMENT`
  - `student_id BIGINT NOT NULL`
  - `teacher_id BIGINT NOT NULL`
  - `bind_time DATETIME DEFAULT CURRENT_TIMESTAMP`
- **约束**：
  - `UNIQUE (student_id, teacher_id)`：防止重复绑定
  - 外键约束关联users表
- **索引**：
  - `idx_student_teacher_student_id`
  - `idx_student_teacher_teacher_id`

### 3. project_reviews表（已存在）
- 保存项目审核历史记录
- 支持多次审核记录

## 代码重构内容

### 1. 模型层更新 (models/project.go)

#### 新增字段和关联
```go
type Project struct {
    // ... 原有字段
    TeacherID   uint       `gorm:"not null;column:teacher_id" json:"teacherId"`
    SubmittedAt *time.Time `gorm:"column:submitted_at" json:"submittedAt"`
    
    // 新增关联
    Teacher *User `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// 新增学生教师绑定模型
type StudentTeacher struct {
    ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
    StudentID uint      `gorm:"not null;column:student_id" json:"studentId"`
    TeacherID uint      `gorm:"not null;column:teacher_id" json:"teacherId"`
    BindTime  time.Time `gorm:"column:bind_time;autoCreateTime" json:"bindTime"`
    
    Student *User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
    Teacher *User `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}
```

#### 新增请求/响应结构体
- `StudentTeacherBindRequest`：绑定请求
- `StudentTeacherBindResponse`：绑定响应
- `TeacherListResponse`：教师列表响应
- 更新现有结构体以包含教师信息

### 2. 服务层扩展 (services/project_service.go)

#### 新增方法
1. **BindStudentTeacher**：绑定学生和教师
2. **GetTeacherList**：获取教师列表
3. **GetStudentTeachers**：获取学生的指导教师
4. **UnbindStudentTeacher**：解绑学生和教师
5. **SubmitProject**：提交项目审核（设置submitted_at时间）

#### 更新现有方法
1. **GetProjectList**：添加教师信息预加载和响应
2. **GetProjectByID**：添加教师信息预加载和响应
3. **CreateProject**：支持teacher_id字段
4. **UpdateProject**：支持teacher_id更新

### 3. 控制器层扩展 (controllers/project_controller.go)

#### 新增控制器方法
1. **BindStudentTeacher**：处理绑定请求
2. **GetTeacherList**：获取教师列表
3. **GetStudentTeachers**：获取学生指导教师
4. **UnbindStudentTeacher**：处理解绑请求

#### 更新现有方法
1. **SubmitProject**：使用新的SubmitProject服务方法

### 4. 路由层更新 (routes/routes.go)

#### 新增路由组
```go
// 教师管理路由
teachers := auth.Group("/teachers")
teachers.Use(middlewares.RoleMiddleware("teacher", "admin"))
{
    teachers.GET("", projectController.GetTeacherList)
    teachers.POST("/bind", projectController.BindStudentTeacher)
    teachers.GET("/students/:studentId", projectController.GetStudentTeachers)
    teachers.DELETE("/students/:studentId/teachers/:teacherId", projectController.UnbindStudentTeacher)
}
```

#### 更新项目路由
- 添加项目提交路由：`POST /api/projects/submit/:id`

## 新增API接口

### 1. 教师管理接口
- `GET /api/teachers`：获取教师列表
- `POST /api/teachers/bind`：绑定学生和教师
- `GET /api/teachers/students/:studentId`：获取学生的指导教师
- `DELETE /api/teachers/students/:studentId/teachers/:teacherId`：解绑学生和教师

### 2. 项目提交接口
- `POST /api/projects/submit/:id`：提交项目审核

### 3. 更新现有接口
- 所有项目相关接口现在都包含教师信息
- 项目创建和更新支持teacher_id字段

## 业务逻辑优化

### 1. 项目状态管理
- **draft**：草稿状态，可以修改
- **pending**：待审核状态，已提交等待审核
- **approved**：已通过审核
- **rejected**：已拒绝

### 2. 权限控制
- **学生**：只能修改自己的草稿项目，必须指定指导教师
- **教师**：可以查看所有项目，审核项目，管理师生关系
- **管理员**：拥有所有权限

### 3. 数据完整性
- 项目必须绑定指导教师
- 学生教师绑定关系唯一性约束
- 提交时间自动记录

## 迁移脚本

### 1. 数据库迁移脚本
- `sql/migrate_projects.sql`：完整的迁移脚本
- 包含表结构更新、索引创建、示例数据插入

### 2. 执行脚本
- `migrate_projects.bat`：Windows批处理脚本
- `migrate_projects.sh`：Linux/Mac shell脚本

## 测试建议

### 1. 数据库迁移测试
```bash
# 执行迁移脚本
./migrate_projects.sh

# 验证表结构
DESCRIBE projects;
DESCRIBE student_teacher;
```

### 2. API接口测试
```bash
# 获取教师列表
curl -X GET http://localhost:8080/api/teachers \
  -H "Authorization: Bearer <token>"

# 绑定学生教师
curl -X POST http://localhost:8080/api/teachers/bind \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"studentId": 3, "teacherId": 2}'

# 创建项目（包含teacherId）
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试项目",
    "description": "项目描述",
    "type": "科研",
    "status": "draft",
    "teacherId": 2,
    "members": [],
    "attachments": []
  }'
```

## 部署注意事项

### 1. 数据库迁移
- 执行迁移脚本前请备份数据库
- 确保数据库连接参数正确
- 检查外键约束是否正常创建

### 2. 代码部署
- 更新所有相关依赖
- 重启服务以加载新的模型和路由
- 验证API接口正常工作

### 3. 数据验证
- 检查现有项目是否已分配指导教师
- 验证师生绑定关系是否正确
- 确认审核记录完整性

## 总结

本次重构成功实现了以下目标：

1. **数据库结构优化**：添加了必要的字段和表，支持师生关系管理
2. **业务逻辑完善**：实现了完整的项目审核流程和师生绑定功能
3. **API接口扩展**：新增了教师管理相关接口，优化了现有接口
4. **权限控制增强**：细化了不同角色的权限范围
5. **数据完整性保障**：通过约束和验证确保数据一致性

重构后的系统更加完善，支持更复杂的业务场景，为后续功能扩展奠定了良好基础。 