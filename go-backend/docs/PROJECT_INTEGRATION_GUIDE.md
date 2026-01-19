# 云梦高校学生科研与竞赛项目管理系统 - 后端整合指南

## 📁 项目结构概览

```
go-backend/
├── 📄 核心文件
│   ├── main.go                    # 主程序入口
│   ├── go.mod                     # Go模块依赖
│   ├── go.sum                     # 依赖校验文件
│   └── README.md                  # 项目说明文档
│
├── 🗂️ 核心目录
│   ├── config/                    # 配置文件
│   ├── controllers/               # 控制器层
│   ├── models/                    # 数据模型
│   ├── services/                  # 业务逻辑层
│   ├── middlewares/               # 中间件
│   ├── routes/                    # 路由配置
│   ├── utils/                     # 工具函数
│   └── docs/                      # API文档
│
├── 🗄️ 数据库相关
│   ├── sql/                       # SQL脚本目录
│   └── uploads/                   # 文件上传目录
│
└── 📋 文档和脚本
    ├── 迁移脚本/                  # 数据库迁移
    ├── 启动脚本/                  # 服务启动
    ├── 测试脚本/                  # 功能测试
    └── 文档/                      # 项目文档
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 检查Go版本 (需要1.16+)
go version

# 检查MySQL版本 (需要5.7+ 或 8.0+)
mysql --version
```

### 2. 数据库初始化
```bash
# Windows
init_database.bat

# Linux/Mac
chmod +x init_database.sh
./init_database.sh
```

### 3. 启动服务
```bash
# Windows
start.bat

# Linux/Mac
chmod +x start.sh
./start.sh
```

## 📋 文件分类索引

### 🔧 核心应用文件

| 文件 | 描述 | 重要性 |
|------|------|--------|
| `main.go` | 应用程序入口点 | ⭐⭐⭐⭐⭐ |
| `go.mod` | Go模块依赖管理 | ⭐⭐⭐⭐⭐ |
| `go.sum` | 依赖校验文件 | ⭐⭐⭐⭐ |
| `README.md` | 项目主要说明文档 | ⭐⭐⭐⭐ |

### 🗄️ 数据库迁移脚本

| 文件 | 描述 | 用途 |
|------|------|------|
| `sql/init_users.sql` | 初始数据库结构 | 首次安装 |
| `sql/add_teacher_id_simple.sql` | TeacherID字段迁移 | 解决500错误 |
| `sql/add_teacher_id_to_projects_mysql8.sql` | MySQL 8.0优化迁移 | 高级功能 |

### 🚀 启动和部署脚本

| 文件 | 平台 | 描述 |
|------|------|------|
| `start.bat` | Windows | 启动服务 |
| `start.sh` | Linux/Mac | 启动服务 |
| `quick_start.bat` | Windows | 快速启动 |
| `quick_start.sh` | Linux/Mac | 快速启动 |

### 🔄 数据库迁移脚本

| 文件 | 平台 | 描述 |
|------|------|------|
| `run_simple_migration.bat` | Windows | 简化版迁移 |
| `run_teacher_id_migration_mysql8.bat` | Windows | MySQL 8.0迁移 |
| `run_teacher_id_migration_mysql8.sh` | Linux/Mac | MySQL 8.0迁移 |

### 🧪 测试和验证脚本

| 文件 | 描述 | 用途 |
|------|------|------|
| `test_new_apis.bat` | API测试脚本 | 验证新功能 |
| `test_refactored_system.bat` | 系统重构测试 | 验证重构 |
| `check_database.bat` | 数据库检查 | 验证数据库状态 |
| `simple_check.bat` | 简单检查 | 基础验证 |

### 📚 项目文档

| 文件 | 描述 | 内容 |
|------|------|------|
| `README.md` | 项目主要说明 | 项目介绍、安装、使用 |
| `API_DOCUMENTATION.md` | API文档 | 完整的API接口说明 |
| `PROJECT_API_DOCUMENTATION.md` | 项目API文档 | 项目相关API详细说明 |
| `ADMIN_BACKEND_SUMMARY.md` | 管理员后端总结 | 管理员功能说明 |
| `PROJECT_MODULE_SUMMARY.md` | 项目模块总结 | 项目功能模块说明 |

### 🔧 迁移和修复文档

