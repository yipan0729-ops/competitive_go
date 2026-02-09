# 🚀 快速开始 - 全流程自动化竞品分析

> **5分钟从输入主题到获得完整报告！**

---

## 📋 前置准备

### 1. 确认服务运行

```powershell
# 检查服务状态
Invoke-WebRequest -Uri http://localhost:8080/health -UseBasicParsing
```

✅ 看到 `"status":"ok"` 表示服务正常

### 2. 确认Ollama运行

```powershell
# 测试Ollama
ollama list
```

✅ 确保已下载 `deepseek-r1:8b` 或 `qwen2.5:7b`

---

## 🎯 方式1：一键自动化（推荐）

### 最简单的方式 - 只需一个API调用！

```powershell
# 定义分析主题
$body = @{
    topic = "项目管理工具"          # 你要分析的主题
    market = "中国"                # 可选：目标市场
    competitor_count = 5          # 可选：竞品数量
} | ConvertTo-Json

# 启动全流程自动化
$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

# 获取任务ID
$result = $response.Content | ConvertFrom-Json
$taskId = $result.task_id

Write-Host "✅ 任务已启动！ID: $taskId"
Write-Host "⏳ 预计5分钟完成..."
Write-Host ""
```

### 监控进度

```powershell
# 自动轮询进度
while ($true) {
    $status = Invoke-WebRequest `
        -Uri "http://localhost:8080/api/discover/status/$taskId" `
        -UseBasicParsing | 
        ConvertFrom-Json
    
    $progress = $status.progress
    $currentStatus = $status.status
    
    # 显示进度条
    $bar = "[" + ("=" * [math]::Floor($progress / 5)) + (" " * (20 - [math]::Floor($progress / 5))) + "]"
    Write-Host "`r$bar $progress% - $currentStatus" -NoNewline
    
    if ($currentStatus -eq "completed") {
        Write-Host "`n`n🎉 分析完成！"
        
        # 显示结果
        Write-Host "`n📊 结果摘要:"
        Write-Host "  竞品数量: $($status.result.competitors.Count)"
        Write-Host "  爬取URL: $($status.result.urls_crawled)"
        Write-Host "  已分析: $($status.result.analyzed_count)"
        Write-Host "  报告位置: $($status.result.report_path)"
        
        # 打开报告
        if ($status.result.report_path) {
            Write-Host "`n📖 打开报告..."
            Start-Process $status.result.report_path
        }
        
        break
    }
    
    if ($currentStatus -eq "failed") {
        Write-Host "`n❌ 任务失败"
        break
    }
    
    Start-Sleep -Seconds 5
}
```

### 完整一键脚本

将以上代码保存为 `auto-analysis.ps1`：

```powershell
# =====================================
# 一键竞品分析脚本
# =====================================

param(
    [string]$Topic = "项目管理工具",
    [string]$Market = "中国",
    [int]$Count = 5
)

Write-Host "🚀 开始竞品分析..." -ForegroundColor Cyan
Write-Host "主题: $Topic" -ForegroundColor Yellow
Write-Host ""

