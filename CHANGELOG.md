# 更新日志

## v2.0.0 (2026-02-09) - 重大更新 🎉

### ✨ 新增功能

#### 1. 批量爬取功能
- **接口**: `POST /api/crawl/batch`
- 支持批量爬取多个URL
- 可配置并发数量（1-10）
- 异步处理，立即返回
- 自动保存到数据库

#### 2. AI智能分析
- **接口**: `POST /api/analyze/competitor`
- 使用LLM提取产品信息
- 自动进行SWOT分析
- 支持自定义市场背景
- 结果结构化存储

**分析内容包括**:
- 产品基本信息（名称、公司、定位）
- 目标用户群体
- 核心功能列表
- 价格策略
- 优势/劣势/机会/威胁

#### 3. 报告自动生成
- **接口**: `POST /api/report/generate`
- 生成专业Markdown格式报告
- 包含功能对比矩阵
- 价格策略分析
- 完整SWOT分析
- 战略建议

**报告章节**:
1. 执行摘要
2. 竞品概览
3. 功能对比分析
4. 价格策略分析
5. SWOT分析
6. 战略建议
7. 附录（数据来源）

#### 4. 全流程自动化 🔥
- **接口**: `POST /api/auto/analysis`
- **一键完成全部流程**！

**工作流程**:
```
输入主题 → 发现竞品 → 搜索数据源 → 批量爬取 → AI分析 → 生成报告
```

**特点**:
- 全自动执行，无需人工干预
- 实时进度跟踪
- 可配置每个步骤
- 异步执行，不阻塞
- 预计完成时间：5分钟

### 🚀 使用示例

#### 场景1：一键自动分析（最简单）

```powershell
$body = @{
    topic = "项目管理工具"
    market = "中国"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing
```

**就这么简单！系统会自动**:
1. ✅ 搜索5个竞品
2. ✅ 为每个竞品找3个数据源
3. ✅ 爬取所有网页（15个URL）
4. ✅ AI分析每个竞品
5. ✅ 生成完整报告

**5分钟后，您将获得**:
- 📊 完整的竞品分析报告（Markdown）
- 💾 所有爬取的原始内容
- 🤖 AI提取的结构化数据
- 📈 SWOT分析结果

---

#### 场景2：分步执行（更灵活）

**步骤1: 发现竞品**
```powershell
POST /api/discover/search
{ "topic": "canvas" }
# 返回: task_id = 1
```

**步骤2: 批量爬取**
```powershell
POST /api/crawl/batch
{
  "urls": [...],  # 从步骤1获取
  "concurrent": 5
}
```

**步骤3: AI分析**
```powershell
POST /api/analyze/competitor
{ "competitor_id": 1 }
```

**步骤4: 生成报告**
```powershell
POST /api/report/generate
{ "competitor_ids": [1,2,3], "topic": "canvas" }
```

---

### 📊 性能对比

| 操作 | v1.0 | v2.0 | 提升 |
|------|------|------|------|
| 发现竞品 | ✅ 手动 | ✅ 自动 | - |
| 爬取5个URL | 5次API调用 | 1次API调用 | **5x** |
| 分析竞品 | ❌ 不支持 | ✅ AI自动 | ∞ |
| 生成报告 | ❌ 不支持 | ✅ 自动生成 | ∞ |
| 完整流程 | 30分钟+手动 | **5分钟全自动** | **6x+** |

---

### 🛠️ 技术改进

#### 1. 并发控制
- 使用信号量（semaphore）控制并发
- 避免同时发起过多请求
- 防止API速率限制

#### 2. 异步处理
- 所有耗时操作异步执行
- 立即返回task_id
- 前端轮询进度

#### 3. 错误处理
- 单个爬取失败不影响整体
- 详细的错误日志
- 自动重试机制

#### 4. 数据结构化
- 产品信息JSON格式
- SWOT分析结构化
- 便于后续查询和展示

---

### 📝 API变更

#### 新增接口
- `POST /api/crawl/batch` - 批量爬取
- `POST /api/analyze/competitor` - AI分析
- `POST /api/report/generate` - 生成报告
- `POST /api/auto/analysis` - 全流程自动化 🔥

#### 保持兼容
- 所有v1.0接口保持不变
- 向后兼容

---

### 🎯 推荐使用方式

#### 新用户
直接使用 **全流程自动化接口**：
```powershell
POST /api/auto/analysis
{ "topic": "你的主题" }
```

#### 高级用户
使用**分步接口**获得更多控制：
1. `POST /api/discover/search` - 发现
2. 人工筛选结果
3. `POST /api/crawl/batch` - 批量爬取
4. `POST /api/analyze/competitor` - 分析
5. `POST /api/report/generate` - 报告

---

### 📚 文档更新

- ✅ `API.md` - 完整API文档更新
- ✅ `CHANGELOG.md` - 本文件
- ✅ 所有示例使用PowerShell和cURL双版本

---

### 🐛 Bug修复

- 修复CGO依赖问题（切换到纯Go SQLite）
- 修复Ollama API调用超时
- 优化LLM响应解析

---

### ⚙️ 配置变更

**新增环境变量**:
- `LLM_TEMPERATURE` - LLM温度参数（默认：0.3）
- `LLM_MAX_TOKENS` - 最大Token数（默认：4000）

**推荐配置** (`.env`):
```env
# Ollama本地LLM（推荐）
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=deepseek-r1:8b  # 或 qwen2.5:7b
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
```

---

### 📦 依赖更新

**新增**:
- `github.com/glebarez/sqlite` - 纯Go SQLite驱动

**移除**:
- `gorm.io/driver/sqlite` - 需要CGO

---

### 🔄 升级指南

#### 从v1.0升级

1. **更新代码**
```bash
git pull origin main
```

2. **更新依赖**
```bash
go mod tidy
```

3. **无需修改配置**
- 所有现有配置保持兼容

4. **开始使用新功能**
```powershell
# 测试自动化接口
POST /api/auto/analysis
{ "topic": "测试主题" }
```

---

### 🎉 总结

**v2.0.0 是一个重大更新**，带来了：

✅ **批量爬取** - 效率提升5倍  
✅ **AI分析** - 自动提取+SWOT  
✅ **报告生成** - 专业Markdown报告  
✅ **全流程自动化** - 一键完成所有步骤  

**从此，竞品分析不再需要手工操作！**

输入主题 → 等5分钟 → 获得完整报告 🚀

---

### 🔮 下一步计划 (v3.0)

- [ ] 定时监控竞品变化
- [ ] 多维度对比图表
- [ ] PDF/HTML报告导出
- [ ] Web UI界面
- [ ] 邮件通知
- [ ] API认证系统

---

**更新时间**: 2026-02-09  
**版本**: v2.0.0  
**状态**: ✅ 已发布
