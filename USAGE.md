# 自动化竞品调研工具 - 使用指南

## 快速上手

### 第一步：配置API密钥

1. 复制环境变量模板：
```bash
cp .env.example .env
```

2. 编辑 `.env` 文件，填入你的API密钥：

```env
# 必需的API密钥
FIRECRAWL_API_KEY=your_key_here     # Firecrawl API（推荐）
SERPER_API_KEY=your_key_here        # Serper搜索API（推荐）
OPENAI_API_KEY=your_key_here        # OpenAI GPT-4（必需）

# 可选的API密钥
GOOGLE_API_KEY=your_key_here        # Google搜索（可选）
BING_API_KEY=your_key_here          # Bing搜索（可选）
```

### 第二步：启动服务

```bash
# 安装依赖
go mod download

# 运行服务
go run main.go
```

服务启动后会自动创建必要的目录：
- `./data/` - 数据库文件
- `./storage/` - 爬取内容存储
- `./reports/` - 生成的报告

### 第三步：使用智能发现功能

#### 方式一：使用API

```bash
# 1. 开始智能发现
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "AI写作助手",
    "market": "中国",
    "competitor_count": 5,
    "depth": "standard"
  }'

# 返回任务ID
{
  "task_id": 1,
  "status": "processing",
  "progress": 0
}

# 2. 查询进度
curl http://localhost:8080/api/discover/status/1

# 3. 任务完成后确认
curl -X POST http://localhost:8080/api/discover/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": 1,
    "selected_competitors": ["Notion AI", "Jasper"],
    "save_as_config": true
  }'
```

#### 方式二：手动爬取

如果你已经知道竞品的URL，可以直接爬取：

```bash
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://notion.so",
    "competitor": "Notion AI",
    "source_type": "官网"
  }'
```

## 完整工作流程

### 场景1：从零开始的竞品分析

假设你想分析"在线协作工具"市场：

```bash
# 1. 智能发现竞品和数据源（约1分钟）
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "在线协作工具",
    "market": "中国",
    "competitor_count": 3,
    "source_types": ["官网", "评价"],
    "depth": "standard"
  }'

# 等待任务完成，查询结果
curl http://localhost:8080/api/discover/status/1

# 2. 系统自动找到：Notion、飞书、语雀
# 3. 确认并开始采集
curl -X POST http://localhost:8080/api/discover/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": 1,
    "selected_competitors": ["Notion", "飞书", "语雀"],
    "save_as_config": true
  }'

# 4. 系统会自动爬取所有数据源
# 5. 内容保存在 ./storage/ 目录
```

### 场景2：已知竞品URL的快速分析

如果你已经知道要分析的URL：

```bash
# 批量爬取多个URL
curl -X POST http://localhost:8080/api/crawl/single \
  -d '{"url": "https://notion.so", "competitor": "Notion"}'

curl -X POST http://localhost:8080/api/crawl/single \
  -d '{"url": "https://www.feishu.cn", "competitor": "飞书"}'

curl -X POST http://localhost:8080/api/crawl/single \
  -d '{"url": "https://www.yuque.com", "competitor": "语雀"}'
```

### 场景3：查看竞品列表和数据源

```bash
# 获取所有竞品
curl http://localhost:8080/api/competitors

# 获取某个竞品的数据源
curl http://localhost:8080/api/competitors/1/sources
```

## 平台支持说明

### 已支持的平台

| 平台 | 支持状态 | 推荐方法 | 注意事项 |
|------|---------|---------|---------|
| 普通网站 | ✅ 完整支持 | Firecrawl | - |
| 微信公众号 | ✅ 完整支持 | Firecrawl | 使用短链接 `/s/xxxxx` |
| 知乎 | ✅ 完整支持 | Firecrawl/Jina | - |
| 淘宝/天猫 | ⚠️ 基础支持 | Jina | 需要处理反爬 |
| 京东 | ⚠️ 基础支持 | Jina | 需要处理反爬 |
| B站 | ✅ 完整支持 | Firecrawl | - |
| 微博 | ⚠️ 基础支持 | Jina | - |
| 小红书 | ❌ 需要登录 | Playwright | 暂未实现 |
| 抖音 | ❌ 需要登录 | Playwright | 暂未实现 |

