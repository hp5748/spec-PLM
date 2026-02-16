# PLM 产品生命周期管理系统

## 环境要求

- Node.js 18+
- Go 1.21+
- Docker Desktop

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd PLM
```

### 2. 一键启动（Windows）

双击运行 `scripts/start-all.bat`，将自动启动：
- Docker 容器（MySQL、Redis、MinIO）
- 后端服务
- 前端服务

### 3. 分步启动

```bash
# 启动 Docker 基础服务
scripts/start-docker.bat

# 启动后端（新终端）
scripts/start-backend.bat

# 启动前端（新终端）
scripts/start-frontend.bat
```

## 服务地址

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:5173 |
| 后端 API | http://localhost:8080 |
| MinIO 控制台 | http://localhost:9001 |

## 默认账号

- 管理员：admin / admin123
- MinIO：minioadmin / minioadmin
- MySQL：root / plm123456
