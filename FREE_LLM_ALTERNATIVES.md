# 🆓 OpenAI API 免费替代方案指南

OpenAI API 需要付费，这里提供多个**免费或超低成本**的替代方案！

---

## 🎯 推荐方案概览

| 方案 | 费用 | 质量 | 难度 | 推荐度 |
|------|------|------|------|--------|
| **DeepSeek** | ⭐⭐⭐⭐⭐ 超便宜 | ⭐⭐⭐⭐ 很好 | ⭐ 简单 | 🔥 强烈推荐 |
| **Ollama本地** | ⭐⭐⭐⭐⭐ 完全免费 | ⭐⭐⭐ 良好 | ⭐⭐ 中等 | 🔥 推荐 |
| **智谱AI** | ⭐⭐⭐⭐ 便宜 | ⭐⭐⭐⭐ 很好 | ⭐ 简单 | ✅ 推荐 |
| **通义千问** | ⭐⭐⭐⭐ 便宜 | ⭐⭐⭐⭐ 很好 | ⭐ 简单 | ✅ 推荐 |
| **Groq** | ⭐⭐⭐⭐⭐ 免费 | ⭐⭐⭐⭐ 很好 | ⭐ 简单 | ✅ 推荐 |
| **Cursor API** | ❌ 不可用 | - | - | ❌ 不支持 |

---

## 方案1: DeepSeek（强烈推荐）🔥

### 💰 费用对比
- **OpenAI GPT-4**: $0.03/1K tokens（输入）
- **DeepSeek**: ¥0.001/1K tokens ≈ **便宜200倍！**

### 🎯 特点
- ✅ 质量接近GPT-4
- ✅ 价格超低（1元=约150万tokens）
- ✅ 兼容OpenAI API格式
- ✅ 无需科学上网
- ✅ 5分钟即可接入

### 📝 获取步骤

**1. 注册账号**
```
官网: https://platform.deepseek.com/
```

**2. 获取API密钥**
- 登录后点击"API Keys"
- 创建新密钥
- 复制密钥（格式：sk-xxx）

**3. 充值**
- 点击"账户余额"
- 充值¥10（够用很久）
- 支持支付宝/微信

**4. 配置到项目**

编辑 `.env` 文件：
```env
# 使用DeepSeek替代OpenAI
OPENAI_API_KEY=sk-你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
```

### 🔧 修改代码

编辑 `ai/llm.go`：

```go
// NewLLMClient 创建LLM客户端
func NewLLMClient(apiKey, model string, temperature float64, maxTokens int) *LLMClient {
	return &LLMClient{
		APIKey:      apiKey,
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		BaseURL:     os.Getenv("OPENAI_BASE_URL"), // 新增
	}
}

// Chat 发送聊天请求
func (c *LLMClient) Chat(messages []ChatMessage) (string, error) {
	if c.APIKey == "" {
		return "", errors.New("API Key未配置")
	}

	// 使用自定义BaseURL或默认OpenAI URL
	apiURL := c.BaseURL
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/chat/completions"
	} else {
		apiURL = apiURL + "/v1/chat/completions"
	}

	// ... 其余代码保持不变
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	// ...
}
```

### 💡 成本对比

**场景：分析3个竞品**
- OpenAI GPT-4: $0.50
- DeepSeek: ¥0.03 ≈ **$0.004**

**充值¥10可以完成**：
- OpenAI: 约20次分析
- DeepSeek: 约**2500次分析**

---

## 方案2: Ollama本地部署（完全免费）🆓

### 💰 费用
- **完全免费！** 无任何API费用
- 只需要电脑有足够的内存（8GB+）

### 🎯 特点
- ✅ 完全免费，无限使用
- ✅ 数据隐私，不发送到外网
- ✅ 离线可用
- ⚠️ 需要一定硬件配置
- ⚠️ 质量略低于GPT-4

### 📝 安装步骤

**1. 安装Ollama**

**Windows**:
```bash
# 下载安装包
https://ollama.com/download/windows
# 双击安装
```