| 文件 | 描述 | 内容 |
|------|------|------|
| `TEACHER_ID_MIGRATION_FIX.md` | TeacherID迁移修复 | 解决500错误的方案 |
| `MYSQL8_MIGRATION_GUIDE.md` | MySQL 8.0迁移指南 | MySQL 8.0优化方案 |
| `TROUBLESHOOTING_GUIDE.md` | 故障排除指南 | 常见问题解决方案 |
| `FINAL_MIGRATION_GUIDE.md` | 最终迁移指南 | 完整的迁移流程 |
| `MANUAL_MIGRATION_GUIDE.md` | 手动迁移指南 | 手动操作步骤 |

### 📊 项目状态和更新文档

| 文件 | 描述 | 内容 |
|------|------|------|
| `API_SUMMARY.md` | API总结 | API功能概览 |
| `UPDATED_API_WITH_NEW_FEATURES.md` | 新功能API更新 | 新增功能说明 |
| `UPDATED_PROJECT_API.md` | 项目API更新 | 项目API更新说明 |
| `API_PATH_FIX_SUMMARY.md` | API路径修复总结 | 路径修复说明 |
| `REFACTOR_README.md` | 重构说明 | 代码重构说明 |
| `STARTUP_SCRIPTS_UPDATE_SUMMARY.md` | 启动脚本更新 | 脚本更新说明 |
| `MIGRATION_ISSUE_SUMMARY.md` | 迁移问题总结 | 迁移过程中的问题 |

## 🎯 核心功能模块

### 1. 用户管理模块
- **文件**: `controllers/user_controller.go`, `models/user.go`
- **功能**: 用户注册、登录、权限管理
- **API**: `/api/users/*`

### 2. 项目管理模块
- **文件**: `controllers/project_controller.go`, `models/project.go`
- **功能**: 项目创建、编辑、审核、查询
- **API**: `/api/projects/*`

### 3. 文件管理模块
- **文件**: `controllers/file_controller.go`
- **功能**: 文件上传、下载、管理
- **API**: `/api/files/*`

### 4. 管理员模块
- **文件**: `controllers/admin_controller.go`
- **功能**: 系统管理、数据统计、用户管理
- **API**: `/api/admin/*`

## 🔧 配置说明

### 数据库配置
- **文件**: `config/database.go`
- **默认配置**:
  - 主机: localhost
  - 端口: 3306
  - 数据库: cloud_dream_system
  - 用户: root
  - 密码: 123456

### 服务器配置
- **端口**: 8080
- **模式**: 开发模式
- **日志**: 控制台输出

## 🚨 重要注意事项

### 1. 数据库迁移
- 首次使用需要运行 `init_database.bat` 或 `init_database.sh`
- 如果遇到TeacherID相关错误，运行 `run_simple_migration.bat`
- 建议在迁移前备份数据库

### 2. 权限配置
- 确保MySQL用户有足够权限
- 检查防火墙设置
- 验证端口8080是否可用

### 3. 依赖管理
- 使用 `go mod tidy` 整理依赖
- 使用 `go mod download` 下载依赖
- 检查 `go.sum` 文件完整性

## 🔍 故障排除

### 常见问题
1. **500错误**: 运行TeacherID迁移脚本
2. **连接错误**: 检查MySQL服务状态
3. **权限错误**: 验证数据库用户权限
4. **端口冲突**: 修改配置文件中的端口

### 调试步骤
1. 检查日志输出
2. 验证数据库连接
3. 测试API接口
4. 查看错误信息

## 📞 技术支持

### 文档资源
- `TROUBLESHOOTING_GUIDE.md` - 故障排除
- `MYSQL8_MIGRATION_GUIDE.md` - MySQL 8.0指南
- `API_DOCUMENTATION.md` - API文档

### 联系信息
- 项目文档: 查看各README文件
- 技术问题: 参考故障排除指南
- 功能问题: 查看API文档

## 🎉 项目特色

### 技术栈
- **后端**: Go + Gin框架
- **数据库**: MySQL 8.0
- **认证**: JWT Token
- **文件处理**: 多文件上传支持

### 功能特色
- 完整的用户权限管理
- 项目全生命周期管理
- 文件上传和管理
- 实时数据统计
- 多角色支持（学生、教师、管理员）

### 部署优势
- 跨平台支持
- 一键启动脚本
- 自动化迁移
- 详细文档支持

---

**最后更新**: 2024年
**版本**: 1.0.0
**维护者**: 系统管理员 