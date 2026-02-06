# 🚀 项目运行指南

## ✅ 项目完整性检查

### 已完成的模块

| 模块 | 状态 | 文件数 | 说明 |
|------|------|--------|------|
| 📦 配置模块 | ✅ 完成 | 1 | config/config.go |
| 🗄️ 数据模型 | ✅ 完成 | 1 | models/models.go (9张表) |
| 💾 数据库层 | ✅ 完成 | 1 | database/database.go |
| 🕷️ 爬虫模块 | ✅ 完成 | 3 | crawler/ (三层策略) |
| 🔍 数据源发现 | ✅ 完成 | 3 | discovery/ (搜索引擎集成) |
| 🤖 AI分析 | ✅ 完成 | 2 | ai/ (LLM客户端+提取器) |
| 📊 报告生成 | ✅ 完成 | 1 | report/generator.go |
| 🌐 API接口 | ✅ 完成 | 1 | handlers/handlers.go |
| 🎯 主程序 | ✅ 完成 | 1 | main.go |

### 文档完整性

| 文档 | 状态 | 说明 |
|------|------|------|
| README.md | ✅ | 项目介绍和功能说明 |
| QUICKSTART.md | ✅ | 5分钟快速上手 |
| USAGE.md | ✅ | 详细使用指南 |
| ARCHITECTURE.md | ✅ | 技术架构文档 |
| VERSION.md | ✅ | 依赖版本说明 |
| GO_1.24_NOTES.md | ✅ | Go 1.24兼容性说明 |

### 辅助工具

| 工具 | 状态 | 说明 |
|------|------|------|
| start.bat | ✅ | Windows启动脚本 |
| start.sh | ✅ | Linux/Mac启动脚本 |
| verify.bat | ✅ | Windows验证脚本 |
| verify.sh | ✅ | Linux/Mac验证脚本 |
| examples/test.go | ✅ | API测试示例 |

**结论：项目100%完整！所有核心功能和文档都已实现。** ✨

---

## 🎯 快速运行（3步）

### 方法一：一键启动（推荐）

#### Windows用户

```bash
# 1. 配置API密钥（首次使用）
copy .env.example .env
notepad .env

# 2. 一键启动
.\start.bat
```

#### Linux/Mac用户

```bash
# 1. 配置API密钥（首次使用）
cp .env.example .env
vi .env

# 2. 一键启动
chmod +x start.sh
./start.sh
```

`start.bat/sh` 会自动：
- ✅ 检查Go环境
- ✅ 创建必要目录（data/、storage/、reports/）
- ✅ 下载依赖
- ✅ 编译项目
- ✅ 启动服务

### 方法二：手动运行

```bash
# 1. 配置环境变量
copy .env.example .env
notepad .env  # 填入API密钥

# 2. 下载依赖（首次运行）
go mod download

# 3. 直接运行（不编译）
go run main.go

# 或者先编译再运行
go build -o competitive-analyzer.exe
.\competitive-analyzer.exe
```

---

## 📋 详细步骤说明

### 第一步：获取API密钥（5分钟）

项目需要以下API密钥（至少2个）：

#### 1. Serper API（推荐，用于搜索）
```
网址：https://serper.dev/
费用：2500次免费搜索
用途：智能发现竞品和数据源
```

#### 2. OpenAI API（必需，用于AI分析）
```
网址：https://platform.openai.com/
费用：按使用量付费
用途：信息提取、SWOT分析
```

#### 3. Firecrawl API（可选，用于爬取）
```
网址：https://firecrawl.dev/
费用：500页/月免费
用途：高质量网页抓取
```

### 第二步：配置项目（2分钟）

编辑 `.env` 文件：

```env
# === 必需配置 ===
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=你的OpenAI密钥

# === 推荐配置 ===
FIRECRAWL_API_KEY=你的Firecrawl密钥

# === 可选配置 ===
GOOGLE_API_KEY=你的Google密钥
BING_API_KEY=你的Bing密钥

# === 服务器配置（默认即可）===
SERVER_PORT=8080
GIN_MODE=release
```

### 第三步：运行项目（1分钟）

**选择A：使用启动脚本（推荐）**
```bash
.\start.bat  # Windows
./start.sh   # Linux/Mac
```

**选择B：手动运行**
```bash
go run main.go
```

看到以下输出表示成功：
```
数据库初始化成功
服务器启动在端口 8080
[GIN-debug] Listening and serving HTTP on :8080
```

---

## 🧪 验证安装

### 方法1：使用验证脚本

