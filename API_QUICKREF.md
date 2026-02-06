# 🚀 API密钥快速获取卡

## ⚡ 最快配置（15分钟）

### 步骤1: Serper API（5分钟）⭐

```
1. 访问: https://serper.dev/
2. 点击 "Sign Up" 注册
3. 登录后进入 Dashboard
4. 复制 API Key
5. 免费额度: 2500次搜索
```

### 步骤2: OpenAI API（10分钟）⭐⭐⭐

```
1. 访问: https://platform.openai.com/
2. 注册账号（需要验证手机号）
3. 充值 $10（最少建议）
4. 进入 "API keys" 创建密钥
5. 复制密钥（只显示一次！）
```

### 步骤3: 配置到项目（1分钟）

```bash
# 1. 复制配置文件
copy .env.example .env

# 2. 编辑配置文件
notepad .env

# 3. 填入密钥
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=你的OpenAI密钥

# 4. 保存并关闭
```

---

## 📋 API网址速查

| API | 注册网址 | 用途 |
|-----|----------|------|
| **Serper** | https://serper.dev/ | 搜索竞品 |
| **OpenAI** | https://platform.openai.com/ | AI分析 |
| **Firecrawl** | https://firecrawl.dev/ | 爬取网页 |
| **Google** | https://console.cloud.google.com/ | 搜索（备选） |
| **Bing** | https://portal.azure.com/ | 搜索（备选） |

---

## 💰 费用速查

| API | 免费额度 | 付费价格 | 推荐预算 |
|-----|----------|----------|----------|
| Serper | 2500次 | $5/1000次 | $10/月 |
| OpenAI | 无 | $0.03-0.06/1K tokens | $20/月 |
| Firecrawl | 500页/月 | $0.02/页 | $20/月 |
| Google | 100次/天 | $5/1000次 | 免费够用 |
| Bing | 1000次/月 | $3/1000次 | 免费够用 |

---

## 🎯 推荐方案

### 方案A: 最小配置（新手）
```env
SERPER_API_KEY=xxx    # 免费2500次
OPENAI_API_KEY=xxx    # 充值$10
```
**成本**: $10  
**够用**: 15-20次分析

### 方案B: 标准配置（推荐）
```env
SERPER_API_KEY=xxx
OPENAI_API_KEY=xxx
FIRECRAWL_API_KEY=xxx  # 免费500页/月
```
**成本**: $10-30/月  
**够用**: 50-100次分析

### 方案C: 完整配置（专业）
```env
SERPER_API_KEY=xxx
GOOGLE_API_KEY=xxx
BING_API_KEY=xxx
OPENAI_API_KEY=xxx
FIRECRAWL_API_KEY=xxx
```
**成本**: $50-100/月  
**够用**: 无限制

---

## ⚠️ 重要提示

### Serper
- ✅ 注册即送2500次免费搜索
- ✅ 无需信用卡
- ✅ 5分钟即可完成

### OpenAI
- ⚠️ 需要验证手机号（可能需要国际号码）
- ⚠️ 必须充值才能使用（最少$5）
- ⚠️ API密钥只显示一次，务必保存
- 💡 建议首次充值$10

### Firecrawl
- ✅ 免费500页/月够日常使用
- ✅ 不配置会自动使用Jina（免费）
- 💡 新手可以先不配置

---

## 🔒 安全提示

### 密钥保护
1. ❌ 不要提交到Git
2. ❌ 不要分享给他人
3. ❌ 不要截图发到网上
4. ✅ 使用密码管理器保存
5. ✅ 定期更换密钥

### 设置限额
- OpenAI: 设置月度预算上限
- 启用使用量提醒
- 定期查看账单

---

## 🐛 常见问题

### Q1: 没有国际手机号怎么办？
**A**: 
- 使用接码平台（不推荐）
- 使用DeepSeek等国产大模型（推荐）
- 找朋友帮忙验证

### Q2: OpenAI太贵了？
**A**: 
- 使用 gpt-3.5-turbo（便宜10倍）
- 使用 DeepSeek（便宜30倍）
- 只分析重要竞品

### Q3: Serper用完了怎么办？
**A**:
- 切换到Google Search（每天100次免费）
- 切换到Bing Search（每月1000次免费）
- 充值Serper（$5/1000次）

### Q4: 密钥填错了？
**A**:
```bash
# 重新编辑 .env 文件
notepad .env

# 检查:
# 1. 密钥完整无误
# 2. 没有多余空格
# 3. 没有引号
# 4. 格式正确: KEY=value
```

---

## 🧪 测试密钥

### 启动项目后测试

```bash
# 1. 启动服务
go run main.go

# 2. 测试健康检查
curl http://localhost:8080/health

# 3. 测试搜索功能
curl -X POST http://localhost:8080/api/discover/search \
  -H "Content-Type: application/json" \
  -d "{\"topic\":\"test\",\"competitor_count\":1}"

# 如果返回task_id，说明配置成功！
```

---

## 📞 获取帮助

### API问题
- 查看详细指南: [API_GUIDE.md](API_GUIDE.md)
- Serper: support@serper.dev
- OpenAI: https://help.openai.com/

### 项目问题
- 查看运行指南: [HOW_TO_RUN.md](HOW_TO_RUN.md)
- 查看快速上手: [QUICKSTART.md](QUICKSTART.md)

---

## ✅ 配置检查清单

完成配置前，请确认：

- [ ] 已注册Serper账号
- [ ] 已复制Serper API密钥
- [ ] 已注册OpenAI账号
- [ ] 已充值OpenAI账户（至少$10）
- [ ] 已创建OpenAI API密钥
- [ ] 已复制 .env.example 为 .env
- [ ] 已填写API密钥到 .env
- [ ] 已保存 .env 文件
- [ ] 准备启动项目！

---

**⏱️ 预计耗时**: 15-20分钟  
**💰 最低成本**: $10（OpenAI充值）  
**🎯 完成后**: 立即可以使用！

---

打印此文档，边看边操作，15分钟搞定！🚀
