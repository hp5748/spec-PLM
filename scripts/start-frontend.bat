@echo off
chcp 65001 >nul
setlocal

echo ========================================
echo    PLM 前端服务启动脚本
echo ========================================
echo.

cd /d "%~dp0..\frontend"

:: 检查Node.js是否安装
node -v >nul 2>&1
if errorlevel 1 (
    echo [错误] Node.js 未安装，请先安装 Node.js 18+
    pause
    exit /b 1
)

:: 杀掉已存在的进程
echo [信息] 检查并关闭已存在的前端进程...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":5173" ^| findstr "LISTENING"') do (
    echo [信息] 发现端口5173被进程 %%a 占用，正在关闭...
    taskkill /F /PID %%a >nul 2>&1
)

:: 等待端口释放
timeout /t 2 /nobreak >nul

:: 检查依赖是否已安装
if not exist "node_modules" (
    echo [信息] 首次运行，正在安装依赖...
    npm install
    if errorlevel 1 (
        echo [错误] 依赖安装失败
        pause
        exit /b 1
    )
)

echo.
echo [信息] 正在启动前端服务...
echo [信息] 服务地址: http://localhost:5173
echo [信息] 按 Ctrl+C 停止服务
echo.

:: 启动服务
npm run dev

pause
