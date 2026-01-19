@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ========================================
echo    管理员用户检查和修复工具
echo ========================================
echo.

:: 设置颜色
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

echo %CYAN%正在检查Go环境...%RESET%

:: 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] Go未安装或未添加到PATH%RESET%
    echo %YELLOW%请先安装Go并确保go命令可用%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] Go环境正常%RESET%

:: 检查是否在正确的目录
if not exist "main.go" (
    echo %RED%[错误] 请在go-backend目录下运行此脚本%RESET%
    echo %YELLOW%当前目录: %CD%%RESET%
    pause
    exit /b 1
)

echo %CYAN%正在检查依赖...%RESET%

:: 检查go.mod文件
if not exist "go.mod" (
    echo %RED%[错误] 未找到go.mod文件%RESET%
    pause
    exit /b 1
)

:: 下载依赖
echo %CYAN%正在下载依赖包...%RESET%
go mod tidy
if errorlevel 1 (
    echo %RED%[错误] 依赖下载失败%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] 依赖检查完成%RESET%

:: 编译修复工具
echo %CYAN%正在编译修复工具...%RESET%
go build -o fix_admin_user.exe fix_admin_user.go
if errorlevel 1 (
    echo %RED%[错误] 编译失败%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] 修复工具编译完成%RESET%

:: 运行修复工具
echo.
echo %CYAN%正在运行管理员用户检查和修复...%RESET%
echo.

fix_admin_user.exe

if errorlevel 1 (
    echo.
    echo %RED%[错误] 修复工具运行失败%RESET%
    echo %YELLOW%请检查数据库连接和配置%RESET%
) else (
    echo.
    echo %GREEN%[成功] 管理员用户检查和修复完成%RESET%
)

:: 清理临时文件
echo %CYAN%正在清理临时文件...%RESET%
if exist "fix_admin_user.exe" del "fix_admin_user.exe"

echo.
echo %CYAN%========================================%RESET%
echo %CYAN%修复工具运行完成%RESET%
echo %CYAN%========================================%RESET%
echo.
echo %YELLOW%如果修复成功，请使用以下账号登录:%RESET%
echo %GREEN%用户名: admin%RESET%
echo %GREEN%密码: 123456%RESET%
echo %GREEN%角色: admin%RESET%
echo.
echo %CYAN%现在可以重新启动后端服务并尝试创建竞赛%RESET%
echo.

pause 