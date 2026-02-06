@echo off
REM 自动化竞品调研工具 - Windows 快速启动脚本

echo === 自动化竞品调研工具 ===
echo.

REM 检查 .env 文件
if not exist .env (
    echo ❌ 未找到 .env 文件
    echo 📝 正在从 .env.example 复制...
    copy .env.example .env
    echo ✅ 已创建 .env 文件，请编辑该文件填入你的API密钥
    echo.
    echo 需要配置的密钥：
    echo   - FIRECRAWL_API_KEY
    echo   - SERPER_API_KEY
    echo   - OPENAI_API_KEY
    echo.
    pause
    exit /b 1
)

REM 检查 Go 环境
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 未安装 Go，请先安装 Go 1.21+
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo ✅ Go 版本: %GO_VERSION%
echo.

REM 创建必要的目录
echo 📁 创建必要的目录...
if not exist data mkdir data
if not exist storage mkdir storage
if not exist reports mkdir reports
echo ✅ 目录创建完成
echo.

REM 下载依赖
echo 📦 下载依赖...
go mod download
if %errorlevel% neq 0 (
    echo ❌ 依赖下载失败
    pause
    exit /b 1
)
echo ✅ 依赖下载完成
echo.

REM 编译项目
echo 🔨 编译项目...
go build -o competitive-analyzer.exe
if %errorlevel% neq 0 (
    echo ❌ 编译失败
    pause
    exit /b 1
)
echo ✅ 编译成功
echo.

REM 启动服务
echo 🚀 启动服务...
echo.
competitive-analyzer.exe
