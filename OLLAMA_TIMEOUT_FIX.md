# 🔧 Ollama连接问题解决方案

## ❌ 问题：Ollama超时

```
请求Ollama失败: context deadline exceeded
```

---

## ✅ 解决方案

### 方案1：切换到更快的模型（推荐）⭐

DeepSeek-R1是推理模型，响应很慢（可能需要5-10分钟）。

#### 下载Qwen2.5（更快）

```powershell
# 下载Qwen2.5（更快，3-5秒响应）
ollama pull qwen2.5:7b

# 测试
ollama run qwen2.5:7b "hello"
```

#### 修改配置

编辑 `.env` 文件：

```env
# 从这个：
LLM_MODEL=deepseek-r1:8b

# 改为：
LLM_MODEL=qwen2.5:7b
```

#### 重启服务

```powershell
# 停止旧服务
taskkill /F /IM competitive-analyzer-v2.exe

# 启动新服务
.\competitive-analyzer-v2.exe
```

---

### 方案2：确保Ollama运行

```powershell
# 1. 检查Ollama
ollama list

# 2. 如果没有运行，启动Ollama
ollama serve

# 3. 测试连接
ollama run deepseek-r1:8b "test"
```

---

### 方案3：增加超时时间（已完成）

代码已修改，现在超时时间：
- 云端API：600秒（10分钟）
- Ollama：600秒（10分钟）

---

## 📡 单独执行AI分析的API接口

### 接口：`POST /api/analyze/competitor`

**功能**：对已爬取的竞品进行AI分析

### 使用步骤

#### 1. 查看现有竞品

```powershell
# 获取竞品列表
$competitors = Invoke-WebRequest `
    -Uri http://localhost:8080/api/competitors `
    -UseBasicParsing |
    ConvertFrom-Json

# 显示
$competitors.competitors | Format-Table id, name, created_at
```

#### 2. 对单个竞品执行分析

```powershell
# 分析竞品ID=1
$body = @{
    competitor_id = 1
    market_context = "中国AI工具市场"
} | ConvertTo-Json

$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

# 查看结果
$result = $response.Content | ConvertFrom-Json
$result | ConvertTo-Json -Depth 5
```

#### 3. 完整示例

```powershell
# =====================================
# 单独执行AI分析脚本
# =====================================

Write-Host "🤖 开始AI分析..." -ForegroundColor Cyan
Write-Host ""

# 1. 获取竞品列表
Write-Host "1️⃣ 获取竞品列表..." -ForegroundColor Yellow
$competitors = Invoke-WebRequest `
    -Uri http://localhost:8080/api/competitors `
    -UseBasicParsing |
    ConvertFrom-Json

if ($competitors.total -eq 0) {
    Write-Host "❌ 没有竞品数据，请先运行爬取任务" -ForegroundColor Red
    exit
}

Write-Host "   找到 $($competitors.total) 个竞品" -ForegroundColor White
Write-Host ""

# 2. 选择第一个竞品
$competitor = $competitors.competitors[0]
Write-Host "2️⃣ 分析竞品: $($competitor.name)" -ForegroundColor Yellow
Write-Host "   ID: $($competitor.id)" -ForegroundColor Gray
Write-Host ""

# 3. 执行AI分析
$body = @{
    competitor_id = $competitor.id
    market_context = "中国市场"
} | ConvertTo-Json

Write-Host "3️⃣ 调用Ollama分析（可能需要1-5分钟）..." -ForegroundColor Yellow
Write-Host "   如果使用deepseek-r1，可能需要5-10分钟" -ForegroundColor Gray
Write-Host ""

