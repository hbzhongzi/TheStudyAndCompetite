@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

:: 设置错误处理
set "ERROR_OCCURRED=false"

:: 创建错误处理函数
:error_handler
if "%ERROR_OCCURRED%"=="true" (
    echo.
    echo %RED%========================================%RESET%
    echo %RED%脚本执行过程中发生错误%RESET%
    echo %RED%========================================%RESET%
    echo.
    echo %YELLOW%请检查上述错误信息并解决问题后重试%RESET%
    echo.
    pause
    exit /b 1
)

echo ========================================
echo    云梦高校科研竞赛管理系统
echo        后端服务启动脚本 v2.0
echo ========================================
echo.

:: 设置颜色代码
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

:: 设置项目路径
set PROJECT_ROOT=%~dp0..
cd /d "!PROJECT_ROOT!" 2>nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法切换到项目目录%RESET%
    set "ERROR_OCCURRED=true"
    goto :error_handler
)

echo %BLUE%项目路径: !PROJECT_ROOT!%RESET%
echo.

:: 检查Go环境
echo %BLUE%检查Go环境...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] Go未安装或未添加到PATH环境变量%RESET%
    echo %YELLOW%请先安装Go: https://golang.org/dl/%RESET%
    echo %YELLOW%安装完成后请重启命令行窗口%RESET%
    pause
    exit /b 1
)

:: 安全获取Go版本
set GO_VERSION=
for /f "tokens=3" %%i in ('go version 2^>nul') do set GO_VERSION=%%i
if "!GO_VERSION!"=="" (
    echo %YELLOW%[警告] 无法获取Go版本信息，继续执行%RESET%
    set GO_VERSION=未知版本
) else (
    echo %GREEN%[成功] Go版本: !GO_VERSION!%RESET%
)

:: 检查MySQL环境
echo %BLUE%检查MySQL环境...%RESET%
mysql --version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] MySQL未安装或未添加到PATH环境变量%RESET%
    echo %YELLOW%请先安装MySQL: https://dev.mysql.com/downloads/mysql/%RESET%
    echo %YELLOW%安装完成后请重启命令行窗口%RESET%
    pause
    exit /b 1
)

:: 安全获取MySQL版本
set MYSQL_VERSION=
for /f "tokens=6" %%i in ('mysql --version 2^>nul') do set MYSQL_VERSION=%%i
if "!MYSQL_VERSION!"=="" (
    echo %YELLOW%[警告] 无法获取MySQL版本信息，继续执行%RESET%
    set MYSQL_VERSION=未知版本
) else (
    echo %GREEN%[成功] MySQL版本: !MYSQL_VERSION!%RESET%
)

echo.
echo %GREEN%========================================%RESET%
echo %GREEN%环境检查完成%RESET%
echo %GREEN%========================================%RESET%
echo.

:: 获取MySQL连接参数
echo %CYAN%请输入MySQL连接信息:%RESET%
set /p DB_USER=用户名 (默认: root): 
if "!DB_USER!"=="" set DB_USER=root

set /p DB_PASSWORD=密码: 
if "!DB_PASSWORD!"=="" (
    echo %RED%[错误] 密码不能为空%RESET%
    pause
    exit /b 1
)

set DB_HOST=localhost
set DB_PORT=3306
set DB_NAME=cloud_dream_system

echo.
echo %BLUE%数据库连接信息:%RESET%
echo   主机: !DB_HOST!
echo   端口: !DB_PORT!
echo   用户: !DB_USER!
echo   数据库: !DB_NAME!
echo.

:: 测试MySQL连接
echo %BLUE%测试MySQL连接...%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "SELECT VERSION() as mysql_version;" 2>nul >nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法连接到MySQL数据库%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. MySQL服务未启动%RESET%
    echo %YELLOW%2. 连接参数错误（主机、端口、用户名、密码）%RESET%
    echo %YELLOW%3. 防火墙阻止连接%RESET%
    echo %YELLOW%4. MySQL配置不允许远程连接%RESET%
    echo.
    echo %CYAN%请检查MySQL服务状态并重试%RESET%
    echo %YELLOW%调试信息: ERRORLEVEL=!ERRORLEVEL!%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] MySQL连接正常%RESET%
echo.

