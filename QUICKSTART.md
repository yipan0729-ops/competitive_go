# 快速入门指南

## 5分钟上手

### 第一步：获取API密钥 (2分钟)

你需要至少以下API密钥之一：

1. **Serper API** (推荐，用于搜索)
   - 访问 https://serper.dev/
   - 注册并获取API密钥
   - 免费额度：2500次搜索

2. **OpenAI API** (必需，用于AI分析)
   - 访问 https://platform.openai.com/
   - 创建API密钥
   - 需要充值使用

3. **Firecrawl API** (可选，用于爬取)
   - 访问 https://firecrawl.dev/
   - 注册并获取API密钥
   - 免费额度：500页/月

### 第二步：配置项目 (1分钟)

```bash
# 1. 进入项目目录
cd D:\Code\Competitive_go

# 2. 复制配置文件
copy .env.example .env

# 3. 编辑 .env 文件，填入你的API密钥
notepad .env
```

最小配置（只需这两个）：
```env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=你的OpenAI密钥
```

完整配置（推荐）：
```env
FIRECRAWL_API_KEY=你的Firecrawl密钥
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=你的OpenAI密钥
```

### 第三步：启动服务 (2分钟)

**Windows用户**：
```bash
# 双击运行
start.bat

# 或者命令行运行
.\start.bat
```

**Linux/Mac用户**：
```bash
# 给脚本执行权限
chmod +x start.sh

# 运行
./start.sh
```

**手动启动**：
```bash
# 下载依赖
go mod download

# 运行
go run main.go
```

看到以下输出表示启动成功：
```
数据库初始化成功
服务器启动在端口 8080
```

## 第一次使用

### 场景1：我想分析"AI写作助手"市场

打开新的命令行窗口，执行：

```bash
# 1. 开始智能发现（约1分钟）
curl -X POST http://localhost:8080/api/discover/search ^
  -H "Content-Type: application/json" ^
  -d "{\"topic\":\"AI写作助手\",\"market\":\"中国\",\"competitor_count\":3,\"depth\":\"quick\"}"

# 会返回任务ID，比如：
# {"task_id":1,"status":"processing","progress":0}

# 2. 等待10秒后查询结果
timeout 10
curl http://localhost:8080/api/discover/status/1

# 3. 看到竞品列表后，确认保存
curl -X POST http://localhost:8080/api/discover/confirm ^
  -H "Content-Type: application/json" ^
  -d "{\"task_id\":1,\"selected_competitors\":[\"Notion AI\",\"Jasper\"],\"save_as_config\":true}"
```

### 场景2：我知道竞品网址，直接爬取

```bash
# 爬取单个网址
curl -X POST http://localhost:8080/api/crawl/single ^
  -H "Content-Type: application/json" ^
  -d "{\"url\":\"https://notion.so\",\"competitor\":\"Notion AI\",\"source_type\":\"官网\"}"

# 成功后内容保存在 storage/ 目录
```

### 场景3：查看已有的竞品

```bash
# 获取竞品列表
curl http://localhost:8080/api/competitors

# 获取某个竞品的数据源
curl http://localhost:8080/api/competitors/1/sources
```

## 使用测试工具

我们提供了一个测试工具，可以快速测试所有功能：

```bash
# 进入examples目录
cd examples

# 运行测试
go run test.go
```

测试工具会自动：
1. 检查服务健康状态
2. 创建一个发现任务
3. 爬取一个测试URL
4. 获取竞品列表

## 查看结果

### 爬取内容位置

```
storage/
└── 20260206_Notion_AI_首页/
    ├── content.md          # Markdown内容
    ├── img_01.jpg          # 图片1
    ├── img_02.jpg          # 图片2
    └── ...
```

### 数据库

使用SQLite浏览器查看 `data/competitive.db`

推荐工具：
- [DB Browser for SQLite](https://sqlitebrowser.org/)
- [DBeaver](https://dbeaver.io/)

## 常见问题

### Q1: 启动失败，提示"未找到 .env 文件"

```bash
# 解决方案：复制配置文件
copy .env.example .env
# 然后编辑 .env 填入你的API密钥
```

### Q2: 搜索失败，提示"API Key未配置"

```bash
# 检查 .env 文件中是否配置了 SERPER_API_KEY
# 至少需要配置以下之一：
# - SERPER_API_KEY
# - GOOGLE_API_KEY
# - BING_API_KEY
```

### Q3: 爬取失败，提示"所有爬虫都失败了"

可能原因：
1. Firecrawl API密钥未配置或额度用完
2. 目标网站有反爬虫保护
3. URL格式不正确

解决方案：
1. 配置 FIRECRAWL_API_KEY（推荐）
2. 使用免费的Jina（会自动尝试）
3. 检查URL是否可访问

### Q4: Windows下双击start.bat闪退

原因：可能是Go未安装或环境变量未配置

解决方案：
1. 安装Go 1.21+
2. 在命令行中运行 `start.bat` 查看错误信息

### Q5: 端口8080被占用

修改 `.env` 文件中的端口：
```env
SERVER_PORT=8081
```

## 下一步

### 学习完整功能

阅读详细文档：
- [README.md](README.md) - 项目介绍
- [USAGE.md](USAGE.md) - 详细使用指南
- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构说明

### 定制开发

1. 添加新的平台支持：编辑 `crawler/platform.go`
2. 添加新的搜索引擎：实现 `SearchEngine` 接口
3. 自定义报告格式：修改 `report/generator.go`

### 社区交流

- 提交Issue报告问题
- 提交PR贡献代码
- 分享你的使用心得

## 提示和技巧

### 1. 节省API成本

- 搜索结果会缓存7天，相同查询不会重复搜索
- 优先使用免费API（Jina、Google/Bing免费额度）
- 使用 `depth: "quick"` 减少搜索次数

### 2. 提高爬取成功率

- 微信公众号使用短链接格式（/s/xxxxx）
- 对于需要登录的平台，优先使用API获取数据
- 设置合理的请求间隔，避免被封IP

### 3. 批量处理

- 使用智能发现功能，一次找到所有竞品
- 编写脚本批量调用 `/api/crawl/single`
- 使用数据库导出功能批量管理

### 4. 监控提醒

- 定期运行爬取任务（使用cron或Windows计划任务）
- 对比历史数据发现变化
- 设置邮件通知（TODO功能）

## 演示视频（待制作）

- [ ] 5分钟快速上手
- [ ] 智能发现功能演示
- [ ] 爬取和分析完整流程
- [ ] 生成分析报告

## 获取帮助

如果遇到问题：

1. 查看日志输出
2. 检查 .env 配置
3. 阅读错误信息
4. 搜索常见问题
5. 提交Issue

---

**祝你使用愉快！** 🎉
