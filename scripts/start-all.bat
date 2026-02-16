@echo off
chcp 65001 >nul
setlocal

echo ========================================
echo    PLM 一键启动所有服务
echo ========================================
echo.

cd /d "%~dp0"

:: 启动Docker容器
echo [步骤1] 启动Docker容器...
call start-docker.bat

:: 在新窗口启动后端
echo.
echo [步骤2] 启动后端服务...
start "PLM Backend" cmd /c start-backend.bat

:: 等待后端启动
timeout /t 5 /nobreak >nul

:: 在新窗口启动前端
echo [步骤3] 启动前端服务...
start "PLM Frontend" cmd /c start-frontend.bat

echo.
echo ========================================
echo [完成] 所有服务已启动
echo ========================================
echo.
echo 服务地址:
echo   - 前端:    http://localhost:5173
echo   - 后端:    http://localhost:8080
echo   - MinIO:   http://localhost:9001
echo.
echo 默认管理员账号: admin / admin123
echo.
pause
