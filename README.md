# 自动化竞品调研工具 v2.0 🚀

> **一键完成竞品分析！从输入主题到生成报告，只需5分钟！**

基于 Go + Gin 框架开发的全自动化竞品分析工具，支持智能数据源发现、批量数据采集、AI智能分析和专业报告生成。

## ✨ v2.0 重大更新

### 🔥 全流程自动化
- **一个API调用完成所有步骤**
- 自动发现竞品 → 批量爬取 → AI分析 → 生成报告
- 5分钟获得完整分析报告
- 零人工干预

### 🆕 新增功能
- ✅ **批量爬取** - 一次爬取多个URL，支持并发控制
- ✅ **AI智能分析** - 自动提取产品信息和SWOT分析
- ✅ **报告自动生成** - 专业Markdown格式报告
- ✅ **全流程自动化** - 一键完成所有步骤

---

## 功能特性

### 1. 智能数据源发现 ✨
- 支持多搜索引擎（Serper、Google、Bing）
- 自动发现竞品和数据源
- 智能链接分类和质量评分
- 搜索结果缓存优化

### 2. 高效数据采集 ⚡
- **批量并发爬取** 🆕
  - 支持1-10个并发
  - 异步执行，立即返回
  - 自动错误处理
- **三层策略自动降级**
  - Firecrawl (AI驱动，首选)
  - Jina Reader (免费，备选)
  - Playwright (浏览器自动化，兜底)
- 支持多平台识别（微信公众号、小红书、知乎、淘宝、京东等）
- 自动下载并本地化图片
- Markdown格式保存

### 3. AI智能分析 🤖 🆕
- **产品信息提取**
  - 公司信息、产品定位
  - 目标用户群体
  - 核心功能列表
  - 价格策略分析
- **SWOT自动分析**
  - 优势/劣势/机会/威胁
  - 证据支持
  - 影响评估
  - 战略建议
- **支持多种LLM**
  - Ollama（本地，免费）⭐ 推荐
  - OpenAI
  - DeepSeek
  - Groq
  - 智谱AI、通义千问

### 4. 专业报告生成 📊 🆕
- 完整的Markdown报告
- 功能对比矩阵
- 价格策略分析表
- 完整SWOT分析
- 战略建议输出
- 数据来源附录

### 5. 数据管理
- SQLite数据库存储（纯Go驱动，无需CGO）
- 历史数据追踪
- 内容变化检测
- 监控预警（v3.0计划）

## 快速开始

### 🎯 最简单方式 - 一键自动化（推荐）

```powershell
# 1. 启动服务
go run main.go

# 2. 运行自动化分析（只需这一个命令！）
$body = @{topic="项目管理工具"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis `
    -Method POST -Body $body -ContentType "application/json" -UseBasicParsing

# 3. 等待5分钟，查看报告
ls reports/
```

**就这么简单！系统会自动**:
- ✅ 发现5个竞品
- ✅ 搜索15个数据源
- ✅ 批量爬取所有网页
- ✅ AI提取产品信息
- ✅ 生成SWOT分析
- ✅ 输出完整报告

---

### 📋 传统方式 - 完整设置

### 1. 环境要求

- Go 1.21+ （已测试：1.21、1.22、1.23、1.24）
- Ollama（本地LLM，免费）⭐ 推荐
- 各API的密钥（可选：Firecrawl、Serper等）

### 2. 安装依赖

```bash
# 克隆项目
git clone <项目地址>
cd Competitive_go

# 设置Go代理（国内必需）
go env -w GOPROXY=https://goproxy.cn,direct

# 安装Go依赖
go mod download
```

### 3. 配置环境变量

复制 `.env.ollama` 为 `.env`：

```bash
cp .env.ollama .env
```

**推荐配置** (Ollama本地LLM)：

```env
# Ollama本地LLM（完全免费）
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b

# 搜索API（可选）
SERPER_API_KEY=你的密钥

# 其他配置
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
SERVER_PORT=8080
```

### 4. 安装Ollama

**Windows**:
```powershell
# 下载安装
# 访问: https://ollama.com/download/windows
# 或使用 winget
winget install Ollama.Ollama

# 下载推荐模型
ollama pull qwen2.5:7b
```

**Linux/Mac**:
```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull qwen2.5:7b
```

### 5. 启动服务

```bash
go run main.go
```

看到 `服务器启动在端口 8080` 表示成功！
FIRECRAWL_API_KEY=your_firecrawl_key_here
SERPER_API_KEY=your_serper_key_here
OPENAI_API_KEY=your_openai_key_here

# 可选
GOOGLE_API_KEY=your_google_key_here
GOOGLE_SEARCH_ENGINE_ID=your_engine_id_here
BING_API_KEY=your_bing_key_here

# 服务器配置
SERVER_PORT=8080
GIN_MODE=release

# 数据库配置
DB_PATH=./data/competitive.db

# 存储配置
STORAGE_PATH=./storage
REPORTS_PATH=./reports
```

### 4. 运行服务

```bash
# 开发模式
go run main.go

# 编译运行
go build -o competitive-analyzer
./competitive-analyzer
```

服务将在 `http://localhost:8080` 启动。

### 5. API使用示例

#### 智能数据源发现

```bash
# 1. 开始发现任务
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "AI写作助手",
    "market": "中国",
    "competitor_count": 5,
    "source_types": ["官网", "评价", "电商"],
    "depth": "standard"
  }'

# 响应
{
  "task_id": 1,
  "status": "processing",
  "progress": 0,
  "estimated_time": 60
}

# 2. 查询任务状态
curl http://localhost:8080/api/discover/status/1

# 响应
{
  "status": "completed",
  "progress": 100,
  "competitors_found": 5,
  "data_sources_found": 42,
  "result": {
    "competitors": ["Notion AI", "Jasper", "Copy.ai"],
    "data_sources": {...}
  }
}

# 3. 确认并保存配置
curl -X POST http://localhost:8080/api/discover/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": 1,
    "selected_competitors": ["Notion AI", "Jasper"],
    "selected_sources": {
      "Notion AI": ["https://notion.so", "https://notion.so/pricing"]
    },
    "save_as_config": true
  }'
```

