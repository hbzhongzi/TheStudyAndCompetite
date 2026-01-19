# 云梦高校科研竞赛管理系统 - 启动指南

## 系统要求

### 必需软件
- **Go 1.20+** - 后端开发语言
- **MySQL 8.0+** - 数据库服务器
- **Git** - 版本控制（可选）

### 推荐配置
- **操作系统**: Windows 10/11, macOS, Linux
- **内存**: 4GB+
- **磁盘空间**: 2GB+

## 快速启动

### 方法一：一键启动（推荐）

1. **双击运行** `quick_start.bat`
2. **等待自动初始化**，脚本会自动：
   - 检查环境
   - 启动MySQL服务
   - 创建数据库
   - 初始化表结构
   - 插入测试数据
   - 启动后端服务

3. **访问系统**：
   - 前端地址：http://localhost:5173
   - 后端API：http://localhost:8080/api

### 方法二：详细启动

1. **运行** `start.bat`
2. **按提示操作**，脚本会逐步检查每个组件
3. **查看详细日志**，了解每个步骤的执行情况

### 方法三：手动启动

1. **检查数据库**：运行 `check_database.bat`
2. **初始化数据库**：运行 `init_database.bat`
3. **启动服务**：运行 `go run main.go`

## 数据库配置

### 默认配置
- **数据库名**: `cloud_dream_system`
- **用户名**: `root`
- **密码**: `123456`
- **主机**: `localhost`
- **端口**: `3306`

### 修改配置
如需修改数据库配置，请编辑 `config/database.go` 文件：

```go
func NewDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "3306"),
        Username: getEnv("DB_USERNAME", "root"),
        Password: getEnv("DB_PASSWORD", "123456"),
        Database: getEnv("DB_DATABASE", "cloud_dream_system"),
        Charset:  getEnv("DB_CHARSET", "utf8mb4"),
    }
}
```

### 环境变量
也可以通过环境变量配置：

```bash
# Windows
set DB_PASSWORD=your_password
set DB_HOST=your_host

# Linux/macOS
export DB_PASSWORD=your_password
export DB_HOST=your_host
```

## 默认用户

系统预置了以下测试用户：

| 用户名 | 密码 | 角色 | 描述 |
|--------|------|------|------|
| admin | 123456 | 管理员 | 系统管理员，拥有所有权限 |
| teacher001 | 123456 | 教师 | 计算机学院教师，可以审核项目 |
| student001 | 123456 | 学生 | 计算机学院学生，可以申报项目 |

## 数据库结构

### 主要表
- **users** - 用户基础信息
- **user_profiles** - 用户详细信息
- **roles** - 角色定义
- **user_roles** - 用户角色关联
- **projects** - 项目信息
- **project_members** - 项目成员
- **project_files** - 项目附件
- **project_reviews** - 项目审核记录
- **login_logs** - 登录日志

### 示例数据
系统包含以下示例数据：
- 3个用户（管理员、教师、学生）
- 5个项目（科研和竞赛项目）
- 完整的项目成员和附件信息
- 项目审核记录

## 故障排除

### 常见问题

#### 1. Go未安装
```
错误: Go未安装，请先安装Go
```
**解决方案**：
- 下载并安装Go：https://golang.org/dl/
- 确保Go已添加到PATH环境变量

#### 2. MySQL未安装
```
错误: MySQL未安装或未添加到PATH环境变量
```
**解决方案**：
- 下载并安装MySQL：https://dev.mysql.com/downloads/
- 确保mysql命令可用
- 启动MySQL服务：`net start mysql`

#### 3. 数据库连接失败
```
错误: 无法连接到MySQL服务器
```
**解决方案**：
- 检查MySQL服务是否启动
- 验证用户名和密码是否正确
- 检查防火墙设置
- 确认MySQL端口（3306）未被占用

#### 4. 数据库初始化失败
```
警告: 数据库初始化失败
```
**解决方案**：
- 检查MySQL权限
- 手动执行SQL脚本：`mysql -u root -p123456 < sql/init_complete.sql`
- 查看MySQL错误日志

#### 5. 端口被占用
```
错误: 端口8080已被占用
```
**解决方案**：
- 修改配置文件中的端口号
- 关闭占用端口的程序
- 使用不同的端口

### 日志查看

#### 后端日志
后端服务启动时会显示详细日志，包括：
- 数据库连接状态
- API路由注册
- 服务启动信息

#### 数据库日志
查看MySQL错误日志：
```bash
# Windows
type "C:\ProgramData\MySQL\MySQL Server 8.0\Data\hostname.err"

# Linux
tail -f /var/log/mysql/error.log
```

## 开发模式

### 热重载
使用 `air` 工具实现热重载：

1. **安装air**：
```bash
go install github.com/cosmtrek/air@latest
```

2. **创建配置文件** `.air.toml`：
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

3. **启动开发模式**：
```bash
air
```

### 调试模式
设置环境变量启用调试模式：
```bash
set DEBUG=true
go run main.go
```

## 部署

### 生产环境
1. **编译二进制文件**：
```bash
go build -o yunmeng-backend main.go
```

2. **配置环境变量**：
```bash
set DB_PASSWORD=production_password
set DB_HOST=production_host
```

3. **启动服务**：
```bash
yunmeng-backend
```

### Docker部署
1. **构建镜像**：
```bash
docker build -t yunmeng-backend .
```

2. **运行容器**：
```bash
docker run -p 8080:8080 yunmeng-backend
```

## 维护

### 数据库备份
```bash
mysqldump -u root -p123456 cloud_dream_system > backup.sql
```

### 数据库恢复
```bash
mysql -u root -p123456 cloud_dream_system < backup.sql
```

### 日志清理
定期清理日志文件以节省磁盘空间。

## 技术支持

如遇到问题，请：
1. 查看本文档的故障排除部分
2. 检查系统日志
3. 联系技术支持团队

---

**注意**：首次启动时，系统会自动创建数据库和表结构，请确保MySQL服务正常运行。 