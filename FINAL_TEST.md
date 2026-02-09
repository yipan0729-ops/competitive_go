# 🎯 现在进行最终测试

## ✅ 准备就绪

- ✅ 服务已重新编译
- ✅ 旧进程已停止
- ✅ 新程序正在运行
- ✅ 修复代码已加载

---

## 🚀 启动测试任务

请在PowerShell中运行以下完整测试：

```powershell
# =====================================
# 最终测试脚本 - 任务6
# =====================================

Write-Host "🚀 开始最终测试..." -ForegroundColor Cyan
Write-Host ""

# 启动新任务
$body = @{
    topic = "markdown编辑器"
    competitor_count = 2
    depth = "quick"
} | ConvertTo-Json

Write-Host "1️⃣ 发送请求..." -ForegroundColor Yellow
$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

$result = $response.Content | ConvertFrom-Json
$taskId = $result.task_id

Write-Host "✅ 任务已启动！" -ForegroundColor Green
Write-Host "   任务ID: $taskId" -ForegroundColor White
Write-Host "   预计2-3分钟完成" -ForegroundColor White
Write-Host ""

# 监控进度
Write-Host "2️⃣ 监控进度..." -ForegroundColor Yellow
Write-Host ""

$completed = $false
for ($i = 1; $i -le 30; $i++) {
    Start-Sleep -Seconds 10
    
    $status = Invoke-WebRequest `
        -Uri "http://localhost:8080/api/discover/status/$taskId" `
        -UseBasicParsing |
        ConvertFrom-Json
    
    $emoji = switch ($status.status) {
        "discovering" { "🔎" }
        "searching_sources" { "🔍" }
        "crawling" { "🕷️" }
        "analyzing" { "🤖" }
        "generating_report" { "📝" }
        "completed" { "✅" }
        "failed" { "❌" }
        default { "⏳" }
    }
    
    Write-Host "[$i] $emoji $($status.progress)% - $($status.status)" -ForegroundColor Cyan
    
    # 关键阶段提示
    if ($status.status -eq "crawling") {
        Write-Host "    >>> 正在爬取网页..." -ForegroundColor Gray
    }
    if ($status.status -eq "analyzing") {
        Write-Host "    >>> 正在调用Ollama分析..." -ForegroundColor Yellow
    }
    if ($status.status -eq "generating_report") {
        Write-Host "    >>> 正在生成报告..." -ForegroundColor Gray
    }
    
    if ($status.status -eq "completed") {
        $completed = $true
        Write-Host ""
        Write-Host "=" * 50 -ForegroundColor Green
        Write-Host "🎉 任务完成！" -ForegroundColor Green
        Write-Host "=" * 50 -ForegroundColor Green
        Write-Host ""
        break
    }
    
    if ($status.status -eq "failed") {
        Write-Host ""
        Write-Host "❌ 任务失败" -ForegroundColor Red
        break
    }
}

if (-not $completed) {
    Write-Host ""
    Write-Host "⏱️ 监控超时，请手动检查" -ForegroundColor Yellow
}

# 获取最终结果
Write-Host "3️⃣ 检查结果..." -ForegroundColor Yellow
Write-Host ""

$finalStatus = Invoke-WebRequest `
    -Uri "http://localhost:8080/api/discover/status/$taskId" `
    -UseBasicParsing |
    ConvertFrom-Json

Write-Host "📊 最终统计：" -ForegroundColor Cyan
Write-Host "   竞品数: $($finalStatus.competitors_found)" -ForegroundColor White
Write-Host "   数据源: $($finalStatus.data_sources_found)" -ForegroundColor White
Write-Host "   爬取URL: $($finalStatus.result.urls_crawled)" -ForegroundColor White
Write-Host "   已分析: $($finalStatus.result.analyzed_count)" -ForegroundColor $(if ($finalStatus.result.analyzed_count -gt 0) { "Green" } else { "Red" })
Write-Host "   报告路径: $($finalStatus.result.report_path)" -ForegroundColor $(if ($finalStatus.result.report_path) { "Green" } else { "Red" })
Write-Host ""

# 验证修复
Write-Host "4️⃣ 验证修复..." -ForegroundColor Yellow
Write-Host ""

$allGood = $true

if ($finalStatus.result.analyzed_count -gt 0) {
    Write-Host "✅ AI分析成功 (调用了Ollama)" -ForegroundColor Green
} else {
    Write-Host "❌ AI分析失败 (未调用Ollama)" -ForegroundColor Red
    $allGood = $false
}

if ($finalStatus.result.report_path -and $finalStatus.result.report_path -ne "") {
    Write-Host "✅ 报告生成成功" -ForegroundColor Green
    
    # 检查文件是否存在
    if (Test-Path $finalStatus.result.report_path) {
        Write-Host "✅ 报告文件存在" -ForegroundColor Green
        Write-Host "   文件: $($finalStatus.result.report_path)" -ForegroundColor White
        
        # 打开报告
        Write-Host ""
        Write-Host "📖 打开报告..." -ForegroundColor Cyan
        Start-Process $finalStatus.result.report_path
    } else {
        Write-Host "⚠️ 报告文件不存在" -ForegroundColor Yellow
        $allGood = $false
    }
} else {
    Write-Host "❌ 报告生成失败" -ForegroundColor Red
    $allGood = $false
}

# 检查目录
Write-Host ""
Write-Host "5️⃣ 检查文件..." -ForegroundColor Yellow
Write-Host ""

if (Test-Path "D:\Code\Competitive_go\reports") {
    $reportFiles = Get-ChildItem "D:\Code\Competitive_go\reports" -File
    Write-Host "   reports/ 目录: $($reportFiles.Count) 个文件" -ForegroundColor White
}

if (Test-Path "D:\Code\Competitive_go\storage\crawled") {
    $crawledDirs = Get-ChildItem "D:\Code\Competitive_go\storage\crawled" -Directory
    Write-Host "   storage/crawled/ 目录: $($crawledDirs.Count) 个竞品文件夹" -ForegroundColor White
}

# 最终结论
Write-Host ""
Write-Host "=" * 50 -ForegroundColor Cyan
if ($allGood) {
    Write-Host "🎊 测试成功！所有功能正常！" -ForegroundColor Green
    Write-Host ""
    Write-Host "修复已生效：" -ForegroundColor Green
    Write-Host "  ✅ 爬虫工作正常" -ForegroundColor Green
    Write-Host "  ✅ Ollama被成功调用" -ForegroundColor Green
    Write-Host "  ✅ 报告成功生成" -ForegroundColor Green
} else {
    Write-Host "⚠️ 测试未完全成功" -ForegroundColor Yellow
    Write-Host "   请检查上面的错误信息" -ForegroundColor Yellow
}
Write-Host "=" * 50 -ForegroundColor Cyan
Write-Host ""

# 保存结果
$finalStatus | ConvertTo-Json -Depth 10 | Out-File -FilePath "final_test_result.json" -Encoding UTF8
Write-Host "📄 完整结果已保存到: final_test_result.json" -ForegroundColor Gray
```

---

## 📋 预期结果

如果修复成功，您应该看到：

```
🎉 任务完成！

📊 最终统计：
   竞品数: 2
   数据源: 6
   爬取URL: 6
   已分析: 2        ← 这个必须 > 0
   报告路径: reports/xxx.md  ← 必须有路径

✅ AI分析成功 (调用了Ollama)
✅ 报告生成成功
✅ 报告文件存在

🎊 测试成功！所有功能正常！
```

---

## ⚠️ 如果还是失败

如果 `analyzed_count` 仍然是 0，请：

1. **查看服务日志**
   - 检查是否看到 `[自动化] 使用默认配置：启用爬取、分析和报告生成`
   - 如果没有这行，说明代码没有正确加载

2. **检查Ollama**
   ```powershell
   ollama list
   ollama run deepseek-r1:8b "test"
   ```

3. **查看完整日志**
   - 查看运行 `competitive-analyzer-new.exe` 的终端输出

---

## 🎯 现在就测试

复制上面的完整脚本，在PowerShell中运行！

这次应该能看到完整的结果了！🚀
