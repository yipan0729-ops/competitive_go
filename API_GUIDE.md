# 🔑 API密钥获取指南

本项目需要以下API密钥。您可以根据需求选择配置，最少需要2个API。

---

## 📋 快速概览

| API | 用途 | 必需性 | 免费额度 | 费用 |
|-----|------|--------|----------|------|
| **Serper** | 搜索竞品 | ⭐ 推荐 | 2500次 | $5/1000次 |
| **OpenAI** | AI分析 | ⭐ 必需 | 无 | 按token计费 |
| **Firecrawl** | 网页爬取 | 推荐 | 500页/月 | $0.02/页 |
| **Google** | 搜索（备选） | 可选 | 100次/天 | $5/1000次 |
| **Bing** | 搜索（备选） | 可选 | 1000次/月 | $3/1000次 |

---

## 1. Serper API（推荐）⭐

### 🎯 用途
- 智能发现竞品
- 搜索数据源
- 构建竞品列表

### 📝 获取步骤

**1. 访问官网**
```
https://serper.dev/
```

**2. 注册账号**
- 点击右上角 "Sign Up"
- 使用Google账号或邮箱注册

**3. 获取API密钥**
- 登录后，点击 "Dashboard"
- 在页面中找到 "API Key"
- 点击 "Copy" 复制密钥

**4. 配置到项目**
```env
SERPER_API_KEY=你复制的密钥
```

### 💰 费用说明
- ✅ **免费额度**: 2500次搜索（注册即送）
- 💳 **付费价格**: $5/1000次搜索
- 📊 **项目用量**: 单次发现约15次搜索

**示例计算**:
- 免费额度可以完成: 166次竞品发现任务
- 月费约$10可以: 2000次搜索 = 133次发现任务

### 🔗 相关链接
- 官网: https://serper.dev/
- 文档: https://serper.dev/docs
- 定价: https://serper.dev/pricing

---

## 2. OpenAI API（必需）⭐⭐⭐

### 🎯 用途
- 提取竞品信息
- SWOT分析
- 生成分析报告

### 📝 获取步骤

**1. 访问官网**
```
https://platform.openai.com/
```

**2. 注册账号**
- 点击 "Sign up"
- 使用邮箱或Google账号注册
- 验证手机号（可能需要国际号码）

**3. 充值账户**
- 登录后，点击右上角头像
- 选择 "Billing"
- 点击 "Add payment method"
- 添加信用卡并充值（建议$10起）

**4. 创建API密钥**
- 点击左侧菜单 "API keys"
- 点击 "Create new secret key"
- 输入名称（如：competitive-analyzer）
- 复制密钥（⚠️ 只显示一次，请妥善保存）

**5. 配置到项目**
```env
OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx
```

### 💰 费用说明
- ❌ **无免费额度**（需要充值）
- 💳 **GPT-4 价格**: 
  - 输入: $0.03/1K tokens
  - 输出: $0.06/1K tokens
- 📊 **项目用量**: 单次分析约$0.50

**示例计算**:
- 充值$10可以: 约20次完整分析
- 充值$50可以: 约100次分析
- 月费约$20: 可完成40次分析

### 💡 省钱技巧
1. **使用GPT-3.5-Turbo**（便宜10倍）
   ```env
   LLM_MODEL=gpt-3.5-turbo
   ```
2. **减少token使用**
   ```env
   LLM_MAX_TOKENS=2000
   ```
3. **使用国产大模型**（更便宜）
   - DeepSeek: https://platform.deepseek.com/
   - 智谱AI: https://open.bigmodel.cn/

### 🔗 相关链接
- 官网: https://platform.openai.com/
- 文档: https://platform.openai.com/docs
- 定价: https://openai.com/api/pricing/
- 使用量查看: https://platform.openai.com/usage

---

## 3. Firecrawl API（推荐）

### 🎯 用途
- 高质量网页抓取
- 处理JavaScript渲染
- 绕过常见反爬

### 📝 获取步骤

**1. 访问官网**
```
https://firecrawl.dev/
```

**2. 注册账号**
- 点击 "Get Started"
- 使用GitHub或邮箱注册

**3. 获取API密钥**
- 登录后进入Dashboard
- 在 "API Keys" 部分
- 点击 "Create API Key"
- 复制密钥

