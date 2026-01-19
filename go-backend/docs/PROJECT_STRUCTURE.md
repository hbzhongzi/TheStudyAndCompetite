# 项目结构说明 
 
## 目录结构 
 
``` 
go-backend/ 
├── README.md                    # 项目主文档 
├── main.go                      # 应用入口 
├── go.mod                       # Go模块定义 
├── go.sum                       # 依赖校验 
├── config/                      # 配置文件 
│   └── database.go              # 数据库配置 
├── controllers/                 # 控制器 
│   ├── login.go                 # 登录控制器 
│   ├── auth_controller.go       # 认证控制器 
│   ├── user_controller.go       # 用户管理 
│   ├── project_controller.go    # 项目管理 
│   └── competition_controller.go # 竞赛管理 
├── models/                      # 数据模型 
│   ├── user.go                  # 用户模型 
│   ├── project.go               # 项目模型 
│   └── competition.go           # 竞赛模型 
├── middlewares/                 # 中间件 
│   └── auth.go                  # 认证中间件 
├── routes/                      # 路由配置 
│   └── routes.go                # 路由定义 
├── services/                    # 业务逻辑 
├── utils/                       # 工具函数 
│   ├── jwt.go                   # JWT工具 
│   ├── password.go              # 密码工具 
│   └── response.go              # 响应工具 
├── sql/                         # SQL脚本 
├── scripts/                     # 脚本文件 
├── debug/                       # 调试工具 
│   ├── debug_token.go           # Token调试工具 
│   └── go.mod                   # 调试模块定义 
└── docs/                        # 文档目录 
    ├── PROJECT_STRUCTURE.md     # 项目结构说明 
    ├── TOKEN_EXPIRATION_FIX.md  # Token过期修复文档 
    ├── TOKEN_FIX_COMPLETE.md    # Token修复完整文档 
    ├── ADMIN_PERMISSION_FIX.md  # 管理员权限修复文档 
    ├── ADMIN_QUICK_FIX.md       # 管理员快速修复文档 
    ├── QUICK_START_README.md    # 快速开始文档 
    ├── START_EXAMPLES.md        # 启动示例文档 
    └── cleanup_and_organize.sh  # 清理脚本 
``` 
