# 版本兼容性说明

## Go版本支持

本项目支持以下Go版本：

- ✅ Go 1.21.x
- ✅ Go 1.22.x
- ✅ Go 1.23.x
- ✅ Go 1.24.x（已在1.24.8上测试）

## 依赖版本

所有依赖都会自动下载，无需手动安装：

### 核心依赖
- `github.com/gin-gonic/gin` v1.10.0 - Web框架
- `gorm.io/gorm` v1.25.7 - ORM
- `gorm.io/driver/sqlite` v1.5.5 - SQLite驱动
- `github.com/joho/godotenv` v1.5.1 - 环境变量加载

### 可选依赖（用于Playwright）
- `github.com/playwright-community/playwright-go` v0.4201.1

### HTTP客户端
- `github.com/go-resty/resty/v2` v2.11.0

### HTML解析
- `github.com/PuerkitoBio/goquery` v1.9.0

## 安装说明

### 方法1：自动安装（推荐）

**Windows**:
```bash
.\start.bat
```

**Linux/Mac**:
```bash
chmod +x start.sh
./start.sh
```

脚本会自动：
1. 检查Go环境
2. 下载所有依赖（`go mod download`）
3. 编译项目
4. 启动服务

### 方法2：手动安装

```bash
# 1. 下载依赖
go mod download

# 2. 验证依赖
go mod tidy

# 3. 编译（可选）
go build -o competitive-analyzer

# 4. 运行
go run main.go
# 或
./competitive-analyzer
```

## 常见问题

### Q1: go mod tidy 很慢或超时

**原因**: 可能是网络问题或Go代理配置

**解决方案**:
```bash
# Windows PowerShell
$env:GOPROXY="https://goproxy.cn,direct"
go mod download

# Linux/Mac
export GOPROXY=https://goproxy.cn,direct
go mod download
```

中国大陆推荐使用的Go代理：
- https://goproxy.cn
- https://goproxy.io
- https://mirrors.aliyun.com/goproxy/

### Q2: 编译报错 "cannot find package"

**解决方案**:
```bash
# 清理缓存并重新下载
go clean -modcache
go mod download
```

### Q3: CGO相关错误（Windows）

mattn/go-sqlite3 需要CGO支持。

**解决方案**:
1. 安装MinGW或TDM-GCC
2. 确保gcc在PATH中
3. 或者使用纯Go的SQLite驱动（性能稍差）

如果不想安装GCC，可以修改 `go.mod`:
```go
// 将
gorm.io/driver/sqlite v1.5.5

// 替换为
github.com/glebarez/sqlite v1.10.0
```

### Q4: Go版本过低

如果你的Go版本 < 1.21，请升级：

**Windows**: 
- 访问 https://go.dev/dl/
- 下载并安装最新版本

**Linux**:
```bash
# 使用官方脚本
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
```

**Mac**:
```bash
brew install go
```

## 验证安装

运行测试命令验证环境：

```bash
# 检查Go版本
go version

# 验证依赖
go mod verify

# 编译测试
go build

# 运行测试
go run examples/test.go
```

## 依赖说明

### 为什么选择这些依赖？

1. **Gin** - 最流行的Go Web框架，性能优秀
2. **GORM** - 功能完整的ORM，支持多种数据库
3. **SQLite** - 零配置，适合单用户场景
4. **Resty** - 简洁的HTTP客户端
5. **Goquery** - jQuery风格的HTML解析

### 可选依赖

- **Playwright-go**: 用于浏览器自动化（暂未使用）
- 如果不需要，可以从 `go.mod` 中移除

## 更新依赖

定期更新依赖以获取安全补丁：

```bash
# 更新所有依赖到最新版本
go get -u ./...
go mod tidy

# 更新特定依赖
go get -u github.com/gin-gonic/gin
```

## 离线安装

如果需要在无网络环境安装：

```bash
# 1. 在有网络的机器上下载依赖
go mod download
go mod vendor  # 将依赖复制到vendor目录

# 2. 将整个项目（包括vendor目录）复制到目标机器

# 3. 使用vendor模式编译
go build -mod=vendor
```

## Docker环境

如果使用Docker，参考以下Dockerfile：

```dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# 安装GCC（用于SQLite）
RUN apk add --no-cache gcc musl-dev

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o competitive-analyzer

# 运行
CMD ["./competitive-analyzer"]
```

---

**更新时间**: 2026-02-06  
**支持的Go版本**: 1.21 - 1.24
