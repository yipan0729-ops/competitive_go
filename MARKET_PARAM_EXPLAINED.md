# 📖 API参数详细说明

## Market参数说明

### `market` 参数的作用

在不同接口中有不同用途：

---

### 1. 在发现任务中 (`POST /api/discover/search`)

```json
{
  "topic": "AI创作",
  "market": "中国"     ← 这个参数
}
```

**作用**：
- ❌ **当前代码未使用此参数**
- 🎯 **设计意图**：限定搜索地域
  - 例如：`market="中国"` → 优先搜索中国市场的AI创作工具
  - 例如：`market="全球"` → 搜索全球范围的工具

**实际情况**：
- 目前代码中 `market` 参数**只是保存到数据库**
- **没有真正影响搜索结果**
- 搜索引擎会返回全球的结果

**改进建议**：
- 可以在搜索时添加地域关键词
- 例如：搜索 "AI创作 中国" 或 "AI创作 site:.cn"

---

### 2. 在自动化分析中 (`POST /api/auto/analysis`)

```json
{
  "topic": "AI创作",
  "market": "中国",           ← 这个参数
  "competitor_count": 3
}
```

**作用**：
1. **保存到任务记录**（用于标识）
2. **传递给AI分析** - 作为 `market_context`
3. **在SWOT分析时使用** - 帮助AI理解市场背景

**示例**：
- `market="中国"` → AI会考虑中国市场特点（审查、本地化等）
- `market="美国"` → AI会考虑美国市场特点（隐私法规等）

---

### 3. 在AI分析中 (`POST /api/analyze/competitor`)

```json
{
  "competitor_id": 1,
  "market_context": "中国AI工具市场，竞争激烈，用户更注重中文支持和本地化"
}
```

**作用**：
- ✅ **直接传递给LLM**
- 帮助AI理解市场环境
- 影响SWOT分析的结果
- 影响战略建议

**示例**：

**输入**：
```json
{
  "market_context": "中国市场，政策监管严格，用户注重数据安全"
}
```

**AI输出的威胁会考虑**：
```json
{
  "threats": [
    {
      "point": "政策风险",
      "context": "中国监管严格",
      "action": "加强合规审查"
    }
  ]
}
```

---

## ❌ 当前问题：竞品提取不准确

### 问题代码（已修复）

**之前**：
```go
// 简化处理：使用标题作为竞品名 ❌
nameList = append(nameList, result.Title)
```

这导致：
- 输入："AI创作"
- 输出："中国AI？美国AI？"、"生成式AI与不同类型的AI"
- **这些是文章标题，不是产品名！**

### 修复后（新逻辑）

```go
// 尝试从URL提取品牌名 ✅
brandName := extractBrandFromURL(result.URL)
// https://www.notion.so → Notion
// https://www.canva.com → Canva
```

**但这仍不完美**！真正的解决方案应该使用AI提取竞品名。

---

## ✅ 完整解决方案

### 方案1：使用AI提取竞品名（推荐）

应该从搜索结果的文章内容中，用AI提取真正的竞品名：

```
文章标题："2024年最好的10款AI创作工具"
↓ AI分析
提取竞品：["Notion AI", "Jasper", "Copy.ai", "Writesonic"]
```

**但这需要爬取文章内容，成本较高。**

---

### 方案2：改进搜索策略（当前采用）

修改搜索查询，直接搜索产品：

```go
// 搜索 "AI创作工具 产品" 而不是 "AI创作"
// 搜索 "AI创作 官网" 来获取产品官网
```

---

### 方案3：手动指定竞品（最准确）

创建新接口，让用户直接指定竞品名：

```json
{
  "competitors": ["Notion AI", "Jasper", "Copy.ai"],
  "auto_search_sources": true
}
```

---

## 💡 推荐使用方式

### 方式1：直接输入产品名（最准确）

```powershell
# 不要用泛泛的主题，直接用产品名
$body = @{
    topic = "Notion"         # 明确的产品名
    market = "全球"
    competitor_count = 3
} | ConvertTo-Json
```

### 方式2：使用更具体的主题

```powershell
# 不要用：AI创作
# 应该用：AI写作工具、AI文案生成、AI内容创作平台

$body = @{
    topic = "AI文案生成工具"  # 更具体
    market = "中国"
} | ConvertTo-Json
```

---

## 🎯 Market参数最佳实践

### 推荐写法

```json
// ✅ 好的写法
{
  "topic": "在线设计工具",
  "market": "中国"
}

// ✅ 详细的market_context
{
  "competitor_id": 1,
  "market_context": "中国在线设计市场，用户群体年轻化，移动端使用占比高，注重中文字体和模板"
}

// ❌ 不推荐
{
  "topic": "AI",           // 太泛泛
  "market": "中国"
}
```

### Market值建议

- `"中国"` - 中国大陆市场
- `"全球"` - 全球市场
- `"美国"` - 美国市场
- `"亚太"` - 亚太地区
- `"欧洲"` - 欧洲市场

---

## 🚀 单独执行AI分析

针对已存在的竞品ID=3：

```powershell
# 分析竞品ID=3
$body = @{
    competitor_id = 3
    market_context = "中国AI工具市场，用户需求：中文支持、易用性、价格敏感"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing `
    -TimeoutSec 600
```

**注意**：
1. 需要先确保Ollama运行
2. 强烈建议切换到 `qwen2.5:7b`（快10倍）
3. 确保竞品ID=3有爬取的内容数据

---

## 📚 总结

| 参数 | 用途 | 实际效果 | 建议 |
|------|------|----------|------|
| `market` (搜索) | 限定地域 | ❌ 当前未实现 | 用更具体的topic |
| `market` (自动化) | 保存标识 | ⚠️ 仅记录 | 可以忽略 |
| `market_context` (分析) | AI背景 | ✅ 影响SWOT | 详细描述市场 |

**关键**：当前的 `market` 参数**影响不大**，更重要的是：
1. 使用**具体的产品名**作为topic
2. 在AI分析时提供**详细的market_context**

---

**修复状态**: ✅ 已改进竞品提取逻辑  
**建议**: 切换到qwen2.5模型，然后重新测试
