# 🆓 使用Ollama本地LLM配置指南

Ollama是**完全免费**的本地LLM方案，无需API密钥，数据完全本地化！

---

## ✨ Ollama的优势

✅ **完全免费** - 无任何API费用  
✅ **隐私保护** - 数据不发送到外网  
✅ **离线可用** - 无需网络连接  
✅ **无限使用** - 没有请求限制  
⚠️ **需要配置** - 需要8GB+内存  

---

## 🚀 3步配置流程

### 第1步：安装Ollama（2分钟）

**Windows用户**：
```
1. 访问: https://ollama.com/download/windows
2. 下载 OllamaSetup.exe
3. 双击安装（会自动安装到后台）
4. 安装完成后，Ollama会自动运行
```

**Linux用户**：
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**Mac用户**：
```bash
# 下载 Mac 安装包
https://ollama.com/download/mac

# 或使用 Homebrew
brew install ollama
```

---

### 第2步：下载模型（5-10分钟）

打开命令行（CMD或PowerShell），执行：

```bash
# 推荐：Qwen2.5 7B（中文好，速度快）
ollama pull qwen2.5:7b
```

**其他可选模型**：

```bash
# DeepSeek R1 8B（推理能力强）
ollama pull deepseek-r1:8b

# Llama 3.1 8B（英文好）
ollama pull llama3.1:8b

# Llama 3.2 3B（更快，质量稍低）
ollama pull llama3.2:3b

# Qwen2.5 14B（质量更高，需要16GB内存）
ollama pull qwen2.5:14b
```

下载进度会显示：
```
pulling manifest
pulling 8934d96d3f08... 100% ▕████████████████▏ 4.7 GB
verifying sha256 digest
writing manifest
success
```

---

### 第3步：验证Ollama运行（1分钟）

**检查Ollama服务**：
```bash
# 测试Ollama是否运行
ollama list

# 应该看到已下载的模型
NAME              SIZE
qwen2.5:7b       4.7GB
```

**测试对话**：
```bash
ollama run qwen2.5:7b "你好"

# 应该返回中文回复
你好！有什么我可以帮助你的吗？
```

**查看API地址**：
```bash
# Ollama API 默认运行在
http://localhost:11434
```

---

## ✅ 配置已完成

您的 `.env` 文件已更新为：

```env
# Ollama本地LLM（完全免费）
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b
```

**无需API密钥！** 🎉

---

## 🧪 启动项目

```bash
# 启动服务
go run main.go

# 或使用启动脚本
.\start.bat  # Windows
./start.sh   # Linux/Mac
```

---

## 🧪 测试配置

### 测试1：健康检查
```bash
curl http://localhost:8080/health
```

### 测试2：爬取功能
```bash
curl -X POST http://localhost:8080/api/crawl/single \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://ollama.com\",\"competitor\":\"Ollama\"}"
```

---

## 🎯 推荐模型对比

### 中文任务（推荐）

| 模型 | 大小 | 内存需求 | 中文能力 | 速度 | 推荐度 |
|------|------|----------|----------|------|--------|
| **qwen2.5:7b** | 4.7GB | 8GB | ⭐⭐⭐⭐⭐ | 快 | 🔥 最推荐 |
| qwen2.5:14b | 9GB | 16GB | ⭐⭐⭐⭐⭐ | 中等 | ✅ 高质量 |
| deepseek-r1:8b | 5GB | 8GB | ⭐⭐⭐⭐ | 中等 | ✅ 推理强 |

### 英文任务

| 模型 | 大小 | 内存需求 | 英文能力 | 速度 | 推荐度 |
|------|------|----------|----------|------|--------|
| llama3.1:8b | 4.7GB | 8GB | ⭐⭐⭐⭐⭐ | 快 | ✅ 推荐 |
| llama3.2:3b | 2GB | 4GB | ⭐⭐⭐⭐ | 很快 | ✅ 快速 |

### 配置不同模型

编辑 `.env` 文件：
```env
# 使用Qwen2.5 7B（默认，推荐）
LLM_MODEL=qwen2.5:7b

# 或使用DeepSeek R1 8B（推理能力强）
LLM_MODEL=deepseek-r1:8b

# 或使用Llama 3.1 8B（英文好）
LLM_MODEL=llama3.1:8b

# 或使用Qwen2.5 14B（更高质量，需16GB内存）
LLM_MODEL=qwen2.5:14b
```

---

## 💻 硬件要求

### 最低配置
- **CPU**: 4核心
- **内存**: 8GB（运行7B模型）
- **硬盘**: 10GB可用空间

### 推荐配置
- **CPU**: 8核心+
- **内存**: 16GB（可运行14B模型）
- **硬盘**: 20GB可用空间
- **GPU**: 可选（NVIDIA显卡会更快）

### 各模型内存需求

| 模型大小 | 最低内存 | 推荐内存 |
|----------|----------|----------|
| 3B | 4GB | 8GB |
| 7B | 8GB | 16GB |
| 8B | 8GB | 16GB |
| 14B | 16GB | 32GB |

---

## ⚙️ Ollama常用命令

### 模型管理
```bash
# 列出已安装的模型
ollama list

# 下载模型
ollama pull qwen2.5:7b

# 删除模型
ollama rm qwen2.5:7b

# 更新模型
ollama pull qwen2.5:7b
```

