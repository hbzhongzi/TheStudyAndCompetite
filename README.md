# 云梦高校学生科研与竞赛项目管理系统

## 项目简介

云梦高校学生科研与竞赛项目管理系统是一个基于Web的综合管理平台，旨在为高校学生、教师和管理员提供科研项目管理、竞赛管理、用户管理等功能的完整解决方案。

## 系统特色

- 🎯 **多角色管理**: 支持管理员、教师、学生三种角色
- 📚 **项目管理**: 完整的科研项目生命周期管理
- 🏆 **竞赛管理**: 竞赛报名、评审、结果管理一体化
- 👥 **用户管理**: 完善的用户权限和角色管理
- 📊 **数据统计**: 丰富的统计分析和报表功能
- 📱 **响应式设计**: 支持PC端和移动端访问
- 🔒 **安全可靠**: JWT认证、权限控制、数据加密

## 技术架构

### 前端技术栈
- **框架**: Vue 3 + Composition API
- **UI库**: Element Plus
- **构建工具**: Vite
- **路由**: Vue Router 4
- **HTTP客户端**: Axios

### 后端技术栈
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT

## 功能模块

### 1. 用户管理模块
- 用户注册、登录、权限管理
- 角色分配和权限控制
- 用户信息维护和状态管理

### 2. 项目管理模块
- **学生端**: 项目创建、编辑、提交、进度跟踪
- **教师端**: 项目指导、审核、质量评估、批量操作
- **管理员端**: 项目监控、审核管理、统计分析

### 3. 竞赛管理模块
- 竞赛信息发布和管理
- 学生报名和作品提交
- 教师评审和结果管理
- 竞赛数据统计和导出

### 4. 系统管理模块
- 系统配置和参数管理
- 日志记录和监控
- 数据备份和恢复
- 系统健康状态监控

## 快速开始

### 环境要求
- Node.js 16+
- Go 1.21+
- MySQL 8.0+
- Git

### 安装步骤

#### 1. 克隆项目
```bash
git clone <repository-url>
cd 云梦高校学生科研与竞赛项目管理系统的设计与实现
```

#### 2. 启动后端服务
```bash
cd go-backend
go mod tidy
go run main.go
```

#### 3. 启动前端服务
```bash
cd yunmeng-frontend
npm install
npm run dev
```

#### 4. 访问系统
- 前端: http://localhost:5173
- 后端: http://localhost:8080

## 项目结构

```
├── go-backend/                    # 后端服务
│   ├── config/                    # 配置管理
│   ├── controllers/               # 控制器
│   ├── models/                    # 数据模型
│   ├── services/                  # 业务逻辑
│   ├── routes/                    # 路由配置
│   ├── middlewares/               # 中间件
│   ├── migrations/                # 数据库迁移
│   └── docs/                      # 后端文档
├── yunmeng-frontend/              # 前端应用
│   ├── src/                       # 源代码
│   │   ├── components/            # 通用组件
│   │   ├── views/                 # 页面组件
│   │   ├── services/              # API服务
│   │   └── utils/                 # 工具函数
│   ├── public/                    # 静态资源
│   └── docs/                      # 前端文档
├── docs/                          # 项目文档
└── README.md                      # 项目说明
```

## 核心功能

### 项目管理
- ✅ 项目创建和编辑
- ✅ 项目审核和状态管理
- ✅ 进度跟踪和里程碑管理
- ✅ 文件上传和管理
- ✅ 延期申请和审批
- ✅ 质量评估和反馈

### 竞赛管理
- ✅ 竞赛信息发布
- ✅ 在线报名系统
- ✅ 作品提交和评审
- ✅ 结果公布和统计
- ✅ 数据导出和报表

### 用户管理
- ✅ 多角色权限控制
- ✅ 用户信息管理
- ✅ 导师学生绑定
- ✅ 批量操作支持
- ✅ 数据统计和分析

### 系统功能
- ✅ 通知消息系统
- ✅ 系统监控和日志
- ✅ 数据备份和恢复
- ✅ 性能优化和缓存
- ✅ 安全防护和审计

## 开发指南

### 代码规范
- 遵循Go和Vue.js官方代码规范
- 使用统一的命名约定
- 添加必要的代码注释
- 完善的错误处理机制

### 测试
```bash
# 后端测试
cd go-backend
go test ./...

# 前端测试
cd yunmeng-frontend
npm run test
```

### 部署
```bash
# 后端构建
cd go-backend
go build -o yunmeng-backend main.go

# 前端构建
cd yunmeng-frontend
npm run build
```

## 文档说明

- **[前端文档](./yunmeng-frontend/FRONTEND_DOCUMENTATION.md)**: 前端功能说明和开发指南
- **[后端文档](./go-backend/BACKEND_DOCUMENTATION.md)**: 后端API文档和部署指南
- **[教师端项目管理](./yunmeng-frontend/TEACHER_PROJECT_MANAGEMENT_ENHANCEMENT.md)**: 教师端功能完善说明

## 更新日志

### v1.2.0 (2024-01-25)
- 🆕 教师端项目管理功能增强
- 🆕 批量操作支持
- 🆕 项目质量评估系统
- 🆕 系统监控和健康检查
- 🔧 项目文件结构优化
- 📚 文档整合和优化

### v1.1.0 (2024-01-20)
- 🆕 竞赛管理模块
- 🆕 用户管理功能
- 🆕 文件上传系统
- 🔧 权限控制优化

### v1.0.0 (2024-01-15)
- 🎉 系统初始版本发布
- 🆕 基础项目管理功能
- 🆕 用户认证系统
- 🆕 基础UI界面

## 贡献指南

欢迎贡献代码和提出建议！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- **项目维护者**: 云梦高校开发团队
- **邮箱**: dev@yunmeng.edu.cn
- **项目地址**: [GitHub Repository](<repository-url>)

## 致谢

感谢所有为这个项目做出贡献的开发者和用户！

---

*最后更新: 2024-01-25*
*版本: v1.2.0* 