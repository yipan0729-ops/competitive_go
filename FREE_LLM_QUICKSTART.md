# 🚀 快速配置免费LLM方案

本项目已支持**多种免费或低成本**的LLM方案！无需OpenAI也能运行。

---

## ⚡ 3分钟配置（推荐方案）

### 方案A: DeepSeek（超便宜）🔥

**成本**: ¥10充值，够用2500次分析  
**质量**: 接近GPT-4

```bash
# 1. 注册并获取密钥
https://platform.deepseek.com/

# 2. 充值¥10

# 3. 编辑 .env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=sk-你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat

# 4. 启动
.\start.bat
```

### 方案B: Ollama（完全免费）🆓

**成本**: $0，完全免费  
**质量**: 良好

```bash
# 1. 下载并安装Ollama
https://ollama.com/download

# 2. 下载模型
ollama pull qwen2.5:7b

# 3. 编辑 .env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b

# 4. 启动
.\start.bat
```

---

## 📋 所有支持的方案

| 方案 | 费用 | 配置难度 | 质量 | 推荐度 |
|------|------|----------|------|--------|
| **DeepSeek** | ¥10/2500次 | ⭐ 简单 | ⭐⭐⭐⭐ | 🔥 强烈推荐 |
| **Ollama** | 免费 | ⭐⭐ 中等 | ⭐⭐⭐ | 🔥 推荐 |
| **智谱AI** | ¥15/100次 | ⭐ 简单 | ⭐⭐⭐⭐ | ✅ 可选 |
| **通义千问** | ¥6/100次 | ⭐ 简单 | ⭐⭐⭐⭐ | ✅ 可选 |
| **Groq** | 免费 | ⭐ 简单 | ⭐⭐⭐⭐ | ✅ 可选 |
| OpenAI | $50/100次 | ⭐ 简单 | ⭐⭐⭐⭐⭐ | 💰 贵 |

---

## 🔧 已完成的代码修改

✅ 项目已支持自定义LLM API地址  
✅ 兼容所有OpenAI格式的API  
✅ 配置文件已更新

**无需额外修改代码，直接配置即可使用！**

---

## 📝 详细配置

### DeepSeek配置

```env
# .env 文件
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=sk-xxxxxxxxxxxx
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
```

### Ollama配置

```env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b
```

### 智谱AI配置

```env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=你的智谱密钥
OPENAI_BASE_URL=https://open.bigmodel.cn/api/paas
LLM_MODEL=glm-4
```

### 通义千问配置

```env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=sk-你的通义密钥
OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode
LLM_MODEL=qwen-plus
```

### Groq配置

```env
SERPER_API_KEY=你的Serper密钥
OPENAI_API_KEY=gsk-你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.1-70b-versatile
```

---

## ❓ 常见问题

### Q: 完全不想花钱，怎么配置？

**A**: 使用Ollama本地方案，完全免费：

1. 安装Ollama
2. 下载模型: `ollama pull qwen2.5:7b`
3. 配置使用Ollama
4. 只需要Serper免费额度（2500次）

**总成本: $0**

### Q: 想要最好的质量，怎么配置？

**A**: 使用DeepSeek，性价比最高：

1. 充值¥10
2. 质量接近GPT-4
3. 可以完成2500次分析

**总成本: ¥10 ≈ $1.4**

### Q: Cursor的API能用吗？

**A**: ❌ **不能**

- Cursor没有公开API
- 只是一个IDE，不提供API服务
- 建议用DeepSeek或Ollama

### Q: 我的电脑配置不够跑Ollama？

**A**: 使用DeepSeek或其他云端API：

- DeepSeek（超便宜）
- 智谱AI（国产）
- 通义千问（阿里）
- Groq（免费但有限制）

---

## 🎯 推荐配置组合

### 组合A: 免费体验
```env
SERPER_API_KEY=xxx（免费2500次）
+ Ollama本地（完全免费）
= 总成本: $0
```

### 组合B: 最佳性价比（推荐）
```env
SERPER_API_KEY=xxx（免费2500次）
+ DeepSeek（¥10充值）
= 总成本: ¥10，够用很久
```

### 组合C: 付费优质
```env
SERPER_API_KEY=xxx（$5/1000次）
+ OpenAI GPT-4（$50/100次）
= 总成本: $55/月
```

---

## 📚 相关文档

- [FREE_LLM_ALTERNATIVES.md](FREE_LLM_ALTERNATIVES.md) - 详细的免费方案指南
- [API_GUIDE.md](API_GUIDE.md) - API获取指南
- [HOW_TO_RUN.md](HOW_TO_RUN.md) - 项目运行指南

---

**更新时间**: 2026-02-06  
**状态**: ✅ 已完成，可直接使用