### 运行和测试
```bash
# 运行模型（交互式）
ollama run qwen2.5:7b

# 运行单个问题
ollama run qwen2.5:7b "Hello"

# 查看模型信息
ollama show qwen2.5:7b
```

### 服务管理
```bash
# 启动Ollama服务（Windows自动启动）
ollama serve

# 停止服务（Windows任务管理器结束进程）
# Linux/Mac: killall ollama
```

---

## 🔧 常见问题

### Q1: 安装后找不到ollama命令？

**Windows**:
```bash
# Ollama安装后会自动添加到PATH
# 如果不行，重启命令行窗口
# 或手动添加: C:\Users\你的用户名\AppData\Local\Programs\Ollama
```

**Linux/Mac**:
```bash
# 检查安装
which ollama

# 如果没有，重新安装
curl -fsSL https://ollama.com/install.sh | sh
```

### Q2: 下载模型很慢？

**解决方案**:
```bash
# 使用国内镜像（如果有）
export OLLAMA_HOST=镜像地址

# 或者换个网络环境
# 或者下载较小的模型
ollama pull llama3.2:3b  # 只有2GB
```

### Q3: 内存不够，运行卡顿？

**解决方案**:
```bash
# 使用更小的模型
ollama pull llama3.2:3b  # 3B模型，需4GB内存
ollama pull qwen2.5:1.5b # 1.5B模型，需2GB内存

# 更新.env
LLM_MODEL=llama3.2:3b
```

### Q4: 想加速运行速度？

**方案1: 使用GPU**
```bash
# 如果有NVIDIA显卡，Ollama会自动使用
# 速度可提升5-10倍
```

**方案2: 使用量化模型**
```bash
# 下载量化版本（更小更快）
ollama pull qwen2.5:7b-q4  # 4-bit量化
```

### Q5: Ollama服务启动失败？

**检查端口占用**:
```bash
# Windows
netstat -ano | findstr :11434

# Linux/Mac
lsof -i :11434

# 如果被占用，修改端口
# 设置环境变量
set OLLAMA_HOST=0.0.0.0:11435  # Windows
export OLLAMA_HOST=0.0.0.0:11435  # Linux/Mac

# 同时修改.env
OPENAI_BASE_URL=http://localhost:11435
```

---

## 🔄 切换回云端API

如果想切换回Groq或其他云端API：

### 切换到Groq
```env
OPENAI_API_KEY=gsk_你的Groq密钥
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile
```

### 切换到DeepSeek
```env
OPENAI_API_KEY=sk_你的DeepSeek密钥
OPENAI_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
```

---

## 📊 性能对比

### Ollama vs 云端API

| 对比项 | Ollama本地 | Groq | DeepSeek |
|--------|-----------|------|----------|
| **费用** | 免费 | 免费 | ¥10/2500次 |
| **速度** | 中等-慢 | 超快 | 快 |
| **质量** | 良好 | 优秀 | 优秀 |
| **隐私** | 完全本地 | 云端 | 云端 |
| **限制** | 无 | 30次/分钟 | 无 |
| **网络** | 不需要 | 需要 | 需要 |

### 响应时间对比（生成100 tokens）

| 方案 | 时间 | 说明 |
|------|------|------|
| Groq (云端) | 0.5秒 | 专用硬件，超快 |
| DeepSeek | 2秒 | 云端API |
| Ollama (CPU) | 5-10秒 | 取决于CPU |
| Ollama (GPU) | 1-2秒 | 有NVIDIA显卡 |

---

## 🎯 使用建议

### 适合使用Ollama的场景
✅ 注重隐私，不想数据上传云端  
✅ 长期大量使用，节省成本  
✅ 有足够的硬件配置（8GB+内存）  
✅ 可以接受稍慢的响应速度  

### 不适合Ollama的场景
❌ 电脑配置较低（<8GB内存）  
❌ 需要极快的响应速度  
❌ 追求最高质量的分析  
❌ 只是偶尔使用  

---

## 📚 相关资源

### Ollama官方
- 官网: https://ollama.com/
- 模型库: https://ollama.com/library
- GitHub: https://github.com/ollama/ollama
- 文档: https://github.com/ollama/ollama/blob/main/docs/api.md

### 模型列表
- Qwen2.5: https://ollama.com/library/qwen2.5
- Llama 3: https://ollama.com/library/llama3.1
- DeepSeek: https://ollama.com/library/deepseek-r1

---

## 🎊 总结

### ✅ 配置完成

您已成功配置Ollama：

- ✅ Ollama已安装并运行
- ✅ 模型已下载（qwen2.5:7b）
- ✅ .env配置已更新
- ✅ 项目可以启动使用

### 📝 配置摘要

```env
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=qwen2.5:7b
```

### 🚀 下一步

```bash
# 1. 确保Ollama运行
ollama list

# 2. 启动项目
go run main.go

# 3. 开始使用
curl http://localhost:8080/health
```

**完全免费，无限使用！** 🎉

---

**更新时间**: 2026-02-06  
**Ollama版本**: 最新  
**推荐模型**: qwen2.5:7b