**Linux/Mac**:
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**2. 下载模型**
```bash
# 推荐：Qwen2.5（中文好）
ollama pull qwen2.5:7b

# 或者：Llama3.1（英文好）
ollama pull llama3.1:8b

# 或者：DeepSeek-R1（推理能力强）
ollama pull deepseek-r1:8b
```

**3. 启动服务**
```bash
# Ollama会自动在后台运行
# API地址：http://localhost:11434
```

**4. 配置到项目**

编辑 `.env`：
```env
# 使用Ollama本地模型
OPENAI_API_KEY=ollama  # 随便填
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b
```

### 🔧 修改代码

编辑 `ai/llm.go`，添加Ollama支持：

```go
// Chat 发送聊天请求
func (c *LLMClient) Chat(messages []ChatMessage) (string, error) {
	// 检测是否使用Ollama
	isOllama := strings.Contains(c.BaseURL, "localhost:11434")
	
	if isOllama {
		return c.chatWithOllama(messages)
	}
	
	// 原有的OpenAI逻辑
	// ...
}

// chatWithOllama Ollama专用请求
func (c *LLMClient) chatWithOllama(messages []ChatMessage) (string, error) {
	type OllamaRequest struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
		Stream   bool          `json:"stream"`
	}

	request := OllamaRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, _ := json.Marshal(request)
	apiURL := c.BaseURL + "/api/chat"
	
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	
	if message, ok := result["message"].(map[string]interface{}); ok {
		if content, ok := message["content"].(string); ok {
			return content, nil
		}
	}
	
	return "", errors.New("Ollama响应格式错误")
}
```

### 💻 硬件要求

| 模型 | 内存需求 | 推荐配置 | 速度 |
|------|----------|----------|------|
| qwen2.5:7b | 8GB | 16GB | 中等 |
| llama3.1:8b | 8GB | 16GB | 中等 |
| deepseek-r1:8b | 8GB | 16GB | 慢 |
| qwen2.5:14b | 16GB | 32GB | 慢 |

---

## 方案3: 智谱AI（国产推荐）

### 💰 费用
- GLM-4: ¥0.05/1K tokens
- 比OpenAI便宜**20倍**

### 📝 获取步骤

**1. 注册**
```
官网: https://open.bigmodel.cn/
```

**2. 获取API Key**
- 进入"API Keys"
- 创建新密钥

**3. 配置**
```env
OPENAI_API_KEY=你的智谱密钥
OPENAI_BASE_URL=https://open.bigmodel.cn/api/paas
LLM_MODEL=glm-4
```

---

## 方案4: 通义千问（阿里）

### 💰 费用
- qwen-plus: ¥0.002/1K tokens
- 比OpenAI便宜**100倍**

### 📝 获取步骤

**1. 注册**
```
官网: https://dashscope.aliyun.com/
```

**2. 获取API Key**
- 进入"API-KEY管理"
- 创建新密钥

**3. 配置**
```env
OPENAI_API_KEY=sk-你的通义密钥
OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode
LLM_MODEL=qwen-plus
```

---

## 方案5: Groq（免费且快速）🚀

### 💰 费用
- **完全免费！**
- 每分钟限制30次请求

### 🎯 特点
- ✅ 完全免费
- ✅ 速度超快（专用硬件）
- ✅ 支持Llama、Mixtral等模型
- ⚠️ 有速率限制

### 📝 获取步骤

**1. 注册**
```
官网: https://console.groq.com/
```

**2. 获取API Key**
- 点击"API Keys"
- 创建新密钥

**3. 配置**
```env
OPENAI_API_KEY=gsk_你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.1-70b-versatile
```

---

## ❌ Cursor API 不可用

### 为什么不能用Cursor的API？

1. **Cursor没有公开API**
   - Cursor是一个IDE，不是API服务
   - 它内部调用Claude/GPT，但不对外提供

2. **技术限制**
   - Cursor的AI功能只能在IDE内使用
   - 无法通过HTTP API访问

3. **替代方案**
   - 使用Claude API（Cursor背后的模型）
   - 或使用上述免费方案

