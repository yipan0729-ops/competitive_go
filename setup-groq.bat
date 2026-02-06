@echo off
REM Groq配置助手 - 快速配置Groq API

echo ============================================
echo Groq 配置助手
echo ============================================
echo.
echo Groq是完全免费的LLM服务，速度超快！
echo.

REM 检查是否已有.env文件
if exist .env (
    echo 发现已有.env文件
    choice /C YN /M "是否覆盖现有配置"
    if errorlevel 2 goto :END
)

echo.
echo 请按照以下步骤操作：
echo.
echo 第1步：获取Groq API密钥
echo ----------------------------------------
echo 1. 访问: https://console.groq.com/
echo 2. 注册/登录账号
echo 3. 点击 "API Keys"
echo 4. 创建新密钥并复制
echo.
pause

echo.
echo 第2步：获取Serper API密钥（用于搜索）
echo ----------------------------------------
echo 1. 访问: https://serper.dev/
echo 2. 注册/登录账号
echo 3. 复制 API Key
echo.
pause

echo.
echo 第3步：输入API密钥
echo ----------------------------------------

REM 输入Groq密钥
set /p GROQ_KEY="请输入Groq API密钥 (gsk_xxx): "
if "%GROQ_KEY%"=="" (
    echo 错误: 密钥不能为空
    pause
    goto :END
)

REM 输入Serper密钥
set /p SERPER_KEY="请输入Serper API密钥 (可选，回车跳过): "

echo.
echo 第4步：生成配置文件
echo ----------------------------------------

REM 生成.env文件
(
echo # API Keys
echo FIRECRAWL_API_KEY=your_firecrawl_key_here
echo SERPER_API_KEY=%SERPER_KEY%
echo.
echo # Groq LLM ^(完全免费^)
echo OPENAI_API_KEY=%GROQ_KEY%
echo OPENAI_BASE_URL=https://api.groq.com/openai
echo LLM_MODEL=llama-3.3-70b-versatile
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

echo.
echo ✅ 配置文件已生成！
echo.
echo 配置信息：
echo - Groq API: %GROQ_KEY:~0,20%...
echo - 模型: llama-3.3-70b-versatile
echo - 完全免费，速度超快！
echo.
echo 下一步：
echo 1. 运行 start.bat 启动服务
echo 2. 或运行 go run main.go
echo.

:END
pause
