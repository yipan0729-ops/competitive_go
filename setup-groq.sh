#!/bin/bash

# Groq配置助手 - 快速配置Groq API

echo "============================================"
echo "Groq 配置助手"
echo "============================================"
echo ""
echo "Groq是完全免费的LLM服务，速度超快！"
echo ""

# 检查是否已有.env文件
if [ -f .env ]; then
    echo "发现已有.env文件"
    read -p "是否覆盖现有配置? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

echo ""
echo "请按照以下步骤操作："
echo ""
echo "第1步：获取Groq API密钥"
echo "----------------------------------------"
echo "1. 访问: https://console.groq.com/"
echo "2. 注册/登录账号"
echo "3. 点击 'API Keys'"
echo "4. 创建新密钥并复制"
echo ""
read -p "按回车继续..."

echo ""
echo "第2步：获取Serper API密钥（用于搜索）"
echo "----------------------------------------"
echo "1. 访问: https://serper.dev/"
echo "2. 注册/登录账号"
echo "3. 复制 API Key"
echo ""
read -p "按回车继续..."

echo ""
echo "第3步：输入API密钥"
echo "----------------------------------------"

# 输入Groq密钥
read -p "请输入Groq API密钥 (gsk_xxx): " GROQ_KEY
if [ -z "$GROQ_KEY" ]; then
    echo "错误: 密钥不能为空"
    exit 1
fi

# 输入Serper密钥
read -p "请输入Serper API密钥 (可选，回车跳过): " SERPER_KEY

echo ""
echo "第4步：生成配置文件"
echo "----------------------------------------"

# 生成.env文件
cat > .env << EOF
# API Keys
FIRECRAWL_API_KEY=your_firecrawl_key_here
SERPER_API_KEY=$SERPER_KEY

# Groq LLM (完全免费)
OPENAI_API_KEY=$GROQ_KEY
OPENAI_BASE_URL=https://api.groq.com/openai
LLM_MODEL=llama-3.3-70b-versatile

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

echo ""
echo "✅ 配置文件已生成！"
echo ""
echo "配置信息："
echo "- Groq API: ${GROQ_KEY:0:20}..."
echo "- 模型: llama-3.3-70b-versatile"
echo "- 完全免费，速度超快！"
echo ""
echo "下一步："
echo "1. 运行 ./start.sh 启动服务"
echo "2. 或运行 go run main.go"
echo ""
