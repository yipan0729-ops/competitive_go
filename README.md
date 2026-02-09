# 自动化竞品分析系统 🚀

> **一键完成竞品分析！从输入主题到生成报告，完全自动化！**

基于 Go + Gin 框架开发的全自动化竞品分析工具，支持智能数据源发现、批量数据采集、AI智能分析和专业报告生成。

---

## ✨ 核心特性

### 🎯 全流程自动化
- **一个API调用完成所有步骤**
- 自动发现竞品 → 批量爬取 → AI分析 → 生成报告
- 零人工干预，自动输出专业报告

### 🔍 智能数据源发现
- 支持多搜索引擎（Serper、Google、Bing）
- 自动发现竞品和数据源
- 智能链接分类和质量评分

### ⚡ 高效数据采集
- **批量并发爬取**
  - 支持1-10个并发
  - 异步执行，立即返回
  - 自动错误处理和重试
- **三层策略自动降级**
  - Firecrawl (AI驱动，首选)
  - Jina Reader (免费，备选)
  - Playwright (浏览器自动化，兜底)
- 支持多平台识别（微信、小红书、知乎等）
- 自动下载并本地化图片
- Markdown格式保存

### 🤖 AI智能分析
- **产品信息提取**
  - 公司信息、产品定位
  - 目标用户群体
  - 核心功能列表
  - 价格策略分析
- **SWOT自动分析**
  - 优势/劣势/机会/威胁
  - 证据支持和影响评估
  - 战略建议
- **支持多种LLM**
  - Ollama（本地，免费）⭐ 推荐
  - OpenAI、DeepSeek、Groq
  - 智谱AI、通义千问

### 📊 专业报告生成
- 完整的Markdown报告
- 功能对比矩阵
- 价格策略分析表
- 完整SWOT分析
- 战略建议输出
- 数据来源附录

### 💾 数据管理
- SQLite数据库存储（纯Go驱动，无需CGO）
- 历史数据追踪
- 内容变化检测

---

## 🚀 快速开始

### 方式1: 一键自动化（推荐）⭐

```powershell
# 1. 启动服务
go run main.go

# 2. 运行自动化分析（只需这一个命令！）
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

# 3. 查询进度
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1

# 4. 查看报告
ls reports/
```

**系统会自动**:
- ✅ 发现竞品
- ✅ 搜索数据源
- ✅ 批量爬取网页
- ✅ AI提取产品信息
- ✅ 生成SWOT分析
- ✅ 输出完整报告

---

### 方式2: 完整安装配置

#### 1. 环境要求

- **Go 1.21+** （已测试：1.21-1.24）
- **Ollama**（本地LLM，完全免费）⭐ 推荐

#### 2. 安装依赖

```bash
# 克隆项目
git clone <项目地址>
cd Competitive_go

# 设置Go代理（国内推荐）
go env -w GOPROXY=https://goproxy.cn,direct

# 安装依赖
go mod download
```

#### 3. 配置环境变量

复制配置文件：

```bash
# Windows
copy .env.ollama .env

# Linux/Mac
cp .env.ollama .env
```

**推荐配置** (使用Ollama本地LLM)：

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

#### 4. 安装Ollama

**Windows**:
```powershell
# 方式1: 从官网下载安装
# https://ollama.com/download/windows

# 方式2: 使用winget
winget install Ollama.Ollama

# 下载推荐模型
ollama pull qwen2.5:7b
```

**Linux/Mac**:
```bash
# 安装Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 下载推荐模型
ollama pull qwen2.5:7b
```

**推荐模型**:
- `qwen2.5:7b` - 中文效果好，速度快 ⭐
- `deepseek-r1:8b` - 推理能力强
- `llama3.2:3b` - 最快，适合测试

#### 5. 启动服务

```bash
go run main.go
```

看到以下输出表示成功：
```
服务器启动在端口 8080
```

---

## 📖 使用示例

### 示例1: 全自动分析