### 三层爬虫策略说明

系统会自动按优先级尝试：

1. **Firecrawl** (首选)
   - AI驱动，智能处理JavaScript
   - 自动绕过常见反爬
   - 成功率最高
   - 需要API密钥

2. **Jina Reader** (备选)
   - 完全免费
   - 适合简单页面
   - 不需要API密钥

3. **Playwright** (兜底)
   - 真实浏览器
   - 支持登录态
   - 暂未实现，将在后续版本支持

## 数据存储说明

### 爬取内容存储

```
storage/
└── 20260206_Notion_AI_首页/
    ├── content.md          # Markdown内容
    ├── img_01.jpg          # 下载的图片
    ├── img_02.jpg
    └── ...
```

### 数据库存储

SQLite数据库文件：`./data/competitive.db`

包含以下表：
- `discovery_tasks` - 发现任务
- `search_caches` - 搜索缓存
- `competitors` - 竞品信息
- `data_sources` - 数据源
- `raw_contents` - 原始内容
- `parsed_data` - 解析结果
- `analysis_reports` - 分析报告
- `change_logs` - 变化日志
- `monitor_tasks` - 监控任务

## 常见错误处理

### 错误1：搜索API返回错误

```
错误信息: Serper返回错误 401
解决方案: 检查 SERPER_API_KEY 是否正确配置
```

### 错误2：爬取失败

```
错误信息: 所有爬虫都失败了
解决方案: 
1. 检查URL是否有效
2. 检查 FIRECRAWL_API_KEY 是否配置
3. 尝试使用Jina（免费）
```

### 错误3：图片下载失败

```
错误信息: 下载图片失败: 403
解决方案: 
- 小红书图片需要设置Referer（已自动处理）
- 部分图片可能需要登录，会被跳过
```

### 错误4：LLM解析失败

```
错误信息: OpenAI返回错误
解决方案:
1. 检查 OPENAI_API_KEY 是否正确
2. 检查API额度是否充足
3. 检查网络连接
```

## 性能优化建议

### 1. 启用搜索缓存

搜索结果会自动缓存7天，避免重复搜索：

```env
SEARCH_CACHE_DAYS=7
```

### 2. 控制并发数

搜索和爬取都支持并发，默认最多5个并发任务。

### 3. 使用免费API优先

- Jina Reader: 完全免费
- Google/Bing: 每日有免费额度
- Serper: 2500次免费搜索

### 4. 批量处理

对于大量URL，建议使用批量爬取接口（TODO）。

## API开发指南

### 添加自定义处理器

1. 在 `handlers/` 目录创建新文件
2. 实现处理函数
3. 在 `main.go` 中注册路由

示例：

```go
// handlers/custom.go
package handlers

import "github.com/gin-gonic/gin"

func CustomHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "custom"})
}

// main.go
api.GET("/custom", handlers.CustomHandler)
```

### 扩展爬虫策略

在 `crawler/crawler.go` 中添加新的爬虫实现：

```go
type CustomCrawler struct {
    // 配置
}

func (c *CustomCrawler) Crawl(url string, platform *PlatformInfo) (*CrawlResult, error) {
    // 实现爬取逻辑
}
```

## 下一步

- [ ] 实现完整的报告生成功能
- [ ] 添加Web前端界面
- [ ] 实现监控和预警功能
- [ ] 支持Playwright浏览器自动化
- [ ] 添加更多数据可视化

## 获取帮助

- 查看 [README.md](README.md) 了解完整功能
- 查看 [Design.md](Design.md) 了解设计思路
- 提交 Issue 反馈问题
