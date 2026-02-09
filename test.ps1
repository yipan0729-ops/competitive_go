# 保存为 continue-from-step4.ps1

Write-Host "继续竞品ID=4的后续流程" -ForegroundColor Cyan
Write-Host ""

# 1. AI分析
Write-Host "1️⃣ 执行AI分析..." -ForegroundColor Yellow
$body = @{
    competitor_id = 4
    market_context = "中国AI工具市场"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -TimeoutSec 600 `
    -OutFile "analysis-4.json"

Write-Host "✅ 分析完成" -ForegroundColor Green
Write-Host ""

# 2. 生成报告
Write-Host "2️⃣ 生成报告..." -ForegroundColor Yellow
$body = @{
    competitor_ids = @(4)
    topic = "AI创作"
    market = "中国"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/report/generate `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -OutFile "report-4.json"

$report = cat report-4.json | ConvertFrom-Json
Write-Host "✅ 报告已生成: $($report.report_path)" -ForegroundColor Green

# 3. 显示报告
Write-Host ""
Write-Host "📄 报告内容:" -ForegroundColor Cyan
cat $report.report_path# 保存为 continue-from-step4.ps1

Write-Host "继续竞品ID=4的后续流程" -ForegroundColor Cyan
Write-Host ""

# 1. AI分析
Write-Host "1️⃣ 执行AI分析..." -ForegroundColor Yellow
$body = @{
    competitor_id = 4
    market_context = "中国AI工具市场"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -TimeoutSec 600 `
    -OutFile "analysis-4.json"

Write-Host "✅ 分析完成" -ForegroundColor Green
Write-Host ""

# 2. 生成报告
Write-Host "2️⃣ 生成报告..." -ForegroundColor Yellow
$body = @{
    competitor_ids = @(4)
    topic = "AI创作"
    market = "中国"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/report/generate `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -OutFile "report-4.json"

$report = cat report-4.json | ConvertFrom-Json
Write-Host "✅ 报告已生成: $($report.report_path)" -ForegroundColor Green

# 3. 显示报告
Write-Host ""
Write-Host "📄 报告内容:" -ForegroundColor Cyan
cat $report.report_path