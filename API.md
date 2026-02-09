# 竞品分析系统 - API接口文档与操作手册

> **版本**: v1.0.0  
> **基础URL**: `http://localhost:8080`  
> **更新日期**: 2026-02-09

---

## 📋 目录

- [快速开始](#快速开始)
- [认证说明](#认证说明)
- [通用响应格式](#通用响应格式)
- [核心模块](#核心模块)
  - [1. 健康检查](#1-健康检查)
  - [2. 数据源发现](#2-数据源发现)
  - [3. 内容爬取](#3-内容爬取)
  - [4. 竞品管理](#4-竞品管理)
- [完整操作流程](#完整操作流程)
- [错误码说明](#错误码说明)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 🚀 快速开始

### 前置条件

1. **服务已启动**：确保服务运行在 `http://localhost:8080`
2. **配置文件**：`.env` 文件已正确配置
3. **Ollama运行**：本地LLM服务已启动

### 快速测试

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

## 🔐 认证说明

当前版本：**无需认证**

未来版本将支持：
- JWT Token认证
- API Key认证
- OAuth 2.0

---

## 📦 通用响应格式

### 成功响应

```json
{
  "success": true,
  "data": {},
  "message": "操作成功"
}
```

### 错误响应

```json
{
  "error": "错误描述信息",
  "code": "ERROR_CODE",
  "details": {}
}
```

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 📡 核心模块

---

## 1. 健康检查

### 1.1 检查服务状态

**接口**: `GET /health`

**说明**: 检查API服务是否正常运行

**请求示例**:

```bash
curl http://localhost:8080/health
```

**响应示例**:

```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

## 2. 数据源发现

智能搜索并发现竞品及其数据源

---

### 2.1 启动发现任务

**接口**: `POST /api/discover/search`

**说明**: 根据主题自动搜索竞品和数据源

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | ✅ | 搜索主题（如："项目管理工具"） |
| market | string | ❌ | 目标市场（如："中国"、"全球"） |
| competitor_count | int | ❌ | 目标竞品数量（默认：5） |
| source_types | array | ❌ | 数据源类型（如：["官网", "定价页"]） |
| depth | string | ❌ | 搜索深度：quick/standard/deep（默认：standard） |

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "项目管理工具",
    "market": "中国",
    "competitor_count": 5,
    "depth": "standard"
  }'
```

**PowerShell示例**:

```powershell
$body = @{
    topic = "项目管理工具"
    market = "中国"
    competitor_count = 5
    depth = "standard"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/discover/search `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应示例**:

```json
{
  "task_id": 1,
  "status": "processing",
  "progress": 0,
  "estimated_time": 60
}
```

---

### 2.2 查询任务状态

**接口**: `GET /api/discover/status/:task_id`

**说明**: 查询发现任务的执行状态和结果

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | int | ✅ | 任务ID |

**请求示例**:

```bash
curl http://localhost:8080/api/discover/status/1
```

**PowerShell示例**:

```powershell
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1
```

**响应示例（进行中）**:

```json
{
  "status": "processing",
  "progress": 60,
  "competitors_found": 3,
  "data_sources_found": 12,
  "created_at": "2026-02-09T10:00:00Z",
  "completed_at": null
}
```

**响应示例（完成）**:

```json
{
  "status": "completed",
  "progress": 100,
  "competitors_found": 5,
  "data_sources_found": 25,
  "result": {
    "competitors": [
      "Notion",
      "飞书",
      "钉钉",
      "Teambition",
      "Worktile"
    ],
    "data_sources": {
      "Notion_官网": [
        {
          "url": "https://www.notion.so",
          "title": "Notion官网",
          "quality_score": 0.95
        }
      ]
    }
  },
  "created_at": "2026-02-09T10:00:00Z",
  "completed_at": "2026-02-09T10:01:30Z"
}
```

**状态说明**:

| 状态 | 说明 | 进度 |
|------|------|------|
| processing | 处理中 | 0-99% |
| completed | 已完成 | 100% |
| failed | 失败 | - |

---

### 2.3 确认并保存配置

**接口**: `POST /api/discover/confirm`

**说明**: 用户确认发现结果并保存到数据库

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | int | ✅ | 任务ID |
| selected_competitors | array | ✅ | 选中的竞品列表 |
| selected_sources | object | ❌ | 选中的数据源（按竞品分组） |
| save_as_config | bool | ❌ | 是否保存为配置 |

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/discover/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": 1,
    "selected_competitors": ["Notion", "飞书", "钉钉"],
    "selected_sources": {
      "Notion": [
        "https://www.notion.so",
        "https://www.notion.so/pricing"
      ],
      "飞书": [
        "https://www.feishu.cn"
      ]
    },
    "save_as_config": true
  }'
```

**PowerShell示例**:

```powershell
$body = @{
    task_id = 1
    selected_competitors = @("Notion", "飞书", "钉钉")
    selected_sources = @{
        "Notion" = @(
            "https://www.notion.so",
            "https://www.notion.so/pricing"
        )
        "飞书" = @("https://www.feishu.cn")
    }
    save_as_config = $true
} | ConvertTo-Json -Depth 5

Invoke-WebRequest -Uri http://localhost:8080/api/discover/confirm `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应示例**:

```json
{
  "success": true,
  "message": "配置已保存"
}
```

---

## 3. 内容爬取

使用三层爬取策略获取网页内容

---

### 3.1 爬取单个URL

**接口**: `POST /api/crawl/single`

**说明**: 爬取指定URL的内容并保存为Markdown

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | string | ✅ | 目标URL |
| competitor | string | ✅ | 竞品名称 |
| source_type | string | ❌ | 数据源类型（如："官网"、"定价页"） |

**爬取策略**:

1. **第一层 - Firecrawl** (最优先)
   - AI驱动，内容质量最高
   - 需要API Key
   
2. **第二层 - Jina Reader** (备选)
   - 免费服务
   - 质量较好
   
3. **第三层 - Playwright** (兜底)
   - 本地浏览器渲染
   - 支持复杂网站

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.notion.so",
    "competitor": "Notion",
    "source_type": "官网"
  }'
```

**PowerShell示例**:

```powershell
$body = @{
    url = "https://www.notion.so"
    competitor = "Notion"
    source_type = "官网"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/crawl/single `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**响应示例**:

```json
{
  "success": true,
  "content_path": "storage/crawled/Notion/2026-02-09_notion-so.md",
  "image_count": 5,
  "title": "Notion – The all-in-one workspace"
}
```

**保存位置**:

```
storage/
  └── crawled/
      └── Notion/
          ├── 2026-02-09_notion-so.md          # Markdown内容
          └── images/
              ├── image_1.png                   # 本地化图片
              └── image_2.png
```

---

### 3.2 批量爬取（规划中）

**接口**: `POST /api/crawl/batch`

**说明**: 批量爬取多个URL

**状态**: 🚧 开发中

---

## 4. 竞品管理

管理竞品和数据源信息

---

### 4.1 获取竞品列表

**接口**: `GET /api/competitors`

**说明**: 分页获取竞品列表

**查询参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | ❌ | 1 | 页码 |
| page_size | int | ❌ | 20 | 每页数量 |

**请求示例**:

```bash
curl "http://localhost:8080/api/competitors?page=1&page_size=10"
```

**PowerShell示例**:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/competitors?page=1&page_size=10"
```

**响应示例**:

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
      "category": "项目管理",
      "confidence": 0.95,
      "status": "active",
      "created_at": "2026-02-09T10:00:00Z",
      "updated_at": "2026-02-09T10:00:00Z"
    },
    {
      "id": 2,
      "name": "飞书",
      "company": "字节跳动",
      "website": "https://www.feishu.cn",
      "category": "协作工具",
      "confidence": 0.92,
      "status": "active",
      "created_at": "2026-02-09T10:01:00Z",
      "updated_at": "2026-02-09T10:01:00Z"
    }
  ]
}
```

---

### 4.2 获取数据源列表

**接口**: `GET /api/competitors/:id/sources`

**说明**: 获取指定竞品的所有数据源

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | ✅ | 竞品ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| competitor_id | int | ❌ | 过滤指定竞品 |

**请求示例**:

```bash
# 获取所有数据源
curl http://localhost:8080/api/competitors/1/sources

# 过滤特定竞品
curl "http://localhost:8080/api/competitors/1/sources?competitor_id=1"
```

**PowerShell示例**:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/competitors/1/sources"
```

**响应示例**:

```json
{
  "data_sources": [
    {
      "id": 1,
      "competitor_id": 1,
      "source_type": "官网",
      "url": "https://www.notion.so",
      "priority": 1,
      "quality_score": 0.95,
      "auto_discovered": true,
      "status": "active",
      "last_crawl_time": "2026-02-09T10:05:00Z"
    },
    {
      "id": 2,
      "competitor_id": 1,
      "source_type": "定价页",
      "url": "https://www.notion.so/pricing",
      "priority": 2,
      "quality_score": 0.90,
      "auto_discovered": true,
      "status": "active",
      "last_crawl_time": "2026-02-09T10:06:00Z"
    }
  ]
}
```

---

## 📖 完整操作流程

### 场景1：从零开始的竞品分析

**步骤1：启动发现任务**

```bash
# 1. 发起搜索
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "在线协作文档",
    "market": "中国",
    "competitor_count": 5,
    "depth": "standard"
  }'

# 响应：{"task_id": 1, "status": "processing"}
```

**步骤2：轮询任务状态**

```bash
# 2. 每5秒查询一次状态
while true; do
  curl http://localhost:8080/api/discover/status/1
  sleep 5
done

# 当 status="completed" 时停止
```

**步骤3：确认并保存**

```bash
# 3. 确认发现结果
curl -X POST http://localhost:8080/api/discover/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": 1,
    "selected_competitors": ["Notion", "飞书", "石墨文档"],
    "selected_sources": {
      "Notion": ["https://www.notion.so", "https://www.notion.so/pricing"],
      "飞书": ["https://www.feishu.cn"],
      "石墨文档": ["https://shimo.im"]
    }
  }'
```

**步骤4：爬取内容**

```bash
# 4. 爬取每个数据源
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.notion.so",
    "competitor": "Notion",
    "source_type": "官网"
  }'

# 重复其他URL...
```

**步骤5：查看结果**

```bash
# 5. 检查竞品列表
curl http://localhost:8080/api/competitors

# 6. 查看爬取的内容
ls storage/crawled/Notion/
```

---

### 场景2：已知竞品，直接爬取

```bash
# 直接爬取指定URL
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.notion.so/product",
    "competitor": "Notion",
    "source_type": "产品页"
  }'
```

---

### 场景3：监控竞品变化（规划中）

```bash
# 创建监控任务
curl -X POST http://localhost:8080/api/monitor/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notion价格监控",
    "competitor_ids": [1],
    "frequency": "daily",
    "alert_rules": {
      "price_change": true,
      "feature_change": true
    }
  }'
```

---

## ⚠️ 错误码说明

| 错误码 | HTTP状态 | 说明 | 解决方案 |
|--------|----------|------|----------|
| INVALID_PARAMS | 400 | 请求参数错误 | 检查JSON格式和必填字段 |
| TASK_NOT_FOUND | 404 | 任务不存在 | 确认task_id是否正确 |
| COMPETITOR_NOT_FOUND | 404 | 竞品不存在 | 检查竞品ID |
| CRAWL_FAILED | 500 | 爬取失败 | 检查URL可访问性和API配置 |
| DATABASE_ERROR | 500 | 数据库错误 | 检查数据库连接 |
| LLM_ERROR | 500 | LLM调用失败 | 检查Ollama服务状态 |

---

## 💡 最佳实践

### 1. 搜索深度选择

| 深度 | 适用场景 | 结果数量 | 耗时 |
|------|----------|----------|------|
| quick | 快速探索 | 5个竞品 | ~30秒 |
| standard | 日常分析 | 10个竞品 | ~60秒 |
| deep | 深度调研 | 20个竞品 | ~120秒 |

### 2. API调用频率

- **发现任务**: 每小时不超过10次
- **爬取请求**: 每分钟不超过30次
- **查询接口**: 无限制

### 3. 数据源优先级

1. **官网** - 最重要，优先爬取
2. **定价页** - 核心信息
3. **产品页** - 功能详情
4. **博客** - 动态信息

### 4. 错误重试策略

```bash
# 示例：带重试的爬取脚本
max_retries=3
retry_count=0

while [ $retry_count -lt $max_retries ]; do
  response=$(curl -X POST http://localhost:8080/api/crawl/single \
    -H "Content-Type: application/json" \
    -d '{"url": "...", "competitor": "..."}')
  
  if [ $? -eq 0 ]; then
    echo "成功"
    break
  fi
  
  retry_count=$((retry_count + 1))
  echo "重试 $retry_count/$max_retries"
  sleep 5
done
```

---

## 🔧 高级配置

### 自定义爬取选项

在 `.env` 文件中配置：

```env
# 爬取策略优先级
CRAWLER_STRATEGY=firecrawl,jina,playwright

# 超时时间（秒）
CRAWLER_TIMEOUT=30

# 最大重试次数
CRAWLER_MAX_RETRIES=3

# 图片下载
CRAWLER_DOWNLOAD_IMAGES=true
```

### LLM配置

```env
# 使用Ollama
OPENAI_BASE_URL=http://localhost:11434
OPENAI_API_KEY=ollama
LLM_MODEL=deepseek-r1:8b

# LLM参数
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
```

---

## ❓ 常见问题

### Q1: 爬取失败怎么办？

**A**: 检查以下几点：

1. URL是否可访问
2. Firecrawl API Key是否配置
3. 网络连接是否正常
4. 目标网站是否需要登录

**解决方案**:

```bash
# 测试URL可访问性
curl -I https://www.notion.so

# 检查API配置
cat .env | grep FIRECRAWL

# 查看详细错误日志
tail -f logs/app.log
```

---

### Q2: Ollama返回错误

**A**: 常见原因：

1. Ollama服务未启动
2. 模型未下载
3. 端口配置错误

**解决方案**:

```bash
# 检查Ollama状态
ollama list

# 测试Ollama连接
curl http://localhost:11434/api/tags

# 重启Ollama
ollama serve

# 下载模型
ollama pull deepseek-r1:8b
```

---

### Q3: 任务一直处于processing状态

**A**: 可能原因：

1. 后台任务卡住
2. 搜索API超时
3. 数据库锁定

**解决方案**:

```bash
# 查看任务详情
curl http://localhost:8080/api/discover/status/1

# 重启服务
# Ctrl+C 停止服务
go run main.go

# 或使用系统命令
pkill -f "go run main.go"
go run main.go
```

---

### Q4: 如何批量爬取多个URL？

**A**: 使用循环脚本：

```bash
# Bash脚本
urls=(
  "https://www.notion.so"
  "https://www.notion.so/pricing"
  "https://www.notion.so/product"
)

for url in "${urls[@]}"; do
  curl -X POST http://localhost:8080/api/crawl/single \
    -H "Content-Type: application/json" \
    -d "{\"url\": \"$url\", \"competitor\": \"Notion\"}"
  sleep 2  # 避免频繁请求
done
```

**PowerShell脚本**:

```powershell
$urls = @(
    "https://www.notion.so",
    "https://www.notion.so/pricing",
    "https://www.notion.so/product"
)

foreach ($url in $urls) {
    $body = @{
        url = $url
        competitor = "Notion"
    } | ConvertTo-Json

    Invoke-WebRequest -Uri http://localhost:8080/api/crawl/single `
        -Method POST `
        -Body $body `
        -ContentType "application/json"
    
    Start-Sleep -Seconds 2
}
```

---

### Q5: 如何导出分析报告？

**A**: 当前版本爬取的内容保存在：

```
storage/crawled/竞品名称/日期_文件名.md
```

可以使用文本编辑器直接查看，或使用脚本合并：

```bash
# 合并所有Notion相关内容
cat storage/crawled/Notion/*.md > Notion完整报告.md
```

---

### Q6: 数据保存在哪里？

**A**: 数据分布：

```
competitive_go/
├── data/
│   └── competitive.db          # SQLite数据库（结构化数据）
├── storage/
│   └── crawled/                # 爬取的Markdown内容
│       └── [竞品名]/
│           ├── *.md            # 网页内容
│           └── images/         # 图片资源
└── reports/                    # 生成的分析报告（规划中）
```

---

## 📚 相关文档

- [README.md](./README.md) - 项目概述
- [HOW_TO_RUN.md](./HOW_TO_RUN.md) - 运行指南
- [OLLAMA_MODEL_GUIDE.md](./OLLAMA_MODEL_GUIDE.md) - 模型推荐
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 架构设计

---

## 🔄 版本历史

### v1.0.0 (2026-02-09)

- ✅ 数据源发现
- ✅ 三层爬取策略
- ✅ 竞品管理
- ✅ Ollama集成
- 🚧 AI分析（开发中）
- 🚧 报告生成（开发中）
- 🚧 变化监控（开发中）

---

## 📞 技术支持

遇到问题？

1. 查看 [常见问题](#常见问题)
2. 检查 [错误码说明](#错误码说明)
3. 查看服务日志：`tail -f logs/app.log`
4. 重启服务试试

---

## 🎯 快速参考卡片

### 核心接口速查

```bash
# 健康检查
GET /health

# 启动发现任务
POST /api/discover/search

# 查询任务状态
GET /api/discover/status/:task_id

# 确认结果
POST /api/discover/confirm

# 爬取单个URL
POST /api/crawl/single

# 获取竞品列表
GET /api/competitors

# 获取数据源
GET /api/competitors/:id/sources
```

### PowerShell快捷命令

```powershell
# 健康检查
iwr http://localhost:8080/health

# 启动发现
$body = @{topic="项目管理"} | ConvertTo-Json
iwr http://localhost:8080/api/discover/search -Method POST -Body $body -ContentType "application/json"

# 查询状态
iwr http://localhost:8080/api/discover/status/1

# 爬取URL
$body = @{url="https://notion.so";competitor="Notion"} | ConvertTo-Json
iwr http://localhost:8080/api/crawl/single -Method POST -Body $body -ContentType "application/json"
```

---

**文档更新日期**: 2026-02-09  
**当前版本**: v1.0.0  
**服务状态**: ✅ 运行中

🎉 开始使用竞品分析系统吧！