:: 检查数据库状态
echo %BLUE%检查数据库状态...%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "USE !DB_NAME!;" 2>nul >nul
if !ERRORLEVEL! NEQ 0 (
    echo %YELLOW%[信息] 数据库不存在，将创建新数据库%RESET%
    mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "CREATE DATABASE !DB_NAME! CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul >nul
    if !ERRORLEVEL! NEQ 0 (
        echo %RED%[错误] 创建数据库失败%RESET%
        echo %YELLOW%可能的原因:%RESET%
        echo %YELLOW%1. 用户权限不足%RESET%
        echo %YELLOW%2. 数据库名称已存在但无法访问%RESET%
        echo %YELLOW%3. MySQL配置问题%RESET%
        pause
        exit /b 1
    )
    echo %GREEN%[成功] 数据库创建完成%RESET%
    set NEED_INIT_DB=true
) else (
    echo %GREEN%[信息] 数据库已存在%RESET%
    
    :: 检查表数量 - 使用更安全的方式
    set TABLE_COUNT=0
    for /f "tokens=1" %%a in ('mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '!DB_NAME!';" 2^>nul ^| findstr /r "^[0-9]"') do (
        set TABLE_COUNT=%%a
    )
    
    if "!TABLE_COUNT!"=="" (
        echo %YELLOW%[警告] 无法检查表数量，假设需要初始化%RESET%
        set NEED_INIT_DB=true
    ) else if !TABLE_COUNT! LSS 15 (
        echo %YELLOW%[警告] 数据库表不完整，需要初始化%RESET%
        set NEED_INIT_DB=true
    ) else (
        echo %GREEN%[信息] 数据库表完整，无需初始化%RESET%
        set NEED_INIT_DB=false
    )
)

:: 显示调试信息
echo %YELLOW%[调试] 数据库检查完成，NEED_INIT_DB=!NEED_INIT_DB!%RESET%

:: 初始化数据库
if "!NEED_INIT_DB!"=="true" (
    echo.
    echo %CYAN%是否初始化数据库? (y/n)%RESET%
    set /p INIT_DB_CHOICE=
    
    if /i "!INIT_DB_CHOICE!"=="y" (
        echo %BLUE%初始化数据库...%RESET%
        
        :: 检查SQL文件
        if exist "sql\database_setup_compatible.sql" (
            set SQL_FILE=sql\database_setup_compatible.sql
            echo %GREEN%[信息] 使用兼容版SQL文件%RESET%
        ) else if exist "sql\database_setup.sql" (
            set SQL_FILE=sql\database_setup.sql
            echo %GREEN%[信息] 使用标准版SQL文件%RESET%
        ) else (
            echo %RED%[错误] 未找到SQL初始化文件%RESET%
            echo %YELLOW%请检查以下文件是否存在:%RESET%
            echo %YELLOW%1. sql\database_setup_compatible.sql%RESET%
            echo %YELLOW%2. sql\database_setup.sql%RESET%
            pause
            exit /b 1
        )
        
        :: 执行SQL脚本
        mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! < "!SQL_FILE!" 2>nul
        if !ERRORLEVEL! EQU 0 (
            echo %GREEN%[成功] 数据库初始化完成%RESET%
        ) else (
            echo %RED%[错误] 数据库初始化失败%RESET%
            echo %YELLOW%可能的原因:%RESET%
            echo %YELLOW%1. SQL文件语法错误%RESET%
            echo %YELLOW%2. 用户权限不足%RESET%
            echo %YELLOW%3. 数据库连接问题%RESET%
            echo %YELLOW%4. MySQL版本兼容性问题%RESET%
            pause
            exit /b 1
        )
    ) else (
        echo %YELLOW%[跳过] 数据库初始化%RESET%
    )
)

:: 安装Go依赖
echo.
echo %BLUE%安装Go依赖...%RESET%
go mod tidy 2>nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] Go依赖安装失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. 网络连接问题%RESET%
    echo %YELLOW%2. Go模块配置错误%RESET%
    echo %YELLOW%3. 依赖包版本冲突%RESET%
    echo %YELLOW%4. 代理设置问题%RESET%
    echo.
    echo %CYAN%请检查网络连接并重试%RESET%
    pause
    exit /b 1
)
echo %GREEN%[成功] Go依赖安装完成%RESET%

