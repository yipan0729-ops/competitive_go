# 竞品提取改进测试脚本

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "🧪 竞品提取改进测试" -ForegroundColor Yellow
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# 1. 停止旧服务
Write-Host "1️⃣ 停止旧服务..." -ForegroundColor Yellow
$process = Get-Process | Where-Object { $_.Name -like "*competitive*" -or ($_.Path -and $_.Path -like "*competitive*.exe") }
if ($process) {
    $process | ForEach-Object { 
        Stop-Process -Id $_.Id -Force
        Write-Host "   ✅ 已停止: $($_.Name) (PID: $($_.Id))" -ForegroundColor Green
    }
    Start-Sleep -Seconds 2
} else {
    Write-Host "   ℹ️  无运行中的服务" -ForegroundColor Gray
}

# 2. 启动新版本
Write-Host ""
Write-Host "2️⃣ 启动v3版本..." -ForegroundColor Yellow
$job = Start-Job -ScriptBlock {
    Set-Location "D:\Code\Competitive_go"
    .\competitive-analyzer-v3.exe
}
Write-Host "   ✅ 服务已启动 (Job ID: $($job.Id))" -ForegroundColor Green
Start-Sleep -Seconds 3

# 3. 健康检查
Write-Host ""
Write-Host "3️⃣ 健康检查..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -UseBasicParsing
    Write-Host "   ✅ 服务运行正常" -ForegroundColor Green
} catch {
    Write-Host "   ❌ 服务启动失败" -ForegroundColor Red
    exit 1
}

# 4. 测试竞品发现（对比改进前后）
Write-Host ""
Write-Host "4️⃣ 测试竞品提取（改进后）..." -ForegroundColor Yellow
Write-Host "   主题: AI创作工具" -ForegroundColor White

$body = @{
    topic = "AI创作工具"
    market = "中国"
    limit = 5
} | ConvertTo-Json -Compress

try {
    Write-Host "   🔍 搜索中..." -ForegroundColor Cyan
    $result = Invoke-RestMethod `
        -Uri "http://localhost:8080/api/discover/search" `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -UseBasicParsing `
        -TimeoutSec 30
    
    Write-Host ""
    Write-Host "   📊 搜索结果:" -ForegroundColor Green
    Write-Host "   任务ID: $($result.task_id)" -ForegroundColor White
    Write-Host "   状态: $($result.status)" -ForegroundColor White
    
    # 等待任务完成
    Start-Sleep -Seconds 5
    
    # 查询任务结果
    Write-Host ""
    Write-Host "5️⃣ 查询任务结果..." -ForegroundColor Yellow
    $task = Invoke-RestMethod `
        -Uri "http://localhost:8080/api/discover/tasks/$($result.task_id)" `
        -UseBasicParsing
    
    Write-Host "   状态: $($task.status)" -ForegroundColor White
    Write-Host "   发现结果数: $($task.results_count)" -ForegroundColor White
    
    if ($task.results -and $task.results.Count -gt 0) {
        Write-Host ""
        Write-Host "   🎯 提取的竞品:" -ForegroundColor Cyan
        $task.results | ForEach-Object {
            Write-Host "      - $($_.title)" -ForegroundColor White
            Write-Host "        URL: $($_.url)" -ForegroundColor Gray
        }
    }
    
    # 6. 对比说明
    Write-Host ""
    Write-Host "==================================" -ForegroundColor Cyan
    Write-Host "📝 改进说明" -ForegroundColor Yellow
    Write-Host "==================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "改进前:" -ForegroundColor Red
    Write-Host "  - 使用文章标题作为竞品名" -ForegroundColor Gray
    Write-Host "  - 结果: '中国AI？美国AI？'、'生成式AI与不同类型的AI'" -ForegroundColor Gray
    Write-Host ""
    Write-Host "改进后:" -ForegroundColor Green
    Write-Host "  - 从URL提取品牌名" -ForegroundColor Gray
    Write-Host "  - https://www.notion.so → Notion" -ForegroundColor Gray
    Write-Host "  - https://www.canva.com → Canva" -ForegroundColor Gray
    Write-Host ""
    Write-Host "⚠️  注意: 仍需要更准确的搜索词" -ForegroundColor Yellow
    Write-Host "   建议使用: 'Notion'、'Jasper'等明确产品名" -ForegroundColor Gray
    
} catch {
    Write-Host "   ❌ 测试失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "测试完成" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "💡 提示:" -ForegroundColor Yellow
Write-Host "  - 查看详细说明: cat MARKET_PARAM_EXPLAINED.md" -ForegroundColor Gray
Write-Host "  - 服务仍在后台运行" -ForegroundColor Gray
Write-Host "  - 停止服务: Stop-Job $($job.Id)" -ForegroundColor Gray
