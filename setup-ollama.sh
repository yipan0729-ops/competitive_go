#!/bin/bash

# Ollama配置助手 - 一键配置本地LLM

echo "============================================"
echo "Ollama 配置助手"
echo "============================================"
echo ""
echo "Ollama是完全免费的本地LLM服务！"
echo ""

# 检查Ollama是否已安装
if ! command -v ollama &> /dev/null; then
    echo "[步骤1/4] Ollama未安装"
    echo ""
    echo "请先安装Ollama:"
    echo "curl -fsSL https://ollama.com/install.sh | sh"
    echo ""
    echo "或访问: https://ollama.com/download"
    echo ""
    exit 1
fi

echo "[步骤1/4] ✅ Ollama已安装"
echo ""

# 检查模型
echo "[步骤2/4] 检查已安装的模型..."
ollama list

echo ""
echo "推荐模型:"
echo "1. qwen2.5:7b     - 中文好，推荐 (4.7GB)"
echo "2. llama3.1:8b    - 英文好 (4.7GB)"
echo "3. llama3.2:3b    - 更快，较小 (2GB)"
echo "4. deepseek-r1:8b - 推理强 (5GB)"
echo ""

read -p "需要下载模型吗? (y/n，默认n): " INSTALL_MODEL
if [[ "$INSTALL_MODEL" == "y" || "$INSTALL_MODEL" == "Y" ]]; then
    echo ""
    echo "选择要下载的模型:"
    echo "1. qwen2.5:7b     (推荐，4.7GB)"
    echo "2. llama3.1:8b    (4.7GB)"
    echo "3. llama3.2:3b    (2GB)"
    echo "4. deepseek-r1:8b (5GB)"
    echo ""
    read -p "请输入数字 (1-4): " MODEL_CHOICE
    
    case $MODEL_CHOICE in
        1) MODEL_NAME="qwen2.5:7b" ;;
        2) MODEL_NAME="llama3.1:8b" ;;
        3) MODEL_NAME="llama3.2:3b" ;;
        4) MODEL_NAME="deepseek-r1:8b" ;;
        *) MODEL_NAME="qwen2.5:7b" ;;
    esac
    
    echo ""
    echo "正在下载 $MODEL_NAME..."
    echo "这可能需要5-10分钟，请耐心等待..."
    ollama pull $MODEL_NAME
else
    MODEL_NAME="qwen2.5:7b"
    echo ""
    echo "使用默认模型: qwen2.5:7b"
    echo "如果未安装，请手动运行: ollama pull qwen2.5:7b"
fi

echo ""
echo "[步骤3/4] 生成配置文件..."

# 获取Serper密钥（可选）
read -p "请输入Serper API密钥（可选，回车跳过）: " SERPER_KEY

# 生成.env文件
cat > .env << EOF
# 搜索API
SERPER_API_KEY=$SERPER_KEY

# Ollama本地LLM (完全免费)
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434
LLM_MODEL=$MODEL_NAME

# LLM参数
LLM_TEMPERATURE=0.3
LLM_MAX_TOKENS=4000

# 服务器配置
SERVER_PORT=8080
GIN_MODE=release

# 数据库配置
DB_PATH=./data/competitive.db

# 存储配置
STORAGE_PATH=./storage
REPORTS_PATH=./reports

# 搜索配置
SEARCH_CACHE_DAYS=7
MAX_SEARCH_RESULTS=10
EOF

echo "✅ 配置文件已生成！"
echo ""

echo "[步骤4/4] 验证配置..."
echo ""

# 测试Ollama
echo "测试Ollama服务..."
if ! ollama list > /dev/null 2>&1; then
    echo "❌ Ollama服务未运行"
    echo "请手动启动: ollama serve"
    exit 1
fi

echo "✅ Ollama服务运行正常"
echo ""

echo "============================================"
echo "✨ 配置完成！"
echo "============================================"
echo ""
echo "配置信息:"
echo "- LLM: Ollama (本地)"
echo "- 模型: $MODEL_NAME"
echo "- API地址: http://localhost:11434"
echo "- 完全免费，无限使用！"
echo ""
echo "下一步:"
echo "1. 运行 ./start.sh 启动服务"
echo "2. 或运行 go run main.go"
echo ""
echo "提示:"
echo "- 首次使用Ollama可能较慢"
echo "- 如需更快速度，考虑使用GPU"
echo "- 可随时切换到云端API（Groq、DeepSeek）"
echo ""
