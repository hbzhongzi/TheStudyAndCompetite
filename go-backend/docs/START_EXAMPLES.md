# 启动脚本使用示例

## 快速开始

### 1. 首次运行（推荐）
```bash
# Windows - 批处理
quick_start_backend.bat

# Windows - PowerShell
.\start_backend.ps1

# Linux/macOS
./quick_start_backend.sh
```

### 2. 日常开发
```bash
# Windows - 批处理
此文件内容已归档并移至 `go-backend/backups/docs_archive/START_EXAMPLES.md`，可安全删除原件以释放空间。
.\start_backend.ps1 -Mode full -Port 8080 -DbHost localhost -DbPort 3306 -DbUser root -DbPassword 123456 -DbName cloud_dream_system
```

## 常见场景

### 场景1：新环境部署
```powershell
# 1. 检查环境
.\start_backend.ps1 -Mode check

# 2. 如果环境正常，启动服务
.\start_backend.ps1 -Mode full
```

### 场景2：开发调试
```powershell
# 快速启动，跳过数据库检查
.\start_backend.ps1 -Mode dev
```

### 场景3：端口冲突
```powershell
# 使用其他端口
.\start_backend.ps1 -Port 9090
```

### 场景4：远程数据库
```powershell
# 连接远程数据库
.\start_backend.ps1 -DbHost 192.168.1.100 -DbUser remote_user -DbPassword remote_pass
```

## 故障排除

### 1. 权限问题
```powershell
# 如果遇到执行策略限制
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 2. 端口被占用
```powershell
# 查看端口占用
netstat -ano | findstr :8080

# 使用其他端口
.\start_backend.ps1 -Port 9090
```

### 3. 数据库连接失败
```powershell
# 检查MySQL服务
net start mysql

# 测试数据库连接
mysql -h localhost -P 3306 -u root -p123456 -e "SELECT 1;"
```

## 批处理脚本使用

### 基本使用
```cmd
# 完整启动
quick_start_backend.bat

# 开发模式
dev_start.bat
```

### 环境变量覆盖
```cmd
# 设置环境变量后运行
set PORT=9090
set DB_PASSWORD=mypassword
dev_start.bat
```

## Linux/macOS脚本使用

### 基本使用
```bash
# 添加执行权限
chmod +x quick_start_backend.sh

# 运行脚本
./quick_start_backend.sh
```

### 自定义环境变量
```bash
# 设置环境变量
export PORT=9090
export DB_PASSWORD=mypassword

# 运行脚本
./quick_start_backend.sh
```

## 服务管理

### 后台运行（Linux/macOS）
```bash
# 后台运行
nohup ./quick_start_backend.sh > backend.log 2>&1 &

# 查看日志
tail -f backend.log

# 停止服务
pkill -f "go run main.go"
```

### Windows服务管理
```powershell
# 使用PowerShell后台运行
Start-Process powershell -ArgumentList "-File", "start_backend.ps1" -WindowStyle Hidden

# 停止服务
Get-Process | Where-Object {$_.ProcessName -eq "go"} | Stop-Process
```

## 性能优化

### 生产环境建议
```powershell
# 1. 编译二进制文件
go build -o yunmeng-backend.exe main.go

# 2. 直接运行二进制文件
.\yunmeng-backend.exe

# 3. 或使用系统服务管理
```

### 开发环境优化
```powershell
# 使用热重载（需要安装air）
go install github.com/cosmtrek/air@latest
air
``` 