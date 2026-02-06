#!/bin/bash

# 自动化竞品调研工具 - 快速启动脚本

echo "=== 自动化竞品调研工具 ==="
echo ""

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "❌ 未找到 .env 文件"
    echo "📝 正在从 .env.example 复制..."
    cp .env.example .env
    echo "✅ 已创建 .env 文件，请编辑该文件填入你的API密钥"
    echo ""
    echo "需要配置的密钥："
    echo "  - FIRECRAWL_API_KEY"
    echo "  - SERPER_API_KEY"
    echo "  - OPENAI_API_KEY"
    echo ""
    exit 1
fi

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 未安装 Go，请先安装 Go 1.21+"
    exit 1
fi

echo "✅ Go 版本: $(go version)"
echo ""

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p data storage reports
echo "✅ 目录创建完成"
echo ""

# 下载依赖
echo "📦 下载依赖..."
go mod download
echo "✅ 依赖下载完成"
echo ""

# 编译项目
echo "🔨 编译项目..."
go build -o competitive-analyzer
if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi
echo "✅ 编译成功"
echo ""

# 启动服务
echo "🚀 启动服务..."
echo ""
./competitive-analyzer
