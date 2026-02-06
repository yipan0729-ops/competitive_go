# ✅ Ollama配置完成！

## 🎉 配置状态

您的项目已成功切换到**Ollama本地LLM**：

```env
✅ OPENAI_API_KEY=ollama
✅ OPENAI_BASE_URL=http://localhost:11434
✅ LLM_MODEL=qwen2.5:7b
```

---

## 📋 接下来的步骤

### 步骤1：安装Ollama（5分钟）

**Windows**:
```
1. 访问: https://ollama.com/download/windows
2. 下载 OllamaSetup.exe（约200MB）
3. 双击安装
4. 安装完成后Ollama会自动后台运行
```

**Linux**:
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**Mac**:
```bash
# 下载安装包
https://ollama.com/download/mac

# 或使用Homebrew
brew install ollama
```

---

### 步骤2：下载模型（5-10分钟）

打开命令行，执行：

```bash
# 下载Qwen2.5 7B模型（推荐，中文好）
ollama pull qwen2.5:7b
```

下载进度显示：
```
pulling manifest
pulling 8934d96d3f08... 100% ▕████████████████▏ 4.7 GB
verifying sha256 digest
writing manifest
success
```

**其他推荐模型**：
```bash
# DeepSeek R1（推理强）
ollama pull deepseek-r1:8b

# Llama 3.1（英文好）
ollama pull llama3.1:8b

# Llama 3.2 3B（更快，需4GB内存）
ollama pull llama3.2:3b
```

---

### 步骤3：验证Ollama（1分钟）

```bash
# 1. 检查已安装的模型
ollama list

# 应该看到
NAME           ID            SIZE     MODIFIED
qwen2.5:7b     8934d96d3f08  4.7 GB   2 minutes ago

# 2. 测试模型
ollama run qwen2.5:7b "你好，请介绍一下自己"

# 应该看到中文回复
```

---

### 步骤4：启动项目（1分钟）

```bash
# 启动服务
go run main.go

# 或使用启动脚本
.\start.bat
```

看到以下输出表示成功：
```
数据库初始化成功
服务器启动在端口 8080
```

---

### 步骤5：测试功能（1分钟）

打开新的命令行窗口：

```bash
# 测试健康检查
curl http://localhost:8080/health

# 测试爬取功能
curl -X POST http://localhost:8080/api/crawl/single ^
  -H "Content-Type: application/json" ^
  -d "{\"url\":\"https://ollama.com\",\"competitor\":\"Ollama\"}"
```

---

## ✅ 优势总结

### Ollama vs 云端API

| 对比项 | Ollama | Groq | DeepSeek | OpenAI |
|--------|--------|------|----------|--------|
| **费用** | 免费 | 免费 | ¥10/2500次 | $50/100次 |
| **隐私** | 完全本地 | 云端 | 云端 | 云端 |
| **网络** | 不需要 | 需要 | 需要 | 需要 |
| **限制** | 无 | 30次/分钟 | 无 | 按量付费 |
| **速度** | 中等 | 超快 | 快 | 中等 |

---

## 💡 使用技巧

### 1. 选择合适的模型

**中文任务（推荐）**:
```env
LLM_MODEL=qwen2.5:7b  # 最推荐
```

**英文任务**:
```env
LLM_MODEL=llama3.1:8b
```

**快速响应**:
```env
LLM_MODEL=llama3.2:3b
```

### 2. 优化性能

**如果有GPU**:
```bash
# Ollama会自动检测并使用NVIDIA GPU
# 速度提升5-10倍
```

**如果只有CPU**:
```bash
# 使用较小的模型
ollama pull llama3.2:3b
```

### 3. 管理模型

```bash
# 查看所有模型
ollama list

# 删除不用的模型（释放空间）
ollama rm 模型名

# 更新模型
ollama pull qwen2.5:7b
```

---

## ⚠️ 常见问题

### Q1: ollama命令找不到？

**Windows**:
```bash
# 重启命令行窗口
# 或检查安装路径
C:\Users\你的用户名\AppData\Local\Programs\Ollama
```

**解决**: 重启电脑或手动添加到PATH

### Q2: 下载模型失败？

```bash
# 检查网络
ping ollama.com

# 重试下载
ollama pull qwen2.5:7b

# 或下载更小的模型
ollama pull llama3.2:3b
```

### Q3: 运行时提示"Ollama服务未运行"？

**Windows**:
```bash
# Ollama应该自动后台运行
# 检查任务管理器是否有ollama进程

# 手动启动
ollama serve
```

**Linux/Mac**:
```bash
# 启动Ollama服务
ollama serve &

# 检查运行状态
ps aux | grep ollama
```

### Q4: 内存不够，卡顿？

```bash
# 使用更小的模型
ollama pull llama3.2:3b  # 只需4GB
ollama pull qwen2.5:1.5b # 只需2GB

# 更新.env
LLM_MODEL=llama3.2:3b
```

### Q5: 想用GPU加速？

**确认GPU可用**:
```bash
# NVIDIA GPU自动支持
# AMD/Intel GPU暂不支持
```

---

## 🔄 切换到其他方案

### 切换到Groq（免费云端）
```env
OPENAI_API_KEY=gsk_你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile
```

### 切换到DeepSeek（超便宜）
```env
OPENAI_API_KEY=sk_你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
```

---

## 📚 更多资源

- Ollama官网: https://ollama.com/
- 模型库: https://ollama.com/library
- 文档: https://github.com/ollama/ollama/blob/main/docs/api.md
- GitHub: https://github.com/ollama/ollama

---

## 🎊 完成！

**您现在拥有**：
- ✅ 完全免费的本地LLM
- ✅ 无限制的使用次数
- ✅ 完全的数据隐私
- ✅ 离线也能工作

**总耗时**: 10-15分钟（含下载模型）  
**总成本**: $0  
**内存需求**: 8GB+

开始使用您的免费AI竞品分析工具吧！🚀
