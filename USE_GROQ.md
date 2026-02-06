# 🎯 使用Groq（推荐免费方案）

## 为什么选择Groq？

✅ **完全免费** - 无需信用卡  
✅ **速度超快** - 比GPT-4快7倍  
✅ **质量优秀** - Llama 3.3性能强大  
✅ **5分钟配置** - 超级简单  

---

## 🚀 三种配置方式

### 方式1：一键配置脚本（最简单）⭐

**Windows**:
```bash
.\setup-groq.bat
```

**Linux/Mac**:
```bash
chmod +x setup-groq.sh
./setup-groq.sh
```

脚本会引导您：
1. 获取Groq API密钥
2. 获取Serper API密钥（可选）
3. 自动生成 `.env` 配置文件

---

### 方式2：手动配置

**1. 获取Groq API密钥**
```
访问: https://console.groq.com/
注册 → API Keys → Create API Key → 复制密钥
```

**2. 创建配置文件**
```bash
copy .env.groq .env  # Windows
# 或
cp .env.groq .env    # Linux/Mac
```

**3. 填入密钥**
```bash
notepad .env  # Windows
# 或
vi .env       # Linux/Mac
```

修改这一行：
```env
OPENAI_API_KEY=gsk_你的Groq密钥
```

---

### 方式3：完全手动

创建 `.env` 文件，填入以下内容：

```env
# 必需：搜索API
SERPER_API_KEY=你的Serper密钥

# 必需：Groq LLM（完全免费）
OPENAI_API_KEY=gsk_你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile

# 其他配置
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
SERVER_PORT=8080
```

---

## ✅ 验证配置

启动服务：
```bash
.\start.bat  # Windows
./start.sh   # Linux/Mac
```

测试API：
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

---

## 🎓 Groq模型选择

| 模型 | 速度 | 质量 | 推荐用途 |
|------|------|------|----------|
| **llama-3.3-70b-versatile** ⭐ | 超快 | 优秀 | 默认推荐 |
| llama-3.1-8b-instant | 极快 | 良好 | 简单任务 |
| mixtral-8x7b-32768 | 快 | 优秀 | 长文本 |

修改模型：
```env
# 在 .env 中修改
LLM_MODEL=llama-3.3-70b-versatile
```

---

## 💰 成本对比

### 分析3个竞品

| 方案 | 费用 | 说明 |
|------|------|------|
| **Groq** | **$0** | 完全免费 |
| DeepSeek | ¥0.03 | 超便宜 |
| OpenAI GPT-4 | $0.50 | 需付费 |

**Groq = 完全免费！**

---

## ⚠️ 使用限制

- ✅ 完全免费
- ⚠️ 每分钟30次请求
- ⚠️ 每天14,400次请求

对本项目影响：
- ✅ 完全够用
- 单次分析约5-10次请求
- 每天可完成1000+次分析

---

## 📚 详细文档

- [GROQ_SETUP.md](GROQ_SETUP.md) - 完整配置指南
- [FREE_LLM_ALTERNATIVES.md](FREE_LLM_ALTERNATIVES.md) - 所有免费方案
- [HOW_TO_RUN.md](HOW_TO_RUN.md) - 项目运行指南

---

## 🎊 开始使用

```bash
# 1. 运行配置脚本
.\setup-groq.bat

# 2. 按提示获取并输入密钥

# 3. 启动服务
.\start.bat

# 完成！开始使用
```

**总耗时**: 5分钟  
**总成本**: $0  

🎉 完全免费，立即可用！
