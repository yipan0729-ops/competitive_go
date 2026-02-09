# 竞品分析系统 API 文档

> **版本**: v1.0.0  
> **基础URL**: `http://localhost:8080`  
> **更新日期**: 2026-02-09

---

## 目录

- [快速开始](#快速开始)
- [完整工作流](#完整工作流)
- [API接口](#api接口)
  - [1. 健康检查](#1-健康检查)
  - [2. 全流程自动化](#2-全流程自动化)
  - [3. 数据源发现](#3-数据源发现)
  - [4. 内容爬取](#4-内容爬取)
  - [5. AI分析](#5-ai分析)
  - [6. 报告生成](#6-报告生成)
  - [7. 竞品管理](#7-竞品管理)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)

---

## 快速开始

### 前置条件

1. **Go服务运行中**: `go run main.go`
2. **LLM服务可用**: Ollama或其他LLM服务
3. **配置文件**: `.env` 已正确配置

### 健康检查

```powershell
# 检查服务状态
Invoke-WebRequest -Uri http://localhost:8080/health
```

**响应示例**:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

## 完整工作流

### 方式1: 全自动流程（推荐）⭐

一个API调用完成所有步骤：发现竞品 → 爬取内容 → AI分析 → 生成报告

```powershell
$body = @{
    topic = "AI创作工具"
    market = "中国"
    competitor_count = 3
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing
```

**响应**:
```json
{
  "success": true,
  "task_id": 1,
  "status": "processing",
  "workflow": "discovery -> crawl -> analysis -> report",
  "estimated_time": 600
}
```

**查询进度**:
```powershell
# 轮询任务状态
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1
```

---

### 方式2: 手动分步流程

适用于需要精细控制的场景。

#### 步骤1: 发现竞品

```powershell
$body = @{
    topic = "项目管理工具"
    market = "中国"
    competitor_count = 5
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/discover/search `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

#### 步骤2: 查询发现结果

```powershell
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1
```

#### 步骤3: 确认并保存配置

```powershell
$body = @{
    task_id = 1
    selected_competitors = @("Notion", "飞书")
    selected_sources = @{
        "Notion" = @("https://www.notion.so")
        "飞书" = @("https://www.feishu.cn")
    }
} | ConvertTo-Json -Depth 5

Invoke-WebRequest `
    -Uri http://localhost:8080/api/discover/confirm `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

#### 步骤4: 爬取内容

```powershell
$body = @{
    url = "https://www.notion.so"
    competitor = "Notion"
    source_type = "官网"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/crawl/single `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

#### 步骤5: AI分析

```powershell
$body = @{
    competitor_id = 1
    market_context = "中国协作工具市场"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -TimeoutSec 600
```

#### 步骤6: 生成报告

```powershell
$body = @{
    competitor_ids = @(1, 2)
    topic = "项目管理工具"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/report/generate `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

---

## API接口

---

## 1. 健康检查

### GET /health

检查API服务是否正常运行。

**响应**:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

## 2. 全流程自动化

### POST /api/auto/analysis

一键完成全流程：发现 → 爬取 → 分析 → 报告。

**请求参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| topic | string | ✅ | - | 分析主题 |
| market | string | ❌ | - | 目标市场 |
| competitor_count | int | ❌ | 5 | 竞品数量 |
| depth | string | ❌ | standard | quick/standard/deep |
| auto_crawl | bool | ❌ | true | 是否自动爬取 |
| auto_analyze | bool | ❌ | true | 是否自动分析 |
| generate_report | bool | ❌ | true | 是否生成报告 |

**请求示例**:
```powershell
$body = @{
    topic = "AI写作助手"
    market = "中国"
    competitor_count = 3
    depth = "standard"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing
```

**响应示例**:
```json
{
  "success": true,
  "task_id": 1,
  "status": "processing",
  "workflow": "discovery -> crawl -> analysis -> report",
  "estimated_time": 600
}
```

**查询进度**:
```powershell
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1
```

**完成后响应**:
```json
{
  "status": "completed",
  "progress": 100,
  "result": {
    "competitors": ["Notion AI", "Jasper", "Copy.ai"],
    "urls_crawled": 9,
    "analyzed_count": 3,
    "report_path": "reports/AI写作助手_自动分析报告_20260209.md"
  }
}
```

---

## 3. 数据源发现

### POST /api/discover/search

智能搜索竞品和数据源。

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | ✅ | 搜索主题 |
| market | string | ❌ | 目标市场 |
| competitor_count | int | ❌ | 目标数量（默认5） |
| depth | string | ❌ | quick/standard/deep |

**请求示例**:
```powershell
$body = @{
    topic = "CRM系统"
    market = "中国"
    competitor_count = 5
    depth = "standard"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/discover/search `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应**:
```json
{
  "task_id": 1,
  "status": "processing",
  "progress": 0,
  "estimated_time": 60
}
```

---

### GET /api/discover/status/:task_id

查询发现任务状态。

**请求示例**:
```powershell
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1
```

**响应（进行中）**:
```json
{
  "status": "processing",
  "progress": 60,
  "competitors_found": 3,
  "data_sources_found": 12
}
```

**响应（完成）**:
```json
{
  "status": "completed",
  "progress": 100,
  "competitors_found": 5,
  "data_sources_found": 25,
  "result": {
    "competitors": ["Salesforce", "HubSpot", "纷享销客"],
    "data_sources": {
      "Salesforce_官网": [
        {
          "url": "https://www.salesforce.com",
          "title": "Salesforce官网",
          "quality_score": 0.95
        }
      ]
    }
  }
}
```

---

### POST /api/discover/confirm

确认并保存发现结果。

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | int | ✅ | 任务ID |
| selected_competitors | array | ✅ | 选中的竞品 |
| selected_sources | object | ❌ | 选中的数据源 |
| save_as_config | bool | ❌ | 保存为配置 |

**请求示例**:
```powershell
$body = @{
    task_id = 1
    selected_competitors = @("Notion", "飞书")
    selected_sources = @{
        "Notion" = @("https://www.notion.so", "https://www.notion.so/pricing")
        "飞书" = @("https://www.feishu.cn")
    }
    save_as_config = $true
} | ConvertTo-Json -Depth 5

Invoke-WebRequest `
    -Uri http://localhost:8080/api/discover/confirm `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应**:
```json
{
  "success": true,
  "message": "配置已保存"
}
```

---

## 4. 内容爬取

### POST /api/crawl/single

爬取单个URL内容。

**爬取策略**（自动选择）:
1. **Firecrawl** - AI驱动，质量最高
2. **Jina Reader** - 免费备选
3. **Playwright** - 本地兜底

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | string | ✅ | 目标URL |
| competitor | string | ✅ | 竞品名称 |
| source_type | string | ❌ | 数据源类型 |

**请求示例**:
```powershell
$body = @{
    url = "https://www.notion.so"
    competitor = "Notion"
    source_type = "官网"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/crawl/single `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应**:
```json
{
  "success": true,
  "content_path": "storage/crawled/Notion/2026-02-09_notion-so.md",
  "image_count": 5,
  "title": "Notion – The all-in-one workspace"
}
```

---

### POST /api/crawl/batch

批量爬取多个URL。

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| urls | array | ✅ | URL列表 |
| concurrent | int | ❌ | 并发数（默认1，最大3） |

**URL项格式**:
```json
{
  "url": "https://example.com",
  "competitor": "竞品名",
  "source_type": "官网"
}
```

**请求示例**:
```powershell
$body = @{
    urls = @(
        @{
            url = "https://www.notion.so"
            competitor = "Notion"
            source_type = "官网"
        },
        @{
            url = "https://www.notion.so/pricing"
            competitor = "Notion"
            source_type = "定价页"
        }
    )
    concurrent = 1
} | ConvertTo-Json -Depth 5

Invoke-WebRequest `
    -Uri http://localhost:8080/api/crawl/batch `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应**:
```json
{
  "success": true,
  "total_urls": 2,
  "concurrent": 1,
  "message": "批量爬取任务已启动"
}
```

---

## 5. AI分析

### POST /api/analyze/competitor

对竞品进行AI分析（产品信息提取 + SWOT分析）。

**注意**: 此接口调用LLM，可能需要较长时间（1-5分钟）。

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| competitor_id | int | ✅ | 竞品ID |
| market_context | string | ❌ | 市场背景 |

**请求示例**:
```powershell
$body = @{
    competitor_id = 1
    market_context = "中国在线协作市场"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -TimeoutSec 600
```

**响应示例**:
```json
{
  "success": true,
  "competitor": "Notion",
  "product_info": {
    "product_name": "Notion",
    "company": "Notion Labs Inc.",
    "tagline": "All-in-one workspace",
    "target_users": ["知识工作者", "创作者"],
    "core_features": [
      {
        "name": "笔记与文档",
        "description": "强大的编辑器",
        "category": "核心功能"
      }
    ],
    "pricing": {
      "model": "订阅制",
      "tiers": [
        {
          "name": "Free",
          "price": 0,
          "features": ["个人使用"]
        }
      ]
    }
  },
  "swot_analysis": {
    "strengths": [
      {
        "point": "功能全面",
        "evidence": "集成多种工具",
        "impact": "高"
      }
    ],
    "weaknesses": [...],
    "opportunities": [...],
    "threats": [...]
  }
}
```

---

## 6. 报告生成

### POST /api/report/generate

生成竞品分析报告（Markdown格式）。

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| competitor_ids | array | ✅ | 竞品ID列表 |
| topic | string | ✅ | 分析主题 |
| report_name | string | ❌ | 报告名称 |

**请求示例**:
```powershell
$body = @{
    competitor_ids = @(1, 2, 3)
    topic = "项目管理工具"
    report_name = "项目管理工具竞品分析"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/report/generate `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应**:
```json
{
  "success": true,
  "report_id": 1,
  "report_name": "项目管理工具竞品分析",
  "report_path": "reports/项目管理工具竞品分析.md",
  "competitors": 3
}
```

**报告包含内容**:
- 📊 执行摘要
- 🏢 竞品概览
- ⚙️ 功能对比矩阵
- 💰 价格策略分析
- 📈 SWOT分析
- 💡 战略建议

---

## 7. 竞品管理

### GET /api/competitors

获取竞品列表（分页）。

**查询参数**:

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 20 | 每页数量 |

**请求示例**:
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/competitors?page=1&page_size=10"
```

**响应**:
```json
{
  "total": 50,
  "page": 1,
  "page_size": 10,
  "competitors": [
    {
      "id": 1,
      "name": "Notion",
      "company": "Notion Labs Inc.",
      "website": "https://www.notion.so",
      "status": "active",
      "created_at": "2026-02-09T10:00:00Z"
    }
  ]
}
```

---

### GET /api/data_sources

获取数据源列表。

**查询参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| competitor_id | int | 过滤指定竞品 |

**请求示例**:
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/data_sources?competitor_id=1"
```

**响应**:
```json
{
  "data_sources": [
    {
      "id": 1,
      "competitor_id": 1,
      "url": "https://www.notion.so",
      "source_type": "官网",
      "priority": 1,
      "quality_score": 0.95,
      "status": "active",
      "last_crawl_time": "2026-02-09T10:05:00Z"
    }
  ]
}
```

---

## 错误处理

### 通用响应格式

**成功响应**:
```json
{
  "success": true,
  "data": {}
}
```

**错误响应**:
```json
{
  "error": "错误描述"
}
```

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

### 常见错误

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| "竞品不存在" | 竞品ID无效 | 检查ID是否正确 |
| "该竞品没有数据源" | 未爬取内容 | 先执行爬取操作 |
| "没有可分析的内容" | 内容文件缺失 | 重新爬取 |
| "SWOT分析失败" | LLM服务异常 | 检查Ollama服务 |
| "请求Ollama失败" | Ollama未运行 | 启动Ollama服务 |

---

## 最佳实践

### 1. 超时设置

AI分析和报告生成可能需要较长时间，建议设置足够的超时时间：

```powershell
# 推荐设置
-TimeoutSec 600  # 10分钟超时
```

### 2. 并发控制

批量爬取时控制并发数，避免触发反爬虫：

```powershell
# 推荐配置
concurrent = 1  # 串行爬取，最稳定
```

### 3. 错误重试

网络请求可能失败，建议添加重试逻辑：

```powershell
$maxRetries = 3
$retryCount = 0

while ($retryCount -lt $maxRetries) {
    try {
        $response = Invoke-WebRequest -Uri $url -Method POST -Body $body
        break
    } catch {
        $retryCount++
        if ($retryCount -lt $maxRetries) {
            Write-Host "重试 $retryCount/$maxRetries..."
            Start-Sleep -Seconds 5
        }
    }
}
```

### 4. 数据源优先级

爬取数据源的推荐顺序：

1. **官网** - 最重要，优先爬取
2. **定价页** - 核心商业信息
3. **产品页** - 功能详情
4. **文档** - 技术细节

### 5. 内容存储

爬取的内容保存在：

```
storage/
└── crawled/
    └── [竞品名]/
        ├── 2026-02-09_notion-so.md
        └── images/
            └── image_1.png
```

报告保存在：

```
reports/
└── [主题]_竞品分析报告_20260209.md
```

---

## 常见问题

### Q1: Ollama连接失败？

**检查步骤**:
```powershell
# 1. 检查Ollama是否运行
ollama list

# 2. 测试连接
curl http://localhost:11434/api/tags

# 3. 检查.env配置
cat .env | Select-String "OPENAI"
```

### Q2: 爬取一直失败？

**可能原因**:
- 网站需要登录
- 反爬虫保护
- URL不可访问

**解决方案**:
```powershell
# 测试URL可访问性
curl -I https://www.notion.so

# 查看详细错误
# 检查Go服务输出日志
```

### Q3: AI分析超时？

**原因**: LLM响应时间过长

**解决方案**:
```powershell
# 增加超时时间
-TimeoutSec 1200  # 20分钟

# 或使用更快的模型
LLM_MODEL=qwen2.5:7b
```

### Q4: 如何查看生成的报告？

```powershell
# 查看reports目录
ls reports/

# 用记事本打开
notepad reports/AI创作_竞品分析报告_20260209.md

# 或用VS Code打开
code reports/AI创作_竞品分析报告_20260209.md
```

---

## 相关文档

- [README.md](./README.md) - 项目概述
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始指南
- [HOW_TO_RUN.md](./HOW_TO_RUN.md) - 详细运行指南
- [OLLAMA_SETUP.md](./OLLAMA_SETUP.md) - Ollama配置
- [FREE_LLM_ALTERNATIVES.md](./FREE_LLM_ALTERNATIVES.md) - 免费LLM选项
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 架构设计

---

## 版本历史

### v1.0.0 (2026-02-09)

- ✅ 数据源发现
- ✅ 三层爬取策略
- ✅ AI分析（产品信息 + SWOT）
- ✅ 报告生成
- ✅ 全流程自动化
- ✅ JSON解析修复

---

**文档更新**: 2026-02-09  
**API版本**: v1.0.0  
**维护状态**: ✅ 活跃维护
