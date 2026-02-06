# 🔧 Groq "Forbidden" 错误排查指南

遇到 `{"error":{"message":"Forbidden"}}` 错误？按以下步骤排查：

---

## 🔍 常见原因及解决方案

### 原因1：API密钥格式错误 ⭐ 最常见

**检查**：
```env
# 错误示例
OPENAI_API_KEY=your_groq_key_here  ❌
OPENAI_API_KEY=gsk_  ❌（不完整）
OPENAI_API_KEY="gsk_xxxx"  ❌（有引号）

# 正确示例
OPENAI_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  ✅
```

**解决**：
1. 检查 `.env` 文件中的密钥
2. 确保密钥完整（约40-50个字符）
3. 确保没有引号、空格等多余字符
4. 确保以 `gsk_` 开头

---

### 原因2：API地址配置错误

**检查**：
```env
# 错误示例
OPENAI_BASE_URL=https://api.groq.com  ❌（缺少/openai）
OPENAI_BASE_URL=https://groq.com/openai  ❌（错误域名）

# 正确示例
OPENAI_BASE_URL=https://api.groq.com/openai  ✅
```

**解决**：
确保完整的API地址：
```env
OPENAI_BASE_URL=https://api.groq.com/openai
```

---

### 原因3：模型名称错误

**检查**：
```env
# 错误示例（旧模型名）
LLM_MODEL=llama-3.1-70b-versatile  ⚠️（可能已废弃）
LLM_MODEL=llama3-70b  ❌（错误格式）

# 正确示例（2026年2月最新）
LLM_MODEL=llama-3.3-70b-versatile  ✅
LLM_MODEL=llama-guard-3-8b  ✅（根据Groq网站）
```

