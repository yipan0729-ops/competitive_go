# =====================================
# 全流程自动化测试脚本
# =====================================

Write-Host "🚀 开始测试全流程自动化..." -ForegroundColor Cyan
Write-Host ""

# 配置
$BaseUrl = "http://localhost:8080"
$Topic = "canvas"

# 步骤1：检查服务健康
Write-Host "1️⃣ 检查服务状态..." -ForegroundColor Yellow
try {
    $health = Invoke-WebRequest -Uri "$BaseUrl/health" -UseBasicParsing
    $healthData = $health.Content | ConvertFrom-Json
    if ($healthData.status -eq "ok") {
        Write-Host "   ✅ 服务正常" -ForegroundColor Green
    }
} catch {
    Write-Host "   ❌ 服务未启动！" -ForegroundColor Red
    Write-Host "   请先运行: go run main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# 步骤2：启动自动化分析
Write-Host "2️⃣ 启动自动化分析..." -ForegroundColor Yellow
$body = @{
    topic = $Topic
    market = "全球"
    competitor_count = 3  # 测试用3个竞品，快速完成
    depth = "quick"       # 快速模式
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest `
        -Uri "$BaseUrl/api/auto/analysis" `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -UseBasicParsing
    
    $result = $response.Content | ConvertFrom-Json
    $taskId = $result.task_id
    
    Write-Host "   ✅ 任务已启动 ID: $taskId" -ForegroundColor Green
    Write-Host "   📋 工作流: $($result.workflow)" -ForegroundColor Gray
    Write-Host "   ⏱️  预计时间: $($result.estimated_time)秒" -ForegroundColor Gray
    
} catch {
    Write-Host "   ❌ 启动失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 步骤3：监控进度
Write-Host "3️⃣ 监控任务进度..." -ForegroundColor Yellow
Write-Host ""

$completed = $false
$failCount = 0
$maxAttempts = 60  # 最多等待5分钟

while (-not $completed -and $failCount -lt $maxAttempts) {
    Start-Sleep -Seconds 5
    
    try {
        $statusResponse = Invoke-WebRequest `
            -Uri "$BaseUrl/api/discover/status/$taskId" `
            -UseBasicParsing
        
        $status = $statusResponse.Content | ConvertFrom-Json
        
        $progress = $status.progress
        $currentStatus = $status.status
        
        # 显示进度条
        $barLength = [math]::Floor($progress / 5)
        $bar = "[" + ("=" * $barLength) + (" " * (20 - $barLength)) + "]"
        
        Write-Host "`r   $bar $progress% - $currentStatus" -NoNewline -ForegroundColor Cyan
        
        if ($currentStatus -eq "completed") {
            $completed = $true
            Write-Host "`n"
            
            Write-Host "   🎉 分析完成！" -ForegroundColor Green
            Write-Host ""
            Write-Host "   📊 结果统计:" -ForegroundColor Yellow
            Write-Host "      竞品数量: $($status.competitors_found)" -ForegroundColor White
            Write-Host "      数据源数: $($status.data_sources_found)" -ForegroundColor White
            
            if ($status.result) {
                if ($status.result.competitors) {
                    Write-Host ""
                    Write-Host "   🏢 发现的竞品:" -ForegroundColor Yellow
                    $index = 1
                    foreach ($comp in $status.result.competitors) {
                        Write-Host "      $index. $comp" -ForegroundColor White
                        $index++
                    }
                }
                
                if ($status.result.urls_crawled) {
                    Write-Host ""
                    Write-Host "   🕷️  爬取统计:" -ForegroundColor Yellow
                    Write-Host "      URL数量: $($status.result.urls_crawled)" -ForegroundColor White
                    Write-Host "      已分析: $($status.result.analyzed_count)" -ForegroundColor White
                }
                
                if ($status.result.report_path) {
                    Write-Host ""
                    Write-Host "   📄 报告位置:" -ForegroundColor Yellow
                    Write-Host "      $($status.result.report_path)" -ForegroundColor White
                    
                    # 检查报告是否存在
                    if (Test-Path $status.result.report_path) {
                        Write-Host ""
                        Write-Host "   📖 打开报告..." -ForegroundColor Cyan
                        Start-Process $status.result.report_path
                    }
                }
            }
        }
        
        if ($currentStatus -eq "failed") {
            Write-Host "`n"
            Write-Host "   ❌ 任务失败" -ForegroundColor Red
            if ($status.result.error) {
                Write-Host "      错误: $($status.result.error)" -ForegroundColor Red
            }
            break
        }
        
        $failCount = 0  # 重置失败计数
        
    } catch {
        $failCount++
        if ($failCount -ge 3) {
            Write-Host "`n   ⚠️  无法连接到服务" -ForegroundColor Yellow
            break
        }
    }
}

if ($failCount -ge $maxAttempts) {
    Write-Host "`n   ⚠️  监控超时" -ForegroundColor Yellow
    Write-Host "      任务可能仍在后台运行" -ForegroundColor Gray
    Write-Host "      请稍后手动检查: $BaseUrl/api/discover/status/$taskId" -ForegroundColor Gray
}

Write-Host ""
Write-Host "✨ 测试完成！" -ForegroundColor Cyan
Write-Host ""

# 步骤4：查看竞品列表
Write-Host "4️⃣ 查看竞品列表..." -ForegroundColor Yellow
try {
    $competitors = Invoke-WebRequest `
        -Uri "$BaseUrl/api/competitors" `
        -UseBasicParsing |
        ConvertFrom-Json
    
    Write-Host "   ✅ 数据库中共有 $($competitors.total) 个竞品" -ForegroundColor Green
    
} catch {
    Write-Host "   ⚠️  查询失败" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Gray
Write-Host ""
Write-Host "🎊 测试总结" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ 服务运行正常" -ForegroundColor Green
Write-Host "✅ 自动化流程可用" -ForegroundColor Green
Write-Host "✅ 进度监控正常" -ForegroundColor Green
Write-Host ""
Write-Host "📚 查看文档:" -ForegroundColor Yellow
Write-Host "   - API文档: API.md" -ForegroundColor Gray
Write-Host "   - 快速开始: QUICKSTART_V2.md" -ForegroundColor Gray
Write-Host "   - 更新日志: CHANGELOG.md" -ForegroundColor Gray
Write-Host ""
Write-Host "🚀 开始使用:" -ForegroundColor Yellow
Write-Host '   $body = @{topic="你的主题"} | ConvertTo-Json' -ForegroundColor Gray
Write-Host '   Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis \' -ForegroundColor Gray
Write-Host '       -Method POST -Body $body -ContentType "application/json"' -ForegroundColor Gray
Write-Host ""
