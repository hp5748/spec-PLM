@echo off
chcp 65001 >nul
setlocal

echo ========================================
echo    PLM 后端服务启动脚本
echo ========================================
echo.

cd /d "%~dp0..\backend"

:: 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] Go 未安装，请先安装 Go 1.21+
    pause
    exit /b 1
)

:: 杀掉已存在的进程
echo [信息] 检查并关闭已存在的后端进程...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":8080" ^| findstr "LISTENING"') do (
    echo [信息] 发现端口8080被进程 %%a 占用，正在关闭...
    taskkill /F /PID %%a >nul 2>&1
)

:: 等待端口释放
timeout /t 2 /nobreak >nul

:: 设置Go代理（国内加速）
set GOPROXY=https://goproxy.cn,direct

:: 下载依赖
echo [信息] 检查并下载依赖...
go mod tidy

if errorlevel 1 (
    echo [错误] 依赖下载失败
    pause
    exit /b 1
)

echo.
echo [信息] 正在启动后端服务...
echo [信息] 服务地址: http://localhost:8080
echo [信息] 按 Ctrl+C 停止服务
echo.

:: 启动服务
go run cmd/server/main.go

pause