#### 单个URL爬取

```bash
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://notion.so",
    "competitor": "Notion AI",
    "source_type": "官网"
  }'

# 响应
{
  "success": true,
  "content_path": "./storage/20260206_Notion_AI_首页/content.md",
  "image_count": 5,
  "title": "Notion – The all-in-one workspace"
}
```

#### 获取竞品列表

```bash
curl http://localhost:8080/api/competitors?page=1&page_size=20

# 响应
{
  "total": 5,
  "page": 1,
  "page_size": 20,
  "competitors": [
    {
      "id": 1,
      "name": "Notion AI",
      "company": "Notion Labs",
      "website": "https://notion.so",
      "status": "active"
    }
  ]
}
```

## 项目结构

```
Competitive_go/
├── main.go                 # 主程序入口
├── config/                 # 配置模块
│   └── config.go
├── models/                 # 数据模型
│   └── models.go
├── database/               # 数据库
│   └── database.go
├── crawler/                # 爬虫模块
│   ├── platform.go         # 平台识别
│   ├── crawler.go          # 三层爬虫策略
│   └── saver.go            # 内容保存
├── discovery/              # 数据源发现模块
│   ├── search.go           # 搜索引擎集成
│   ├── manager.go          # 搜索管理器
│   └── classifier.go       # 链接分类和评分
├── ai/                     # AI分析模块
│   ├── llm.go              # LLM客户端
│   └── extractor.go        # 信息提取器
├── report/                 # 报告生成模块
│   └── generator.go
├── handlers/               # HTTP处理器
│   └── handlers.go
├── storage/                # 存储目录（自动创建）
├── reports/                # 报告目录（自动创建）
├── data/                   # 数据库目录（自动创建）
├── go.mod                  # Go模块依赖
├── .env.example            # 环境变量示例
└── README.md               # 项目说明
```

## API文档

### 数据源发现模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/discover/search` | POST | 开始智能发现任务 |
| `/api/discover/status/:task_id` | GET | 查询任务状态 |
| `/api/discover/confirm` | POST | 确认并保存配置 |

### 爬取模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/crawl/single` | POST | 爬取单个URL |
| `/api/crawl/batch` | POST | 批量爬取（TODO） |
| `/api/crawl/status/:task_id` | GET | 查询爬取进度（TODO） |

### 竞品管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/competitors` | GET | 获取竞品列表 |
| `/api/competitors/:id/sources` | GET | 获取数据源列表 |

### 分析模块（TODO）

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/analyze/extract` | POST | 信息提取 |
| `/api/analyze/compare` | POST | 对比分析 |
| `/api/analyze/swot` | POST | SWOT分析 |

### 报告模块（TODO）

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/report/generate` | POST | 生成报告 |
| `/api/report/:id` | GET | 查看报告 |
| `/api/report/export/:id` | GET | 导出报告 |

## 成本估算

### API成本（单次分析3个竞品）

#### 🆓 使用Groq（推荐）
- **数据源发现**: $0.03（Serper）
- **数据采集**: $0（Jina免费）
- **AI分析**: $0（Groq免费）
- **总计**: **$0.03/次** ✨

#### 使用OpenAI
- **数据源发现**: $0.03-0.08
- **数据采集**: $0-0.06
- **AI分析**: $0.50
- **总计**: ~$0.55-0.65/次

### 月度成本对比

| 方案 | 月费（100次分析） |
|------|------------------|
| **Groq** | **$3**（只需Serper） |
| DeepSeek | $7-10 |
| OpenAI | $55-65 |

## 开发计划

### Phase 1: MVP核心功能 ✅
- [x] 三层爬虫策略
- [x] 智能数据源发现
- [x] 基础信息提取
- [x] 简单对比分析
- [x] Markdown报告生成

### Phase 2: 增强分析能力（进行中）
- [ ] 完善信息提取维度
- [ ] SWOT自动生成
- [ ] 功能对比矩阵
- [ ] 用户口碑分析
- [ ] 图表可视化

### Phase 3: 监控与自动化
- [ ] 定时监控调度
- [ ] 变化检测算法
- [ ] 预警通知系统
- [ ] 历史数据对比

### Phase 4: 产品化
- [ ] Web前端界面
- [ ] 配置管理
- [ ] 报告管理和多格式导出
- [ ] 监控面板

## 注意事项

1. **API密钥安全**: 不要将 `.env` 文件提交到版本控制
2. **爬虫合规**: 遵守 robots.txt，尊重网站服务条款
3. **频率控制**: 避免过于频繁的请求，建议使用搜索缓存
4. **成本控制**: 优先使用免费API额度，设置请求限制

## 常见问题

### Q1: Firecrawl返回验证码怎么办？
A: 系统会自动降级到Jina或Playwright。对于微信公众号，建议使用短链接格式。

### Q2: 如何降低API成本？
A: 
- 优先使用免费API（Google/Bing每日免费额度）
- 启用搜索缓存（默认7天）
- 使用本地LLM模型（DeepSeek/Llama）

### Q3: 支持哪些平台？
A: 目前支持：微信公众号、小红书、知乎、抖音、淘宝、京东、B站、微博等，以及普通网站。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或联系维护者。

---

**版本**: v1.0.0  
**更新时间**: 2026-02-06
