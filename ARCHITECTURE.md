# 项目架构说明

## 目录结构

```
Competitive_go/
├── main.go                     # 主程序入口，Gin路由配置
├── go.mod                      # Go模块依赖管理
├── .env.example                # 环境变量模板
├── .gitignore                  # Git忽略文件
├── README.md                   # 项目说明文档
├── USAGE.md                    # 使用指南
├── Design.md                   # 设计文档（原始需求）
├── source.md                   # 灵感来源
├── start.sh                    # Linux/Mac启动脚本
├── start.bat                   # Windows启动脚本
│
├── config/                     # 配置模块
│   └── config.go               # 应用配置加载
│
├── models/                     # 数据模型
│   └── models.go               # 数据库模型定义
│
├── database/                   # 数据库层
│   └── database.go             # 数据库初始化和连接
│
├── crawler/                    # 数据采集模块
│   ├── platform.go             # 平台识别（微信/小红书/知乎等）
│   ├── crawler.go              # 三层爬虫策略实现
│   └── saver.go                # 内容保存和图片下载
│
├── discovery/                  # 智能数据源发现模块
│   ├── search.go               # 搜索引擎集成（Serper/Google/Bing）
│   ├── manager.go              # 搜索管理器和查询生成
│   └── classifier.go           # 链接分类和质量评分
│
├── ai/                         # AI分析模块
│   ├── llm.go                  # LLM客户端（OpenAI）
│   └── extractor.go            # 信息提取器（竞品/产品/SWOT）
│
├── report/                     # 报告生成模块
│   └── generator.go            # 报告生成器（Markdown格式）
│
├── handlers/                   # HTTP处理器
│   └── handlers.go             # API请求处理逻辑
│
└── examples/                   # 示例代码
    └── test.go                 # API测试示例
```

## 核心模块说明

### 1. 配置模块 (config/)

**职责**: 加载和管理应用配置

- 从环境变量读取API密钥
- 设置服务器参数
- 初始化存储路径

**关键文件**: `config.go`

### 2. 数据模型 (models/)

**职责**: 定义数据库表结构

- `DiscoveryTask`: 数据源发现任务
- `SearchCache`: 搜索结果缓存
- `Competitor`: 竞品信息
- `DataSource`: 数据源链接
- `RawContent`: 原始爬取内容
- `ParsedData`: AI解析结果
- `AnalysisReport`: 分析报告
- `ChangeLog`: 变化日志
- `MonitorTask`: 监控任务

**关键文件**: `models.go`

### 3. 数据库层 (database/)

**职责**: 数据库初始化和管理

- SQLite连接
- 自动表迁移
- GORM ORM封装

**关键文件**: `database.go`

### 4. 爬虫模块 (crawler/)

**职责**: 实现三层爬虫策略

**文件说明**:
- `platform.go`: 识别URL所属平台（微信/小红书/知乎等）
- `crawler.go`: 实现Firecrawl、Jina、Playwright三层策略
- `saver.go`: 保存内容为Markdown，下载图片到本地

**核心逻辑**:
```
识别平台 → 选择策略 → 爬取内容 → 保存文件
```

### 5. 数据源发现模块 (discovery/)

**职责**: 智能发现竞品和数据源

**文件说明**:
- `search.go`: 集成Serper/Google/Bing搜索引擎
- `manager.go`: 管理搜索任务，生成查询语句
- `classifier.go`: 对搜索结果分类和质量评分

**核心逻辑**:
```
生成查询 → 并发搜索 → 提取竞品 → 搜索数据源 → 分类评分
```

### 6. AI分析模块 (ai/)

**职责**: 使用LLM进行智能分析

**文件说明**:
- `llm.go`: OpenAI客户端封装
- `extractor.go`: 实现竞品提取、产品信息提取、SWOT分析

**核心功能**:
- 从搜索结果中提取竞品名称
- 从网页内容中提取产品信息
- 自动生成SWOT分析
- 支持结构化JSON输出

### 7. 报告生成模块 (report/)

**职责**: 生成竞品分析报告

**文件说明**:
- `generator.go`: Markdown报告生成器

**报告结构**:
- 执行摘要
- 竞品概览
- 功能对比矩阵
- 价格策略分析
- SWOT分析
- 战略建议

### 8. HTTP处理器 (handlers/)

**职责**: 处理API请求

**主要接口**:
- `/api/discover/search` - 开始智能发现
- `/api/discover/status/:task_id` - 查询任务状态
- `/api/discover/confirm` - 确认并保存
- `/api/crawl/single` - 爬取单个URL
- `/api/competitors` - 获取竞品列表

## 数据流设计

### 智能发现模式流程