### 如果想用Claude

**直接使用Claude API**:
```
官网: https://console.anthropic.com/
费用: 类似OpenAI，需要付费
```

配置：
```env
# 注意：需要修改代码适配Claude API格式
CLAUDE_API_KEY=sk-ant-你的密钥
```

---

## 🔧 统一配置方案

为了支持多种LLM，我为您创建了一个**统一配置方案**：

### 1. 更新 `.env.example`

```env
# ========== LLM配置 ==========
# 选择一个LLM提供商，填写对应的配置

# 方案1: DeepSeek（推荐，超便宜）
# OPENAI_API_KEY=sk-你的DeepSeek密钥
# OPENAI_BASE_URL=https://api.deepseek.com
# LLM_MODEL=deepseek-chat

# 方案2: Ollama本地（免费）
# OPENAI_API_KEY=ollama
# OPENAI_BASE_URL=http://localhost:11434
# LLM_MODEL=qwen2.5:7b

# 方案3: 智谱AI
# OPENAI_API_KEY=你的智谱密钥
# OPENAI_BASE_URL=https://open.bigmodel.cn/api/paas
# LLM_MODEL=glm-4

# 方案4: 通义千问
# OPENAI_API_KEY=sk-你的通义密钥
# OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode
# LLM_MODEL=qwen-plus

# 方案5: Groq（免费）
# OPENAI_API_KEY=gsk_你的Groq密钥
# OPENAI_BASE_URL=https://api.groq.com/openai
# LLM_MODEL=llama-3.1-70b-versatile

# 方案6: OpenAI（原版，需付费）
# OPENAI_API_KEY=sk-proj-你的OpenAI密钥
# OPENAI_BASE_URL=
# LLM_MODEL=gpt-4

# LLM参数
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
```

### 2. 更新配置模块

我会为您创建一个支持多种LLM的版本。

---

## 📊 方案对比总结

### 成本对比（分析3个竞品）

| 方案 | 单次成本 | 月费（100次） |
|------|----------|---------------|
| **OpenAI GPT-4** | $0.50 | $50 |
| **DeepSeek** | ¥0.03 ($0.004) | ¥3 ($0.40) |
| **智谱AI** | ¥0.15 ($0.02) | ¥15 ($2) |
| **通义千问** | ¥0.06 ($0.008) | ¥6 ($0.80) |
| **Ollama** | $0 | $0 |
| **Groq** | $0 | $0 |

### 质量对比

| 方案 | 中文能力 | 英文能力 | 推理能力 | 速度 |
|------|----------|----------|----------|------|
| OpenAI GPT-4 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 中等 |
| DeepSeek | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 快 |
| 智谱AI | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | 中等 |
| 通义千问 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | 快 |
| Ollama | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | 慢 |
| Groq | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 极快 |

---

## 🎯 推荐选择

### 🏆 最佳性价比: DeepSeek
- 质量接近GPT-4
- 价格是OpenAI的1/200
- 充值¥10够用很久

### 🆓 完全免费: Ollama
- 一次性投入（下载模型）
- 无限使用
- 适合隐私敏感场景

### ⚡ 快速体验: Groq
- 免费且快速
- 适合轻度使用
- 有请求限制

---

## 🚀 快速开始

### 使用DeepSeek（推荐）

```bash
# 1. 注册获取密钥
https://platform.deepseek.com/

# 2. 充值¥10

# 3. 配置
copy .env.example .env
notepad .env

# 4. 填入配置
OPENAI_API_KEY=sk-你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat

# 5. 启动
.\start.bat
```

### 使用Ollama（免费）

```bash
# 1. 安装Ollama
https://ollama.com/download

# 2. 下载模型
ollama pull qwen2.5:7b

# 3. 配置
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b

# 4. 启动
.\start.bat
```

---

**推荐**: 先用DeepSeek试试，¥10充值，够用很久！🚀

**免费**: 用Ollama，完全免费，但需要好点的电脑配置。

**不推荐**: Cursor API不可用，无法使用。