**解决**：
根据 [Groq官网](https://console.groq.com/) 最新模型列表：

```env
# 推荐使用这些模型（2026年2月）
LLM_MODEL=llama-3.3-70b-versatile
# 或
LLM_MODEL=kimi-k2
# 或
LLM_MODEL=gpt-oss-120b
```

---

### 原因4：API密钥被删除或失效

**检查**：
1. 登录 https://console.groq.com/keys
2. 查看密钥状态
3. 确认密钥未被删除

**解决**：
1. 删除旧密钥
2. 创建新密钥
3. 更新 `.env` 文件

---

### 原因5：账户或地区限制

**检查**：
- Groq账户是否正常
- 是否在受限地区

**解决**：
1. 确认账户状态正常
2. 如有地区限制，考虑：
   - 使用VPN
   - 或切换到其他免费方案（DeepSeek、Ollama）

---

## ✅ 正确配置示例

### 完整的 .env 配置

```env
# 搜索API（可选）
SERPER_API_KEY=99f29bd66d1a22564a3cbe2f56f3a0981eba5993

# Groq LLM配置
OPENAI_API_KEY=gsk_你的完整密钥（约40-50个字符）
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile

# 其他配置
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000
SERVER_PORT=8080
GIN_MODE=release
DB_PATH=./data/competitive.db
STORAGE_PATH=./storage
REPORTS_PATH=./reports
SEARCH_CACHE_DAYS=7
MAX_SEARCH_RESULTS=10
```

---

## 🧪 测试步骤

### 步骤1：验证API密钥

使用 curl 直接测试 Groq API：

```bash
curl https://api.groq.com/openai/v1/models \
  -H "Authorization: Bearer gsk_你的密钥"
```

**成功返回**：
```json
{
  "object": "list",
  "data": [
    {
      "id": "llama-3.3-70b-versatile",
      ...
    }
  ]
}
```

**失败返回**：
```json
{
  "error": {
    "message": "Forbidden"
  }
}
```

如果直接测试也失败，说明密钥本身有问题。

---

### 步骤2：测试简单请求

```bash
curl https://api.groq.com/openai/v1/chat/completions \
  -H "Authorization: Bearer gsk_你的密钥" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama-3.3-70b-versatile",
    "messages": [{"role": "user", "content": "Hello"}],
    "max_tokens": 10
  }'
```

---

### 步骤3：验证项目配置

```bash
# 1. 启动服务
go run main.go

# 2. 测试健康检查
curl http://localhost:8080/health

# 3. 查看日志
# 如果有"Forbidden"错误，说明配置有问题
```

---

## 🔧 快速修复

### 方法1：重新获取密钥（推荐）

```bash
# 1. 访问 Groq 控制台
https://console.groq.com/keys

# 2. 删除现有密钥

# 3. 创建新密钥
点击 "Create API Key"
输入名称: competitive-analyzer
点击 "Create"

# 4. 复制完整密钥（gsk_开头）

# 5. 更新 .env
OPENAI_API_KEY=gsk_你刚复制的完整密钥
```

---

### 方法2：检查配置文件

```bash
# Windows
notepad .env

# 检查以下内容：
# 1. OPENAI_API_KEY 是否完整
# 2. OPENAI_BASE_URL 是否正确
# 3. LLM_MODEL 是否是有效模型
# 4. 没有多余的引号或空格
```

---

### 方法3：使用配置脚本重新配置

```bash
# 运行配置脚本
.\setup-groq.bat

# 按提示重新输入密钥
```

---

## 🔄 备选方案

如果Groq问题持续，可以立即切换到其他免费方案：

### 方案A：DeepSeek（超便宜）

```env
OPENAI_API_KEY=sk_你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
```

**获取**：https://platform.deepseek.com/
**费用**：¥10充值，够用很久

---

### 方案B：Ollama（本地免费）

```bash
# 1. 安装 Ollama
https://ollama.com/download

# 2. 下载模型
ollama pull qwen2.5:7b

# 3. 配置
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b
```

---

## 📊 Groq最新模型列表（2026年2月）

根据 [Groq官网](https://console.groq.com/)，以下是可用的模型：

### 文本生成
- `llama-3.3-70b-versatile` ⭐ 推荐
- `gpt-oss-120b` - GPT OSS 120B
- `gpt-oss-20b` - GPT OSS 20B
- `kimi-k2` - Kimi K2
- `qwen-3-32b` - Qwen 3 32B

### 多模态
- `llama-4-scout` - Llama 4 Scout（支持视觉）
- `llama-4-maverick` - Llama 4 Maverick（支持视觉）

### 语音
- `whisper-large-v3` - 语音转文本
- `orpheus-english` - 文本转语音

---

## 💡 调试技巧

### 1. 启用详细日志

在 `.env` 中：
```env
GIN_MODE=debug  # 改为debug模式
```

### 2. 查看完整错误

修改 `ai/llm.go`，添加日志：
```go
if resp.StatusCode != http.StatusOK {
    log.Printf("API Error: Status=%d, Body=%s", resp.StatusCode, string(body))
    return "", fmt.Errorf("API返回错误 %d: %s", resp.StatusCode, string(body))
}
```

### 3. 测试不同模型

```env
# 尝试不同模型
LLM_MODEL=gpt-oss-120b  # 试试这个
LLM_MODEL=kimi-k2       # 或这个
```

---

## 📞 获取帮助

### Groq支持
- 官网：https://groq.com/
- 控制台：https://console.groq.com/
- 文档：https://console.groq.com/docs

### 其他方案
- DeepSeek：https://platform.deepseek.com/
- Ollama：https://ollama.com/

---

## 🎯 检查清单

解决"Forbidden"错误前，请确认：

- [ ] API密钥完整（gsk_开头，约40-50字符）
- [ ] API密钥没有引号或空格
- [ ] API地址正确：`https://api.groq.com/openai`
- [ ] 模型名称正确（检查最新列表）
- [ ] 密钥未过期或被删除
- [ ] .env 文件在项目根目录
- [ ] 重启了服务（go run main.go）

---

**最快解决方案**：
1. 重新创建Groq API密钥
2. 完整复制密钥（不要有空格）
3. 更新 `.env` 文件
4. 重启服务

**如果还是不行**：
切换到DeepSeek或Ollama（5分钟配置）

祝您顺利解决！🚀
