@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ========================================
echo    快速启动后端服务
echo ========================================
echo.

:: 设置颜色
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

echo %CYAN%正在启动后端服务...%RESET%
echo.

:: 检查Go环境
echo %YELLOW%检查Go环境...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] Go未安装或未配置到PATH%RESET%
    pause
    exit /b 1
)
echo %GREEN%[成功] Go环境正常%RESET%

:: 检查依赖
echo %YELLOW%检查依赖...%RESET%
go mod tidy
if errorlevel 1 (
    echo %RED%[错误] 依赖检查失败%RESET%
    pause
    exit /b 1
)
echo %GREEN%[成功] 依赖检查通过%RESET%

:: 启动服务
echo %YELLOW%启动后端服务...%RESET%
echo %CYAN%服务将在 http://localhost:8080 启动%RESET%
echo %CYAN%按 Ctrl+C 停止服务%RESET%
echo.

go run main.go

echo.
echo %CYAN%服务已停止%RESET%
pause 