:: 编译后端
echo %BLUE%编译后端服务...%RESET%
go build -o yunmeng-backend.exe . 2>nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 后端编译失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. 代码语法错误%RESET%
    echo %YELLOW%2. 依赖包缺失%RESET%
    echo %YELLOW%3. Go版本不兼容%RESET%
    echo %YELLOW%4. 文件权限问题%RESET%
    echo.
    echo %CYAN%请检查代码并修复错误后重试%RESET%
    pause
    exit /b 1
)
echo %GREEN%[成功] 后端编译完成%RESET%

:: 检查端口占用
echo %BLUE%检查端口占用...%RESET%
netstat -an 2>nul | findstr ":8080" >nul
if !ERRORLEVEL! EQU 0 (
    echo %YELLOW%[警告] 端口8080已被占用%RESET%
    echo %CYAN%是否强制释放端口? (y/n)%RESET%
    set /p FORCE_RELEASE=
    
    if /i "!FORCE_RELEASE!"=="y" (
        echo %BLUE%释放端口8080...%RESET%
        for /f "tokens=5" %%a in ('netstat -ano 2^>nul ^| findstr ":8080"') do (
            taskkill /f /pid %%a >nul 2>&1
        )
        echo %GREEN%[成功] 端口已释放%RESET%
    ) else (
        echo %YELLOW%[跳过] 端口释放%RESET%
    )
) else (
    echo %GREEN%[信息] 端口8080可用%RESET%
)

:: 启动后端服务
echo.
echo %BLUE%启动后端服务...%RESET%

:: 设置环境变量
set "DB_HOST=!DB_HOST!"
set "DB_PORT=!DB_PORT!"
set "DB_USERNAME=!DB_USER!"
set "DB_PASSWORD=!DB_PASSWORD!"
set "DB_DATABASE=!DB_NAME!"
set "DB_CHARSET=utf8mb4"
set "PORT=8080"

echo %CYAN%环境变量已设置%RESET%
echo %CYAN%数据库: !DB_HOST!:!DB_PORT!/!DB_NAME!%RESET%
echo %CYAN%服务端口: !PORT!%RESET%
echo.

:: 尝试启动服务
echo %CYAN%正在启动Go服务...%RESET%
go run main.go
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 后端服务启动失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. 数据库连接失败%RESET%
    echo %YELLOW%2. 端口被占用%RESET%
    echo %YELLOW%3. 依赖包问题%RESET%
    echo %YELLOW%4. 代码编译错误%RESET%
    echo.
    echo %CYAN%请检查上述问题后重试%RESET%
    pause
    exit /b 1
)
echo %GREEN%[成功] 后端服务已启动%RESET%

:: 等待后端服务启动
echo %BLUE%等待后端服务启动...%RESET%
timeout /t 5 /nobreak >nul 2>&1

:: 检查后端服务状态
echo %BLUE%检查后端服务状态...%RESET%
echo %CYAN%服务已启动，请检查以下地址:%RESET%
echo %CYAN%前端地址: http://localhost:5173%RESET%
echo %CYAN%后端API地址: http://localhost:8080/api%RESET%
echo.
echo %YELLOW%如果服务启动失败，请检查:%RESET%
echo %YELLOW%1. 数据库连接是否正常%RESET%
echo %YELLOW%2. 端口8080是否被占用%RESET%
echo %YELLOW%3. 依赖包是否正确安装%RESET%
echo %YELLOW%4. 代码是否有编译错误%RESET%

:: 显示服务信息
echo.
echo %GREEN%========================================%RESET%
echo %GREEN%后端服务启动完成%RESET%
echo %GREEN%========================================%RESET%
echo.
echo %CYAN%服务信息:%RESET%
echo   后端服务: http://localhost:8080
echo   数据库: !DB_HOST!:!DB_PORT!/!DB_NAME!
echo.
echo %CYAN%默认用户:%RESET%
echo   管理员: admin / 123456
echo   教师: teacher001 / 123456
echo   学生: student001 / 123456
echo.
echo %CYAN%API文档:%RESET%
echo   后端API: http://localhost:8080/api/docs
echo.
echo %YELLOW%提示: 在后端服务窗口中按Ctrl+C可以停止服务%RESET%
echo.

:: 等待用户操作
echo %CYAN%服务启动完成，按任意键退出...%RESET%
pause >nul

echo.
echo %GREEN%========================================%RESET%
echo %GREEN%启动脚本执行完成%RESET%
echo %GREEN%========================================%RESET%
echo.
pause 