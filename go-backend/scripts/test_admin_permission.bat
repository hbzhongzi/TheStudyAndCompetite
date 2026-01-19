@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ========================================
echo    管理员权限测试工具
echo ========================================
echo.

:: 设置颜色
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

echo %CYAN%正在检查后端服务状态...%RESET%

:: 检查后端服务是否运行
curl -s http://localhost:8080/api/health >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] 后端服务未运行%RESET%
    echo %YELLOW%请先启动后端服务: go run main.go%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] 后端服务运行正常%RESET%

:: 检查是否在正确的目录
if not exist "main.go" (
    echo %RED%[错误] 请在go-backend目录下运行此脚本%RESET%
    echo %YELLOW%当前目录: %CD%%RESET%
    pause
    exit /b 1
)

:: 编译测试工具
echo %CYAN%正在编译测试工具...%RESET%
go build -o test_admin_permission.exe test_admin_permission.go
if errorlevel 1 (
    echo %RED%[错误] 编译失败%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] 测试工具编译完成%RESET%

:: 运行测试
echo.
echo %CYAN%正在运行管理员权限测试...%RESET%
echo.

test_admin_permission.exe

if errorlevel 1 (
    echo.
    echo %RED%[错误] 测试失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo   1. 管理员用户不存在
    echo   2. 密码错误
    echo   3. 角色未正确分配
    echo   4. 权限中间件配置问题
    echo.
    echo %CYAN%建议运行修复工具: scripts\fix_admin_user.bat%RESET%
) else (
    echo.
    echo %GREEN%[成功] 所有测试通过！%RESET%
    echo %CYAN%管理员权限正常工作%RESET%
)

:: 清理临时文件
echo %CYAN%正在清理临时文件...%RESET%
if exist "test_admin_permission.exe" del "test_admin_permission.exe"

echo.
echo %CYAN%========================================%RESET%
echo %CYAN%测试完成%RESET%
echo %CYAN%========================================%RESET%
echo.

pause 