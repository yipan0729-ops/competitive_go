# 🚀 快速使用指南

## ✅ 服务已启动并修复

当前服务版本：**v2.0.1（已修复）**  
服务地址：`http://localhost:8080`  
状态：✅ 运行中

---

## 📊 启动新任务（推荐）

之前的任务4因为是旧版本处理的，所以没有生成报告。  
现在使用修复后的版本启动新任务：

### 方式1：快速测试（2分钟）

```powershell
# 启动新任务
$body = @{
    topic = "在线绘图工具"
    competitor_count = 2
    depth = "quick"
} | ConvertTo-Json

$response = Invoke-WebRequest `
    -Uri http://localhost:8080/api/auto/analysis `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

$result = $response.Content | ConvertFrom-Json
$taskId = $result.task_id

Write-Host "✅ 新任务已启动，ID: $taskId"
Write-Host ""
Write-Host "监控进度（每10秒检查一次）："
Write-Host ""

# 监控循环
for ($i = 0; $i -lt 20; $i++) {
    Start-Sleep -Seconds 10
    
    $status = Invoke-WebRequest `
        -Uri "http://localhost:8080/api/discover/status/$taskId" `
        -UseBasicParsing |
        ConvertFrom-Json
    
    Write-Host "[$i] 进度: $($status.progress)% | 状态: $($status.status)" -ForegroundColor Cyan
    
    if ($status.status -eq "completed") {
        Write-Host ""
        Write-Host "🎉 任务完成！" -ForegroundColor Green
        Write-Host "   爬取: $($status.result.urls_crawled) 个URL" -ForegroundColor White
        Write-Host "   分析: $($status.result.analyzed_count) 个竞品" -ForegroundColor White
        Write-Host "   报告: $($status.result.report_path)" -ForegroundColor Yellow
        Write-Host ""
        
        # 如果有报告，打开它
        if ($status.result.report_path -and (Test-Path $status.result.report_path)) {
            Write-Host "📖 打开报告..." -ForegroundColor Cyan
            Start-Process $status.result.report_path
        } else {
            Write-Host "⚠️ 报告路径为空或文件不存在" -ForegroundColor Yellow
        }
        
        break
    }
    
    if ($status.status -eq "failed") {
        Write-Host "❌ 任务失败" -ForegroundColor Red
        break
    }
}
```

### 方式2：使用测试脚本

```powershell
.\test-auto-workflow.ps1
```

---

## 🔍 检查修复是否生效

新任务应该会看到以下日志（说明修复生效）：

```
[自动化] 使用默认配置：启用爬取、分析和报告生成  ← 这行是新的！
[自动化] 开始执行任务 #5
[自动化] 发现竞品: [...]
[自动化] 找到数据源: X 个
[自动化] 开始爬取 X 个URL              ← 会执行爬取
[爬取成功: ...]                        ← 爬取日志
[自动化] 开始AI分析                    ← 会执行分析  
[自动化] 生成报告                      ← 会生成报告
[自动化] 任务完成 #5
```

---

## 📁 文件位置

修复后，文件会保存在：

### 报告
```
D:\Code\Competitive_go\reports\
  ├─ {主题}_自动分析报告_{时间}.md
  └─ ...
```

### 爬取的内容
```
D:\Code\Competitive_go\storage\crawled\
  ├─ 竞品1\
  │   ├─ 官网.md
  │   └─ ...
  ├─ 竞品2\
  └─ ...
```

---

## 🎯 验证修复的步骤

1. **启动新任务**（任务ID会是5或更大）
2. **等待完成**（2-3分钟）
3. **检查结果**：
   ```powershell
   # 查看报告目录
   ls D:\Code\Competitive_go\reports\
   
   # 查看爬取内容
   ls D:\Code\Competitive_go\storage\crawled\
   
   # 查看任务详情
   Invoke-WebRequest -Uri http://localhost:8080/api/discover/status/5 -UseBasicParsing
   ```

4. **预期结果**：
   - ✅ `analyzed_count` > 0
   - ✅ `report_path` 不为空
   - ✅ `reports/` 目录有文件
   - ✅ `storage/crawled/` 目录有文件

---

## 🛑 停止服务

### 方法1：Ctrl+C
在运行程序的终端按 `Ctrl + C`

### 方法2：命令停止
```powershell
# 查找进程
$pid = netstat -ano | findstr :8080 | ForEach-Object { 
    if ($_ -match '\s+(\d+)$') { $matches[1] } 
} | Select-Object -First 1

# 停止
if ($pid) { taskkill /F /PID $pid }
```

---

## 🔄 重启服务

```powershell
# 1. 停止（如果在运行）
$pid = netstat -ano | findstr :8080 | ForEach-Object { 
    if ($_ -match '\s+(\d+)$') { $matches[1] } 
} | Select-Object -First 1
if ($pid) { taskkill /F /PID $pid }

# 2. 等待1秒
Start-Sleep -Seconds 1

# 3. 启动
.\competitive-analyzer.exe
```

---

## 💡 常见问题

### Q: 为什么任务4没有报告？
**A**: 任务4是在修复前启动的，所以没有执行爬取和分析。请启动新任务（任务5+）。

### Q: 如何确认修复生效？
**A**: 查看日志，应该看到 `[自动化] 使用默认配置：启用爬取、分析和报告生成` 这行。

### Q: 报告保存在哪？
**A**: `D:\Code\Competitive_go\reports\` 目录。

### Q: 可以查看之前的任务吗？
**A**: 可以，但任务4及之前的任务因为是旧版本，不会有完整结果。

---

## 📚 相关文档

- [完整API文档](./API.md)
- [问题修复说明](./BUGFIX_AUTO_ANALYSIS.md)
- [新功能指南](./NEW_FEATURES.md)
- [快速开始](./QUICKSTART_V2.md)

---

## 🎉 现在开始测试！

复制上面的PowerShell代码，启动一个新任务，这次应该能看到完整的报告了！

**服务状态**: ✅ 运行中  
**版本**: v2.0.1（已修复）  
**端口**: 8080
