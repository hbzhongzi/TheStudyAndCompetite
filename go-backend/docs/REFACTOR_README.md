此文件内容已归档并移至 `go-backend/backups/docs_archive/REFACTOR_README.md`，可安全删除原件以释放空间。

#### 1.1 projects表新增字段
- `teacher_id BIGINT NOT NULL`：指导老师ID（必填）
- `submitted_at DATETIME NULL`：项目提交时间
- 新增外键约束：`fk_projects_teacher`

#### 1.2 新增student_teacher表
- 管理学生和教师的绑定关系
- 支持一对多的师生关系
- 包含唯一性约束防止重复绑定

### 2. 代码重构

#### 2.1 模型层 (models/project.go)
- 新增 `TeacherID` 和 `SubmittedAt` 字段
- 新增 `StudentTeacher` 模型
- 更新所有请求/响应结构体

#### 2.2 服务层 (services/project_service.go)
- 新增师生绑定相关方法
- 更新项目查询方法以包含教师信息
- 新增项目提交方法（设置submitted_at时间）

#### 2.3 控制器层 (controllers/project_controller.go)
- 新增教师管理相关控制器方法
- 更新项目提交逻辑

#### 2.4 路由层 (routes/routes.go)
- 新增教师管理路由组
- 更新项目路由

### 3. 新增API接口

#### 3.1 教师管理接口
- `GET /api/teachers`：获取教师列表
- `POST /api/teachers/bind`：绑定学生和教师
- `GET /api/teachers/students/:studentId`：获取学生的指导教师
- `DELETE /api/teachers/students/:studentId/teachers/:teacherId`：解绑学生和教师

#### 3.2 项目提交接口
- `POST /api/projects/submit/:id`：提交项目审核

## 快速开始

### 1. 执行数据库迁移

#### Windows
```bash
migrate_projects.bat
```

#### Linux/Mac
```bash
chmod +x migrate_projects.sh
./migrate_projects.sh
```

### 2. 启动服务
```bash
go run main.go
```

### 3. 测试系统
```bash
# Windows
test_refactored_system.bat
```

## 文件结构

```
go-backend/
├── models/
│   └── project.go              # 更新：新增字段和模型
├── services/
│   └── project_service.go      # 更新：新增师生管理方法
├── controllers/
│   └── project_controller.go   # 更新：新增控制器方法
├── routes/
│   └── routes.go               # 更新：新增路由
├── sql/
│   └── migrate_projects.sql    # 新增：数据库迁移脚本
├── migrate_projects.bat        # 新增：Windows迁移脚本
├── migrate_projects.sh         # 新增：Linux迁移脚本
├── test_refactored_system.bat  # 新增：系统测试脚本
├── UPDATED_PROJECT_API.md      # 更新：API文档
├── PROJECT_MODULE_SUMMARY.md   # 新增：重构总结
└── REFACTOR_README.md          # 新增：本文件
```

## 主要功能

### 1. 师生关系管理
- 学生和教师的绑定/解绑
- 教师列表查询
- 学生指导教师查询

### 2. 项目审核流程
- 项目必须指定指导教师
- 支持项目提交时间记录
- 完整的审核历史记录

### 3. 权限控制
- 学生：只能修改自己的草稿项目
- 教师：可以审核项目和管理师生关系
- 管理员：拥有所有权限

## 测试用例

### 1. 师生绑定测试
```bash
# 获取教师列表
curl -X GET http://localhost:8080/api/teachers \
  -H "Authorization: Bearer <token>"

# 绑定学生和教师
curl -X POST http://localhost:8080/api/teachers/bind \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"studentId": 3, "teacherId": 2}'
```

### 2. 项目创建测试
```bash
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

### 3. 项目提交测试
```bash
# 提交项目审核
curl -X POST http://localhost:8080/api/projects/submit/1 \
  -H "Authorization: Bearer <token>"
```

## 注意事项

### 1. 数据库迁移
- 执行迁移前请备份数据库
- 确保MySQL服务正在运行
- 检查数据库连接参数

### 2. 代码部署
- 更新所有Go依赖
- 重启服务以加载新代码
- 验证API接口正常工作

### 3. 数据验证
- 检查现有项目是否已分配指导教师
- 验证师生绑定关系
- 确认审核记录完整性

## 常见问题

### Q1: 迁移脚本执行失败
**A**: 检查数据库连接参数，确保MySQL服务正在运行，用户有足够权限。

### Q2: 项目创建时提示teacherId必填
**A**: 这是新功能要求，所有项目必须指定指导教师。可以先绑定师生关系，再创建项目。

### Q3: API接口返回401错误
**A**: 检查JWT token是否有效，确保用户已正确登录。

## 联系支持

如有问题，请查看以下文档：
- `UPDATED_PROJECT_API.md`：完整的API文档
- `PROJECT_MODULE_SUMMARY.md`：详细的重构总结
- `sql/migrate_projects.sql`：数据库迁移脚本

## 版本信息

- **重构版本**: v2.0
- **数据库版本**: v2.0
- **API版本**: v2.0
- **更新日期**: 2024年1月 