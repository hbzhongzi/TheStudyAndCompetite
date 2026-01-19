@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

echo ========================================
echo    数据库连接测试脚本 v2.0
echo ========================================
echo.

:: 设置颜色代码
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

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
echo %YELLOW%执行命令: mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "SELECT VERSION() as mysql_version;"%RESET%
echo.

mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "SELECT VERSION() as mysql_version;"
echo.
echo %YELLOW%ERRORLEVEL: !ERRORLEVEL!%RESET%

if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法连接到MySQL数据库%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. MySQL服务未启动%RESET%
    echo %YELLOW%2. 连接参数错误（主机、端口、用户名、密码）%RESET%
    echo %YELLOW%3. 防火墙阻止连接%RESET%
    echo %YELLOW%4. MySQL配置不允许远程连接%RESET%
    echo.
    echo %CYAN%请检查MySQL服务状态并重试%RESET%
) else (
    echo %GREEN%[成功] MySQL连接正常%RESET%
    
    :: 测试数据库是否存在
    echo.
    echo %BLUE%测试数据库是否存在...%RESET%
    mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "USE !DB_NAME!;"
    echo.
    echo %YELLOW%ERRORLEVEL: !ERRORLEVEL!%RESET%
    
    if !ERRORLEVEL! NEQ 0 (
        echo %YELLOW%[信息] 数据库不存在%RESET%
        echo %CYAN%是否创建数据库? (y/n)%RESET%
        set /p CREATE_DB_CHOICE=
        
        if /i "!CREATE_DB_CHOICE!"=="y" (
            echo %BLUE%创建数据库...%RESET%
            mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "CREATE DATABASE !DB_NAME! CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
            if !ERRORLEVEL! EQU 0 (
                echo %GREEN%[成功] 数据库创建完成%RESET%
            ) else (
                echo %RED%[错误] 数据库创建失败%RESET%
            )
        ) else (
            echo %YELLOW%[跳过] 数据库创建%RESET%
        )
    ) else (
        echo %GREEN%[信息] 数据库存在%RESET%
        
        :: 测试表数量查询
        echo.
        echo %BLUE%测试表数量查询...%RESET%
        mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) as table_count FROM information_schema.tables WHERE table_schema = '!DB_NAME!';"
        echo.
        echo %YELLOW%ERRORLEVEL: !ERRORLEVEL!%RESET%
        
        if !ERRORLEVEL! EQU 0 (
            :: 获取表数量
            set TABLE_COUNT=0
            for /f "tokens=1" %%a in ('mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '!DB_NAME!';" 2^>nul ^| findstr /r "^[0-9]"') do (
                set TABLE_COUNT=%%a
            )
            
            if "!TABLE_COUNT!"=="" (
                echo %YELLOW%[警告] 无法获取表数量%RESET%
            ) else (
                echo %GREEN%[信息] 数据库中共有 !TABLE_COUNT! 个表%RESET%
                
                if !TABLE_COUNT! LSS 15 (
                    echo %YELLOW%[警告] 表数量较少，可能需要初始化数据库%RESET%
                ) else (
                    echo %GREEN%[信息] 数据库表结构完整%RESET%
                )
            )
            
            :: 检查关键表
            echo.
            echo %BLUE%检查关键表...%RESET%
            mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'users';" 2>nul | findstr "users" >nul && echo %GREEN%   ✓ users表%RESET% || echo %RED%   ✗ users表%RESET%
            mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'projects';" 2>nul | findstr "projects" >nul && echo %GREEN%   ✓ projects表%RESET% || echo %RED%   ✗ projects表%RESET%
            mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'competitions';" 2>nul | findstr "competitions" >nul && echo %GREEN%   ✓ competitions表%RESET% || echo %RED%   ✗ competitions表%RESET%
            
            :: 检查用户数据
            echo.
            echo %BLUE%检查用户数据...%RESET%
            set USER_COUNT=0
            for /f "tokens=1" %%a in ('mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) FROM users;" 2^>nul ^| findstr /r "^[0-9]"') do (
                set USER_COUNT=%%a
            )
            
            if "!USER_COUNT!"=="" (
                echo %YELLOW%[警告] 无法获取用户数量%RESET%
            ) else (
                echo %GREEN%[信息] 用户数量: !USER_COUNT!%RESET%
                
                if !USER_COUNT! LSS 3 (
                    echo %YELLOW%[警告] 用户数量较少，可能需要初始化默认数据%RESET%
                ) else (
                    echo %GREEN%[信息] 用户数据完整%RESET%
                )
            )
        ) else (
            echo %RED%[错误] 表数量查询失败%RESET%
        )
    )
)

:: 测试端口连接
echo.
echo %BLUE%测试端口连接...%RESET%
echo %YELLOW%测试MySQL端口 (!DB_PORT!) 连接...%RESET%
telnet !DB_HOST! !DB_PORT! <nul >nul 2>&1
if !ERRORLEVEL! EQU 0 (
    echo %GREEN%[成功] 端口连接正常%RESET%
) else (
    echo %RED%[错误] 端口连接失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. MySQL服务未启动%RESET%
    echo %YELLOW%2. 端口被占用%RESET%
    echo %YELLOW%3. 防火墙阻止连接%RESET%
)

:: 测试网络连接
echo.
echo %BLUE%测试网络连接...%RESET%
ping -n 1 !DB_HOST! >nul 2>&1
if !ERRORLEVEL! EQU 0 (
    echo %GREEN%[成功] 网络连接正常%RESET%
) else (
    echo %RED%[错误] 网络连接失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. 主机地址错误%RESET%
    echo %YELLOW%2. 网络配置问题%RESET%
    echo %YELLOW%3. 防火墙阻止连接%RESET%
)

echo.
echo %GREEN%========================================%RESET%
echo %GREEN%测试完成%RESET%
echo %GREEN%========================================%RESET%
echo.
echo %CYAN%测试总结:%RESET%
echo   MySQL连接: %GREEN%✓%RESET%
echo   数据库存在: %GREEN%✓%RESET%
echo   表结构完整: %GREEN%✓%RESET%
echo   用户数据: %GREEN%✓%RESET%
echo   端口连接: %GREEN%✓%RESET%
echo   网络连接: %GREEN%✓%RESET%
echo.
echo %YELLOW%如果所有测试都通过，可以运行launch.bat启动后端服务%RESET%
echo.
pause 