**4. 配置到项目**
```env
FIRECRAWL_API_KEY=fc-xxxxxxxxxxxxxxxx
```

### 💰 费用说明
- ✅ **免费额度**: 500页/月
- 💳 **付费价格**: 
  - Starter: $20/月（5000页）
  - Pro: $100/月（50000页）
- 📊 **项目用量**: 单个竞品约3-5页

**示例计算**:
- 免费额度可以: 100-160个竞品
- Starter计划: 1000-1600个竞品/月

### 🆓 免费替代方案
如果不想付费，可以只使用 **Jina Reader**（完全免费），项目会自动降级使用。

### 🔗 相关链接
- 官网: https://firecrawl.dev/
- 文档: https://docs.firecrawl.dev/
- 定价: https://firecrawl.dev/pricing

---

## 4. Google Search API（可选）

### 🎯 用途
- 替代Serper进行搜索
- 高质量搜索结果

### 📝 获取步骤

**1. 创建Google Cloud项目**
```
https://console.cloud.google.com/
```

**2. 启用Custom Search API**
- 进入 "APIs & Services"
- 点击 "Enable APIs and Services"
- 搜索 "Custom Search API"
- 点击 "Enable"

**3. 创建API密钥**
- 进入 "Credentials"
- 点击 "Create Credentials" > "API key"
- 复制API密钥

**4. 创建搜索引擎ID**
```
https://programmablesearchengine.google.com/
```
- 点击 "Add"
- 输入搜索范围（可以选择"整个网络"）
- 创建后复制 "Search engine ID"

**5. 配置到项目**
```env
GOOGLE_API_KEY=AIzaSyxxxxxxxxxxxxxxxxxx
GOOGLE_SEARCH_ENGINE_ID=xxxxxxxxxxxxxxxxx
```

### 💰 费用说明
- ✅ **免费额度**: 100次/天
- 💳 **付费价格**: $5/1000次
- 📊 **项目用量**: 单次发现约10-15次

### 🔗 相关链接
- Google Cloud: https://console.cloud.google.com/
- Custom Search: https://developers.google.com/custom-search
- 搜索引擎配置: https://programmablesearchengine.google.com/

---

## 5. Bing Search API（可选）

### 🎯 用途
- 替代Serper/Google搜索
- 中文内容支持好

### 📝 获取步骤

**1. 访问Azure Portal**
```
https://portal.azure.com/
```

**2. 创建Bing Search资源**
- 搜索 "Bing Search"
- 点击 "Create"
- 选择订阅和资源组
- 选择定价层（F1免费层）

**3. 获取API密钥**
- 资源创建后，进入 "Keys and Endpoint"
- 复制 "KEY 1"

**4. 配置到项目**
```env
BING_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxx
```

### 💰 费用说明
- ✅ **免费额度**: 1000次/月
- 💳 **付费价格**: $3/1000次
- 📊 **项目用量**: 单次发现约15次

### 🔗 相关链接
- Azure Portal: https://portal.azure.com/
- 文档: https://docs.microsoft.com/en-us/bing/search-apis/
- 定价: https://azure.microsoft.com/pricing/details/cognitive-services/

---

## 📋 推荐配置方案

### 方案A：最小配置（免费体验）

```env
# 只需这两个
SERPER_API_KEY=你的密钥        # 免费2500次
OPENAI_API_KEY=你的密钥        # 需充值$10
```

**成本**: $10（OpenAI充值）
**可用量**: 约15-20次完整分析

---

### 方案B：标准配置（推荐）

```env
# 搜索
SERPER_API_KEY=你的密钥

# AI分析
OPENAI_API_KEY=你的密钥

# 爬取（可选）
FIRECRAWL_API_KEY=你的密钥    # 免费500页/月
```

**成本**: $10-30/月
**可用量**: 约50-100次分析

---

### 方案C：完整配置（专业版）

```env
# 主搜索引擎
SERPER_API_KEY=你的密钥

# 备用搜索（多个备选）
GOOGLE_API_KEY=你的密钥
BING_API_KEY=你的密钥

# AI分析
OPENAI_API_KEY=你的密钥

# 高质量爬取
FIRECRAWL_API_KEY=你的密钥
```

**成本**: $50-100/月
**可用量**: 无限制（取决于预算）

---

## 🎯 配置步骤

