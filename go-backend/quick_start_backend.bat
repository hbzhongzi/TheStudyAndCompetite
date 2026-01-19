@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

:: 设置颜色代码
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

echo %BLUE%========================================%RESET%
echo %BLUE%  云梦高校科研竞赛管理系统 - 快速启动%RESET%
echo %BLUE%========================================%RESET%
echo.

:: 切换到项目目录
cd /d "%~dp0" 2>nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法切换到项目目录%RESET%
    pause
    exit /b 1
)

echo %GREEN%[信息] 当前目录: %CD%%RESET%
echo.

:: 检查Go环境
echo %BLUE%[步骤1] 检查Go环境...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] Go未安装或未添加到PATH环境变量%RESET%
    echo %YELLOW%请先安装Go: https://golang.org/dl/%RESET%
    echo %YELLOW%安装完成后请重启命令行窗口%RESET%
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version 2^>nul') do set GO_VERSION=%%i
echo %GREEN%[成功] Go版本: !GO_VERSION!%RESET%
echo.

:: 检查MySQL环境
echo %BLUE%[步骤2] 检查MySQL环境...%RESET%
mysql --version >nul 2>&1
if errorlevel 1 (
    echo %YELLOW%[警告] MySQL未安装或未添加到PATH环境变量%RESET%
    echo %YELLOW%请确保MySQL服务正在运行%RESET%
    echo.
) else (
    for /f "tokens=6" %%i in ('mysql --version 2^>nul') do set MYSQL_VERSION=%%i
    echo %GREEN%[成功] MySQL版本: !MYSQL_VERSION!%RESET%
)
echo.

:: 设置默认环境变量
echo %BLUE%[步骤3] 设置环境变量...%RESET%
set "DB_HOST=localhost"
set "DB_PORT=3306"
set "DB_USERNAME=root"
set "DB_PASSWORD=123456"
set "DB_DATABASE=cloud_dream_system"
set "DB_CHARSET=utf8mb4"
set "PORT=8080"

echo %GREEN%[成功] 环境变量设置完成%RESET%
echo %CYAN%数据库配置:%RESET%
echo   主机: %DB_HOST%
echo   端口: %DB_PORT%
echo   用户名: %DB_USERNAME%
echo   数据库: %DB_DATABASE%
echo   服务端口: %PORT%
echo.

:: 检查并安装依赖
echo %BLUE%[步骤4] 检查并安装Go依赖...%RESET%
if not exist "go.mod" (
    echo %RED%[错误] 未找到go.mod文件%RESET%
    pause
    exit /b 1
)

echo %CYAN%正在下载依赖包...%RESET%
go mod download
if errorlevel 1 (
    echo %RED%[错误] 依赖下载失败%RESET%
    pause
    exit /b 1
)

echo %CYAN%正在整理依赖...%RESET%
go mod tidy
if errorlevel 1 (
    echo %RED%[错误] 依赖整理失败%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] 依赖安装完成%RESET%
echo.

:: 检查数据库连接
echo %BLUE%[步骤5] 检查数据库连接...%RESET%
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USERNAME% -p%DB_PASSWORD% -e "SELECT 1;" >nul 2>&1
if errorlevel 1 (
    echo %YELLOW%[警告] 无法连接到数据库%RESET%
    echo %YELLOW%请确保MySQL服务正在运行，并且连接参数正确%RESET%
    echo %CYAN%尝试使用以下命令启动MySQL服务:%RESET%
    echo   net start mysql
    echo.
    echo %YELLOW%是否继续启动后端服务？(y/n):%RESET%
    set /p CONTINUE=
    if /i not "!CONTINUE!"=="y" (
        echo %YELLOW%启动已取消%RESET%
        pause
        exit /b 0
    )
    echo %YELLOW%[警告] 数据库连接失败，但将继续启动服务%RESET%
    echo %YELLOW%[警告] 如果后端启动失败，请先解决数据库连接问题%RESET%
) else (
    echo %GREEN%[成功] 数据库连接正常%RESET%
)
echo.

:: 启动后端服务
echo %BLUE%[步骤6] 启动后端服务...%RESET%
echo %CYAN%服务启动中，请稍候...%RESET%
echo %CYAN%前端地址: http://localhost:5173%RESET%
echo %CYAN%后端API地址: http://localhost:%PORT%/api%RESET%
echo %CYAN%按 Ctrl+C 停止服务%RESET%
echo.

:: 设置环境变量并启动服务
set "DB_HOST=%DB_HOST%"
set "DB_PORT=%DB_PORT%"
set "DB_USERNAME=%DB_USERNAME%"
set "DB_PASSWORD=%DB_PASSWORD%"
set "DB_DATABASE=%DB_DATABASE%"
set "DB_CHARSET=%DB_CHARSET%"
set "PORT=%PORT%"

echo %CYAN%正在启动Go服务...%RESET%
go run main.go

if errorlevel 1 (
    echo.
    echo %RED%[错误] 后端服务启动失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo   1. 数据库连接失败
    echo   2. 端口被占用
    echo   3. 依赖包问题
    echo   4. 代码编译错误
    echo.
    echo %CYAN%请检查上述问题后重试%RESET%
    pause
    exit /b 1
)

echo.
echo %YELLOW%[信息] 服务已停止%RESET%
pause 