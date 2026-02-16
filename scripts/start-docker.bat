@echo off
chcp 65001 >nul
setlocal

echo ========================================
echo    PLM Docker 容器启动脚本
echo ========================================
echo.

cd /d "%~dp0.."

:: 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker 未运行，请先启动 Docker Desktop
    pause
    exit /b 1
)

echo [信息] 正在启动 Docker 容器...
docker compose up -d

if errorlevel 1 (
    echo [错误] Docker 容器启动失败
    pause
    exit /b 1
)

echo.
echo [成功] Docker 容器已启动
echo.
echo 容器状态:
docker compose ps
echo.
echo 服务地址:
echo   - MySQL:    localhost:3306  (用户: root, 密码: plm123456)
echo   - Redis:    localhost:6379
echo   - MinIO:    localhost:9000 (API) / localhost:9001 (Console)
echo              (用户: minioadmin, 密码: minioadmin)
echo.
echo 等待MySQL初始化完成...
timeout /t 10 /nobreak >nul

echo.
echo [完成] 开发环境已就绪
pause