```
用户输入主题
    ↓
生成搜索查询 (discovery/manager.go)
    ↓
并发搜索多个引擎 (discovery/search.go)
    ↓
提取搜索结果 (discovery/manager.go)
    ↓
【可选】LLM提取竞品名称 (ai/extractor.go)
    ↓
为每个竞品搜索数据源 (discovery/manager.go)
    ↓
链接分类和质量评分 (discovery/classifier.go)
    ↓
展示给用户确认
    ↓
保存到数据库 (database/)
    ↓
【后续】批量爬取 (crawler/)
```

### 爬取流程

```
接收URL请求
    ↓
识别平台类型 (crawler/platform.go)
    ↓
选择爬虫策略 (crawler/crawler.go)
    ├─ 尝试 Firecrawl
    ├─ 失败则尝试 Jina
    └─ 失败则尝试 Playwright
    ↓
获取Markdown内容
    ↓
提取图片URL (crawler/saver.go)
    ↓
下载图片到本地
    ↓
替换图片链接
    ↓
保存Markdown文件
    ↓
保存到数据库
```

## 技术选型说明

### 为什么选择Go？

1. **高性能**: 并发爬取和搜索
2. **简单部署**: 单个二进制文件
3. **强类型**: 减少运行时错误
4. **丰富的库**: Gin、GORM等成熟框架

### 为什么选择Gin？

1. **性能优秀**: 路由速度快
2. **中间件丰富**: CORS、日志等
3. **文档完善**: 易于上手
4. **社区活跃**: 问题易解决

### 为什么选择SQLite？

1. **零配置**: 无需安装数据库服务
2. **文件存储**: 易于备份和迁移
3. **性能足够**: 单用户场景下性能充足
4. **可升级**: 后续可迁移到PostgreSQL

### 为什么选择三层爬虫策略？

1. **Firecrawl**: AI驱动，成功率最高，但收费
2. **Jina**: 免费，适合简单页面
3. **Playwright**: 兜底方案，真实浏览器

**优势**: 自动降级，确保最大成功率

## 扩展性设计

### 1. 添加新的搜索引擎

实现 `SearchEngine` 接口：

```go
type CustomSearchEngine struct {
    APIKey string
}

func (s *CustomSearchEngine) Search(query string, numResults int) ([]SearchResult, error) {
    // 实现搜索逻辑
}
```

### 2. 添加新的爬虫策略

实现 `Crawler` 接口：

```go
type CustomCrawler struct {}

func (c *CustomCrawler) Crawl(url string, platform *PlatformInfo) (*CrawlResult, error) {
    // 实现爬取逻辑
}
```

### 3. 添加新的平台支持

在 `platform.go` 中添加识别规则：

```go
if strings.Contains(host, "newplatform.com") {
    return &PlatformInfo{
        Name:       "新平台",
        NeedsLogin: false,
        Priority:   1,
    }, nil
}
```

## 性能优化

### 1. 并发控制

使用 `sync.WaitGroup` 和信号量控制并发数：

```go
semaphore := make(chan struct{}, maxWorkers)
```

### 2. 搜索缓存

搜索结果缓存7天，避免重复搜索：

```go
type SearchCache struct {
    Query     string
    Results   JSONB
    ExpiresAt time.Time
}
```

### 3. 连接池

使用 `http.Client` 复用连接：

```go
client := &http.Client{Timeout: 30 * time.Second}
```

## 安全考虑

### 1. API密钥管理

- 使用环境变量存储密钥
- `.env` 文件不提交到Git
- 支持多种密钥源（环境变量/配置文件）

### 2. 输入验证

- 使用 Gin 的 `binding` 标签验证请求
- 检查URL格式
- 防止SQL注入（使用GORM）

### 3. 错误处理

- 统一的错误响应格式
- 详细的日志记录
- 优雅的错误降级

## 部署建议

### 开发环境

```bash
go run main.go
```

### 生产环境

```bash
# 编译
go build -o competitive-analyzer

# 使用systemd管理（Linux）
sudo systemctl start competitive-analyzer

# 使用Docker
docker build -t competitive-analyzer .
docker run -d -p 8080:8080 competitive-analyzer
```

## 监控和日志

### 日志

使用 Gin 的默认日志中间件：

```go
r := gin.Default()  // 包含日志和恢复中间件
```

### 健康检查

```bash
curl http://localhost:8080/health
```

## 后续优化方向

1. **前端界面**: Vue.js或React
2. **实时通知**: WebSocket
3. **任务队列**: Redis + Celery
4. **分布式**: 多节点部署
5. **缓存层**: Redis缓存
6. **数据库**: 迁移到PostgreSQL
7. **监控**: Prometheus + Grafana
8. **日志**: ELK Stack

## 参考资料

- [Gin文档](https://gin-gonic.com/)
- [GORM文档](https://gorm.io/)
- [Firecrawl API](https://docs.firecrawl.dev/)
- [Serper API](https://serper.dev/)
- [OpenAI API](https://platform.openai.com/docs)