```bash
# Windows
.\verify.bat

# Linux/Mac
chmod +x verify.sh
./verify.sh
```

### 方法2：手动测试

**1. 健康检查**
```bash
curl http://localhost:8080/health
```

应该返回：
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

**2. 运行测试程序**
```bash
cd examples
go run test.go
```

---

## 📖 使用示例

### 场景1：智能发现竞品

```bash
# 发起发现任务
curl -X POST http://localhost:8080/api/discover/search ^
  -H "Content-Type: application/json" ^
  -d "{\"topic\":\"AI写作助手\",\"market\":\"中国\",\"competitor_count\":5}"

# 返回
{
  "task_id": 1,
  "status": "processing",
  "progress": 0
}

# 查询进度
curl http://localhost:8080/api/discover/status/1

# 确认并保存
curl -X POST http://localhost:8080/api/discover/confirm ^
  -H "Content-Type: application/json" ^
  -d "{\"task_id\":1,\"selected_competitors\":[\"Notion AI\"],\"save_as_config\":true}"
```

### 场景2：爬取单个URL

```bash
curl -X POST http://localhost:8080/api/crawl/single ^
  -H "Content-Type: application/json" ^
  -d "{\"url\":\"https://notion.so\",\"competitor\":\"Notion AI\"}"
```

### 场景3：查看竞品列表

```bash
curl http://localhost:8080/api/competitors
```

---

## 📁 查看结果

### 爬取的内容

```
storage/
└── 20260206_Notion_AI_首页/
    ├── content.md          # Markdown内容
    ├── img_01.jpg          # 下载的图片
    └── img_02.jpg
```

### 数据库

位置：`data/competitive.db`

推荐工具：
- [DB Browser for SQLite](https://sqlitebrowser.org/)
- [DBeaver](https://dbeaver.io/)

---

## ⚠️ 常见问题

### Q1: 启动失败，提示"go: no such file"

**原因**: Go未安装或不在PATH中

**解决**:
```bash
# 检查Go版本
go version

# 如果没有输出，需要安装Go或添加到PATH
```

### Q2: 依赖下载很慢

**原因**: 网络问题

**解决**: 配置Go代理
```bash
# Windows
set GOPROXY=https://goproxy.cn,direct
go mod download

# Linux/Mac
export GOPROXY=https://goproxy.cn,direct
go mod download
```

### Q3: 编译错误 "cannot find package"

**解决**:
```bash
go clean -modcache
go mod download
go mod tidy
```

### Q4: CGO错误（mattn/go-sqlite3）

**原因**: Windows缺少GCC

**快速解决**: 安装 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)

**或者使用纯Go驱动**（性能稍差）

### Q5: 端口8080被占用

**解决**: 修改 `.env` 中的端口
```env
SERVER_PORT=8081
```

---

## 🎓 下一步

### 学习使用

1. **快速上手**: 阅读 [QUICKSTART.md](QUICKSTART.md)
2. **详细教程**: 阅读 [USAGE.md](USAGE.md)
3. **API测试**: 运行 `examples/test.go`

### 开发扩展

1. **架构了解**: 阅读 [ARCHITECTURE.md](ARCHITECTURE.md)
2. **添加平台**: 编辑 `crawler/platform.go`
3. **自定义分析**: 修改 `ai/extractor.go`

### 获取帮助

- 📖 查看文档目录
- 🐛 提交Issue
- 💬 社区讨论

---

## 📊 系统要求

| 项目 | 要求 | 推荐 |
|------|------|------|
| Go版本 | 1.21+ | 1.24.8 |
| 内存 | 512MB | 1GB+ |
| 磁盘 | 100MB | 500MB+ |
| 网络 | 需要 | 需要 |

---

## 🎉 快速检查清单

运行前确认：

- [ ] Go已安装（`go version`）
- [ ] 已复制 `.env.example` 为 `.env`
- [ ] 已填写至少2个API密钥
- [ ] 端口8080未被占用
- [ ] 网络连接正常

全部勾选？执行：`.\start.bat` 或 `./start.sh`

---

## 💡 提示

1. **首次运行**会自动创建 `data/`、`storage/`、`reports/` 目录
2. **依赖下载**可能需要2-5分钟，请耐心等待
3. **API密钥**务必保密，不要提交到Git
4. **测试环境**可以只配置Serper和OpenAI两个密钥

---

**项目状态**: ✅ 100%完整，可以立即使用

**建议用时**: 
- 快速运行：5分钟
- 完整配置：10分钟
- 熟悉功能：30分钟

**祝你使用愉快！** 🎊