# 启动分析
$body = @{
    topic = $Topic
    market = $Market
    competitor_count = $Count
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest `
        -Uri http://localhost:8080/api/auto/analysis `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -UseBasicParsing
    
    $result = $response.Content | ConvertFrom-Json
    $taskId = $result.task_id
    
    Write-Host "✅ 任务 #$taskId 已启动！" -ForegroundColor Green
    Write-Host ""
    
    # 监控进度
    $completed = $false
    while (-not $completed) {
        Start-Sleep -Seconds 5
        
        $status = Invoke-WebRequest `
            -Uri "http://localhost:8080/api/discover/status/$taskId" `
            -UseBasicParsing | 
            ConvertFrom-Json
        
        $progress = $status.progress
        $bar = "[" + ("=" * [math]::Floor($progress / 5)) + (" " * (20 - [math]::Floor($progress / 5))) + "]"
        Write-Host "`r$bar $progress% - $($status.status)" -NoNewline
        
        if ($status.status -eq "completed") {
            $completed = $true
            Write-Host "`n"
            Write-Host "🎉 分析完成！" -ForegroundColor Green
            Write-Host ""
            Write-Host "📊 结果:" -ForegroundColor Cyan
            Write-Host "  竞品: $($status.result.competitors.Count) 个"
            Write-Host "  URL: $($status.result.urls_crawled) 个"
            Write-Host "  已分析: $($status.result.analyzed_count) 个"
            Write-Host "  报告: $($status.result.report_path)"
            Write-Host ""
            
            if ($status.result.report_path) {
                Write-Host "📖 打开报告..." -ForegroundColor Yellow
                Start-Process $status.result.report_path
            }
        }
        
        if ($status.status -eq "failed") {
            Write-Host "`n❌ 失败" -ForegroundColor Red
            break
        }
    }
    
} catch {
    Write-Host "❌ 错误: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n✨ 完成！" -ForegroundColor Cyan
```

**使用方式**:

```powershell
# 方式1：使用默认参数
.\auto-analysis.ps1

# 方式2：指定主题
.\auto-analysis.ps1 -Topic "在线设计工具"

# 方式3：完整参数
.\auto-analysis.ps1 -Topic "CRM系统" -Market "全球" -Count 8
```

---

## 🎨 方式2：分步执行（更多控制）

### 步骤1：发现竞品

```powershell
$body = @{
    topic = "canvas"
    market = "中国"
    competitor_count = 5
} | ConvertTo-Json

$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/discover/search `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

$taskId = ($response.Content | ConvertFrom-Json).task_id
```

### 步骤2：等待发现完成

```powershell
# 等待30秒
Start-Sleep -Seconds 30

# 查看结果
$status = Invoke-WebRequest `
    -Uri "http://localhost:8080/api/discover/status/$taskId" `
    -UseBasicParsing |
    ConvertFrom-Json

# 显示发现的竞品
$status.result.competitors
```

### 步骤3：批量爬取

```powershell
# 准备URL列表（从发现结果中获取）
$urls = @(
    @{ url="https://www.canva.com"; competitor="Canva"; source_type="官网" },
    @{ url="https://www.adobe.com/express"; competitor="Adobe Express"; source_type="官网" },
    @{ url="https://www.figma.com"; competitor="Figma"; source_type="官网" }
)

$crawlBody = @{
    urls = $urls
    concurrent = 3
} | ConvertTo-Json -Depth 5

Invoke-WebRequest `
    -Uri http://localhost:8080/api/crawl/batch `
    -Method POST `
    -Body $crawlBody `
    -ContentType "application/json" `
    -UseBasicParsing

# 等待爬取完成（约2分钟）
Start-Sleep -Seconds 120
```

### 步骤4：AI分析

```powershell
# 获取竞品列表
$competitors = Invoke-WebRequest `
    -Uri "http://localhost:8080/api/competitors" `
    -UseBasicParsing |
    ConvertFrom-Json

# 分析每个竞品
foreach ($competitor in $competitors.competitors) {
    Write-Host "分析: $($competitor.name)"
    
    $analyzeBody = @{
        competitor_id = $competitor.id
        market_context = "中国在线设计工具市场"
    } | ConvertTo-Json
    
    Invoke-WebRequest `
        -Uri http://localhost:8080/api/analyze/competitor `
        -Method POST `
        -Body $analyzeBody `
        -ContentType "application/json" `
        -UseBasicParsing
    
    Start-Sleep -Seconds 10  # 避免频繁调用LLM
}
```

### 步骤5：生成报告

```powershell
# 收集所有竞品ID
$competitorIds = $competitors.competitors | Select-Object -ExpandProperty id

$reportBody = @{
    competitor_ids = $competitorIds
    topic = "在线设计工具"
    report_name = "Canvas竞品分析"
} | ConvertTo-Json

$report = Invoke-WebRequest `
    -Uri http://localhost:8080/api/report/generate `
    -Method POST `
    -Body $reportBody `
    -ContentType "application/json" `
    -UseBasicParsing |
    ConvertFrom-Json

# 打开报告
Start-Process $report.report_path
```

---

## 📊 实战示例

### 示例1：分析Canva竞品

```powershell
POST /api/auto/analysis
{
  "topic": "canvas",
  "market": "全球",
  "competitor_count": 5
}
```

**5分钟后获得**:
- ✅ 5个竞品（Canva, Adobe Express, Figma, Visme, VistaCreate）
- ✅ 15个爬取的网页
- ✅ 完整产品信息提取
- ✅ SWOT分析
- ✅ 专业报告

### 示例2：分析项目管理工具

```powershell
POST /api/auto/analysis
{
  "topic": "项目管理工具",
  "market": "中国",
  "competitor_count": 8,
  "depth": "deep"  # 深度搜索
}
```

### 示例3：快速探索（只发现不分析）

```powershell
POST /api/auto/analysis
{
  "topic": "AI写作工具",
  "auto_crawl": false,
  "auto_analyze": false,
  "generate_report": false
}
```

只获得竞品列表和数据源，不执行爬取和分析。

---

## 🔍 查看结果

### 查看报告

```powershell
# 报告保存在 reports/ 目录
ls reports/

# 打开最新报告
$latest = Get-ChildItem reports/ | Sort-Object LastWriteTime -Descending | Select-Object -First 1
Start-Process $latest.FullName
```

### 查看爬取的内容

```powershell
# 内容保存在 storage/crawled/ 目录
ls storage/crawled/

# 查看某个竞品的内容
ls storage/crawled/Canva/
```

### 查看数据库

```powershell
# 使用SQLite查看
sqlite3 data/competitive.db

# 查询竞品
SELECT * FROM competitors;

# 查询分析结果
SELECT * FROM parsed_data;
```

---

## ⚡ 性能优化

### 提升速度

1. **增加并发数**
```powershell
# 批量爬取时设置更高并发
{ "concurrent": 5 }  # 最大10
```

2. **减少分析深度**
```powershell
# 快速模式
{ "depth": "quick", "competitor_count": 3 }
```

3. **使用更快的LLM模型**
```env
# .env文件
LLM_MODEL=qwen2.5:7b  # 比deepseek-r1快
```

---

## 🐛 常见问题

### Q: 任务卡在某个进度？

**A**: 查看服务日志

```powershell
# 查看最新日志
tail -f logs/app.log  # Linux/Mac
Get-Content logs/app.log -Tail 50  # Windows
```

### Q: Ollama报错？

**A**: 确认Ollama服务运行

```powershell
# 重启Ollama
ollama serve

# 测试模型
ollama run deepseek-r1:8b "hello"
```

### Q: 爬取失败？

**A**: 检查网络和API配置

```powershell
# 测试URL可访问
curl -I https://www.canva.com

# 检查API Key（如果使用Firecrawl）
cat .env | grep FIRECRAWL
```

---

## 📚 下一步

- 📖 查看 [API完整文档](./API.md)
- 🔧 查看 [配置指南](./HOW_TO_RUN.md)
- 📝 查看 [更新日志](./CHANGELOG.md)
- 🤖 查看 [Ollama配置](./OLLAMA_MODEL_GUIDE.md)

---

## 🎉 开始使用

```powershell
# 一条命令开始你的第一次分析
$body = @{topic="你感兴趣的主题"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis `
    -Method POST -Body $body -ContentType "application/json" -UseBasicParsing
```

**5分钟后，享受完整的竞品分析报告！** 🚀

---

**最后更新**: 2026-02-09  
**版本**: v2.0.0
