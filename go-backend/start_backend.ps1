# 云梦高校科研竞赛管理系统 - PowerShell快速启动脚本
# 需要PowerShell 5.0或更高版本

param(
    [string]$Mode = "full",  # full, dev, check
    [string]$Port = "8080",
    [string]$DbHost = "localhost",
    [string]$DbPort = "3306",
    [string]$DbUser = "root",
    [string]$DbPassword = "123456",
    [string]$DbName = "cloud_dream_system"
)

# 设置控制台编码
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# 颜色函数
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "[成功] $Message" "Green"
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "[信息] $Message" "Cyan"
}

function Write-Warning {
    param([string]$Message)
    Write-ColorOutput "[警告] $Message" "Yellow"
}

function Write-Error {
    param([string]$Message)
    Write-ColorOutput "[错误] $Message" "Red"
}

# 显示标题
Write-ColorOutput "========================================" "Blue"
Write-ColorOutput "  云梦高校科研竞赛管理系统 - 快速启动" "Blue"
Write-ColorOutput "========================================" "Blue"
Write-Host ""

# 切换到脚本所在目录
$ScriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptPath
Write-Info "当前目录: $(Get-Location)"
Write-Host ""

# 环境检查函数
function Test-GoEnvironment {
    Write-Info "检查Go环境..."
    try {
        $goVersion = go version 2>$null
        if ($LASTEXITCODE -eq 0) {
            $version = ($goVersion -split " ")[2]
            Write-Success "Go版本: $version"
            return $true
        }
    }
    catch {
        Write-Error "Go未安装或未添加到PATH环境变量"
        Write-Warning "请先安装Go: https://golang.org/dl/"
        Write-Warning "安装完成后请重启PowerShell"
        return $false
    }
}

function Test-MySQLEnvironment {
    Write-Info "检查MySQL环境..."
    try {
        $mysqlVersion = mysql --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            $version = ($mysqlVersion -split " ")[5] -replace ",", ""
            Write-Success "MySQL版本: $version"
            return $true
        }
    }
    catch {
        Write-Warning "MySQL未安装或未添加到PATH环境变量"
        Write-Warning "请确保MySQL服务正在运行"
        return $false
    }
}

function Test-DatabaseConnection {
    Write-Info "检查数据库连接..."
    try {
        $testQuery = "SELECT 1;"
        $result = mysql -h$DbHost -P$DbPort -u$DbUser -p$DbPassword -e $testQuery 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "数据库连接正常"
            return $true
        }
    }
    catch {
        Write-Warning "无法连接到数据库"
        Write-Warning "请确保MySQL服务正在运行，并且连接参数正确"
        Write-Info "尝试使用以下命令启动MySQL服务:"
        Write-Host "  net start mysql" -ForegroundColor Cyan
        Write-Host ""
        
        $continue = Read-Host "是否继续启动后端服务？(y/n)"
        if ($continue -eq "y" -or $continue -eq "Y") {
            return $true
        }
        return $false
    }
}

function Install-Dependencies {
    Write-Info "检查并安装Go依赖..."
    
    if (-not (Test-Path "go.mod")) {
        Write-Error "未找到go.mod文件"
        return $false
    }
    
    Write-Info "正在下载依赖包..."
    go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Error "依赖下载失败"
        return $false
    }
    
    Write-Info "正在整理依赖..."
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "依赖整理失败"
        return $false
    }
    
    Write-Success "依赖安装完成"
    return $true
}

function Start-BackendService {
    Write-Info "启动后端服务..."
    Write-Info "服务启动中，请稍候..."
    Write-Info "前端地址: http://localhost:5173"
    Write-Info "后端API地址: http://localhost:$Port/api"
    Write-Info "按 Ctrl+C 停止服务"
    Write-Host ""
    
    # 设置环境变量
    $env:DB_HOST = $DbHost
    $env:DB_PORT = $DbPort
    $env:DB_USERNAME = $DbUser
    $env:DB_PASSWORD = $DbPassword
    $env:DB_DATABASE = $DbName
    $env:DB_CHARSET = "utf8mb4"
    $env:PORT = $Port
    
    # 启动服务
    Write-Info "正在启动Go服务..."
    try {
        go run main.go
        if ($LASTEXITCODE -ne 0) {
            Write-Host ""
            Write-Error "后端服务启动失败"
            Write-Warning "可能的原因:"
            Write-Host "  1. 数据库连接失败" -ForegroundColor Yellow
            Write-Host "  2. 端口被占用" -ForegroundColor Yellow
            Write-Host "  3. 依赖包问题" -ForegroundColor Yellow
            Write-Host "  4. 代码编译错误" -ForegroundColor Yellow
            Write-Host ""
            Write-Info "请检查上述问题后重试"
            exit 1
        }
    }
    catch {
        Write-Host ""
        Write-Error "后端服务启动失败: $($_.Exception.Message)"
        Write-Info "请检查上述问题后重试"
        exit 1
    }
}

# 主执行逻辑
switch ($Mode.ToLower()) {
    "check" {
        Write-Info "执行环境检查模式..."
        $goOk = Test-GoEnvironment
        $mysqlOk = Test-MySQLEnvironment
        if ($goOk -and $mysqlOk) {
            Write-Success "环境检查完成，所有依赖正常"
        } else {
            Write-Error "环境检查失败，请解决上述问题后重试"
            exit 1
        }
    }
    "dev" {
        Write-Info "执行开发模式快速启动..."
        if (-not (Test-GoEnvironment)) { exit 1 }
        if (-not (Install-Dependencies)) { exit 1 }
        Start-BackendService
    }
    "full" {
        Write-Info "执行完整模式启动..."
        if (-not (Test-GoEnvironment)) { exit 1 }
        if (-not (Test-MySQLEnvironment)) { 
            Write-Warning "MySQL检查失败，但继续执行"
        }
        if (-not (Install-Dependencies)) { exit 1 }
        if (-not (Test-DatabaseConnection)) { exit 1 }
        Start-BackendService
    }
    default {
        Write-Error "未知的启动模式: $Mode"
        Write-Info "可用模式: full, dev, check"
        exit 1
    }
}

Write-Host ""
Write-Warning "服务已停止" 