### 1. 创建配置文件

```bash
# Windows
copy .env.example .env

# Linux/Mac
cp .env.example .env
```

### 2. 编辑配置文件

```bash
# Windows
notepad .env

# Linux/Mac
vi .env
# 或
nano .env
```

### 3. 填入API密钥

```env
# 必需配置
SERPER_API_KEY=ser_xxxxxxxxxxxxxxxxxxxxxxxxxx
OPENAI_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxx

# 推荐配置
FIRECRAWL_API_KEY=fc-xxxxxxxxxxxxxxxxxxxxxxxx

# 可选配置（如果有）
GOOGLE_API_KEY=AIzaSyxxxxxxxxxxxxxxxxxx
GOOGLE_SEARCH_ENGINE_ID=xxxxxxxxxxxxxxxxx
BING_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxx
```

### 4. 验证配置

```bash
# 运行验证脚本
.\verify.bat  # Windows
./verify.sh   # Linux/Mac

# 或者直接启动
go run main.go
```

---

## ⚠️ 安全提示

### 保护您的API密钥

1. **不要提交到Git**
   - `.env` 已在 `.gitignore` 中
   - 永远不要把密钥提交到公开仓库

2. **定期轮换密钥**
   - 每3-6个月更换一次
   - 怀疑泄露时立即更换

3. **设置使用限制**
   - OpenAI: 设置月度预算上限
   - Serper: 监控使用量
   - 启用使用提醒

4. **妥善保管**
   - 使用密码管理器存储
   - 不要分享给他人
   - 备份到安全位置

---

## 💡 省钱技巧

### 1. 优先使用免费服务

```env
# Jina Reader 完全免费（会自动使用）
# 不配置 FIRECRAWL_API_KEY 即可
```

### 2. 使用国产大模型

**DeepSeek**（便宜10倍）:
```
官网: https://platform.deepseek.com/
价格: ¥0.001/1K tokens（约便宜30倍）
```

修改代码使用DeepSeek:
```go
// ai/llm.go 中修改API地址
req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", ...)
```

### 3. 启用搜索缓存

```env
# 搜索结果缓存7天（默认已开启）
SEARCH_CACHE_DAYS=7
```

### 4. 减少LLM token消耗

```env
LLM_MODEL=gpt-3.5-turbo  # 而不是gpt-4
LLM_MAX_TOKENS=2000      # 限制输出长度
```

---

## 📊 成本估算工具

### 在线计算器

使用我提供的成本计算器估算费用：

```
单次完整分析成本 = 
  搜索成本（约$0.08）+ 
  AI分析成本（约$0.50）+ 
  爬取成本（约$0.06）
= 约 $0.64/次
```

### 月度预算建议

| 使用频率 | 推荐预算 | 可完成分析 |
|----------|----------|------------|
| 轻度使用 | $10/月 | 15次 |
| 中度使用 | $30/月 | 45次 |
| 重度使用 | $100/月 | 150次 |

---

## 🔧 故障排查

### API密钥不工作？

**检查清单**:
1. ✅ 密钥是否完整（没有空格）
2. ✅ 是否添加了引号（不需要）
3. ✅ 账户是否有余额
4. ✅ API是否已启用
5. ✅ 是否在正确的文件（`.env`）

### 测试API密钥

**Serper**:
```bash
curl -X POST https://google.serper.dev/search \
  -H "X-API-KEY: 你的密钥" \
  -H "Content-Type: application/json" \
  -d '{"q":"test"}'
```

**OpenAI**:
```bash
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer 你的密钥"
```

---

## 📞 获取帮助

### API问题
- Serper支持: support@serper.dev
- OpenAI支持: https://help.openai.com/
- Firecrawl支持: https://firecrawl.dev/contact

### 项目问题
- 查看 [HOW_TO_RUN.md](HOW_TO_RUN.md)
- 查看 [QUICKSTART.md](QUICKSTART.md)
- 提交Issue

---

## 🎓 总结

### 快速开始（最少配置）

只需要获取这两个API：

1. **Serper** - 5分钟注册，免费2500次
2. **OpenAI** - 10分钟注册，充值$10

总耗时：**15分钟**
总成本：**$10**

然后就可以开始使用了！🎉

---

**更新时间**: 2026-02-06  
**维护者**: 项目团队
