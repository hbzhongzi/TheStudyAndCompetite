@echo off
chcp 65001 >nul 2>&1

echo ========================================
echo    云梦高校科研竞赛管理系统
echo        开发环境快速启动
echo ========================================
echo.

:: 切换到项目目录
cd /d "%~dp0"

:: 设置环境变量
set DB_HOST=localhost
set DB_PORT=3306
set DB_USERNAME=root
set DB_PASSWORD=123456
set DB_DATABASE=cloud_dream_system
set DB_CHARSET=utf8mb4
set PORT=8080

echo [信息] 环境变量已设置
echo [信息] 数据库: %DB_HOST%:%DB_PORT%/%DB_DATABASE%
echo [信息] 服务端口: %PORT%
echo.

:: 快速检查依赖
echo [信息] 检查依赖...
go mod tidy >nul 2>&1
if errorlevel 1 (
    echo [错误] 依赖检查失败
    pause
    exit /b 1
)

echo [成功] 依赖检查完成
echo.

:: 启动服务
echo [信息] 启动后端服务...
echo [信息] 前端地址: http://localhost:5173
echo [信息] 后端API地址: http://localhost:%PORT%/api
echo [信息] 按 Ctrl+C 停止服务
echo.

echo [信息] 正在启动Go服务...
go run main.go

if errorlevel 1 (
    echo.
    echo [错误] 后端服务启动失败
    echo [警告] 可能的原因:
    echo   1. 数据库连接失败
    echo   2. 端口被占用
    echo   3. 依赖包问题
    echo   4. 代码编译错误
    echo.
    echo [信息] 请检查上述问题后重试
    pause
    exit /b 1
)

echo.
echo [信息] 服务已停止
pause 