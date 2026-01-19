# 用户表扩展迁移脚本
Write-Host "正在执行用户表扩展迁移..." -ForegroundColor Green
Write-Host ""

# 检查MySQL是否运行
Write-Host "检查MySQL服务状态..." -ForegroundColor Yellow
$mysqlService = Get-Service -Name "*MySQL*" -ErrorAction SilentlyContinue
if ($mysqlService -and $mysqlService.Status -eq "Running") {
    Write-Host "MySQL服务正在运行" -ForegroundColor Green
} else {
    Write-Host "MySQL服务未运行，请先启动MySQL服务" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

# 执行迁移脚本
Write-Host "执行数据库迁移脚本..." -ForegroundColor Yellow
try {
    $migrationPath = Join-Path $PSScriptRoot "..\migrations\extend_user_table.sql"
    Get-Content $migrationPath | mysql -u root -p
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "用户表扩展迁移完成！" -ForegroundColor Green
        Write-Host "新添加的字段：" -ForegroundColor Cyan
        Write-Host "- department: 部门/院系" -ForegroundColor White
        Write-Host "- title: 职称/职位" -ForegroundColor White
        Write-Host "- grade: 年级" -ForegroundColor White
        Write-Host "- major: 专业" -ForegroundColor White
        Write-Host "- created_at: 创建时间" -ForegroundColor White
        Write-Host "- updated_at: 更新时间" -ForegroundColor White
    } else {
        Write-Host ""
        Write-Host "迁移失败，请检查错误信息" -ForegroundColor Red
    }
} catch {
    Write-Host "执行迁移时发生错误: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Read-Host "按回车键退出" 