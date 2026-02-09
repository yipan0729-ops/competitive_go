# 🎯 Ollama连接失败 - 完整解决方案

## ❌ 问题

```
请求Ollama失败: context deadline exceeded
```

---

## 🔍 原因分析

### 1. DeepSeek-R1模型太慢
- **DeepSeek-R1:8B** 是推理模型
- 首次加载：30-60秒
- 单次推理：5-10分钟 ⚠️
- 每个竞品需要调用2次（产品信息 + SWOT）
- **总耗时：10-20分钟/竞品**

### 2. Ollama可能没运行
- 服务未启动
- 模型未加载

---

## ✅ 解决方案

### 方案1：切换到Qwen2.5（强烈推荐）⭐⭐⭐

#### 下载模型
```powershell
ollama pull qwen2.5:7b
```

#### 修改配置
编辑 `D:\Code\Competitive_go\.env`：

```env
# 从这个：
LLM_MODEL=deepseek-r1:8b

# 改为：
LLM_MODEL=qwen2.5:7b
```

#### 重启服务
```powershell
# 方式1：按Ctrl+C停止，然后重新运行
.\competitive-analyzer-v2.exe

# 方式2：强制停止并重启
taskkill /F /IM competitive-analyzer-v2.exe
Start-Sleep -Seconds 1
.\competitive-analyzer-v2.exe
```

#### 性能对比
| 模型 | 加载时间 | 单次推理 | 适用场景 |
|------|----------|----------|----------|
| qwen2.5:7b | 5秒 | 5-10秒 | ✅ 推荐 |
| deepseek-r1:8b | 30秒 | 5-10分钟 | ❌ 太慢 |

---

### 方案2：确保Ollama运行

```powershell
# 1. 检查Ollama
ollama list

# 2. 如果没有运行，启动
ollama serve

# 3. 在另一个终端测试
ollama run qwen2.5:7b "你好"
```

---

### 方案3：使用云端API（最快，需付费）

如果需要最快速度，使用Groq：

#### 配置Groq
```env
# .env
OPENAI_API_KEY=你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile
```

#### 性能
- 单次推理：**1-2秒** ⚡
- 免费额度：每天6000次

---

## 📡 单独执行AI分析的方法

### API接口

```
POST http://localhost:8080/api/analyze/competitor
```

### 完整示例

```powershell
# 1. 获取竞品列表
$competitors = Invoke-WebRequest `
    -Uri http://localhost:8080/api/competitors `
    -UseBasicParsing |
    ConvertFrom-Json

Write-Host "找到 $($competitors.total) 个竞品"

# 2. 选择第一个竞品
$firstId = $competitors.competitors[0].id
$firstName = $competitors.competitors[0].name

Write-Host "分析: $firstName (ID: $firstId)"

# 3. 执行AI分析
$body = @{
    competitor_id = $firstId
    market_context = "中国市场"
} | ConvertTo-Json

Write-Host "⏳ 调用Ollama，请等待..."

$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -TimeoutSec 600

# 4. 查看结果
$result = $response.Content | ConvertFrom-Json

Write-Host "✅ 分析完成！"
Write-Host "产品: $($result.product_info.product_name)"
Write-Host "公司: $($result.product_info.company)"
Write-Host "优势: $($result.swot_analysis.strengths.Count) 个"

# 保存结果
$result | ConvertTo-Json -Depth 10 | Out-File "analysis.json" -Encoding UTF8
```

### 或使用脚本

```powershell
# 使用提供的测试脚本
.\test-ai-analysis.ps1
```

---

## 🎯 推荐工作流

### 完整流程（推荐使用Qwen2.5）

```powershell
# 1. 切换模型（编辑.env）
LLM_MODEL=qwen2.5:7b

# 2. 下载模型
ollama pull qwen2.5:7b

# 3. 重启服务
taskkill /F /IM competitive-analyzer-v2.exe
.\competitive-analyzer-v2.exe

# 4. 运行完整分析
$body = @{
    topic = "在线工具"
    competitor_count = 2
    depth = "quick"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis `
    -Method POST -Body $body -ContentType "application/json" -UseBasicParsing

# 5. 监控进度
# 等待2-3分钟完成
```

---

## 🐛 故障排查清单

### ✅ 检查项目

- [ ] Ollama服务运行中 (`ollama list`)
- [ ] 模型已下载 (`qwen2.5:7b` 或 `deepseek-r1:8b`)
- [ ] 模型配置正确（`.env` 文件）
- [ ] 服务已重启（加载新配置）
- [ ] 已有爬取数据（`storage/crawled/` 不为空）
- [ ] 使用正确的竞品ID

---

## 💡 关键改进

已实现的优化：

1. ✅ **超时时间增加**：从120秒 → 600秒（10分钟）
2. ✅ **爬虫重试机制**：失败后自动重试3次
3. ✅ **渐进式延迟**：避免并发触发403
4. ✅ **降低并发数**：从3 → 1，避免被封
5. ✅ **Firecrawl集成**：使用您的API Key

---

## 🚀 立即行动

**最重要的一步**：切换到Qwen2.5！

```powershell
# 1. 下载
ollama pull qwen2.5:7b

# 2. 修改.env
notepad D:\Code\Competitive_go\.env
# 改为: LLM_MODEL=qwen2.5:7b

# 3. 重启服务

# 4. 测试
.\test-ai-analysis.ps1
```

**Qwen2.5响应速度快10-20倍！** ⚡

---

**文档**: [OLLAMA_TIMEOUT_FIX.md](./OLLAMA_TIMEOUT_FIX.md)  
**测试脚本**: `test-ai-analysis.ps1`
