@echo off
REM Ollama配置助手 - 一键配置本地LLM

echo ============================================
echo Ollama 配置助手
echo ============================================
echo.
echo Ollama是完全免费的本地LLM服务！
echo.

REM 检查Ollama是否已安装
where ollama >nul 2>nul
if %errorlevel% neq 0 (
    echo [步骤1/4] Ollama未安装
    echo.
    echo 请先安装Ollama:
    echo 1. 访问: https://ollama.com/download/windows
    echo 2. 下载并安装 OllamaSetup.exe
    echo 3. 安装完成后重新运行此脚本
    echo.
    pause
    exit /b 1
)

echo [步骤1/4] ✅ Ollama已安装
echo.

REM 检查模型
echo [步骤2/4] 检查已安装的模型...
ollama list

echo.
echo 推荐模型:
echo 1. qwen2.5:7b     - 中文好，推荐 (4.7GB)
echo 2. llama3.1:8b    - 英文好 (4.7GB)
echo 3. llama3.2:3b    - 更快，较小 (2GB)
echo 4. deepseek-r1:8b - 推理强 (5GB)
echo.

set /p INSTALL_MODEL="需要下载模型吗? (y/n，默认n): "
if /i "%INSTALL_MODEL%"=="y" (
    echo.
    echo 选择要下载的模型:
    echo 1. qwen2.5:7b     (推荐，4.7GB)
    echo 2. llama3.1:8b    (4.7GB)
    echo 3. llama3.2:3b    (2GB)
    echo 4. deepseek-r1:8b (5GB)
    echo.
    set /p MODEL_CHOICE="请输入数字 (1-4): "
    
    if "%MODEL_CHOICE%"=="1" set MODEL_NAME=qwen2.5:7b
    if "%MODEL_CHOICE%"=="2" set MODEL_NAME=llama3.1:8b
    if "%MODEL_CHOICE%"=="3" set MODEL_NAME=llama3.2:3b
    if "%MODEL_CHOICE%"=="4" set MODEL_NAME=deepseek-r1:8b
    
    echo.
    echo 正在下载 %MODEL_NAME%...
    echo 这可能需要5-10分钟，请耐心等待...
    ollama pull %MODEL_NAME%
) else (
    set MODEL_NAME=qwen2.5:7b
    echo.
    echo 使用默认模型: qwen2.5:7b
    echo 如果未安装，请手动运行: ollama pull qwen2.5:7b
)

echo.
echo [步骤3/4] 生成配置文件...

REM 获取Serper密钥（可选）
set /p SERPER_KEY="请输入Serper API密钥（可选，回车跳过）: "

REM 生成.env文件
(
echo # 搜索API
echo SERPER_API_KEY=%SERPER_KEY%
echo.
echo # Ollama本地LLM ^(完全免费^)
echo OPENAI_API_KEY=ollama
echo OPENAI_BASE_URL=http://localhost:11434
echo LLM_MODEL=%MODEL_NAME%
echo.
echo # LLM参数
echo LLM_TEMPERATURE=0.3
echo LLM_MAX_TOKENS=4000
echo.
echo # 服务器配置
echo SERVER_PORT=8080
echo GIN_MODE=release
echo.
echo # 数据库配置
echo DB_PATH=./data/competitive.db
echo.
echo # 存储配置
echo STORAGE_PATH=./storage
echo REPORTS_PATH=./reports
echo.
echo # 搜索配置
echo SEARCH_CACHE_DAYS=7
echo MAX_SEARCH_RESULTS=10
) > .env

echo ✅ 配置文件已生成！
echo.

echo [步骤4/4] 验证配置...
echo.

REM 测试Ollama
echo 测试Ollama服务...
ollama list
if %errorlevel% neq 0 (
    echo ❌ Ollama服务未运行
    echo 请手动启动: ollama serve
    pause
    exit /b 1
)

echo ✅ Ollama服务运行正常
echo.

echo ============================================
echo ✨ 配置完成！
echo ============================================
echo.
echo 配置信息:
echo - LLM: Ollama (本地)
echo - 模型: %MODEL_NAME%
echo - API地址: http://localhost:11434
echo - 完全免费，无限使用！
echo.
echo 下一步:
echo 1. 运行 start.bat 启动服务
echo 2. 或运行 go run main.go
echo.
echo 提示:
echo - 首次使用Ollama可能较慢
echo - 如需更快速度，考虑使用GPU
echo - 可随时切换到云端API（Groq、DeepSeek）
echo.

pause