try {
    $response = Invoke-WebRequest `
        -Uri http://localhost:8080/api/analyze/competitor `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -UseBasicParsing `
        -TimeoutSec 600  # 10分钟超时
    
    $result = $response.Content | ConvertFrom-Json
    
    Write-Host "✅ 分析成功！" -ForegroundColor Green
    Write-Host ""
    
    # 显示产品信息
    if ($result.product_info) {
        Write-Host "📊 产品信息:" -ForegroundColor Cyan
        Write-Host "   产品名: $($result.product_info.product_name)" -ForegroundColor White
        Write-Host "   公司: $($result.product_info.company)" -ForegroundColor White
        Write-Host "   定位: $($result.product_info.tagline)" -ForegroundColor White
        Write-Host ""
    }
    
    # 显示SWOT
    if ($result.swot_analysis) {
        Write-Host "📈 SWOT分析:" -ForegroundColor Cyan
        Write-Host "   优势: $($result.swot_analysis.strengths.Count) 个" -ForegroundColor Green
        Write-Host "   劣势: $($result.swot_analysis.weaknesses.Count) 个" -ForegroundColor Yellow
        Write-Host "   机会: $($result.swot_analysis.opportunities.Count) 个" -ForegroundColor Cyan
        Write-Host "   威胁: $($result.swot_analysis.threats.Count) 个" -ForegroundColor Red
        Write-Host ""
    }
    
    # 保存完整结果
    $result | ConvertTo-Json -Depth 10 | Out-File "analysis_result.json" -Encoding UTF8
    Write-Host "📄 完整结果已保存到: analysis_result.json" -ForegroundColor Gray
    
} catch {
    Write-Host "❌ 分析失败: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "💡 可能原因：" -ForegroundColor Yellow
    Write-Host "   1. Ollama服务未运行" -ForegroundColor Gray
    Write-Host "   2. 模型响应太慢（deepseek-r1需要5-10分钟）" -ForegroundColor Gray
    Write-Host "   3. 没有爬取的内容数据" -ForegroundColor Gray
}

Write-Host ""
Write-Host "✨ 完成" -ForegroundColor Cyan
```

---

### 4. 批量分析所有竞品

```powershell
# 获取所有竞品
$competitors = Invoke-WebRequest `
    -Uri http://localhost:8080/api/competitors `
    -UseBasicParsing |
    ConvertFrom-Json

# 逐个分析
foreach ($comp in $competitors.competitors) {
    Write-Host "分析: $($comp.name)" -ForegroundColor Cyan
    
    $body = @{
        competitor_id = $comp.id
        market_context = "中国市场"
    } | ConvertTo-Json
    
    try {
        Invoke-WebRequest `
            -Uri http://localhost:8080/api/analyze/competitor `
            -Method POST `
            -Body $body `
            -ContentType "application/json" `
            -UseBasicParsing `
            -TimeoutSec 600
        
        Write-Host "  ✅ 完成" -ForegroundColor Green
    } catch {
        Write-Host "  ❌ 失败: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # 等待一下，避免过载
    Start-Sleep -Seconds 5
}
```

---

## 🎯 API接口完整说明

### 接口地址
```
POST http://localhost:8080/api/analyze/competitor
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| competitor_id | int | ✅ | 竞品ID（从竞品列表获取） |
| market_context | string | ❌ | 市场背景（用于SWOT分析） |

### 请求示例

```json
{
  "competitor_id": 1,
  "market_context": "中国在线设计工具市场，竞争激烈，用户需求多样化"
}
```

### 响应示例

```json
{
  "success": true,
  "competitor": "Canva",
  "product_info": {
    "product_name": "Canva",
    "company": "Canva Pty Ltd",
    "tagline": "设计变得简单",
    "target_users": ["设计师", "营销人员", "教育工作者"],
    "core_features": [
      {
        "name": "模板库",
        "description": "海量设计模板",
        "category": "核心功能",
        "unique": false
      }
    ],
    "pricing": {
      "model": "订阅制",
      "tiers": [...]
    }
  },
  "swot_analysis": {
    "strengths": [...],
    "weaknesses": [...],
    "opportunities": [...],
    "threats": [...]
  }
}
```

---

## 💡 推荐流程

### 完整流程

1. **切换到Qwen2.5模型**（快）
2. **运行爬取任务**（获取数据）
3. **单独执行AI分析**（分析每个竞品）
4. **生成报告**

### 快速命令

```powershell
# 1. 切换模型
# 编辑 .env: LLM_MODEL=qwen2.5:7b

# 2. 重启服务
taskkill /F /IM competitive-analyzer-v2.exe
.\competitive-analyzer-v2.exe

# 3. 分析竞品
$body = @{competitor_id=1; market_context="中国市场"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST -Body $body -ContentType "application/json" -UseBasicParsing
```

---

## 🐛 故障排查

### 1. 检查Ollama

```powershell
# 查看模型列表
ollama list

# 测试连接
curl http://localhost:11434/api/version

# 测试模型
ollama run qwen2.5:7b "hello"
```

### 2. 查看服务日志

运行 `competitive-analyzer-v2.exe` 的终端会显示详细日志。

### 3. 检查是否有爬取数据

```powershell
# 查看竞品
Invoke-WebRequest -Uri http://localhost:8080/api/competitors -UseBasicParsing

# 查看数据源
Invoke-WebRequest -Uri http://localhost:8080/api/competitors/1/sources -UseBasicParsing
```

---

**建议**：立即切换到 `qwen2.5:7b`，响应速度快10倍！🚀