```powershell
# 一键分析AI创作工具市场
$body = @{
    topic = "AI创作工具"
    market = "中国"
    competitor_count = 3
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

### 示例2: 手动分步操作

```powershell
# 1. 发现竞品
$body = @{topic="项目管理工具"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/discover/search `
    -Method POST -Body $body -ContentType "application/json"

# 2. 查询状态
Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/1

# 3. 爬取内容
$body = @{
    url = "https://www.notion.so"
    competitor = "Notion"
} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/crawl/single `
    -Method POST -Body $body -ContentType "application/json"

# 4. AI分析
$body = @{competitor_id=1} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/analyze/competitor `
    -Method POST -Body $body -ContentType "application/json" -TimeoutSec 600

# 5. 生成报告
$body = @{competitor_ids=@(1);topic="项目管理"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8080/api/report/generate `
    -Method POST -Body $body -ContentType "application/json"
```

---

## 📁 项目结构

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
│   └── classifier.go       # 链接分类
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
└── go.mod                  # Go模块依赖
```

---

## 📡 API接口

详细API文档请查看 [API.md](./API.md)

### 核心接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/auto/analysis` | POST | 全流程自动化⭐ |
| `/api/discover/search` | POST | 开始发现任务 |
| `/api/discover/status/:id` | GET | 查询任务状态 |
| `/api/crawl/single` | POST | 爬取单个URL |
| `/api/crawl/batch` | POST | 批量爬取 |
| `/api/analyze/competitor` | POST | AI分析竞品 |
| `/api/report/generate` | POST | 生成报告 |
| `/api/competitors` | GET | 获取竞品列表 |

---

## 💰 成本估算

### 使用Ollama（推荐）🆓

- **数据源发现**: $0.03（Serper）或 $0（不用搜索API）
- **数据采集**: $0（Jina免费）
- **AI分析**: $0（Ollama本地免费）
- **总计**: **$0/次** ✨

### 使用Groq（云端免费）🆓

- **数据源发现**: $0.03（Serper）
- **数据采集**: $0（Jina免费）
- **AI分析**: $0（Groq免费）
- **总计**: **$0.03/次** ✨

### 使用OpenAI

- **数据源发现**: $0.03-0.08
- **数据采集**: $0-0.06
- **AI分析**: $0.50
- **总计**: ~$0.55-0.65/次

---

## 🛠️ 常见问题

### Q1: Ollama连接失败？

**检查步骤**:
```powershell
# 1. 确认Ollama运行中
ollama list

# 2. 测试连接
curl http://localhost:11434/api/tags

# 3. 检查.env配置
OPENAI_BASE_URL=http://localhost:11434
OPENAI_API_KEY=ollama
```

### Q2: AI分析返回JSON解析错误？

**解决方案**: 已在v1.0.0中修复。确保使用最新代码。

### Q3: 爬取失败？

**可能原因**:
- 网站需要登录
- 反爬虫保护
- URL不可访问

**解决方案**:
- 系统会自动重试3次
- 自动降级到其他爬取策略
- 检查URL是否可访问

### Q4: 超时错误？

**解决方案**:
```powershell
# AI分析需要较长时间，增加超时设置
-TimeoutSec 600  # 10分钟
```

---

## 📚 相关文档

- [API.md](./API.md) - 完整API文档⭐
- [QUICKSTART.md](./QUICKSTART.md) - 快速开始指南
- [HOW_TO_RUN.md](./HOW_TO_RUN.md) - 详细运行指南
- [OLLAMA_SETUP.md](./OLLAMA_SETUP.md) - Ollama配置说明
- [FREE_LLM_ALTERNATIVES.md](./FREE_LLM_ALTERNATIVES.md) - 免费LLM选项
- [GROQ_SETUP.md](./GROQ_SETUP.md) - Groq配置说明
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 架构设计文档
- [CHANGELOG.md](./CHANGELOG.md) - 更新日志

---

## 🎯 开发计划

### ✅ Phase 1: MVP核心功能
- [x] 三层爬虫策略
- [x] 智能数据源发现
- [x] 批量并发爬取
- [x] AI产品信息提取
- [x] SWOT自动分析
- [x] Markdown报告生成
- [x] 全流程自动化

### 🚧 Phase 2: 功能增强（进行中）
- [ ] 功能对比矩阵优化
- [ ] 用户口碑分析
- [ ] 图表可视化
- [ ] 多语言支持

### 📋 Phase 3: 监控与自动化
- [ ] 定时监控调度
- [ ] 变化检测算法
- [ ] 预警通知系统
- [ ] 历史数据对比

### 🌟 Phase 4: 产品化
- [ ] Web前端界面
- [ ] 配置管理面板
- [ ] 报告多格式导出
- [ ] 监控仪表板

---

## ⚠️ 注意事项

1. **API密钥安全**: 不要将 `.env` 文件提交到版本控制
2. **爬虫合规**: 遵守 robots.txt，尊重网站服务条款
3. **频率控制**: 避免过于频繁的请求
4. **成本控制**: 优先使用免费API，设置请求限制

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 贡献方式

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

MIT License

---

## 📞 联系方式

如有问题，请提交 Issue。

---

**版本**: v1.0.0  
**更新时间**: 2026-02-09  
**维护状态**: ✅ 活跃维护
