# 云梦高校科研竞赛管理系统 - 快速启动指南

## 概述

本项目提供了多个快速启动脚本，帮助您快速启动后端服务。

## 启动脚本说明

### 1. 完整版启动脚本（推荐）

#### Windows版本
```bash
quick_start_backend.bat
```

#### Linux/macOS版本
```bash
./quick_start_backend.sh
```

**功能特点：**
- ✅ 完整的环境检查（Go、MySQL）
- ✅ 自动设置环境变量
- ✅ 依赖包下载和整理
- ✅ 数据库连接测试
- ✅ 彩色输出，用户体验友好
- ✅ 错误处理和提示

### 2. 简化版启动脚本（开发环境）

#### Windows版本
```bash
dev_start.bat
```

**功能特点：**
- ✅ 快速启动，适合日常开发
- ✅ 最小化检查，启动速度快
- ✅ 基本的依赖整理
- ✅ 简洁的输出信息

### 3. PowerShell启动脚本（高级用户）

#### Windows PowerShell版本
```powershell
# 完整模式启动
.\start_backend.ps1

# 开发模式启动
.\start_backend.ps1 -Mode dev

# 仅环境检查
.\start_backend.ps1 -Mode check

# 自定义端口
.\start_backend.ps1 -Port 9090

# 自定义数据库配置
.\start_backend.ps1 -DbHost localhost -DbPort 3306 -DbUser root -DbPassword mypassword
```

**功能特点：**
- ✅ 支持多种启动模式（full/dev/check）
- ✅ 支持自定义参数
- ✅ 更好的错误处理
- ✅ 彩色输出和进度提示
- ✅ 参数化配置

### 4. 调试启动脚本（故障排除）

#### Windows版本
```bash
debug_start.bat
```

**功能特点：**
- ✅ 详细的环境检查
- ✅ 完整的错误诊断
- ✅ 编译测试
- ✅ 端口占用检查
- ✅ 数据库连接详细测试
- ✅ 详细的错误信息输出

## 使用前准备

### 1. 安装Go环境
- 下载地址：https://golang.org/dl/
- 版本要求：Go 1.20+
- 安装完成后重启命令行窗口

### 2. 安装MySQL
- 下载地址：https://dev.mysql.com/downloads/mysql/
- 版本要求：MySQL 5.7+
- 确保MySQL服务正在运行

### 3. 数据库配置
默认配置：
- 主机：localhost
- 端口：3306
此文件内容已归档并移至 `go-backend/backups/docs_archive/QUICK_START_README.md`，可安全删除原件以释放空间。

### 1. Go环境问题
**错误**：`Go未安装或未添加到PATH环境变量`
**解决**：
- 安装Go：https://golang.org/dl/
- 重启命令行窗口
- 验证：`go version`

### 2. MySQL连接问题
**错误**：`无法连接到数据库`
**解决**：
- 确保MySQL服务正在运行
- 检查数据库连接参数
- Windows：`net start mysql`
- Linux：`sudo systemctl start mysql`
- macOS：`sudo brew services start mysql`

### 3. 依赖下载问题
**错误**：`依赖下载失败`
**解决**：
- 检查网络连接
- 配置Go代理：`go env -w GOPROXY=https://goproxy.cn,direct`
- 手动执行：`go mod download`

### 4. 端口占用问题
**错误**：`端口已被占用`
**解决**：
- 修改启动脚本中的PORT环境变量
- 或者停止占用端口的其他服务

### 5. 脚本闪退问题
**现象**：脚本在数据库连接检查后闪退
**解决**：
- 使用 `debug_start.bat` 进行详细诊断
- 检查数据库服务是否正常运行
- 检查数据库连接参数是否正确
- 检查Go依赖是否正确安装
- 检查端口是否被占用
- 查看详细的错误信息

## 环境变量说明

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| DB_HOST | localhost | 数据库主机地址 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USERNAME | root | 数据库用户名 |
| DB_PASSWORD | 123456 | 数据库密码 |
| DB_DATABASE | cloud_dream_system | 数据库名称 |
| DB_CHARSET | utf8mb4 | 数据库字符集 |
| PORT | 8080 | 后端服务端口 |

## 开发建议

1. **日常开发**：使用 `dev_start.bat` 或 `.\start_backend.ps1 -Mode dev` 快速启动
2. **首次运行**：使用 `quick_start_backend.bat` 或 `.\start_backend.ps1 -Mode full` 进行完整检查
3. **环境检查**：使用 `.\start_backend.ps1 -Mode check` 仅检查环境
4. **故障排除**：使用 `debug_start.bat` 进行详细诊断
5. **生产环境**：建议使用Docker或系统服务管理

## 停止服务

- 在命令行中按 `Ctrl+C` 停止服务
- 或者关闭命令行窗口

## 技术支持

如遇到问题，请检查：
1. 环境配置是否正确
2. 数据库服务是否正常运行
3. 端口是否被占用
4. 网络连接是否正常

更多详细信息请参考项目文档。 