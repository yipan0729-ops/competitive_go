# =====================================
# AI分析单独测试脚本
# =====================================

Write-Host "🤖 AI分析独立测试" -ForegroundColor Cyan
Write-Host ""

# 步骤1：检查Ollama
Write-Host "1️⃣ 检查Ollama状态..." -ForegroundColor Yellow
try {
    $ollamaVersion = Invoke-WebRequest -Uri http://localhost:11434/api/version -UseBasicParsing -TimeoutSec 2
    Write-Host "   ✅ Ollama运行正常" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Ollama未运行！" -ForegroundColor Red
    Write-Host "   请运行: ollama serve" -ForegroundColor Yellow
    exit 1
}

# 步骤2：检查竞品
Write-Host ""
Write-Host "2️⃣ 查找竞品..." -ForegroundColor Yellow

$competitors = Invoke-WebRequest `
    -Uri http://localhost:8080/api/competitors `
    -UseBasicParsing |
    ConvertFrom-Json

if ($competitors.total -eq 0) {
    Write-Host "   ❌ 没有竞品数据！" -ForegroundColor Red
    Write-Host ""
    Write-Host "   请先运行爬取任务：" -ForegroundColor Yellow
    Write-Host '   $body = @{topic="测试"} | ConvertTo-Json' -ForegroundColor Gray
    Write-Host '   Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis \' -ForegroundColor Gray
    Write-Host '       -Method POST -Body $body -ContentType "application/json"' -ForegroundColor Gray
    exit 1
}

Write-Host "   ✅ 找到 $($competitors.total) 个竞品" -ForegroundColor Green
Write-Host ""

# 显示竞品列表
Write-Host "   竞品列表：" -ForegroundColor Cyan
$index = 1
foreach ($comp in $competitors.competitors) {
    Write-Host "   $index. ID=$($comp.id) - $($comp.name)" -ForegroundColor White
    $index++
}

# 步骤3：选择竞品
Write-Host ""
$firstCompetitor = $competitors.competitors[0]
Write-Host "3️⃣ 分析竞品: $($firstCompetitor.name)" -ForegroundColor Yellow
Write-Host "   ID: $($firstCompetitor.id)" -ForegroundColor Gray
Write-Host ""

# 步骤4：执行AI分析
$body = @{
    competitor_id = $firstCompetitor.id
    market_context = "中国市场，竞争激烈"
} | ConvertTo-Json

Write-Host "4️⃣ 调用Ollama分析..." -ForegroundColor Yellow
Write-Host "   ⏳ 这可能需要1-10分钟，取决于模型..." -ForegroundColor Gray
Write-Host "   💡 提示: qwen2.5快，deepseek-r1慢" -ForegroundColor Gray
Write-Host ""

$startTime = Get-Date

try {
    $response = Invoke-WebRequest `
        -Uri http://localhost:8080/api/analyze/competitor `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -UseBasicParsing `
        -TimeoutSec 600  # 10分钟超时
    
    $elapsed = (Get-Date) - $startTime
    
    Write-Host ""
    Write-Host "=" * 60 -ForegroundColor Green
    Write-Host "✅ AI分析成功！" -ForegroundColor Green
    Write-Host "=" * 60 -ForegroundColor Green
    Write-Host ""
    Write-Host "⏱️ 耗时: $([math]::Round($elapsed.TotalSeconds, 1)) 秒" -ForegroundColor White
    Write-Host ""
    
    $result = $response.Content | ConvertFrom-Json
    
    # 显示产品信息
    if ($result.product_info) {
        Write-Host "📊 产品信息:" -ForegroundColor Cyan
        Write-Host "   产品名: $($result.product_info.product_name)" -ForegroundColor White
        Write-Host "   公司: $($result.product_info.company)" -ForegroundColor White
        Write-Host "   定位: $($result.product_info.tagline)" -ForegroundColor White
        
        if ($result.product_info.target_users) {
            Write-Host "   目标用户: $($result.product_info.target_users -join ', ')" -ForegroundColor White
        }
        
        if ($result.product_info.core_features) {
            Write-Host "   核心功能: $($result.product_info.core_features.Count) 个" -ForegroundColor White
        }
        
        Write-Host ""
    }
    
    # 显示SWOT
    if ($result.swot_analysis) {
        Write-Host "📈 SWOT分析:" -ForegroundColor Cyan
        Write-Host "   优势 (S): $($result.swot_analysis.strengths.Count) 个" -ForegroundColor Green
        Write-Host "   劣势 (W): $($result.swot_analysis.weaknesses.Count) 个" -ForegroundColor Yellow
        Write-Host "   机会 (O): $($result.swot_analysis.opportunities.Count) 个" -ForegroundColor Cyan
        Write-Host "   威胁 (T): $($result.swot_analysis.threats.Count) 个" -ForegroundColor Red
        Write-Host ""
        
        # 显示详细优势
        if ($result.swot_analysis.strengths.Count -gt 0) {
            Write-Host "   优势示例:" -ForegroundColor Green
            $result.swot_analysis.strengths | Select-Object -First 2 | ForEach-Object {
                Write-Host "   • $($_.point)" -ForegroundColor White
            }
            Write-Host ""
        }
    }
    
    # 保存结果
    $result | ConvertTo-Json -Depth 10 | Out-File "ai_analysis_result.json" -Encoding UTF8
    Write-Host "💾 完整结果已保存到: ai_analysis_result.json" -ForegroundColor Gray
    Write-Host ""
    
    Write-Host "🎉 测试成功！Ollama被成功调用！" -ForegroundColor Green
    
} catch {
    $elapsed = (Get-Date) - $startTime
    
    Write-Host ""
    Write-Host "❌ AI分析失败" -ForegroundColor Red
    Write-Host "   耗时: $([math]::Round($elapsed.TotalSeconds, 1)) 秒" -ForegroundColor Gray
    Write-Host "   错误: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    
    Write-Host "💡 故障排查：" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "1. 检查Ollama是否运行：" -ForegroundColor White
    Write-Host "   ollama list" -ForegroundColor Gray
    Write-Host "   ollama run qwen2.5:7b `"test`"" -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. 切换到更快的模型：" -ForegroundColor White
    Write-Host "   编辑 .env: LLM_MODEL=qwen2.5:7b" -ForegroundColor Gray
    Write-Host "   重启服务" -ForegroundColor Gray
    Write-Host ""
    Write-Host "3. 检查是否有爬取数据：" -ForegroundColor White
    Write-Host "   ls storage\crawled\" -ForegroundColor Gray
}

Write-Host ""
Write-Host "✨ 脚本执行完毕" -ForegroundColor Cyan
