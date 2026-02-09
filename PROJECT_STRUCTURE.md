# 项目文件清单

> 整理后的项目结构说明

---

## 📁 核心代码

### 源代码目录
```
├── main.go                      # 程序入口
├── config/
│   └── config.go               # 配置管理
├── models/
│   └── models.go               # 数据模型
├── database/
│   └── database.go             # 数据库操作
├── crawler/
│   ├── crawler.go              # 爬虫核心逻辑
│   ├── platform.go             # 平台识别
│   └── saver.go                # 内容保存
├── discovery/
│   ├── search.go               # 搜索引擎集成
│   ├── manager.go              # 搜索管理器
│   └── classifier.go           # 链接分类
├── ai/
│   ├── llm.go                  # LLM客户端
│   └── extractor.go            # AI提取器
├── report/
│   └── generator.go            # 报告生成器
└── handlers/
    └── handlers.go             # HTTP处理器
```

---

## 📖 文档

### 核心文档（必读）
- **README.md** - 项目概述和快速开始 ⭐
- **API.md** - 完整API文档 ⭐
- **QUICKSTART.md** - 快速开始指南
- **HOW_TO_RUN.md** - 详细运行指南
- **CHANGELOG.md** - 更新日志

### 配置文档
- **OLLAMA_SETUP.md** - Ollama配置说明
- **GROQ_SETUP.md** - Groq配置说明
- **FREE_LLM_ALTERNATIVES.md** - 免费LLM选项
- **FREE_LLM_QUICKSTART.md** - 免费方案快速开始

### 技术文档
- **ARCHITECTURE.md** - 架构设计文档

---

## ⚙️ 配置文件

### 环境配置
- **.env.example** - 配置模板
- **.env.ollama** - Ollama配置示例
- **.env.groq** - Groq配置示例
- **.gitignore** - Git忽略规则

### 启动脚本
- **start.sh** / **start.bat** - 快速启动脚本
- **setup-ollama.sh** / **setup-ollama.bat** - Ollama安装脚本
- **setup-groq.sh** / **setup-groq.bat** - Groq配置脚本

---

## 📦 依赖管理

- **go.mod** - Go模块定义
- **go.sum** - 依赖版本锁定

---

## 🗂️ 运行时目录

这些目录在首次运行时自动创建：

```
├── data/                        # 数据库文件
│   └── competitive.db          # SQLite数据库
├── storage/                     # 爬取内容
│   └── crawled/
│       └── [竞品名]/
│           ├── *.md            # Markdown内容
│           └── images/         # 图片资源
└── reports/                     # 生成的报告
    └── *.md                    # 分析报告
```

---

## 🗑️ 已清理的文件

以下临时和过时文件已删除：

### 临时测试文件
- test.ps1
- test-competitor-extraction.ps1
- test-ai-analysis.ps1
- test-auto-workflow.ps1
- analysis-4.json

### 过时文档
- BUGFIX_AUTO_ANALYSIS.md
- COMPLETION_SUMMARY.md
- FINAL_TEST.md
- HOW_TO_FIX_OLLAMA.md
- OLLAMA_TIMEOUT_FIX.md
- PROJECT_COMPLETED.md
- QUICK_START_FIXED.md
- QUICKSTART_V2.md
- MARKET_PARAM_EXPLAINED.md
- GO_1.24_NOTES.md
- GROQ_FORBIDDEN_FIX.md
- GROQ_CHECKLIST.md
- API_GUIDE.md
- API_QUICKREF.md
- source.md
- USAGE.md
- VERSION.md
- Design.md
- USE_GROQ.md

---

## 📊 项目统计

### 代码文件
- Go源代码: 13个文件
- 总行数: ~2500行

### 文档文件
- 核心文档: 5个
- 配置文档: 4个
- 技术文档: 1个

### 配置文件
- 环境配置: 4个
- 启动脚本: 6个

---

## 🎯 使用建议

### 新用户
1. 阅读 **README.md** 了解项目
2. 按照 **QUICKSTART.md** 快速开始
3. 参考 **API.md** 使用接口

### 开发者
1. 阅读 **ARCHITECTURE.md** 了解架构
2. 查看 **CHANGELOG.md** 了解更新
3. 参考源代码注释进行开发

### 部署运维
1. 配置 **.env** 文件
2. 使用 **start.sh/bat** 启动服务
3. 查看 **HOW_TO_RUN.md** 解决问题

---

## 📝 文档维护

### 更新原则
- README.md: 项目有重大变更时更新
- API.md: 接口变更时立即更新
- CHANGELOG.md: 每个版本发布时更新
- 其他文档: 按需更新

### 版本管理
- 所有文档包含更新日期
- 重要变更记录在CHANGELOG中
- 保持文档与代码同步

---

**文档更新**: 2026-02-09  
**项目版本**: v1.0.0  
**维护状态**: ✅ 活跃维护
